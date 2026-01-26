package main

import (
	"fmt"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/term"
)

type WTConfig struct {
	Profiles struct {
		List []struct {
			Name string `json:"name"`
			Hidden bool `json:"hidden"`
		} `json:"list"`
	} `json:"profiles"`
}

func main() {
	configPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		`Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json`)

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

	var names []string
	for _, p := range config.Profiles.List {
		if !p.Hidden {
			names = append(names, p.Name)
		}
	}

	fmt.Println("=== Windows Terminal Launcher===")
	for i, name := range names {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	////fmt.Println("--------------------------------")
	////fmt.Print("Input a number: ")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}

	term.Restore(int(os.Stdin.Fd()), oldState)
	////fmt.Printf("%s\n", string(b))

	// Logic processing
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
}

// vim:foldmethod=marker:
