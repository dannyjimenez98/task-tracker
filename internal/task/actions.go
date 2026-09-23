package task

import (
	"fmt"
	"slices"
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
			return tasklist.Tasks[len(tasklist.Tasks)-1].ID + 1
		}(),
		Description: taskDescription,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
}

func (tasklist *Tasks) Delete(taskID int) error {
	// find the index of the task with taskID in tasklist slice
	i := slices.IndexFunc(tasklist.Tasks, func(n Task) bool {
		return n.ID == taskID
	})
	if i == -1 {
		return fmt.Errorf("could not delete: task with id %d not in tasklist", taskID)
	}

	tasklist.Tasks = slices.Delete(tasklist.Tasks, i, i+1)

	return nil
}

func (tasklist *Tasks) Update(taskID int, newTaskDescription string) error {
	// find the index of the target taskID
	i := slices.IndexFunc(tasklist.Tasks, func(n Task) bool {
		return n.ID == taskID
	})
	if i == -1 {
		return fmt.Errorf("could not update: task with id %d not in tasklist", taskID)
	}

	tasklist.Tasks[i].Description = newTaskDescription

	return nil
}

func (tasklist *Tasks) MarkDone(taskID int) error {
	// find the index of the task with taskID in tasklist slice
	i := slices.IndexFunc(tasklist.Tasks, func(n Task) bool {
		return n.ID == taskID
	})
	if i == -1 {
		return fmt.Errorf("could not mark done: task with id %d not in tasklist", taskID)
	}

	tasklist.Tasks[i].Status = "done"

	return nil
}
