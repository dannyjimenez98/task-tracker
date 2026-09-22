// Package cli handles the logic for the CLI commands and actions
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/dannyjimenez98/task-tracker/internal/task"
)

// Start parses the command line arguments and calls the appropriate function
// to handle the requested action.
func Start(tasklist *task.Tasks, args []string) error {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	switch args[0] {
	// "add" expects 1 argument: taskDescription (string)
	case "add":
		if err := addCmd.Parse(args[1:]); err != nil {
			fmt.Printf("parsing error: %s\n", err)
			return err
		}

		if addCmd.NArg() == 0 {
			return errors.New("task description not provided")
		}

		taskDescription := addCmd.Arg(0)
		tasklist.Add(taskDescription)

	// "delete" expects 1 argument: taskID (int)
	case "delete":
		if err := deleteCmd.Parse(args[1:]); err != nil {
			fmt.Printf("parsing error: %s\n", err)
			return err
		}

		if deleteCmd.NArg() == 0 {
			return errors.New("task ID not provided")
		}

		taskID, err := strconv.Atoi(deleteCmd.Arg(0))
		if err != nil {
			return fmt.Errorf("invalid task ID %q: %w", deleteCmd.Arg(0), err)
		}

		return tasklist.Delete(taskID)

	// "update" expects 2 arguments: taskID (int), newTaskDescription (string)
	case "update":
		if err := updateCmd.Parse(args[1:]); err != nil {
			fmt.Printf("parsing error: %s\n", err)
			return err
		}

		if updateCmd.NArg() < 2 {
			return errors.New("not enough arguments passed")
		}

		taskID, err := strconv.Atoi(updateCmd.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
		newTaskDescription := updateCmd.Arg(1)

		return tasklist.Update(taskID, newTaskDescription)

	default:
		fmt.Println("unknown subcommand")
		os.Exit(1)
	}

	return nil
}
