package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	machines = []string{
		"node0",
		"node1",
		"node2",
		"node3",
		"node4",
		"node5",
		"node6",
		"node7",
		"worker0",
		"zstation",
	}

	version = "v0.2.0" // Update this to the latest version of your script

	commands = []string{
		`sudo sh -c "groupadd admin"`,
		`sudo sh -c "usermod -aG admin michael"`,
		`sudo sh -c "echo >> /etc/sudoers && echo 'michael ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers"`,
		`sudo sh -c "apt install -y zsh"`,
		`sudo sh -c "apt install -y btop iperf iperf3"`,
		`sudo sh -c "snap install --classic go"`,
		`sudo sh -c "snap install --classic kubectl"`,
		`sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`,
		`sh -c "curl -fsSL https://tailscale.com/install.sh | sh"`,
		`sh -c "$(curl -fsSL https://github.com/mat285/linklan/releases/download/` + version + `/install.sh)"`,
	}
)

func main() {
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	for _, machine := range machines {
		wg.Add(1)
		func(machine string) {
			defer wg.Done()
			lock.Lock()
			fmt.Println("Setting up machine:", machine)
			lock.Unlock()
			for _, cmdStr := range commands {
				lock.Lock()
				fmt.Println("Running cmd:", cmdStr, "for machine:", machine)
				lock.Unlock()
				cmd := exec.Command("ssh",
					"-A",
					machine,
					"-T",
					cmdStr,
				)
				cmd.Env = append(os.Environ(), `SUDO_OPTS="-S"`)
				output, err := cmd.CombinedOutput()
				lock.Lock()
				fmt.Println(string(output))
				lock.Unlock()
				if err != nil {
					fmt.Printf("Error running command on %s: %v\n", machine, err)
					continue
				}
			}
			lock.Lock()
			fmt.Printf("Successfully updated %s\n", machine)
			lock.Unlock()
		}(machine)
	}
	wg.Wait()
}
