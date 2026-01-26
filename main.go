package main

import (
	"fmt"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/term"
)

// 对应 settings.json 的结构
type WTConfig struct {
	Profiles struct {
		List []struct {
			Name string `json:"name"`
			Hidden bool `json:"hidden"`
		} `json:"list"`
	} `json:"profiles"`
}

func main() {
	// 1. 定义配置文件路径
	configPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		`Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json`)

	// 2. 读取并解析文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Faild to read config file: %v\n", err)
		return
	}

	var config WTConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Faild to analysis JSON: %v\n", err)
		return
	}

	// 过滤掉隐藏的 Profile
	var names []string
	for _, p := range config.Profiles.List {
		if !p.Hidden {
			names = append(names, p.Name)
		}
	}

	// 3. 打印菜单
	fmt.Println("=== Windows Terminal Launcher===")
	for i, name := range names {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	////fmt.Println("--------------------------------")
	////fmt.Print("Input a number: ")

	// 4. 设置终端为 Raw 模式，实现无需回车获取输入
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// 5. 读取单个字符
	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}

	// 恢复终端状态以便后续输出
	term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Printf("%s\n", string(b))

	// 6. 处理逻辑
	choice := int(b[0] - '1') // 将 ASCII 转换为索引
	if choice >= 0 && choice < len(names) {
		selected := names[choice]
		////fmt.Printf("Starting: %s...\n", selected)
		
		// 执行 wt -p "Profile Name"
		cmd := exec.Command("wt.exe", "new-tab", "--profile", selected)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Faild to execute: %v\n", err)
		}
	} else {
		fmt.Println("Invalid, exit.")
	}
}
