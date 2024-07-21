package setuptasks

import "github.com/mat285/setup-home-server/task"

func InstallTailscale() *task.Task {
	return &task.Task{
		Name:  "install Tailscale",
		Check: isTailscaleInstalled(),
		Job:   installTailscale(),
	}
}

func isTailscaleInstalled() *task.CheckCommand {
	return task.NewCheckCommand(
		"command",
		[]string{
			"-v",
			"tailscale",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}

func installTailscale() *task.Command {
	return task.NewCommand(
		"sh",
		[]string{
			"-c",
			"\"$(curl -fsSL https://tailscale.com/install.sh | sh)\"",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}
