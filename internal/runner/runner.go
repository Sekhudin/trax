package runner

import (
	"fmt"
	"io"
	"os/exec"

	appErr "trax/internal/errors"
)

type Runner interface {
	Run(command map[string]any) error
}

type CommandRunner struct {
	Stdout io.Writer
	Stderr io.Writer
}

func NewRunner(stdout, stderr io.Writer) Runner {
	return &CommandRunner{
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (r *CommandRunner) Run(command map[string]any) error {
	exe, args, err := r.parseCommand(command)
	if err != nil {
		return err
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr

	if err := cmd.Run(); err != nil {
		return appErr.NewExecutionError("runner", err.Error(), err)
	}

	return nil
}

func (r *CommandRunner) parseCommand(command map[string]any) (string, []string, error) {
	if command == nil {
		return "", nil, appErr.NewInvalidConfigError("runner", "command configuration is nil")
	}

	exe, ok := command["exec"].(string)
	if !ok || exe == "" {
		return "", nil, appErr.NewInvalidConfigError("runner", "executable not defined or invalid")
	}

	var args []string
	if rawArgs := command["args"]; rawArgs != nil {
		if sliceAny, ok := rawArgs.([]any); ok {
			for _, v := range sliceAny {
				args = append(args, fmt.Sprintf("%v", v))
			}
		} else if sliceString, ok := rawArgs.([]string); ok {
			args = sliceString
		}
	}

	return exe, args, nil
}
