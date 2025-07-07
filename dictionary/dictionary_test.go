// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package dictionary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDictionary(t *testing.T) {
	// Ensure the directory exists
	dir := "test_dictionaries"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	// Create and write the test dictionary file
	path := filepath.Join(dir, "test_dict.json")
	tmp := `{"STKE": "stack"}`
	err = os.WriteFile(path, []byte(tmp), 0644)
	if err != nil {
		t.Fatalf("failed to write test dict: %v", err)
	}
	defer os.Remove(path)

	// Load and test dictionary
	dict, err := LoadDictionaries(dir)
	if err != nil {
		t.Fatalf("failed to load dictionary: %v", err)
	}

	got := dict["STKE"]
	want := "stack"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
