package cli_test

import (
	"testing"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
	"github.com/dannyjimenez98/task-tracker/internal/task"
)

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
