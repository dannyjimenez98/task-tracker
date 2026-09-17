package main

import (
	// "fmt"

	"fmt"
	"os"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("invalid subcommand entry")
		os.Exit(1)
	}
	cli.Start()
}
