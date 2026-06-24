package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type State struct {
	Sessions map[string]*Session `json:"sessions"`
}

type Session struct {
	Name     string   `json:"name"`
	Owners   []string `json:"owners"`
	StartCmd string   `json:"start_cmd"`
	StopCmd  string   `json:"stop_cmd"`
	CheckCmd string   `json:"check_cmd"`
}

var statePath string

func init() {
	dir := os.Getenv("RAII_STATE_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "/tmp"
		}
		dir = filepath.Join(home, ".local", "share", "raii")
	}
	os.MkdirAll(dir, 0755)
	statePath = filepath.Join(dir, "state.json")
}

func loadState() (*State, error) {
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Sessions: make(map[string]*Session)}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}
	if s.Sessions == nil {
		s.Sessions = make(map[string]*Session)
	}
	return &s, nil
}

func saveState(s *State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0644)
}

func withState(f func(*State) error) error {
	lockPath := statePath + ".lock"
	lockFile, err := os.Create(lockPath)
	if err != nil {
		return fmt.Errorf("failed to create lock: %w", err)
	}
	defer os.Remove(lockPath)
	defer lockFile.Close()

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	state, err := loadState()
	if err != nil {
		return err
	}
	if err := f(state); err != nil {
		return err
	}
	return saveState(state)
}
