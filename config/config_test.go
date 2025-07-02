package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	data := `{
		"serial_port": "/dev/testport",
		"layout": { "S1-": "S", "T-": "T" }
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

	if Layout["S1-"] != "S" {
		t.Errorf("expected S1- to map to S")
	}
	if Layout["T-"] != "T" {
		t.Errorf("expected T- to map to T")
	}
}
