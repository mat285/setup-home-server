package directlan

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	SecondaryInterfacePrefix = "enxa0"
	SecondaryLanIpPrefix     = "192.168.0."
	SecondaryLanCidr         = SecondaryLanIpPrefix + "0/24"
)

func SetupDirectLan(ctx context.Context, node int) error {
	iface, err := FindSecondaryNetworkInterface(ctx)
	if err != nil {
		return fmt.Errorf("failed to find secondary network interface: %w", err)
	}
	if err := SetInterfaceDown(ctx, iface); err != nil {
		return fmt.Errorf("failed to set interface down: %w", err)
	}
	if err := AssignLanIp(ctx, iface, node); err != nil {
		return fmt.Errorf("failed to assign lan ip: %w", err)
	}
	if err := AssignLanRoutes(ctx, iface); err != nil {
		return fmt.Errorf("failed to assign lan routes: %w", err)
	}
	if err := SetInterfaceUp(ctx, iface); err != nil {
		return fmt.Errorf("failed to set interface up: %w", err)
	}
	return nil
}

func FindSecondaryNetworkInterface(ctx context.Context) (string, error) {
	output, err := exec.CommandContext(ctx,
		"ip",
		"link",
		"list",
	).CombinedOutput()
	if err != nil {
		return "", err
	}
	idx := strings.Index(string(output), SecondaryInterfacePrefix)
	if idx < 0 {
		return "", fmt.Errorf("no secondary interface found")
	}
	cut := string(output)[idx:]
	idx = strings.Index(cut, ":")
	if idx < 0 {
		return "", fmt.Errorf("no secondary interface found")
	}
	return cut[:idx], nil
}

func AssignLanIp(ctx context.Context, interfaceName string, node int) error {
	cmd := exec.CommandContext(ctx,
		"ip",
		"addr",
		"add",
		fmt.Sprintf("%s%d", SecondaryLanIpPrefix, node),
		"dev",
		interfaceName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AssignLanRoutes(ctx context.Context, interfaceName string) error {
	cmd := exec.CommandContext(ctx,
		"ip",
		"route",
		"add",
		SecondaryLanCidr,
		"dev",
		interfaceName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetInterfaceDown(ctx context.Context, interfaceName string) error {
	cmd := exec.CommandContext(ctx,
		"ip",
		"link",
		"set",
		interfaceName,
		"down",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetInterfaceUp(ctx context.Context, interfaceName string) error {
	cmd := exec.CommandContext(ctx,
		"ip",
		"link",
		"set",
		interfaceName,
		"up",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
