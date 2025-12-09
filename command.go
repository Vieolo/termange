package termange

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/vieolo/termange/internal"
)

type CommandConfig struct {
	Command        string            // The main command to be run
	Args           []string          // The list of arguments passes to the main command
	Env            map[string]string // The key-value pair of values added to the env for running the command
	StartText      string            // The text to be printed before running the command
	ClearStartText bool              // Whether to clear the start text after the command is completed
}

func (c CommandConfig) HasStartText() bool {
	return len(strings.TrimSpace(c.StartText)) > 0
}

type CommandResult struct {
	Stdout bytes.Buffer // stdout of the command
	Stderr bytes.Buffer // stderr of the command
}

// Runs the command and returns the stdout and strerr
//
// Args:
//   - config: an instance of CommandConfig
//
// Returns:
//   - Instance of CommandResult
//   - error
//
// Errors:
//   - *exec.ExitError: when the command starts but does not complete successfully
func RunCommand(config CommandConfig) (CommandResult, error) {
	// Creating the command
	c := exec.Command(config.Command, config.Args...)

	// Addind stdout and stderr to the command
	var sout bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &sout
	c.Stderr = &stderr

	// Setting up the env variables
	if len(config.Env) > 0 {
		c.Env = os.Environ()
		for k, v := range config.Env {
			c.Env = append(c.Env, fmt.Sprintf("%v=%v", k, v))
		}
	}

	// Printing the start text
	if config.HasStartText() {
		PrintInfoln(config.StartText)

		// Clearing the start text after the command is run
		if config.ClearStartText {
			defer internal.IT{}.CursorUp().ClearLine()
		}
	}

	// Running the command
	cErr := c.Run()

	return CommandResult{
		Stdout: sout,
		Stderr: stderr,
	}, cErr
}

// Runs a raw command without the need to break it down into different parts
//
// Args:
//   - command: the raw command to be run
//
// Returns:
//   - Instance of CommandResult
//   - error
//
// Errors:
//   - *exec.ExitError: when the command starts but does not complete successfully
func RunRawCommand(command string) (CommandResult, error) {
	return RunCommand(CommandConfig{
		Command: "sh",
		Args: []string{
			"-c",
			command,
		},
	})
}
