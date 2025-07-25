// Copyright (c) 2025 Garrett result.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"sten/dictionary"
	"sten/output"
	"sten/stroke"
	"strings"
)

type Translation struct {
	result   Result
	outline  stroke.Outline
	prev     *Translation // previous
	replaced *Translation // previous stroke in multi
}

type Result struct {
	raw      string
	text     string
	replaced string // multi
}

// Translator is the main engine for converting strokes to translations.
type Translator struct {
	dict       dictionary.Dict
	latest     *Translation
	outlineCap int
	in         chan stroke.Stroke
	out        chan output.Output
}

func newResult(raw, replaced string) Result {
	return Result{
		raw:      raw,
		text:     raw + " ",
		replaced: replaced,
	}
}

func newCommand(result Result, outline stroke.Outline) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     nil,
		replaced: nil,
	}
}

func newSingleStroke(result Result, outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

func newMultiStroke(result Result, outline stroke.Outline, prev, replaced *Translation) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     prev,
		replaced: replaced,
	}
}

func newUntranslatable(outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   newResult(outline.String(), ""),
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

// NewTranslator creates a new Translator instance.
func NewTranslator(dict dictionary.Dict, outlineCap int, in chan stroke.Stroke) *Translator {
	t := &Translator{
		dict: dict,
		latest: &Translation{
			result:   newResult("", ""),
			outline:  stroke.Outline{},
			prev:     nil,
			replaced: nil,
		},
		outlineCap: outlineCap,
		in:         in,
		out:        make(chan output.Output, 16),
	}
	return t
}

// provides the longest possible match
func (tr *Translator) translate(outline stroke.Outline, prev *Translation, replacing string) *Translation {
	if len(outline) > tr.outlineCap {
		return nil // too deep to match
	}

	if prev.prev != nil {
		latest := tr.translate(append(prev.outline, outline...), prev.prev, prev.result.text+replacing)
		if latest != nil {
			return latest // return the longest possible match
		}
	}

	if text, ok := tr.dict.Lookup(outline); ok {
		result := newResult(text, replacing)
		if strings.HasPrefix(text, "=") {
			return newCommand(result, outline)
		} else if len(outline) == 1 {
			return newSingleStroke(result, outline, prev)
		} else {
			return newMultiStroke(result, outline, prev, tr.latest)
		}
	}

	if len(outline) == 1 {
		return newUntranslatable(outline, prev)
	}

	return nil // unreachable
}

// pops most recent translation returns result()
func (t *Translator) undo() output.Output {
	if t.latest.prev == nil {
		return output.NewOutput("", "") // noop
	}

	write := t.latest.result.replaced
	undo := t.latest.result.text

	if t.latest.replaced != nil {
		t.latest = t.latest.replaced
	} else {
		t.latest = t.latest.prev
	}

	return output.NewOutput(write, undo)
}

func (t *Translation) isCommand() bool {
	if strings.HasPrefix(t.result.raw, "=") {
		return true
	} else {
		return false
	}
}

func (tr *Translator) appendHistory(latest *Translation) {
	if !latest.isCommand() {
		tr.latest = latest
	}
}

func (t *Translator) Out() chan output.Output {
	return t.out
}

func (t *Translator) Run() {
	for stroke := range t.in {
		latest := t.translate(stroke.Outline(), t.latest, "")
		t.appendHistory(latest)
		if latest.isCommand() {
			if latest.result.raw == "=undo" {
				t.out <- t.undo()
			}
		} else {
			t.out <- output.NewOutput(latest.result.text, latest.result.replaced)
		}
	}
	close(t.out)
}
