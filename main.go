package main

import (
	"context"
	"fmt"

	setuptasks "github.com/mat285/setup-home-server/setup-tasks"
	"github.com/mat285/setup-home-server/task"
)

func main() {
	runner, err := task.NewRunner(setuptasks.All()...)
	if err != nil {
		panic(err)
	}
	success, err := runner.Run(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Command succesful:", success)
}
