// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"gplover/dictionary"
)

type State []Translation

type Translation struct {
	english  string
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

// NewTranslator creates a new Translator instance.
func NewTranslator(dict dictionary.Dictionary, maxOutlineLen int) *Translator {
	return &Translator{
		dict: dict,
		latest: &Translation{
			english:  "",
			outline:  "",
			prev:     nil,
			replaced: nil,
		},
		// MaxHistory:    maxHistory, // Probably shouldn't let history grow forever
		maxOutlineLen: maxOutlineLen,
	}
}

func (tr *Translator) Translate(outline string) *Translation {
	newNode := &Translation{
		english:  "",
		outline:  outline,
		prev:     tr.latest,
		replaced: nil,
	}
	tr.latest = tr.getLatest(outline, newNode, 0)
	return tr.latest
}

func (tr *Translator) getLatest(outline string, node *Translation, depth int) *Translation {
	if depth < tr.maxOutlineLen {
		if node.prev != nil {
			latest := tr.getLatest(node.outline+"/"+outline, node.prev, depth+1)
			if latest != nil {
				return latest // return the longest match possible
			}
		}
		if english, ok := tr.dict.Lookup(outline); ok {
			t := &Translation{
				english:  english,
				outline:  outline,
				prev:     node.prev,
				replaced: tr.latest,
			}
			return t
		} else if depth == 0 {
			return node // return new node
		}
	}
	return nil // dont seek longer than possible matches
}

func (t *Translation) Text() string {
	if t.english != "" {
		return t.english
	} else {
		return t.outline
	}

}
