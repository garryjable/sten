// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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

func (s *Stroke) String() string {
	prefix := ""
	if s.IsCorrection() {
		prefix = "*"
	}
	return prefix + "Stroke(" + s.Steno() + " : [" + strings.Join(s.Keys, ", ") + "])"
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
