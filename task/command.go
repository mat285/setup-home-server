package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExitCheck func(ExitCondition) (bool, error)

type ExitCondition struct {
	Code   int
	StdErr []byte
	StdOut []byte
}

type Command struct {
	Cmd        string
	Args       []string
	WorkingDir string
	Success    ExitCheck
	Task       string
}

type CheckCommand struct {
	Command
	Exists ExitCheck
}

func NewCommand(cmd string, args []string, wd string, success ExitCheck) *Command {
	return &Command{
		Cmd:        cmd,
		Args:       args,
		WorkingDir: wd,
		Success:    success,
	}
}

func NewCheckCommand(cmd string, args []string, wd string, exists ExitCheck) *CheckCommand {
	return &CheckCommand{
		Command: Command{
			Cmd:        cmd,
			Args:       args,
			WorkingDir: wd,
		},
		Exists: exists,
	}
}

func (c *CheckCommand) Run(ctx context.Context) (bool, error) {
	ec, err := c.RunToExit(ctx)
	if err != nil {
		return false, err
	}
	return c.Exists(*ec)
}

func (c *Command) Run(ctx context.Context) (bool, error) {
	ec, err := c.RunToExit(ctx)
	if err != nil {
		return false, err
	}
	return c.Success(*ec)
}

func (c *Command) RunToExit(ctx context.Context) (*ExitCondition, error) {
	fmt.Println()
	cmd := exec.CommandContext(ctx, c.Cmd, c.Args...)
	cmd.Dir = c.WorkingDir
	cmd.Env = os.Environ()

	prepend := "[ " + c.cmdStr() + " ]  "
	if len(c.Task) > 0 {
		prepend = fmt.Sprintf("[ Task: %q ] %s", c.Task, prepend)
	}

	capOut := NewCaptureWriter(os.Stdout, []byte(prepend))
	capErr := NewCaptureWriter(os.Stderr, []byte(prepend))

	cmd.Stdout = capOut
	cmd.Stderr = capErr

	err := cmd.Run()
	if err != nil && !errors.Is(err, &exec.ExitError{}) {
		return nil, err
	}
	fmt.Println()
	ec := &ExitCondition{
		StdOut: capOut.Captured(),
		StdErr: capErr.Captured(),
	}
	if cmd.ProcessState != nil {
		ec.Code = cmd.ProcessState.ExitCode()
	}
	return ec, nil
}

func (c *Command) cmdStr() string {
	return strings.Join(append([]string{c.Cmd}, c.Args...), " ")
}
