package cli_test

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
	"github.com/dannyjimenez98/task-tracker/internal/task"
)

func TestList(t *testing.T) {
	createdAt := task.CustomTime{
		Time: time.Date(2026, time.September, 20, 10, 0, 0, 0, time.UTC),
	}
	updatedAt := task.CustomTime{
		Time: time.Date(2026, time.September, 21, 11, 30, 0, 0, time.UTC),
	}

	inputTasklist := []task.Task{
		{
			ID: 1, Description: "task one", Status: "todo",
			CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
		{
			ID: 2, Description: "task two", Status: "in-progress",
			CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
		{
			ID: 3, Description: "task three", Status: "done",
			CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
		{
			ID: 4, Description: "task four", Status: "todo",
			CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
	}

	// Expectations use fixed timestamps rather than the production formatter.
	const timestamps = " 2026-09-20 10:00:00 +0000 UTC" +
		" 2026-09-21 11:30:00 +0000 UTC"
	const header = "ID DESCRIPTION STATUS CREATED UPDATED"

	rows := []string{
		"1 task one todo" + timestamps,
		"2 task two in-progress" + timestamps,
		"3 task three done" + timestamps,
		"4 task four todo" + timestamps,
	}

	table := func(rows ...string) string {
		return header + "\n" + strings.Join(rows, "\n") + "\n"
	}

	type testCase struct {
		name       string
		tasks      []task.Task
		args       []string
		wantOutput string
		wantErr    bool
	}
	tests := []testCase{
		{
			name:       "list defaults to all",
			tasks:      inputTasklist,
			args:       []string{"list"},
			wantOutput: table(rows...),
		},
		{
			name:       "list explicit all",
			tasks:      inputTasklist,
			args:       []string{"list", "all"},
			wantOutput: table(rows...),
		},
		{
			name:       "list todo includes every matching task",
			tasks:      inputTasklist,
			args:       []string{"list", "todo"},
			wantOutput: table(rows[0], rows[3]),
		},
		{
			name:       "list in-progress",
			tasks:      inputTasklist,
			args:       []string{"list", "in-progress"},
			wantOutput: table(rows[1]),
		},
		{
			name:       "list done",
			tasks:      inputTasklist,
			args:       []string{"list", "done"},
			wantOutput: table(rows[2]),
		},
		{
			name: "list escapes table control characters",
			tasks: []task.Task{
				{
					ID: 5, Description: "first\tsecond\nthird\rfourth",
					Status:    "todo",
					CreatedAt: createdAt, UpdatedAt: updatedAt,
				},
			},
			args: []string{"list"},
			wantOutput: table(
				`5 first\tsecond\nthird\rfourth todo` + timestamps,
			),
		},
	}

	// Cover nil and initialized empty task slices with every supported filter.
	emptyInputs := []struct {
		name  string
		tasks []task.Task
	}{
		{name: "nil task slice", tasks: nil},
		{name: "empty task slice", tasks: []task.Task{}},
	}
	filters := []struct {
		name string
		args []string
	}{
		{name: "default", args: []string{"list"}},
		{name: "all", args: []string{"list", "all"}},
		{name: "todo", args: []string{"list", "todo"}},
		{name: "in-progress", args: []string{"list", "in-progress"}},
		{name: "done", args: []string{"list", "done"}},
	}

	for _, input := range emptyInputs {
		for _, filter := range filters {
			wantOutput := "no tasks to print\n"
			if filter.name != "default" && filter.name != "all" {
				wantOutput = fmt.Sprintf(
					"no tasks found with a current status of \"%s\"\n",
					filter.name,
				)
			}
			tests = append(tests, testCase{
				name:       input.name + "/" + filter.name,
				tasks:      input.tasks,
				args:       filter.args,
				wantOutput: wantOutput,
			})
		}
	}

	// A nonempty list can also have no tasks matching a valid filter.
	for _, status := range []string{"todo", "in-progress", "done"} {
		var nonmatching []task.Task
		for _, entry := range inputTasklist {
			if entry.Status != status {
				nonmatching = append(nonmatching, entry)
			}
		}
		tests = append(tests, testCase{
			name:  "no matching tasks/" + status,
			tasks: nonmatching,
			args:  []string{"list", status},
			wantOutput: fmt.Sprintf(
				"no tasks found with a current status of \"%s\"\n", status,
			),
		})
	}

	invalidInputs := []struct {
		name string
		args []string
	}{
		{name: "unknown status", args: []string{"list", "pending"}},
		{name: "uppercase status", args: []string{"list", "TODO"}},
		{name: "empty status", args: []string{"list", ""}},
		{name: "padded status", args: []string{"list", " todo "}},
		{name: "multiple statuses", args: []string{"list", "todo", "done"}},
		{name: "extra argument after all", args: []string{"list", "all", "extra"}},
		{name: "unknown flag", args: []string{"list", "--invalid"}},
	}
	for _, input := range invalidInputs {
		tests = append(tests, testCase{
			name:    input.name,
			tasks:   inputTasklist,
			args:    input.args,
			wantErr: true,
		})
	}

	// Ignore tabwriter's padding, but preserve line boundaries and row order.
	normalizeOutput := func(output string) string {
		output = strings.TrimSuffix(output, "\n")
		if output == "" {
			return ""
		}
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			lines[i] = strings.Join(strings.Fields(line), " ")
		}
		return strings.Join(lines, "\n") + "\n"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tasklist task.Tasks
			if tt.tasks != nil {
				tasklist.Tasks = make([]task.Task, len(tt.tasks))
				copy(tasklist.Tasks, tt.tasks)
			}

			output, err := captureListOutput(t, func() error {
				return cli.Start(&tasklist, tt.args)
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if got := normalizeOutput(output); got != tt.wantOutput {
				t.Errorf("Start() output:\n%q\nwant:\n%q", got, tt.wantOutput)
			}

			// Verify every attribute, ordering, and nil-versus-empty slice state.
			if !reflect.DeepEqual(tasklist.Tasks, tt.tasks) {
				t.Errorf(
					"list modified tasks:\ngot:  %#v\nwant: %#v",
					tasklist.Tasks, tt.tasks,
				)
			}
		})
	}
}

// Do not use this helper in parallel tests: os.Stdout is process-wide.
func captureListOutput(t *testing.T, run func() error) (string, error) {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "list-output-*")
	if err != nil {
		t.Fatalf("create output file: %v", err)
	}
	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Errorf("close output file: %v", err)
		}
	})

	originalStdout := os.Stdout
	os.Stdout = file
	defer func() {
		os.Stdout = originalStdout
	}()

	runErr := run()
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("rewind output file: %v", err)
	}
	output, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	return string(output), runErr
}

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
