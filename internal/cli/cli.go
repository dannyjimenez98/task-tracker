// Package cli handles the logic for the CLI commands and actions
package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Start parses the command line arguments and calls the appropriate function
// to handle the requested action.
func Start() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	switch os.Args[1] {
	// "add" expects 1 argument: taskDescription (string)
	case "add": 
		if err := addCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		taskDescription := os.Args[2]
		fmt.Printf("adding task %s\n", taskDescription)

	// "delete" expects 1 argument: taskID (int)
	case "delete": 
		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("deleting task %d\n", taskID)

	// "update" expects 2 arguments: taskID (int), newTaskDescription (string)
	case "update": 
		if err := updateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
		newTaskDescription := os.Args[3]
		fmt.Printf("updating task %d with description %s\n", taskID, newTaskDescription)
	default:
		fmt.Println("unknown subcommand")
		os.Exit(1)
	} 
}
