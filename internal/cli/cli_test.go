package cli_test

import (
	"testing"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
	"github.com/dannyjimenez98/task-tracker/internal/task"
)

func TestMarkInProgress(t *testing.T) {
	inputTasklist := []task.Task{
		{
			ID:          1,
			Description: "task 1",
			Status:      "todo",
		},
		{
			ID:          2,
			Description: "task 2",
			Status:      "in-progress",
		},
		{
			ID:          3,
			Description: "task 3",
			Status:      "done",
		},
	}

	tests := []struct {
		name       string
		tasks      []task.Task
		args       []string
		wantOutput []task.Task
		wantErr    bool
	}{
		{
			name:  "mark-in-progress on todo task",
			tasks: inputTasklist,
			args:  []string{"mark-in-progress", "1"},
			wantOutput: []task.Task{
				{ID: 1, Description: "task 1", Status: "in-progress"},
				inputTasklist[1],
				inputTasklist[2],
			},
			wantErr: false,
		},
		{
			name:       "mark-in-progress on in-progress task",
			tasks:      inputTasklist,
			args:       []string{"mark-in-progress", "2"},
			wantOutput: inputTasklist,
			wantErr:    false,
		},
		{
			name:  "mark-in-progress on done task",
			tasks: inputTasklist,
			args:  []string{"mark-in-progress", "3"},
			wantOutput: []task.Task{
				inputTasklist[0],
				inputTasklist[1],
				{ID: 3, Description: "task 3", Status: "in-progress"},
			},
			wantErr: false,
		},
		{
			name:       "mark-in-progress nonexistent task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-in-progress", "99"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "mark-in-progress in empty task list",
			tasks:      []task.Task{},
			args:       []string{"mark-in-progress", "1"},
			wantOutput: []task.Task{},
			wantErr:    true,
		},
		{
			name:       "mark-in-progress without task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-in-progress"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "mark-in-progress with invalid task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-in-progress", "abc"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasklist := task.Tasks{
				Tasks: make([]task.Task, len(tt.tasks)),
			}
			copy(tasklist.Tasks, tt.tasks)

			err := cli.Start(&tasklist, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Check that IDs remain unchanged and descriptions
			// match the expected result for every task.
			for i, want := range tt.wantOutput {
				got := tasklist.Tasks[i]

				if got.ID != want.ID {
					t.Errorf(
						"task[%d].ID = %d, want %d",
						i, got.ID, want.ID,
					)
				}

				if got.Description != want.Description {
					t.Errorf(
						"task[%d].Description = %q, want %q",
						i, got.Description, want.Description,
					)
				}

				if got.Status != want.Status {
					t.Errorf(
						"task[%d].Status = %q, want %q",
						i, got.Status, want.Status,
					)
				}
			}
		})
	}
}

