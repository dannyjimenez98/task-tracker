package cli_test

import (
	"testing"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
	"github.com/dannyjimenez98/task-tracker/internal/task"
)

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
