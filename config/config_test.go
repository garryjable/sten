// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	data := `{
		"serial_port": "/dev/testport",
		"machine": "geminipr"
	}`

	err := os.WriteFile("test_config.json", []byte(data), 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	defer os.Remove("test_config.json")

	cfg, err := Load("test_config.json")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Port != "/dev/testport" {
		t.Errorf("config port failed to parse")
	}
}
