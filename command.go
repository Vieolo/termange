package termange

import (
	"bytes"
	"errors"
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

type CommandError struct {
	RawError  error
	ExitError *exec.ExitError // exit error, returned when command starts but does not complete successfully
}

func (c CommandError) Error() string {
	return c.RawError.Error()
}

// Runs the command and returns the stdout and strerr
func RunCommand(config CommandConfig) (CommandResult, *CommandError) {
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
	}

	// Running the command
	cErr := c.Run()

	// Clearing the start text after the command is run
	if config.ClearStartText && config.HasStartText() {
		internal.IT{}.CursorUp().ClearLine()
	}

	// Preparing the result
	result := CommandResult{
		Stdout: sout,
		Stderr: stderr,
	}

	// Type conversion of the error (if applicable)
	// and add it to the result
	//
	// One of the error types returned is `*ExitError` which is given
	// when the command starts but does not complete successfully
	if cErr != nil {
		ce := CommandError{}
		ce.RawError = cErr

		var exitError *exec.ExitError
		if errors.As(cErr, &exitError) {
			ce.ExitError = exitError
		}
		return result, &ce
	}

	return result, nil
}

// Runs a raw command without the need to break it down into different parts
//
// It returns the stdout and stderr
func RunRawCommand(command string) (CommandResult, *CommandError) {
	return RunCommand(CommandConfig{
		Command: "sh",
		Args: []string{
			"-c",
			command,
		},
	})
}
