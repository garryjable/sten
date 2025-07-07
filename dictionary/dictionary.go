// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package dictionary

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Dictionary map[string]string

func LoadDictionaries(folder string) (map[string]string, error) {
	combined := make(map[string]string)

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(folder, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			log.Printf("failed to open %s: %v", entry.Name(), err)
			continue
		}
		defer f.Close()

		var dict map[string]string
		if err := json.NewDecoder(f).Decode(&dict); err != nil {
			log.Printf("failed to decode %s: %v", entry.Name(), err)
			continue
		}

		for k, v := range dict {
			combined[k] = v
		}
	}

	return combined, nil
}

func (d Dictionary) Lookup(stroke string) (string, bool) {
	result, ok := d[stroke]
	return result, ok
}
