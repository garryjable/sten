package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Layout map[string]string

type Config struct {
	Port        string            `json:"serial_port"`
	Baud        int               `json:"baud_rate"`
	ReadTimeout int               `json:"timeout"`
	Layout      map[string]string `json:"layout"`
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
	Layout = cfg.Layout
	return &cfg, nil
}
