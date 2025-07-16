// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"gplover/dictionary"
	"strings"
)

type State []Translation

type Translation struct {
	entry    string
	outline  string
	prev     *Translation
	replaced *Translation
}

// Translator is the main engine for converting strokes to translations.
type Translator struct {
	dict   dictionary.Dictionary
	latest *Translation
	// UndoBuffer     []Translation
	// Listeners      []func([]Translation, []Translation, *Translation)
	// MaxHistory    int
	maxOutlineLen int
}

func newCommand(entry string, outline string) *Translation {
	return &Translation{
		entry:    entry,
		outline:  outline,
		prev:     nil,
		replaced: nil,
	}
}

func newWord(entry string, outline string, prev *Translation, latest *Translation) *Translation {
	return &Translation{
		entry:    entry,
		outline:  outline,
		prev:     prev,
		replaced: latest,
	}
}

func newBlank(outline string, prev *Translation) *Translation {
	return &Translation{
		entry:    "",
		outline:  outline,
		prev:     prev,
		replaced: nil,
	}
}

// NewTranslator creates a new Translator instance.
func NewTranslator(dict dictionary.Dictionary, maxOutlineLen int) *Translator {
	return &Translator{
		dict: dict,
		latest: &Translation{
			entry:    "",
			outline:  "",
			prev:     nil,
			replaced: nil,
		},
		// MaxHistory:    maxHistory, // Probably shouldn't let history grow forever
		maxOutlineLen: maxOutlineLen,
	}
}

func (tr *Translator) Translate(outline string) *Translation {
	translation := newBlank(outline, tr.latest)
	latest := tr.getLatest(outline, translation, 0)
	if !latest.isCommand() {
		tr.latest = latest
	}
	return latest
}

func (tr *Translator) getLatest(outline string, node *Translation, depth int) *Translation {
	if depth < tr.maxOutlineLen {
		if node.prev != nil {
			latest := tr.getLatest(node.outline+"/"+outline, node.prev, depth+1)
			if latest != nil {
				return latest // return the longest match possible
			}
		}
		if entry, ok := tr.dict.Lookup(outline); ok {
			if strings.HasPrefix(entry, "=") {
				return newCommand(entry, outline)
			} else {
				return newWord(entry, outline, node.prev, tr.latest)
			}
		} else if depth == 0 {
			return node // return blank translation
		}
	}
	return nil // dont seek longer than possible matches
}

func (t *Translation) Text() string {
	if t.entry != "" {
		return t.entry
	} else {
		return t.outline
	}

}

func (t *Translation) isCommand() bool {
	if strings.HasPrefix(t.entry, "=") {
		return true
	} else {
		return false
	}

}
