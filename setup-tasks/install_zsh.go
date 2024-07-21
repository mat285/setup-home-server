package setuptasks

import "github.com/mat285/setup-home-server/task"

func InstallZsh() *task.Task {
	return &task.Task{
		Name:  "install zsh",
		Check: isZshInstalled(),
		Job:   installZsh(),
	}
}

func isZshInstalled() *task.CheckCommand {
	return task.NewCheckCommand(
		"command",
		[]string{
			"-v",
			"zsh",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}

func installZsh() *task.Command {
	return task.NewCommand(
		"sudo",
		[]string{
			"apt",
			"install",
			"zsh",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}
