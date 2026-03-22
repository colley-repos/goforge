// Package config provides game configuration loading from JSON/TOML files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load reads a JSON config file into the given struct.
func Load(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: read %q: %w", path, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("config: parse %q: %w", path, err)
	}

	return nil
}

// Save writes a struct to a JSON config file with indentation.
func Save(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("config: write %q: %w", path, err)
	}

	return nil
}

// LoadWithDefaults reads a JSON config, falling back to defaults if the file doesn't exist.
// If the file doesn't exist, it creates it with the default values.
func LoadWithDefaults(path string, target any) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist — write defaults and return
		return Save(path, target)
	}
	return Load(path, target)
}
