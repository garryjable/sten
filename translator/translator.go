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
	result   string
	outline  stroke.Outline
	prev     *Translation // previous
	replaced *Translation // latest prior to multistroke absorb
}

// Translator is the main engine for converting strokes to translations.
type Translator struct {
	dict       dictionary.Dict
	latest     *Translation
	outlineCap int
	in         chan stroke.Stroke
	out        chan output.Output
}

func newCommand(result string, outline stroke.Outline) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     nil,
		replaced: nil,
	}
}

func newWord(result string, outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

func newMultiWord(result string, outline stroke.Outline, prev *Translation, replaced *Translation) *Translation {
	return &Translation{
		result:   result,
		outline:  outline,
		prev:     prev,
		replaced: replaced,
	}
}

func newUntranslatable(outline stroke.Outline, prev *Translation) *Translation {
	return &Translation{
		result:   "",
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
			result:   "",
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

func (tr *Translation) PrintHistory() {
	if tr.prev != nil {
		tr.prev.PrintHistory()
	}
}

// provides the longest possible match
func (tr *Translator) getLatest(outline stroke.Outline, prev *Translation) *Translation {
	if len(outline) > tr.outlineCap {
		return nil // too deep to match
	}

	if prev.prev != nil {
		latest := tr.getLatest(append(prev.outline, outline...), prev.prev)
		if latest != nil {
			return latest // return the longest possible match
		}
	}

	if result, ok := tr.dict.Lookup(outline); ok {
		if strings.HasPrefix(result, "=") {
			return newCommand(result, outline)
		} else if len(outline) == 1 {
			return newWord(result, outline, prev)
		} else {
			return newMultiWord(result, outline, prev, tr.latest)
		}
	}

	if len(outline) == 1 {
		return newUntranslatable(outline, prev)
	}

	return nil // unreachable
}

func (t *Translator) gatherReplaced(current, until *Translation) string {
	if current == nil || current == until {
		return ""
	}
	return t.gatherReplaced(current.Prev(), until) + current.Text()
}

func (t *Translation) Text() string {
	if t.result != "" {
		return t.result + " "
	} else {
		return t.outline.String() + " "
	}

}

func (t *Translation) Prev() *Translation {
	if t.prev == nil {
		return newUntranslatable(stroke.Outline{0}, nil)
	}
	return t.prev
}

func (t *Translation) Replaced() *Translation {
	if t.replaced == nil {
		return t.prev
	}
	return t.replaced
}

func (t *Translator) Latest() *Translation {
	if t.latest == nil {
		return newUntranslatable(stroke.Outline{0}, nil)
	}
	return t.latest
}

// pops most recent translation
func (t *Translator) Undo() *Translation {
	if t.latest.prev == nil {
		return newUntranslatable(stroke.Outline{0}, nil)
	}
	latest := t.latest
	t.latest = latest.Replaced()
	return latest
}

func (t *Translation) IsCommand() bool {
	if strings.HasPrefix(t.result, "=") {
		return true
	} else {
		return false
	}
}

func (t *Translation) IsMulti() bool {
	if t.replaced != nil {
		return true
	} else {
		return false
	}
}

func (tr *Translator) appendHistory(latest *Translation) {
	if !latest.IsCommand() {
		tr.latest = latest
	}
}

// For engine to send strokes:
func (t *Translator) Translate(stroke stroke.Stroke) {
	t.in <- stroke
}

func (t *Translator) Out() chan output.Output {
	return t.out
}

func (t *Translator) Run() {
	for stroke := range t.in {
		latest := t.getLatest(stroke.Outline(), t.latest)
		t.appendHistory(latest)
		if latest.IsCommand() {
			if latest.result == "=undo" {
				deleted := t.Undo()
				t.out <- output.NewUndo(deleted.Text())
				if deleted.replaced != nil {
					t.out <- output.NewWrite(t.gatherReplaced(deleted.replaced, deleted.prev))
				}
			}
		} else {
			if latest.replaced != nil {
				t.out <- output.NewUndo(t.gatherReplaced(latest.replaced, latest.prev))
			}
			t.out <- output.NewWrite(latest.Text())
		}
	}
	close(t.out)
}
