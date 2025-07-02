package test

import (
	"testing"

	"gplover/config"
	"gplover/stroke"
)

func TestStrokeTranslation(t *testing.T) {
	config.Layout = map[string]string{
		"S1-": "S",
		"T-":  "T",
		"K-":  "K",
		"-E":  "E",
	}

	s := &stroke.Stroke{
		Keys: []string{"S1-", "T-", "K-", "-E"},
	}
	got := s.Steno()
	want := "STKE"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
