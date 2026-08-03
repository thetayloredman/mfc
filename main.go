package main

import (
	"fmt"

	"github.com/thetayloredman/mfc/cli"
	"github.com/thetayloredman/mfc/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	root := cli.NewRootCommand(config)
	if err := root.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
