// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"sten/machine"
	"sten/output"
	"sten/translator"
)

type Engine struct {
	output     output.Output
	translator *translator.Translator
}

func NewEngine(output output.Output) *Engine {
	e := &Engine{
		output: output,
	}
	return e
}

func (e *Engine) Run(machine machine.Machine, translator *translator.Translator) {
	for stroke := range machine.Strokes() {
		translation := translator.Translate(stroke.Steno())
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
