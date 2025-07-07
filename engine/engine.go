// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package engine

import (
	"gplover/dictionary"
	"gplover/stroke"
)

type Engine struct {
	Dict dictionary.Dictionary
}

func NewEngine(dict dictionary.Dictionary) *Engine {
	return &Engine{Dict: dict}
}

func (e *Engine) TranslateSteno(keys []string) string {
	stroke, err := stroke.NewStroke(keys)
	if err != nil {
		return "[error]"
	}
	word, ok := e.Dict.Lookup(stroke.Steno())
	if !ok {
		return "[" + stroke.Steno() + "]"
	}
	return word
}
