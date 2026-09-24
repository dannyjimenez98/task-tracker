package task

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
)

var timestamp = CustomTime{time.Now()}

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
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
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
	tasklist.Tasks[i].UpdatedAt = timestamp

	return nil
}

func (tasklist *Tasks) MarkInProgress(taskID int) error {
	// find the index of the task with taskID in tasklist slice
	i := slices.IndexFunc(tasklist.Tasks, func(n Task) bool {
		return n.ID == taskID
	})
	if i == -1 {
		return fmt.Errorf("could not mark in-progress: task with id %d not in tasklist", taskID)
	}

	tasklist.Tasks[i].Status = "in-progress"
	tasklist.Tasks[i].UpdatedAt = timestamp

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
	tasklist.Tasks[i].UpdatedAt = timestamp

	return nil
}

func (tasklist *Tasks) List(statusFilter string) error {
	switch statusFilter {
	case "all", "todo", "in-progress", "done":
		// Valid filters.
	default:
		return fmt.Errorf("invalid status filter")
	}

	// check if the applied status filter returns at least one task
	// if no filter applied ("all"), check if at least one task exists in task list
	hasMatch := slices.ContainsFunc(tasklist.Tasks, func(t Task) bool {
		return statusFilter == "all" || t.Status == statusFilter
	})

	if !hasMatch {
		if statusFilter == "all" {
			_, err := fmt.Fprintln(os.Stdout, "no tasks to print")
			return err
		}

		_, err := fmt.Fprintf(
			os.Stdout,
			"no tasks found with a current status of \"%s\"\n",
			statusFilter,
		)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	// table header
	if _, err := fmt.Fprintln(w, "ID\tDESCRIPTION\tSTATUS\tCREATED\tUPDATED\t"); err != nil {
		return err
	}

	// prevent escape characters in task description from being applied to table
	// literal escape characters will be printed within task's description
	escapeLiterals := strings.NewReplacer(
		"\t", `\t`,
		"\n", `\n`,
		"\r", `\r`,
	)

	// table body
	for _, taskEntry := range tasklist.Tasks {
		if statusFilter != "all" && taskEntry.Status != statusFilter {
			continue
		}

		if _, err := fmt.Fprintf(
			w,
			"%v\t%v\t%v\t%v\t%v\t\n",
			taskEntry.ID,
			escapeLiterals.Replace(taskEntry.Description),
			taskEntry.Status,
			taskEntry.CreatedAt,
			taskEntry.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return w.Flush()
}
