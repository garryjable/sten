// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"strings"
)

type Outline []Stroke

func (o Outline) Prepend(s Stroke) Outline {
	return append([]Stroke{s}, o...)
}

func (o Outline) Steno() string {
	if len(o) == 0 {
		return ""
	}
	parts := make([]string, len(o))
	for i, stroke := range o {
		parts[i] = stroke.Steno()
	}
	return strings.Join(parts, "/")
}

func (o Outline) String() string {
	return o.Steno()
}

func (s Stroke) Outline() Outline {
	return Outline{s}
}
