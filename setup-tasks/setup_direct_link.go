package setuptasks

import "github.com/mat285/setup-home-server/task"

func SetupDirectLink() *task.Task {
	return &task.Task{
		Name:  "setup direct lan link",
		Check: isDirectLanSetup(),
		Job:   setupDirectLanLinks(),
	}
}

func isDirectLanSetup() *task.CheckCommand {
	return task.NewCheckCommand(
		"sudo",
		[]string{
			"systemctl",
			"status",
			"linklandaemon.service",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}

func setupDirectLanLinks() *task.Command {
	return task.NewCommand(
		"sh",
		[]string{
			"-c",
			"\"$(curl -fsSL https://github.com/mat285/linklan/releases/download/v0.1.0/install.sh)\"",
		},
		".",
		func(ec task.ExitCondition) (bool, error) {
			return ec.Code == 0, nil
		},
	)
}
