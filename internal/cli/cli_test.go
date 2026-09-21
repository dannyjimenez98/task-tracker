package cli_test

import (
	"testing"

	"github.com/dannyjimenez98/task-tracker/internal/cli"
)

func TestAddTask(t *testing.T) {
	tests := []struct{
		name string
		args []string
		wantErr bool
	}{
		{
			name: "add task", 
			args: []string{"add", "return books"}, 
			wantErr: false,
		},
		{
			name: "add without description", 
			args: []string{"add"}, 
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			err := cli.Start(tt.args)
			if tt.wantErr && err == nil {
				t.Error(err)
			}
			if !tt.wantErr && err != nil {
				t.Error(err)
			}
		})
	}
}
