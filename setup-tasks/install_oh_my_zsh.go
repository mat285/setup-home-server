package setuptasks

import "github.com/mat285/setup-home-server/task"

func InstallOhMyZsh() *task.Task {
	return &task.Task{
		Name:     "install oh my zsh",
		Required: []string{InstallZsh().Name},
		Check:    isOhMyZshInstalled(),
		Job:      installOhMyZsh(),
	}
}

func isOhMyZshInstalled() *task.CheckCommand {
	return task.NewCheckCommand(
		"sh",
		[]string{
			"-c",
			"stat",
			"$HOME/.oh-my-zsh",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}

func installOhMyZsh() *task.Command {
	return task.NewCommand(
		"sh",
		[]string{
			"-c",
			"\"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\"",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}
