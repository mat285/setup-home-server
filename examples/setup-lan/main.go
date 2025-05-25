package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	directlan "github.com/mat285/setup-home-server/setup-tasks/direct-lan"
)

func main() {
	host, err := os.Hostname()
	if err != nil || host == "" {
		fmt.Println("Error: Unable to retrieve hostname.")
		os.Exit(1)
	}
	node, err := strconv.Atoi(strings.TrimPrefix(host, "node"))
	if err != nil {
		fmt.Printf("Error: Unable to parse node number from hostname '%s': %v\n", host, err)
		os.Exit(1)
	}
	if err := directlan.SetupDirectLan(context.Background(), node); err != nil {
		fmt.Printf("Error setting up direct LAN: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Setup completed successfully.")
}
