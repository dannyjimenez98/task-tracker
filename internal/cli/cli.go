// Package cli handles the logic for the CLI commands and actions
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dannyjimenez98/task-tracker/internal/task"
)

// Add appends a new task to the slice of tasks and writes it to the tasklist.json file.
func add(taskDescription string) {

	tasklist := task.GetTaskList()

	tasks := append(tasklist, task.Task{
		// set ID to ID of last entry of tasklist incremented by 1
		// sets ID to 1 if tasklist is empty
		ID: func(tasklist []task.Task) int {
				if len(tasklist) == 0 {
					return 1
				}
				return tasklist[len(tasklist) - 1].ID + 1
			}(tasklist),
		Description: taskDescription,
		Status: "todo",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	tasklistJSON, err := json.MarshalIndent(tasks, "", "  ") 
	if err != nil {
		fmt.Println("marshal error: ", err)
	}

	if err := os.WriteFile("tasklist.json", tasklistJSON, 0644); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully added task.")
}

// Start parses the command line arguments and calls the appropriate function
// to handle the requested action.
func Start() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	switch os.Args[1] {
	// "add" expects 1 argument: taskDescription (string)
	case "add": 
		if err := addCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("parsing error: %s\n", err)
			return
		}
		taskDescription := os.Args[2]
		fmt.Printf("adding task %s\n", taskDescription)
		add(taskDescription)
	
	// "delete" expects 1 argument: taskID (int)
	case "delete": 
		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("parsing error: %s\n", err)
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
			fmt.Printf("parsing error: %s\n", err)
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

