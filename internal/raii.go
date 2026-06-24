package internal

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func Start(configPath, owner string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	if owner == "" {
		owner = fmt.Sprintf("%d", os.Getppid())
	}

	err = withState(func(state *State) error {
		sess, exists := state.Sessions[cfg.Session]
		if !exists {
			running, _ := CheckCommand(cfg.Commands.Check)
			if !running {
				if err := RunCommand(cfg.Commands.Start); err != nil {
					return fmt.Errorf("start command failed: %w", err)
				}
			}
			state.Sessions[cfg.Session] = &Session{
				Name:     cfg.Session,
				Owners:   []string{owner},
				StartCmd: cfg.Commands.Start,
				StopCmd:  cfg.Commands.Stop,
				CheckCmd: cfg.Commands.Check,
			}
		} else {
			for _, o := range sess.Owners {
				if o == owner {
					return nil
				}
			}
			sess.Owners = append(sess.Owners, owner)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Spawn a watchdog that will call stop when the owner process dies.
	go spawnWatchdog(configPath, owner)

	return nil
}

func Stop(configPath, owner string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	if owner == "" {
		owner = fmt.Sprintf("%d", os.Getppid())
	}

	return withState(func(state *State) error {
		sess, exists := state.Sessions[cfg.Session]
		if !exists {
			return nil
		}

		newOwners := make([]string, 0, len(sess.Owners))
		for _, o := range sess.Owners {
			if o != owner {
				newOwners = append(newOwners, o)
			}
		}
		sess.Owners = newOwners

		if len(sess.Owners) == 0 {
			if err := RunCommand(cfg.Commands.Stop); err != nil {
				return fmt.Errorf("stop command failed: %w", err)
			}
			delete(state.Sessions, cfg.Session)
		}
		return nil
	})
}

func Check(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	running, err := CheckCommand(cfg.Commands.Check)
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("check failed: resource is not running")
	}
	return nil
}

func Watch(configPath, owner string) error {
	pid := 0
	fmt.Sscanf(owner, "%d", &pid)
	if pid <= 0 {
		return fmt.Errorf("invalid owner pid: %s", owner)
	}

	for {
		time.Sleep(2 * time.Second)
		process, err := os.FindProcess(pid)
		if err != nil {
			break
		}
		if err := process.Signal(syscall.Signal(0)); err != nil {
			break
		}
	}

	cmd := exec.Command(os.Args[0], "stop", "--config", configPath, "--owner", owner)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func spawnWatchdog(configPath, owner string) {
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer devNull.Close()

	cmd := exec.Command(os.Args[0], "watch", "--config", configPath, "--owner", owner)
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Stdin = devNull
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Start()
}
