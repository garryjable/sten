// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"fmt"
	"sten/stroke"
	"testing"
)

// MockDictionary is a basic stub for dictionary.Dictionary
type MockDictionary struct {
	entries map[string]string
}

func (m *MockDictionary) Lookup(outline fmt.Stringer) (string, bool) {
	val, ok := m.entries[outline.String()]
	return val, ok
}

func TestSingleWordTranslation(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"STPH": "hello",
	}}
	tr := NewTranslator(dict, 1)

	tr.Translate(stroke.ParseSteno("STPH"))
	result := <-tr.Out()

	if result.Text() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result.Text())
	}
}

func TestCommandTranslation(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"*": "=undo",
	}}
	tr := NewTranslator(dict, 1)

	tr.Translate(stroke.ParseSteno("*"))
	result := <-tr.Out()

	if !result.isCommand() {
		t.Error("Expected command translation")
	}
	if result.Text() != "=undo" {
		t.Errorf("Expected '=undo', got '%s'", result.Text())
	}
}

func TestBlankFallback(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"TPHOPB":               "known",
		"TPHOPB/TPHOPB":        "known",
		"TPHOPB/TPHOPB/TPHOPB": "known",
	}}
	tr := NewTranslator(dict, 3)

	unknown := "TPHO"
	tr.Translate(stroke.ParseSteno(unknown))
	result := <-tr.Out()

	if result.Text() != unknown {
		t.Errorf("Expected fallback to raw, got '%s'", result.Text())
	}
}

func TestMultiStrokeTranslation(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"U":                      "you",
		"R":                      "are",
		"EUPB":                   "in",
		"TE":                     "the",
		"EUPB/TE/HREB/TWAL":      "intellectual",
		"EUPB/TE/HREB/TWAL/TWAL": "Not Reachable",
	}}
	tr := NewTranslator(dict, 4)

	tr.Translate(stroke.ParseSteno("U"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("R"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("EUPB"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("TE"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("HREB"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("TWAL"))
	result := <-tr.Out()

	if result.Text() != "intellectual" {
		t.Errorf("Expected 'intellectual' from multi-stroke, got '%s'", result.Text())
	}

	tr.Translate(stroke.ParseSteno("TWAL"))
	result = <-tr.Out()
	if result.Text() != "TWAL" {
		t.Errorf("Expected 'TWAL' from exceed max stroke length, got '%s'", result.Text())
	}
}

func TestCommandDoesNotUpdateLatest(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"WOPB": "1",
		"*":    "=undo",
		"TWO":  "2",
	}}
	tr := NewTranslator(dict, 1)

	tr.Translate(stroke.ParseSteno("WOPB"))
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("*")) // should not change latest
	<-tr.Out()
	tr.Translate(stroke.ParseSteno("TWO"))
	result := <-tr.Out()

	if result.Text() != "2" {
		t.Errorf("Expected '2', got '%s'", result.Text())
	}
	if result.prev == nil || result.prev.Text() != "1" {
		t.Errorf("Expected previous translation to be '1', got '%v'", result.prev.Text())
	}
}
