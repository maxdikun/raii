package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(cmdStr string) error {
	if cmdStr == "" {
		return fmt.Errorf("command is empty")
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func CheckCommand(cmdStr string) (bool, error) {
	if cmdStr == "" {
		return false, fmt.Errorf("no check command configured")
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err == nil, nil
}
