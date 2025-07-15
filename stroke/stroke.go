// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"errors"
	"gplover/config"
	"sort"
	"strings"
)

var undoStroke = "*" // you can change this

type Stroke struct {
	Keys []string
}

func NewStroke(keys []string) (*Stroke, error) {
	if len(keys) == 0 {
		return nil, errors.New("empty stroke")
	}
	return &Stroke{Keys: keys}, nil
}

func (s *Stroke) Steno() string {
	var parts []string
	for _, key := range s.Keys {
		if val, ok := config.Layout[key]; ok {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, "")
}

func (s *Stroke) IsCorrection() bool {
	return s.Steno() == undoStroke
}

func SortStrokes(strokes [][]*Stroke) {
	sort.Slice(strokes, func(i, j int) bool {
		// Sort by number of strokes, then number of keys
		if len(strokes[i]) != len(strokes[j]) {
			return len(strokes[i]) < len(strokes[j])
		}
		return totalKeys(strokes[i]) < totalKeys(strokes[j])
	})
}

func totalKeys(strokeSeq []*Stroke) int {
	count := 0
	for _, s := range strokeSeq {
		count += len(s.Keys)
	}
	return count
}

func ToRtfcre(strokes []Stroke) string {
	chords := make([]string, len(strokes))
	for i, stroke := range strokes {
		chords[i] = stroke.Steno() // or however you serialize a stroke
	}
	return strings.Join(chords, "/")
}
