package termange

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type CommandConfig struct {
	Command string            // The main command to be run
	Args    []string          // The list of arguments passes to the main command
	Env     map[string]string // The key-value pair of values added to the env for running the command
}

type CommandResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

// Runs the command and returns the stdout and strerr
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

	// Running the command
	cErr := c.Run()

	return CommandResult{
		Stdout: sout,
		Stderr: stderr,
	}, cErr
}

// Runs a raw command without the need to break it down into different parts
//
// It returns the stdout and stderr
func RunRawCommand(command string) (CommandResult, error) {
	return RunCommand(CommandConfig{
		Command: "sh",
		Args: []string{
			"-c",
			command,
		},
	})
}
