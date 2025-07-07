// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
