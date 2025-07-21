// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sten/machine"
)

type Config struct {
	Port        string            `json:"serial_port"`
	Baud        int               `json:"baud_rate"`
	ReadTimeout int               `json:"timeout"`
	Machine     string            `json:"machine"`
	CustomKeys  map[string]string `json:"custom_keys"`
}

func (cfg *Config) setCustomKeys() map[string]string {
	switch cfg.Machine {
	case "geminipr":
		for k, v := range cfg.CustomKeys {
			machine.GeminiDefaults[k] = v
		}
		return machine.GeminiDefaults
	// case "other":
	//     return SomeOtherLayout
	default:
		panic("Unknown machine type: " + cfg.Machine)
	}
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode config: %w", err)
	}
	cfg.setCustomKeys()
	return &cfg, nil
}
