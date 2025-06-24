package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()
    cmd.Start()

    outScanner := bufio.NewScanner(stdout)
    errScanner := bufio.NewScanner(stderr)

    outScanner.Split(bufio.ScanLines)
	errScanner.Split(bufio.ScanLines)

    for outScanner.Scan() {
		currentTime := time.Now()
		d := currentTime.Format(time.RFC3339)
        m := outScanner.Text()
		fmt.Printf("%v %v\n", color.CyanString(fmt.Sprintf("[%v]", d)), m)
    }

	for errScanner.Scan() {
		currentTime := time.Now()
		d := currentTime.Format(time.RFC3339)
        m := errScanner.Text()
		fmt.Printf("%v %v\n", color.RedString(fmt.Sprintf("[%v]", d)), m)
	}
    cmd.Wait()

	exitCode := cmd.ProcessState.ExitCode()
	os.Exit(exitCode)
}
