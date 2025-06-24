package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func StdPrint(args ...string) {
	currentTime := time.Now()
	d := currentTime.Format(time.RFC3339)
	fmt.Println(color.CyanString(d), strings.Join(args, " "))
}

func ErrPrint(args ...string) {
	currentTime := time.Now()
	d := currentTime.Format(time.RFC3339)
	fmt.Println(color.RedString(d), strings.Join(args, " "))
}

func WarnPrint(args ...string) {
	currentTime := time.Now()
	d := currentTime.Format(time.RFC3339)
	fmt.Println(color.YellowString(d), strings.Join(args, " "))
}

func reader(scanner *bufio.Scanner, channel chan string) {
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		channel <- scanner.Text()
	}
}

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	defer stdin.Close()
	defer stdout.Close()
	defer stderr.Close()

	cmd.Start()

	inputScanner := bufio.NewScanner(os.Stdin)
	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	stdinChan := make(chan string)
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	go reader(inputScanner, stdinChan)
	go reader(outScanner, stdoutChan)
	go reader(errScanner, stderrChan)

	finished := false
	done := make(chan bool)
	go func() {
		for !finished {
			select {
			case line := <-stdinChan:
				io.WriteString(stdin, fmt.Sprintf("%v\n", line))
			case line := <-stdoutChan:
				StdPrint(line)
			case line := <-stderrChan:
				ErrPrint(line)
			case receivedSignal := <-sig:
				WarnPrint(fmt.Sprintf("Signal Received: %v\n", receivedSignal.String()))
				cmd.Process.Signal(receivedSignal)
			case <-done:
				finished = true
			}
		}
	}()

	cmd.Wait()
	done <- true

	exitCode := cmd.ProcessState.ExitCode()
	WarnPrint(fmt.Sprintf("Process exited: %d\n", exitCode))
	os.Exit(exitCode)
}
