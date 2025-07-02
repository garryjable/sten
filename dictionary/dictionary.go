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
