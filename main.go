package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/term"
)

// WTConfig matches Windows Terminal settings.json structure
// WTConfig 对应 Windows Terminal 的 settings.json 结构
type WTConfig struct {
	Profiles struct {
		List []struct {
			Name   string `json:"name"`   // Profile name / 配置项名称
			Hidden bool   `json:"hidden"` // Whether it is hidden / 是否隐藏
		} `json:"list"`
	} `json:"profiles"`
}

func main() {
	// 1. Get config path and read file
	// 1. 获取配置文件路径并读取文件
	// {{{
	configPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		`Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json`)

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Faild to read config file: %v\n", err)
		return
	}
	// }}}

	// 2. Parse JSON and filter profiles
	// 2. 解析 JSON 并过滤配置项
	// {{{
	var config WTConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Faild to analysis JSON: %v\n", err)
		return
	}

	var names []string
	for _, p := range config.Profiles.List {
		if !p.Hidden {
			names = append(names, p.Name)
		}
	}
	// }}}

	// 3. Display Menu
	// 3. 显示菜单
	// {{{
	fmt.Println("=== Windows Terminal Launcher===")
	for i, name := range names {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	////fmt.Println("--------------------------------")
	////fmt.Print("Input a number: ")
	// }}}

	// 4. Set Raw Mode and Read Key
	// 4. 设置原始模式并读取按键
	// {{{
	// term.MakeRaw allows reading input without pressing Enter
	// term.MakeRaw 允许在不按回车的情况下读取输入
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	// Ensure terminal state is restored on exit
	// 确保退出时恢复终端状态
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Create a byte slice as a buffer to store 1 byte of input
	// 创建一个字节切位作为缓冲区，用于存储 1 字节的输入
	b := make([]byte, 1)

	// Read exactly one byte from the raw terminal
	// 从原始模式的终端中精确读取一个字节
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}

	// Restore state immediately after reading to allow normal output
	// 读取后立即恢复状态以允许正常输出
	term.Restore(int(os.Stdin.Fd()), oldState)
	////fmt.Printf("%s\n", string(b))
	// }}}

	// 5. Logic processing
	// 5. 逻辑处理
	// {{{

	// Old 5
	/*
	// {{{
	choice := int(b[0] - '1')
	if choice >= 0 && choice < len(names) {
		selected := names[choice]
		////fmt.Printf("Starting: %s...\n", selected)
		
		cmd := exec.Command("wt.exe", "new-tab", "--profile", selected)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Faild to execute: %v\n", err)
		}
	} else {
		fmt.Println("Invalid, exit.")
	}
	// }}}
	*/

	var selected string

	// Handle numeric keys 1-9
	// 处理数字键 1-9
	if b[0] >= '1' && b[0] <= '9' {
		choice := int(b[0] - '1')
		if choice < len(names) {
			selected = names[choice]
		} else {
			// Fallback to first if out of range
			// 如果超出范围则回退到第一个
			selected = names[0]
		}
	} else if b[0] == 'q' {
		// Quit if 'q' is pressed
		// 如果按下 'q' 则退出
		return
	} else {
		// Default to first profile for any other key
		// 任意其他键默认启动第一个配置
		if len(names) > 0 {
			selected = names[0]
		}
	}

	// Execute Windows Terminal command
	// 执行 Windows Terminal 命令
	if selected != "" {
		// nt (new-tab) opens in the current window if possible
		// nt (new-tab) 尽可能在当前窗口打开新标签
		cmd := exec.Command("wt.exe", "new-tab", "--profile", selected)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Faild to execute: %v\n", err)
		}
		////time.Sleep(0.2 * 1000 * 1000 * 1000)
		////fmt.Printf("Starting profile: %s...\n", selected)
	} else {
		fmt.Println("Invalid, exit.")
	}
	// }}}
}

// vim:foldmethod=marker:
