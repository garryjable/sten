package test

import (
	"testing"

	"gplover/config"
	"gplover/engine"
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

	e := engine.NewEngine(dict)

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
