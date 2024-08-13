package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the structure of our configuration file
type Config struct {
	Processes []ProcessConfig `json:"processes"`
	Commands  []ProcessConfig `json:"commands"`
}

// ProcessConfig represents the configuration for a single process
type ProcessConfig struct {
	Shortname   string `json:"shortname"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// Load reads and parses the configuration file at the given path
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validate checks if the loaded configuration is valid
func validate(cfg *Config) error {
	if len(cfg.Processes) == 0 {
		return fmt.Errorf("no processes defined in configuration")
	}

	for i, proc := range cfg.Processes {
		if proc.Shortname == "" {
			return fmt.Errorf("process %d is missing a shortname", i+1)
		}
		if proc.Command == "" {
			return fmt.Errorf("process %d (%s) is missing a command", i+1, proc.Shortname)
		}
	}

	return nil
}
