package stroke

import (
	"testing"

	"gplover/config"
)

func TestStrokeTranslation(t *testing.T) {
	// Test steno translation from key codes to standard steno notation determined by layout config
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
