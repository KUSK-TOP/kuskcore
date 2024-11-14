package main

import (
	"runtime"

	cmd "kuskcore/cmd/kuskcli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
