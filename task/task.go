package task

import (
	"context"
	"fmt"
)

type Task struct {
	Name     string
	Required []string
	Check    *CheckCommand
	Job      *Command
}

func (t *Task) Run(ctx context.Context) (bool, error) {
	fmt.Println("Running task", t.Name)
	if t.Job == nil {
		fmt.Println("no job configured")
		return true, nil
	}
	should, err := t.shouldRun(ctx)
	if err != nil {
		return false, err
	}
	if !should {
		fmt.Println("Task", t.Name, "already run successfully, skipping run")
		return true, nil
	}
	t.Job.Task = t.Name
	return t.Job.Run(ctx)
}

func (t *Task) shouldRun(ctx context.Context) (bool, error) {
	if t.Check == nil {
		return true, nil
	}
	t.Check.Task = t.Name
	exists, err := t.Check.Run(ctx)
	if err != nil {
		return false, err
	}
	return !exists, nil
}
