package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// プロセス情報を格納する構造体
type ProcessInfo struct {
	PID  int
	Name string
	Cmd  string
}

// 実行中のプロセス一覧を取得
func getRunningProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	// /procディレクトリから数字のディレクトリ（PID）を取得
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// ディレクトリ名がPIDかチェック
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		// プロセス名を取得
		cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
		cmdlineBytes, err := ioutil.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		cmdline := string(cmdlineBytes)
		cmdline = strings.ReplaceAll(cmdline, "\x00", " ")
		cmdline = strings.TrimSpace(cmdline)

		if cmdline == "" {
			continue
		}

		// プロセス名（最初の引数）を抽出
		parts := strings.Fields(cmdline)
		if len(parts) == 0 {
			continue
		}

		processName := parts[0]
		// パスからファイル名のみを抽出
		if lastSlash := strings.LastIndex(processName, "/"); lastSlash != -1 {
			processName = processName[lastSlash+1:]
		}

		processes = append(processes, ProcessInfo{
			PID:  pid,
			Name: processName,
			Cmd:  cmdline,
		})
	}

	return processes, nil
}

// プロセス名でプロセスを検索
func findProcessByName(name string) ([]ProcessInfo, error) {
	processes, err := getRunningProcesses()
	if err != nil {
		return nil, err
	}

	var matches []ProcessInfo
	for _, proc := range processes {
		if strings.Contains(strings.ToLower(proc.Name), strings.ToLower(name)) {
			matches = append(matches, proc)
		}
	}

	return matches, nil
}

// 指定されたPIDのプロセスに標準入力を送信
func sendInputToProcess(pid int, input string) error {
	// 方法1: /proc/PID/fd/0 に直接書き込み（権限が必要）
	stdinPath := fmt.Sprintf("/proc/%d/fd/0", pid)

	// ファイルが存在するかチェック
	if _, err := os.Stat(stdinPath); os.IsNotExist(err) {
		return fmt.Errorf("プロセス %d の標準入力が見つかりません", pid)
	}

	// 標準入力に書き込み
	file, err := os.OpenFile(stdinPath, os.O_WRONLY, 0)
	if err != nil {
		// 権限がない場合は、echoコマンドとパイプを使用
		return sendInputUsingEcho(pid, input)
	}
	defer file.Close()

	_, err = file.WriteString(input + "\n")
	return err
}

// echoコマンドとパイプを使用してプロセスに入力を送信
func sendInputUsingEcho(pid int, input string) error {
	// プロセスの標準入力ファイルディスクリプタを見つける
	stdinPath := fmt.Sprintf("/proc/%d/fd/0", pid)

	// echoコマンドを使って入力を送信
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -ne '%s\n\n' > %s", input, stdinPath))
	return cmd.Run()
}

// シグナルを使ってプロセスに入力を送信（代替方法）
func sendSignalToProcess(pid int, signal syscall.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Signal(signal)
}

func main() {
	fmt.Println("=== 既存プロセスへの標準入力送信ツール ===")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n選択してください:")
		fmt.Println("1. プロセス名で検索")
		fmt.Println("2. PIDを直接指定")
		fmt.Println("3. 全プロセス一覧表示")
		fmt.Println("4. 終了")
		fmt.Print("選択 (1-4): ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("プロセス名を入力してください: ")
			if !scanner.Scan() {
				continue
			}
			processName := strings.TrimSpace(scanner.Text())

			processes, err := findProcessByName(processName)
			if err != nil {
				fmt.Printf("エラー: %v\n", err)
				continue
			}

			if len(processes) == 0 {
				fmt.Println("該当するプロセスが見つかりませんでした")
				continue
			}

			fmt.Println("見つかったプロセス:")
			for i, proc := range processes {
				fmt.Printf("%d. PID: %d, 名前: %s, コマンド: %s\n", i+1, proc.PID, proc.Name, proc.Cmd)
			}

			fmt.Print("選択番号を入力してください: ")
			if !scanner.Scan() {
				continue
			}

			selection, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || selection < 1 || selection > len(processes) {
				fmt.Println("無効な選択です")
				continue
			}

			selectedProcess := processes[selection-1]

			fmt.Print("送信する文字列を入力してください: ")
			if !scanner.Scan() {
				continue
			}
			input := scanner.Text()

			fmt.Printf("PID %d のプロセスに '%s' を送信しています...\n", selectedProcess.PID, input)
			if err := sendInputToProcess(selectedProcess.PID, input); err != nil {
				fmt.Printf("エラー: %v\n", err)
			} else {
				fmt.Println("入力を送信しました")
			}

		case "2":
			fmt.Print("PIDを入力してください: ")
			if !scanner.Scan() {
				continue
			}

			pid, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("無効なPIDです")
				continue
			}

			fmt.Print("送信する文字列を入力してください: ")
			if !scanner.Scan() {
				continue
			}
			input := scanner.Text()

			fmt.Printf("PID %d のプロセスに '%s' を送信しています...\n", pid, input)
			if err := sendInputToProcess(pid, input); err != nil {
				fmt.Printf("エラー: %v\n", err)
			} else {
				fmt.Println("入力を送信しました")
			}

		case "3":
			processes, err := getRunningProcesses()
			if err != nil {
				fmt.Printf("エラー: %v\n", err)
				continue
			}

			fmt.Printf("実行中のプロセス（最初の20個）:\n")
			count := 0
			for _, proc := range processes {
				if count >= 20 {
					fmt.Println("... (さらに多くのプロセスがあります)")
					break
				}
				fmt.Printf("PID: %d, 名前: %s, コマンド: %s\n", proc.PID, proc.Name, proc.Cmd)
				count++
			}

		case "4":
			fmt.Println("終了します")
			return

		default:
			fmt.Println("無効な選択です")
		}
	}
}
