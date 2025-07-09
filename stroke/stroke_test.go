// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"testing"

	"gplover/config"
)

func TestStrokeTranslation(t *testing.T) {
	config.Layout = map[string]string{
		"S1-": "S",
		"T-":  "T",
		"K-":  "K",
		"-E":  "E",
	}

	s := &Stroke{
		Keys: []string{"S1-", "T-", "K-", "-E"},
	}
	got := s.Steno()
	want := "STKE"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
