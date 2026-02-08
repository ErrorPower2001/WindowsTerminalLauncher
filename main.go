package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"time"

	// Added to detect OS for clear screen
	// 用于检测系统以清屏
	"runtime"

	"golang.org/x/term"
)


// WTConfig matches Windows Terminal settings.json structure
// WTConfig 对应 Windows Terminal 的 settings.json 结构
// {{{
type WTConfig struct {
	Profiles struct {
		List []struct {
			Name   string `json:"name"`   // Profile name / 配置项名称
			Hidden bool   `json:"hidden"` // Whether it is hidden / 是否隐藏
		} `json:"list"`
	} `json:"profiles"`
}
// }}}


// clearScreen cleans the terminal display
// clearScreen 清除终端屏幕内容
// {{{
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
// }}}


func main() {
	// Get config path and read file
	// 获取配置文件路径并读取文件
	// {{{
	configPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		`Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json`)

	data, err := os.ReadFile(configPath)
	if err != nil {
		////fmt.Printf("\x1b[1;31mFaild to read config file: %v\x1b[0m\n", err)
		fmt.Print("\x1b[1;31m")
		fmt.Print("Faild to read config file: ")
		fmt.Print("\x1b[0m")
		fmt.Print("\n")
		fmt.Print(err)
		fmt.Print("\n")
		return
	}
	// }}}


	// Parse JSON and filter profiles
	// 解析 JSON 并过滤配置项
	// {{{
	var config WTConfig
	if err := json.Unmarshal(data, &config); err != nil {
		////fmt.Printf("\x1b[1;31mFaild to analysis JSON: %v\x1b[0m\n", err)
		fmt.Print("\x1b[1;31m")
		fmt.Print("Faild to analysis JSON: ")
		fmt.Print("\x1b[0m")
		fmt.Print("\n")
		fmt.Print(err)
		fmt.Print("\n")
		return
	}

	var names []string
	for _, p := range config.Profiles.List {
		if !p.Hidden {
			names = append(names, p.Name)
		}
	}
	// }}}


	// Display Menu
	// 显示菜单
	// {{{
	////fmt.Println("\x1b[1;7m=== Windows Terminal Launcher ===\x1b[0m")
	fmt.Print("\x1b[1;7m")
	fmt.Print("=== Windows Terminal Launcher ===")
	fmt.Print("\x1b[0m")
	fmt.Print("\n")

	for i, name := range names {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	fmt.Println("[q] Quit / Exit")

	////fmt.Print("\x1b[1;7mSelect a profile\x1b[0m: ")
	fmt.Print("\x1b[1;7m")
	fmt.Print("Select a profile")
	fmt.Print("\x1b[0m: ")
	// }}}


	// Set Raw Mode and Read Key
	// 设置原始模式并读取按键
	// {{{
	// term.MakeRaw allows reading input without pressing Enter
	// term.MakeRaw 允许在不按回车的情况下读取输入 
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	// Create a byte slice as a buffer to store 1 byte of input
	// 创建一个字节切位作为缓冲区，用于存储 1 字节的输入 
	b := make([]byte, 1)
	// Read exactly one byte from the raw terminal
	// 从原始模式的终端中精确读取一个字节 
	_, err = os.Stdin.Read(b)
	if err != nil {
	}

	// Restore state immediately to allow clear output for the next loop
	// 立即恢复状态以便为下一次循环提供正常的输出环境
	term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Printf("%s\n", string(b))
	// }}}


	// Logic processing
	// 逻辑处理
	// {{{
	var selected string

	// Handle numeric keys 1-9
	// 处理数字键 1-9
	if b[0] >= '1' && b[0] <= '9' {
		choice := int(b[0] - '1')
		if choice < len(names) {
			selected = names[choice]
		} else if len(names) > 0 {
			selected = names[0]
		}
	} else if b[0] == 'q' || b[0] == 3 { // 'q' or Ctrl+C
		// Exit loop
		// 退出循环
		fmt.Println("\n\x1b[1mExiting...\x1b[0m")
		return
	} /* else {
		// Default to first profile for any other key
		// 任意其他键默认启动第一个配置
		if len(names) > 0 {
			selected = names[0]
		}
	} */
	// }}}


	// Execute Windows Terminal command
	// 执行 Windows Terminal 命令
	// {{{
	if selected != "" {
		// Get current working directory to pass to Windows Terminal
		// 获取当前工作目录并传递给 Windows Terminal
		cwd, err := os.Getwd()
		if err != nil {
			// Fallback to current if error
			// 出错时回退至当前点
			cwd = "."
		}
		cmd := exec.Command("wt.exe", "--window", "0", "new-tab", "--profile", selected, "--startingDirectory", cwd)
		////cmd := exec.Command("wt.exe", "--window", "0", "split-pane", "horizontal", "--profile", selected, "--startingDirectory", cwd)
		err = cmd.Start()
		if err != nil {
			////fmt.Printf("\n\x1b[1;31mFaild to execute: %v\x1b[0m\n", err)
			fmt.Print("\n")
			fmt.Print("\x1b[1;31m")
			fmt.Print("Faild to execute: ")
			fmt.Print("\x1b[0m")
			fmt.Print(err)
		} else {
			////fmt.Printf("\n\x1b[1;32mLaunched [%s] at [%s]\x1b[0m\n", selected, cwd)
			fmt.Printf("\n")
			fmt.Print("\x1b[1;32m")
			fmt.Printf("Launched [%s] at [%s]", selected, cwd)
			fmt.Print("\x1b[0m")
			fmt.Print("\n")

			// Important: Wait for wt.exe to handover the IPC (Inter-Process Communication) 
			// commands to the main Windows Terminal server process. 
			// If this launcher exits immediately, the terminal will destroy the 
			// current pane before it has time to process the "new-pane" request.
			// 重要：等待 wt.exe 通过 IPC（进程间通信）将指令移交给 Windows Terminal 
			// 服务进程。如果启动器立即退出，终端服务可能还没来得及处理“新建窗格”请求，
			// 当前窗格就已经因为进程结束而被销毁了。
			time.Sleep(50 * time.Millisecond)
		}
	}
	////fmt.Print("\x1b[1;33mInvalid input. \x1b[0m\x1b[1;7mSelect again\x1b[0m: ")
	fmt.Print("\x1b[1;33m")
	fmt.Print("Invalid input. ")
	fmt.Print("\x1b[0m")
	fmt.Print("\x1b[1;7m")
	fmt.Print("Select again")
	fmt.Print("\x1b[0m")
	fmt.Print(": ")
	// }}}
}


// vim:foldmethod=marker:
