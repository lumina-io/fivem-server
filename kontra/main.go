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

	"github.com/lumina-io/kontra/internal/logging"
)

func forwardOutput(scanner *bufio.Scanner, logFunc func(string), wg *sync.WaitGroup) {
	defer wg.Done()
	for scanner.Scan() {
		logFunc(scanner.Text())
	}
}

func forwardInput(stdin io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(stdin, "%s\n", scanner.Text())
	}
}

func main() {
	handler := logging.NewSimpleHandler(os.Stdout, slog.LevelInfo)
	logger := slog.New(handler)

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

	var wg sync.WaitGroup
	wg.Add(2)

	go forwardOutput(bufio.NewScanner(stdout), func(msg string) { logger.Info(msg) }, &wg)
	go forwardOutput(bufio.NewScanner(stderr), func(msg string) { logger.Error(msg) }, &wg)

	go forwardInput(stdin)

	// プロセス終了またはシグナル受信を待機
	done := make(chan bool, 1)

	// プロセス終了を監視
	go func() {
		cmd.Wait()
		done <- true
	}()

	// シグナルまたはプロセス終了を待つ
	select {
	case sig := <-sigChan:
		logger.Warn(fmt.Sprintf("Signal received: %s", sig.String()))
		logger.Info("Forwarding signal to child process")

		// 子プロセスにシグナルを転送
		if err := cmd.Process.Signal(sig); err != nil {
			logger.Error(fmt.Sprintf("Failed to forward signal: %v", err))
			// シグナル転送に失敗した場合は強制終了
			cmd.Process.Kill()
		}

		// 子プロセスの終了を待つ
		<-done

	case <-done:
		logger.Debug("Process finished normally")
	}

	// 出力処理完了まで待機
	wg.Wait()

	exitCode := cmd.ProcessState.ExitCode()
	logger.Warn(fmt.Sprintf("Process exited: %d", exitCode))
	os.Exit(exitCode)
}
