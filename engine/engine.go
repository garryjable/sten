// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"gplover/output"
	"gplover/translator"
)

type Engine struct {
	output output.Output
}

func NewEngine(output output.Output) *Engine {
	e := &Engine{
		output: output,
	}
	return e
}

func (e *Engine) Execute(translation *translator.Translation) {
	e.output.Type(translation.Text())
}
