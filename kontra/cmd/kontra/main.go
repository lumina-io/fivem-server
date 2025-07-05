package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lumina-io/kontra/internal/logging"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration constants
const (
	// Buffer configuration
	initialBufSize = 64 * 1024        // 64KB initial buffer - efficient for most logs
	maxTokenSize   = 10 * 1024 * 1024 // 10MB max token size - handles large lines

	// Timeout configuration
	outputProcessingTimeout = 10 * time.Second // Max wait for output processing

	// Log rotation configuration
	maxLogSizeMB  = 1    // 1MB max log file size
	maxLogBackups = 10   // Keep 10 backup files
	maxLogAgeDays = 30   // Keep logs for 30 days
	compressLogs  = true // Compress old logs
)

// forwardOutput efficiently processes and logs output from a reader
// Memory is released immediately after each line is logged
func forwardOutput(reader io.Reader, logger *slog.Logger, isError bool, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, initialBufSize)
	scanner.Buffer(buf, maxTokenSize)

	for scanner.Scan() {
		lineText := scanner.Text()

		if isError {
			logger.Error(lineText)
		} else {
			logger.Info(lineText)
		}
	}

	// Log scanner errors, ignoring expected ones
	if err := scanner.Err(); err != nil && !isFileClosedError(err) {
		logger.Error(fmt.Sprintf("Scanner error: %v", err))
	}
}

// isFileClosedError checks if the error is an expected file closure
func isFileClosedError(err error) bool {
	return strings.Contains(err.Error(), "file already closed")
}

// handleProcessCompletion manages process termination and signal forwarding
func handleProcessCompletion(cmd *exec.Cmd, sigChan <-chan os.Signal, done <-chan bool, logger *slog.Logger) {
	select {
	case sig := <-sigChan:
		logger.Warn(fmt.Sprintf("Signal received: %s", sig.String()))
		logger.Info("Forwarding signal to child process")

		if err := cmd.Process.Signal(sig); err != nil {
			logger.Error(fmt.Sprintf("Failed to forward signal: %v", err))
			cmd.Process.Kill()
		}
		<-done

	case <-done:
		logger.Debug("Process finished normally")
	}
}

// waitForOutputCompletion waits for all output processing with timeout
func waitForOutputCompletion(outputWg *sync.WaitGroup, logger *slog.Logger) {
	outputDone := make(chan struct{})
	go func() {
		outputWg.Wait()
		close(outputDone)
	}()

	select {
	case <-outputDone:
		// All output processed successfully
	case <-time.After(outputProcessingTimeout):
		logger.Warn("Output processing timeout - some output may be lost")
	}
}

func forwardInput(stdin io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(stdin, "%s\n", scanner.Text())
	}
}

// createLogger sets up structured logging with rotation
func createLogger() *slog.Logger {
	logRotator := &lumberjack.Logger{
		Filename:   "./logs/server.log",
		MaxSize:    maxLogSizeMB,
		MaxBackups: maxLogBackups,
		MaxAge:     maxLogAgeDays,
		Compress:   compressLogs,
	}

	return slog.New(slogmulti.Fanout(
		logging.NewSimpleHandler(os.Stdout),
		logging.NewFileHandler(logRotator, slog.LevelInfo),
	))
}

// validateArgs checks command line arguments
func validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kontra <command> [args...]")
		os.Exit(1)
	}
}

func main() {
	validateArgs()
	logger := createLogger()

	// Signal Handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// Executor
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	defer stdin.Close()

	done := make(chan bool, 1)
	var outputWg sync.WaitGroup

	// Set up goroutines before starting the process
	outputWg.Add(2)
	go forwardOutput(stdout, logger, false, &outputWg)
	go forwardOutput(stderr, logger, true, &outputWg)
	go forwardInput(stdin)

	// Start the process after goroutines are ready
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		cmd.Wait()
		done <- true
	}()

	handleProcessCompletion(cmd, sigChan, done, logger)

	waitForOutputCompletion(&outputWg, logger)

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		logger.Warn(fmt.Sprintf("Process exited: %d", exitCode))
	} else {
		logger.Debug("Process exited normally")
	}

	os.Exit(exitCode)
}
