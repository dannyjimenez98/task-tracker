package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
	"github.com/dannyjimenez98/task-tracker/internal/task"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("invalid subcommand entry")
	}

	const filename = "tasklist.json"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("cannot open or create file: ", err)
	}
	defer file.Close()

	var tasklist *task.Tasks

	tasklistJSON, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if len(tasklistJSON) == 0 {
		tasklist = &task.Tasks{
			Tasks: []task.Task{},
		}
	} else {
		err := json.Unmarshal(tasklistJSON, &tasklist)
		if err != nil {
			log.Fatalln("could not unmarshal json data: ", err)
		}
	}

	if err := cli.Start(tasklist, os.Args[1:]); err != nil {
		panic(err)
	}

	// prevent JSON file from being rewritten when 'list' command is called
	// nothing is written or changed to JSON file for this action
	if os.Args[1] == "list" {
		return
	}

	updatedTasklistJSON, err := json.MarshalIndent(tasklist, "", "  ")
	if err != nil {
		log.Fatalln("could not marshal data to json: ", err)
	}

	if err := file.Truncate(0); err != nil {
		log.Fatalln(err)
	}

	if _, err := file.Write(updatedTasklistJSON); err != nil {
		log.Fatalln(err)
	}
}
