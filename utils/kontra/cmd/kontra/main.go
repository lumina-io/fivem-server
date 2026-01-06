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
	"syscall"

	"github.com/lumina-io/kontra/internal/logging"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration constants
const (
	// Log rotation configuration
	maxLogSizeMB  = 1    // 1MB max log file size
	maxLogBackups = 10   // Keep 10 backup files
	maxLogAgeDays = 30   // Keep logs for 30 days
	compressLogs  = true // Compress old logs
)

// forwardOutput efficiently processes and logs output from a reader
// Memory is released immediately after each line is logged
func forwardOutput(reader io.Reader, logger *slog.Logger, isError bool) {
	scanner := bufio.NewScanner(reader)

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
	for {
		if cmd.Process != nil {
			sig := <-sigChan
			logger.Warn(fmt.Sprintf("Signal received: %s", sig.String()))
			logger.Info("Forwarding signal to child process")

			if err := cmd.Process.Signal(sig); err != nil {
				logger.Error(fmt.Sprintf("Failed to forward signal: %v", err))
				cmd.Process.Kill()
				return
			}
		}
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

	// Executor
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	defer stderr.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	done := make(chan bool, 2)

	go func() {
		forwardOutput(stdout, logger, false)
		done <- true
	}()

	go func() {
		forwardOutput(stderr, logger, true)
		done <- true
	}()

	go forwardInput(stdin)

	// Signal Handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		handleProcessCompletion(cmd, sigChan, done, logger)
	}()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Wait stdout/stderr
	<-done
	<-done

	err = cmd.Wait()

	signal.Stop(sigChan)
	close(sigChan)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			logger.Warn(fmt.Sprintf("Process exited: %d", exitCode))
			os.Exit(exitCode)
		}
		logger.Error(fmt.Sprintf("Command execution failed: %v", err))
		os.Exit(1)
	}
}
