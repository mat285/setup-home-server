package setuptasks

import "github.com/mat285/setup-home-server/task"

func InstallGo() *task.Task {
	return &task.Task{
		Name:  "install Go",
		Check: isGoInstalled(),
		Job:   installGo(),
	}
}

func isGoInstalled() *task.CheckCommand {
	return task.NewCheckCommand(
		"command",
		[]string{
			"-v",
			"go",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}

func installGo() *task.Command {
	return task.NewCommand(
		"sudo",
		[]string{
			"snap",
			"install",
			"go",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}
