package main

import (
	"os"
	"fmt"
	"github.com/Lavos/gbvideo/commands/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
