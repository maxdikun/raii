package internal

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Session  string   `toml:"session"`
	Commands Commands `toml:"commands"`
}

type Commands struct {
	Start string `toml:"start"`
	Stop  string `toml:"stop"`
	Check string `toml:"check"`
}

func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(absPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config %s: %w", path, err)
	}

	if cfg.Session == "" {
		cfg.Session = filepath.Dir(absPath)
	}

	return &cfg, nil
}
