package terminalutils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type CommandConfig struct {
	Command string
	Args    []string
	Env     map[string]string
}

// Runs the command and returns the stdout and strerr
func RunCommand(config CommandConfig) (bytes.Buffer, bytes.Buffer, error) {
	c := exec.Command(config.Command, config.Args...)
	var sout bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &sout
	c.Stderr = &stderr
	if len(config.Env) > 0 {
		c.Env = os.Environ()
		for k, v := range config.Env {
			c.Env = append(c.Env, fmt.Sprintf("%v=%v", k, v))
		}
	}
	cErr := c.Run()

	if cErr != nil {
		return sout, stderr, cErr
	} else {
		return sout, stderr, nil
	}
}

// Runs a raw command without the need to break it down into different parts
//
// It returns the stdout and stderr
func RunRawCommand(command string) (bytes.Buffer, bytes.Buffer, error) {
	return RunCommand(CommandConfig{
		Command: "sh",
		Args: []string{
			"-c",
			command,
		},
	})
}
