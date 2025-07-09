// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"testing"

	"gplover/config"
)

func TestEngineTranslate(t *testing.T) {
	// Manually set the layout used for translation
	config.Layout = map[string]string{
		"S2-": "S",
		"T-":  "T",
		"K-":  "K",
		"E":   "E",
		"P-":  "P",
		"W-":  "W",
	}

	dict := map[string]string{
		"STKE": "stack",
		"TKPW": "go",
	}

	e := NewEngine(dict)

	if word := e.TranslateSteno([]string{"S2-", "T-", "K-", "E"}); word != "stack" {
		t.Errorf("expected stack, got %q", word)
	}

	if word := e.TranslateSteno([]string{"T-", "K-", "P-", "W-"}); word != "go" {
		t.Errorf("expected go, got %q", word)
	}

	if word := e.TranslateSteno([]string{"X", "Y", "Z"}); word != "[]" {
		t.Errorf("expected empty string, got %q", word)
	}
}
