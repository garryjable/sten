// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package engine

import (
	"testing"

	"gplover/config"
)

func TestEngineTranslate(t *testing.T) {
	// Manually set the layout used for translation
	config.Layout = map[string]string{
		"S2-": "S",
		"T-":  "T",
		"K-":  "K",
		"E":   "E",
		"P-":  "P",
		"W-":  "W",
	}

	dict := map[string]string{
		"STKE": "stack",
		"TKPW": "go",
	}

	e := NewEngine(dict)

	if word := e.TranslateSteno([]string{"S2-", "T-", "K-", "E"}); word != "stack" {
		t.Errorf("expected stack, got %q", word)
	}

	if word := e.TranslateSteno([]string{"T-", "K-", "P-", "W-"}); word != "go" {
		t.Errorf("expected go, got %q", word)
	}

	if word := e.TranslateSteno([]string{"X", "Y", "Z"}); word != "[]" {
		t.Errorf("expected empty string, got %q", word)
	}
}
