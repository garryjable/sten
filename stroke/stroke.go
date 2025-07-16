// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"sten/config"
	"strings"
)

type Stroke []string

func (s *Stroke) Steno() string {
	var parts []string
	for _, key := range *s {
		if val, ok := config.Layout[key]; ok {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, "")
}

func ToRtfcre(strokes []Stroke) string {
	chords := make([]string, len(strokes))
	for i, stroke := range strokes {
		chords[i] = stroke.Steno() // or however you serialize a stroke
	}
	return strings.Join(chords, "/")
}