func TestMarkDone(t *testing.T) {
	inputTasklist := []task.Task{
		{
			ID:          1,
			Description: "task 1",
			Status:      "todo",
		},
		{
			ID:          2,
			Description: "task 2",
			Status:      "in-progress",
		},
		{
			ID:          3,
			Description: "task 3",
			Status:      "done",
		},
	}

	tests := []struct {
		name       string
		tasks      []task.Task
		args       []string
		wantOutput []task.Task
		wantErr    bool
	}{
		{
			name:  "mark-done on todo task",
			tasks: inputTasklist,
			args:  []string{"mark-done", "1"},
			wantOutput: []task.Task{
				{ID: 1, Description: "task 1", Status: "done"},
				inputTasklist[1],
				inputTasklist[2],
			},
			wantErr: false,
		},
		{
			name:  "mark-done on in-progress task",
			tasks: inputTasklist,
			args:  []string{"mark-done", "2"},
			wantOutput: []task.Task{
				inputTasklist[0],
				{ID: 2, Description: "task 2", Status: "done"},
				inputTasklist[2],
			},
			wantErr: false,
		},
		{
			name:       "mark-done on done task",
			tasks:      inputTasklist,
			args:       []string{"mark-done", "3"},
			wantOutput: inputTasklist,
			wantErr:    false,
		},
		{
			name:       "mark-done nonexistent task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-done", "99"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "mark-done in empty task list",
			tasks:      []task.Task{},
			args:       []string{"mark-done", "1"},
			wantOutput: []task.Task{},
			wantErr:    true,
		},
		{
			name:       "mark-done without task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-done"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "mark-done with invalid task ID",
			tasks:      inputTasklist,
			args:       []string{"mark-done", "abc"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasklist := task.Tasks{
				Tasks: make([]task.Task, len(tt.tasks)),
			}
			copy(tasklist.Tasks, tt.tasks)

			err := cli.Start(&tasklist, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Check that IDs remain unchanged and descriptions
			// match the expected result for every task.
			for i, want := range tt.wantOutput {
				got := tasklist.Tasks[i]

				if got.ID != want.ID {
					t.Errorf(
						"task[%d].ID = %d, want %d",
						i, got.ID, want.ID,
					)
				}

				if got.Description != want.Description {
					t.Errorf(
						"task[%d].Description = %q, want %q",
						i, got.Description, want.Description,
					)
				}

				if got.Status != want.Status {
					t.Errorf(
						"task[%d].Status = %q, want %q",
						i, got.Status, want.Status,
					)
				}
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	inputTasklist := []task.Task{
		{
			ID:          1,
			Description: "task 1",
		},
		{
			ID:          2,
			Description: "task 2",
		},
		{
			ID:          3,
			Description: "task 3",
		},
	}

	tests := []struct {
		name       string
		tasks      []task.Task
		args       []string
		wantOutput []task.Task
		wantErr    bool
	}{
		{
			name:  "update task",
			tasks: inputTasklist,
			args:  []string{"update", "2", "updated task 2"},
			wantOutput: []task.Task{
				inputTasklist[0],
				{ID: 2, Description: "updated task 2"},
				inputTasklist[2],
			},
			wantErr: false,
		},
		{
			name:       "update nonexistent task ID",
			tasks:      inputTasklist,
			args:       []string{"update", "99", "updated task"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "update in empty task list",
			tasks:      []task.Task{},
			args:       []string{"update", "1", "updated task"},
			wantOutput: []task.Task{},
			wantErr:    true,
		},
		{
			name:       "update without task ID and description",
			tasks:      inputTasklist,
			args:       []string{"update"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "update without description",
			tasks:      inputTasklist,
			args:       []string{"update", "2"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
		{
			name:       "update with invalid task ID",
			tasks:      inputTasklist,
			args:       []string{"update", "abc", "updated task"},
			wantOutput: inputTasklist,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasklist := task.Tasks{
				Tasks: make([]task.Task, len(tt.tasks)),
			}
			copy(tasklist.Tasks, tt.tasks)

			err := cli.Start(&tasklist, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Check that IDs remain unchanged and descriptions
			// match the expected result for every task.
			for i, want := range tt.wantOutput {
				got := tasklist.Tasks[i]

				if got.ID != want.ID {
					t.Errorf(
						"task[%d].ID = %d, want %d",
						i, got.ID, want.ID,
					)
				}

				if got.Description != want.Description {
					t.Errorf(
						"task[%d].Description = %q, want %q",
						i, got.Description, want.Description,
					)
				}
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	inputTasklist := []task.Task{
		{
			ID:          1,
			Description: "task 1",
		},
		{
			ID:          2,
			Description: "task 2",
		},
		{
			ID:          3,
			Description: "task 3",
		},
	}

	tests := []struct {
		name          string
		tasks         []task.Task
		args          []string
		wantOutput    []task.Task
		wantTaskCount int
		wantErr       bool
	}{
		{
			name:          "delete task in middle of list",
			tasks:         inputTasklist,
			args:          []string{"delete", "2"},
			wantOutput:    []task.Task{inputTasklist[0], inputTasklist[2]},
			wantTaskCount: 2,
			wantErr:       false,
		},
		{
			name:          "delete first task",
			tasks:         inputTasklist,
			args:          []string{"delete", "1"},
			wantOutput:    []task.Task{inputTasklist[1], inputTasklist[2]},
			wantTaskCount: 2,
			wantErr:       false,
		},
		{
			name:          "delete last task",
			tasks:         inputTasklist,
			args:          []string{"delete", "3"},
			wantOutput:    []task.Task{inputTasklist[0], inputTasklist[1]},
			wantTaskCount: 2,
			wantErr:       false,
		},
		{
			name:          "delete nonexistent task ID",
			tasks:         inputTasklist,
			args:          []string{"delete", "99"},
			wantOutput:    inputTasklist,
			wantTaskCount: 3,
			wantErr:       true,
		},
		{
			name:          "delete from empty task list",
			tasks:         []task.Task{},
			args:          []string{"delete", "1"},
			wantOutput:    []task.Task{},
			wantTaskCount: 0,
			wantErr:       true,
		},
		{
			name:          "delete without task ID",
			tasks:         inputTasklist,
			args:          []string{"delete"},
			wantOutput:    inputTasklist,
			wantTaskCount: 3,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasklist := task.Tasks{
				Tasks: make([]task.Task, len(tt.tasks)),
			}
			copy(tasklist.Tasks, tt.tasks)

			err := cli.Start(&tasklist, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if got := len(tasklist.Tasks); got != tt.wantTaskCount {
				t.Fatalf("task count = %d, want %d", got, tt.wantTaskCount)
			}

			// check both ID and description of each task in output
			// ID and description pairing should remain unchanged
			for i, want := range tt.wantOutput {
				got := tasklist.Tasks[i]

				if got.ID != want.ID {
					t.Errorf(
						"task[%d].ID = %d, want %d",
						i, got.ID, want.ID,
					)
				}

				if got.Description != want.Description {
					t.Errorf(
						"task[%d].Description = %q, want %q",
						i, got.Description, want.Description,
					)
				}
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantDescription string
		wantTaskCount   int
		wantErr         bool
	}{
		{
			name:            "add task",
			args:            []string{"add", "return books"},
			wantDescription: "return books",
			wantTaskCount:   1,
			wantErr:         false,
		},
		{
			name:            "add without description",
			args:            []string{"add"},
			wantDescription: "",
			wantTaskCount:   0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tasklist task.Tasks
			err := cli.Start(&tasklist, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if got := len(tasklist.Tasks); got != tt.wantTaskCount {
				t.Fatalf("task count = %d, want %d", got, tt.wantTaskCount)
			}

			if tt.wantTaskCount > 0 {
				got := tasklist.Tasks[0].Description
				if got != tt.wantDescription {
					t.Errorf(
						"task description = %q, want %q",
						got,
						tt.wantDescription,
					)
				}
			}
		})
	}
}
