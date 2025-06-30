package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lumina-io/kontra/internal/logging"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logMessage struct {
	message string
	isError bool
}

func forwardOutput(scanner *bufio.Scanner, ch chan<- logMessage, isError bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for scanner.Scan() {
		select {
		case ch <- logMessage{
			message: scanner.Text(),
			isError: isError,
		}:
		default:
			return
		}
	}
}

func logProcessor(ch <-chan logMessage, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range ch {
		if msg.isError {
			logger.Error(msg.message)
		} else {
			logger.Info(msg.message)
		}
	}
}

func forwardInput(stdin io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(stdin, "%s\n", scanner.Text())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kontra <command> [args...]")
		os.Exit(1)
	}

	logRotator := &lumberjack.Logger{
		Filename:   "./logs/server.log",
		MaxSize:    1,    // Max size in MB
		MaxBackups: 10,   // Number of backups
		MaxAge:     30,   // Days
		Compress:   true, // Enable compression
	}

	logger := slog.New(slogmulti.Fanout(
		logging.NewSimpleHandler(os.Stdout),
		logging.NewFileHandler(logRotator, slog.LevelInfo),
	))

	// Signal Handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// Executor
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdin, _ := cmd.StdinPipe()

	defer stdout.Close()
	defer stderr.Close()
	defer stdin.Close()

	cmd.Start()

	logCh := make(chan logMessage, 100)
	done := make(chan bool, 1)

	var outputWg sync.WaitGroup
	var logWg sync.WaitGroup

	outputWg.Add(2)
	logWg.Add(1)

	go forwardOutput(bufio.NewScanner(stdout), logCh, false, &outputWg)
	go forwardOutput(bufio.NewScanner(stderr), logCh, true, &outputWg)
	go logProcessor(logCh, logger, &logWg)

	go forwardInput(stdin)

	go func() {
		cmd.Wait()
		done <- true
	}()

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

	outputDone := make(chan struct{})
	go func() {
		outputWg.Wait()
		close(outputDone)
	}()

	select {
	case <-outputDone:
		// No any actions
	case <-time.After(5 * time.Second):
		logger.Warn("Output forwarding timeout")
	}

	close(logCh)
	logWg.Wait()

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		logger.Warn(fmt.Sprintf("Process exited: %d", exitCode))
	} else {
		logger.Debug("Process exited normally")
	}

	os.Exit(exitCode)
}
