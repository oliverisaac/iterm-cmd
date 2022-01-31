package generate

import (
	"log"
	"os"
	"testing"

	"github.com/oliverisaac/go-levellogger"
	"github.com/pkg/errors"
)

func TestCommandGenerator_CommandFile(t *testing.T) {
	type fields struct {
		logger    levellogger.Logger
		directory string
	}

	workDir, err := os.MkdirTemp(os.TempDir(), "iterm-cmd-")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create workingdir for tests"))
	}
	defer os.RemoveAll(workDir)

	tests := []struct {
		name           string
		fields         fields
		args           []string
		wantOutputFile string
		wantErr        bool
	}{
		{
			name: "Test that a file is created",
			fields: fields{
				directory: workDir,
			},
			args: []string{
				"/bin/bash",
				"-c",
				"sleep 1",
			},
			wantOutputFile: workDir + "/7",
			wantErr:        false,
		},
		{
			name: "Test that a with the same command generates the same name",
			fields: fields{
				directory: workDir,
			},
			args: []string{
				"/bin/bash",
				"-c",
				"sleep 1",
			},
			wantOutputFile: workDir + "/7",
			wantErr:        false,
		},
		{
			name: "Test that a command that has the same first hash generates a different filename",
			fields: fields{
				directory: workDir,
			},
			args: []string{
				"/bin/bash",
				"-c",
				"sleep 13011", // Found this value of sleep 17 through trial and error
			},
			wantOutputFile: workDir + "/7c",
			wantErr:        false,
		},
		{
			name: "Test that a command with no args works",
			fields: fields{
				directory: workDir,
			},
			args: []string{
				"/bin/echo",
			},
			wantOutputFile: workDir + "/c",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := NewCommandGenerator(tt.fields.logger, tt.fields.directory)
			exe := tt.args[0]
			var args []string
			if len(tt.args) > 1 {
				args = tt.args[1:]
			}
			gotOutputFile, err := cg.CommandFile(exe, args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandGenerator.CommandFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutputFile != tt.wantOutputFile {
				t.Errorf("\n got:  %v\n want: %v", gotOutputFile, tt.wantOutputFile)
			}
		})
	}
}
