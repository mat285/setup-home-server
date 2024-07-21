package setuptasks

import "github.com/mat285/setup-home-server/task"

func All() []*task.Task {
	return []*task.Task{
		InstallTailscale(),
		InstallZsh(),
		InstallOhMyZsh(),
		InstallGo(),
	}
}
