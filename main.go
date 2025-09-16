package main

import (
	"fmt"
	"os"

	"github.com/Programmer-Bugs-Bunny/GreenFishCli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
