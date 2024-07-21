package task

import (
	"context"
	"fmt"
)

type Runner struct {
	tasks    map[string]*Task
	run      map[string]bool
	children map[string][]*Task
	start    []*Task
}

func NewRunner(tasks ...*Task) (*Runner, error) {
	r := &Runner{
		tasks:    make(map[string]*Task),
		run:      make(map[string]bool),
		children: make(map[string][]*Task),
	}

	for _, task := range tasks {
		r.tasks[task.Name] = task
		if len(task.Required) == 0 {
			r.start = append(r.start, task)
			continue
		}
		for _, dep := range task.Required {
			r.children[dep] = append(r.children[dep], task)
		}
	}

	if len(r.start) == 0 {
		return nil, fmt.Errorf("no tasks can be started")
	}
	// TODO check for cycles

	return r, nil
}

func (r *Runner) Run(ctx context.Context) (bool, error) {
	stack := append([]*Task{}, r.start...)
	status := make(map[string]bool)
	errs := make(map[string]error)

	for len(stack) > 0 {
		task := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Println()
		success, err := task.Run(ctx)
		fmt.Println()
		status[task.Name] = success
		if err != nil {
			errs[task.Name] = err
			continue
		}
		if !success {
			continue
		}
		children := r.children[task.Name]
		for _, child := range children {
			if canRunTask(status, child) {
				stack = append(stack, child)
			}
		}
	}

	if len(errs) > 0 {
		return false, fmt.Errorf("error running tasks %q", errs)
	}

	for task := range r.tasks {
		if !status[task] {
			return false, fmt.Errorf("not all tasks succeeded %q", status)
		}
	}
	return true, nil
}

func canRunTask(status map[string]bool, task *Task) bool {
	for _, dep := range task.Required {
		if !status[dep] {
			return false
		}
	}
	return true
}
