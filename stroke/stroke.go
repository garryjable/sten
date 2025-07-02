package stroke

import (
	"errors"
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
	parts := make([]string, len(s.Keys))
	for i, key := range s.Keys {
		parts[i] = strings.ReplaceAll(key, "-", "")
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
