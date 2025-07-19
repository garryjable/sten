// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"testing"
)

// MockDictionary is a basic stub for dictionary.Dictionary
type MockDictionary struct {
	entries map[string]string
}

func (m *MockDictionary) Lookup(outline string) (string, bool) {
	val, ok := m.entries[outline]
	return val, ok
}

func TestSingleWordTranslation(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"STROKE": "hello",
	}}
	tr := NewTranslator(dict, 1)

	tr.Translate("STROKE")
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

	tr.Translate("*")
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
		"KNOWN":             "known",
		"KNOWN/KNOWN":       "known",
		"KNOWN/KNOWN/KNOWN": "known",
	}}
	tr := NewTranslator(dict, 3)

	tr.Translate("UNKNOWN")
	result := <-tr.Out()

	if result.Text() != "UNKNOWN" {
		t.Errorf("Expected fallback to 'UNKNOWN', got '%s'", result.Text())
	}
}

func TestMultiStrokeTranslation(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"U":                            "you",
		"R":                            "are",
		"EUPB":                         "in",
		"TE":                           "the",
		"EUPB/TE/HREB/TWAL":            "intellectual",
		"EUPB/TE/HREB/TWAL/MAXOUTLINE": "Not Reachable",
	}}
	tr := NewTranslator(dict, 4)

	tr.Translate("U")
	<-tr.Out()
	tr.Translate("R")
	<-tr.Out()
	tr.Translate("EUPB")
	<-tr.Out()
	tr.Translate("TE")
	<-tr.Out()
	tr.Translate("HREB")
	<-tr.Out()
	tr.Translate("TWAL")
	result := <-tr.Out()

	if result.Text() != "intellectual" {
		t.Errorf("Expected 'intellectual' from multi-stroke, got '%s'", result.Text())
	}

	tr.Translate("MAXOUTLINE")
	result = <-tr.Out()
	if result.Text() != "MAXOUTLINE" {
		t.Errorf("Expected 'MAXOUTLINE' from exceed max stroke length, got '%s'", result.Text())
	}
}

func TestCommandDoesNotUpdateLatest(t *testing.T) {
	dict := &MockDictionary{entries: map[string]string{
		"ONE": "1",
		"CMD": "=do",
		"TWO": "2",
	}}
	tr := NewTranslator(dict, 1)

	tr.Translate("ONE")
	<-tr.Out()
	tr.Translate("CMD") // should not change latest
	<-tr.Out()
	tr.Translate("TWO")
	result := <-tr.Out()

	if result.Text() != "2" {
		t.Errorf("Expected '2', got '%s'", result.Text())
	}
	if result.prev == nil || result.prev.Text() != "1" {
		t.Errorf("Expected previous translation to be '1', got '%v'", result.prev)
	}
}
