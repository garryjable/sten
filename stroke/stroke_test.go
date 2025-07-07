// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package stroke

import (
	"testing"

	"gplover/config"
)

func TestStrokeTranslation(t *testing.T) {
	config.Layout = map[string]string{
		"S1-": "S",
		"T-":  "T",
		"K-":  "K",
		"-E":  "E",
	}

	s := &Stroke{
		Keys: []string{"S1-", "T-", "K-", "-E"},
	}
	got := s.Steno()
	want := "STKE"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
