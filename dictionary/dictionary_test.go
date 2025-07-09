// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

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
