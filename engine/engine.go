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

func (e *Engine) TranslateSteno(strokeText string) string {
	stroke, err := stroke.NewStroke(strokeText)
	if err != nil {
		return "[error]"
	}
	word, ok := e.Dict.Lookup(stroke.Steno())
	if !ok {
		return "[" + stroke.Steno() + "]"
	}
	return word
}
