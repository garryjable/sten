// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"testing"
)

func TestStrokeTranslation(t *testing.T) {
	s := ParseSteno("STKES")
	got := s.Steno()
	want := "STKES"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}

	s = ParseSteno("ES")
	got = s.Steno()
	want = "ES"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}

	s = ParseSteno("-S")
	got = s.Steno()
	want = "-S"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}

}
