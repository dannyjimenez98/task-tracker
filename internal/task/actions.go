package task

import (
	"time"
)

// Add appends a new task to the tasklist
func (tasklist *Tasks) Add(taskDescription string) {
	tasklist.Tasks = append(tasklist.Tasks, Task{
		// set ID to ID of last entry of tasklist incremented by 1
		// sets ID to 1 if tasklist is empty
		ID: func() int {
				if len(tasklist.Tasks) == 0 {
					return 1
				}
				return tasklist.Tasks[len(tasklist.Tasks) - 1].ID + 1
			}(),
		Description: taskDescription,
		Status: "todo",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}


