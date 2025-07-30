// Copyright (c) 2025 Garrett result.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"fmt"
	"sten/dictionary"
	"sten/output"
	"sten/stroke"
	"strings"
)

type Translation struct {
	result   Result
	outline  stroke.Outline
	prev     *Translation // previous
	replaced *Translation // store replaced translations
}

type Result struct {
	raw      string
	text     string
	replaced string // cache of replaced result
}

func (r *Result) write() output.Output {
	return output.NewOutput(r.text, r.replaced)
}

func (tr *Translator) newTranslation(raw, replaced string, outline stroke.Outline, prev *Translation) *Translation {
	fmt.Printf("translated %v \n", raw)
	if raw == "=undo" {
		return newUndo(raw, outline, prev)
	} else if strings.HasPrefix(raw, "{^}") {
		suffix := raw[3:]
		return newSuffix(raw, suffix, outline, prev)
	} else if len(outline) == 1 {
		return newSingleStroke(raw, outline, prev)
	} else {
		return newMultiStroke(raw, replaced, outline, prev, tr.latest)
	}
}

// Translator is the main engine for converting strokes to translations.
type Translator struct {
	dict       dictionary.Dict
	latest     *Translation
	outlineCap int
	in         chan stroke.Stroke
	out        chan output.Output
}

func newSingleStroke(raw string, outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   Result{raw, raw + " ", ""},
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

func newMultiStroke(raw, replaced string, outline stroke.Outline, prev, replacedPtr *Translation) *Translation {
	return &Translation{
		result:   Result{raw, raw + " ", replaced},
		outline:  outline,
		prev:     prev,
		replaced: replacedPtr,
	}
}

func newSuffix(raw, suffix string, outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   Result{raw, suffix + " ", " "},
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

func newUntranslatable(outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   Result{outline.String(), outline.String() + " ", ""},
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

func newUndo(raw string, outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   Result{raw, prev.result.replaced, prev.result.text},
		outline:  outline,
		prev:     nil,
		replaced: nil,
	}
}

func newBlank() *Translation {
	return &Translation{
		result:   Result{"", "", ""},
		outline:  stroke.Outline{},
		prev:     nil,
		replaced: nil,
	}
}

// NewTranslator creates a new Translator instance.
func NewTranslator(dict dictionary.Dict, outlineCap int, in chan stroke.Stroke) *Translator {
	t := &Translator{
		dict: dict,
		latest: &Translation{
			result:   Result{"", "", ""},
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
	fmt.Printf("outline %v, replacing %s \n", outline, replacing)

	if len(outline) > tr.outlineCap {
		return nil // too deep to match
	}

	if prev.prev != nil {
		translation := tr.translate(append(prev.outline, outline...), prev.prev, prev.result.text+replacing)
		if translation != nil {
			return translation // return the longest possible translation
		}
	}

	if entry, ok := tr.dict.Lookup(outline); ok {
		return tr.newTranslation(entry, replacing, outline, prev)
	}

	if len(outline) == 1 {
		return newUntranslatable(outline, prev)
	}

	return nil // unreachable
}

func (tr *Translator) updateHistory(latest *Translation) {
	if latest.result.raw == "=undo" {
		tr.undo()
	} else {
		tr.latest = latest
	}
}

func (tr *Translator) undo() {
	if tr.latest.replaced != nil {
		tr.latest = tr.latest.replaced
	} else if tr.latest.prev != nil {
		tr.latest = tr.latest.prev
	}
}

func (tr *Translator) Out() chan output.Output {
	return tr.out
}

func (tr *Translator) Run() {
	for stroke := range tr.in {
		latest := tr.translate(stroke.Outline(), tr.latest, "")
		tr.updateHistory(latest)
		tr.out <- latest.result.write()
	}
	close(tr.out)
}
