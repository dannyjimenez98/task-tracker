package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)


type Task struct {
	ID int `json:"id"`
	Description string `json:"description"`
	Status string `json:"status"`// todo, in-progress, done
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetTaskList returns the slice of tasks from the tasklist.json file.
// Creates the file with an empty Task slice if it does not exist.
func GetTaskList() []Task {
	const filename = "tasklist.json"

	var tasks []Task
	
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			panic("could not create file")
		}
	} else {
		tasklistFile, err := os.ReadFile(filename)	
		if err != nil {
			fmt.Println("error reading file: ", err)
		}

		if err := json.Unmarshal(tasklistFile, &tasks); err != nil {
			fmt.Println("unmarshal error: ", err)
		}
	}
	
	return tasks
}
