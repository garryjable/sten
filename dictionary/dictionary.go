// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.
package dictionary

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Dictionary is the interface Translator depends on.
type Dict interface {
	Lookup(outline string) (string, bool)
}

type Dictionary struct {
	entries map[string]string
}

func LoadDictionaries(folder string) (Dict, int, error) {
	entries := make(map[string]string)
	combined := &Dictionary{
		entries: entries,
	}
	longestOutline := 0

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, longestOutline, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range files {
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
			combined.entries[k] = v

			// Count strokes: number of slashes + 1
			count := strings.Count(k, "/") + 1
			if count > longestOutline {
				longestOutline = count
			}
		}
	}

	fmt.Printf("Loaded %d entries across dictionaries. Max outline length: %d strokes.\n", len(combined.entries), longestOutline)

	return combined, longestOutline, nil
}

func (d *Dictionary) Lookup(stroke string) (string, bool) {
	result, ok := d.entries[stroke]
	return result, ok
}
