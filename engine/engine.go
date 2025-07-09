// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"gplover/dictionary"
	"gplover/stroke"
)

type Engine struct {
	Dict dictionary.Dictionary
}

func NewEngine(dict dictionary.Dictionary) *Engine {
	return &Engine{Dict: dict}
}

func (e *Engine) TranslateSteno(keys []string) string {
	stroke, err := stroke.NewStroke(keys)
	if err != nil {
		return "[error]"
	}
	word, ok := e.Dict.Lookup(stroke.Steno())
	if !ok {
		return "[" + stroke.Steno() + "]"
	}
	return word
}
