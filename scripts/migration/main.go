package main

import (
	"os"
)

func main() {
	cli := NewDefaultCLI()
	exitCode := cli.Run()
	os.Exit(int(exitCode))
}
