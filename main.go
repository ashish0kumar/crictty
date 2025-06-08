package main

import (
	"github.com/ashish0kumar/crictty/cmd"
	"os"
)

// main is the entry point of the application
func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
