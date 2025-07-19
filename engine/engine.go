// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"log"
	"sten/dictionary"
	"sten/machine"
	"sten/output"
	"sten/translator"
)

type Engine struct {
	output     output.Output
	translator *translator.Translator
}

func NewEngine() *Engine {
	// Load your dictionary
	dict, longestOutline, err := dictionary.LoadDictionaries("dictionaries")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}

	// Create virtual Output
	output := output.NewDevOutputService()
	translator := translator.NewTranslator(dict, longestOutline)

	e := &Engine{
		output:     output,
		translator: translator,
	}
	return e
}

func (e *Engine) Run(machine machine.Machine) {
	// Start a goroutine to feed strokes into translator
	go func() {
		for stroke := range machine.Strokes() {
			e.translator.Translate(stroke.Steno())
		}
		e.translator.Close() // Close when machine input is done
	}()

	// Main loop: read translation outputs and handle them
	for translation := range e.translator.Out() {
		e.Execute(translation)
	}
}

func (e *Engine) Execute(newTranslation *translator.Translation) {
	// Get all translations being replaced
	// replaced := getReplacedTranslations(newTranslation)
	// Sum the output length
	// backspaces := getBackspaceCount(replaced)
	// if backspaces > 0 {
	// 	e.output.Backspace(backspaces)
	// }
	// Output the new translation's text
	outText := newTranslation.Text()
	e.output.Type(outText)
	// Optionally update engine's history or pointer if needed
}
