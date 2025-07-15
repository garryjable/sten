package translator

import (
	"gplover/dictionary"
	"gplover/stroke"
)

// Translation represents a successful dictionary lookup for one or more strokes.
type Translation struct {
	Strokes  []stroke.Stroke
	Outline  string
	English  string
	Replaced State
}

type State []Translation

//type Translation struct {
//	English      string
//	Outline		 string
//}
//
//
//type Translation struct {
//	Entries []*TranslationStrokes // actual stroke objects
//	Outline string
//	English string
//	Maybe Replaces []*Translation if you want undo
//}

///type State struct {
///	History []*StrokeEntry
///}

// Translator is the main engine for converting strokes to translations.
type Translator struct {
	Dict  dictionary.Dictionary
	State State
	// UndoBuffer     []Translation
	// Listeners      []func([]Translation, []Translation, *Translation)
	MaxHistory    int
	MaxOutlineLen int
}

// NewTranslator creates a new Translator instance.
func NewTranslator(dict dictionary.Dictionary, maxHistory int, maxOutlineLen int) *Translator {
	return &Translator{
		Dict:          dict,
		State:         State{},
		MaxHistory:    maxHistory,
		MaxOutlineLen: maxOutlineLen,
	}
}

func (tr *Translator) Translate(s *stroke.Stroke) Translation {
	// Build recent strokes including the current one (max length = MaxOutlineLen)
	recentState := tr.recentState() // need to append stroke
	recents := recentState.toStrokes()

	// Try matching longest possible stroke sequence to shortest (greedy match)
	for i := 0; i <= len(recents); i++ {
		candidate := recents[i:] // last N-i strokes
		outline := stroke.ToRtfcre(append(candidate, *s))

		if eng, ok := tr.Dict.Lookup(outline); ok {
			t := Translation{
				Strokes:  candidate,
				Outline:  outline,
				English:  eng,
				Replaced: recentState,
			}
			tr.applyTranslation(t)
			return t
		}
	}

	// No match: emit raw stroke
	t := Translation{
		Strokes: []stroke.Stroke{*s},
		Outline: stroke.ToRtfcre([]stroke.Stroke{*s}),
	}
	tr.applyTranslation(t)
	return t
}

func (tr *Translator) applyTranslation(t Translation) {
	// Remove replaced entries from history and append new entry
	trimLen := len(tr.State) - len(t.Replaced)
	tr.State = append(tr.State[:trimLen], t)
}

// recentStrokes returns the most recent strokes from translation history,
// plus the current stroke, trimmed to at most n strokes total.
func (tr *Translator) recentState() State {
	var result []Translation
	strokeCount := 1

	for i := len(tr.State) - 1; i >= 0; i-- {
		t := tr.State[i]
		strokeCount += len(t.Strokes)
		if strokeCount >= tr.MaxOutlineLen {
			break
		}

		result = append(State{t}, result...)
	}

	return result
}

func (s *State) toStrokes() []stroke.Stroke {
	var strokes []stroke.Stroke
	for _, t := range *s {
		strokes = append(strokes, t.Strokes...)
	}
	return strokes
}
