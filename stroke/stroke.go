// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package stroke

import (
	"strings"
)

// Steno key order: # S T K P W H R A O * E U F R P B L G T S D Z
const (
	StenoDash  uint32 = 0         // - noop
	StenoHash  uint32 = 1 << iota // #
	StenoS1                       // S-
	StenoT                        // T-
	StenoK                        // K-
	StenoP                        // P-
	StenoW                        // W-
	StenoH                        // H-
	StenoR                        // R-
	StenoA                        // A-
	StenoO                        // O-
	StenoSplat                    // *
	StenoE                        // -E
	StenoU                        // -U
	StenoF                        // -F
	StenoR2                       // -R
	StenoP2                       // -P
	StenoB                        // -B
	StenoL                        // -L
	StenoG                        // -G
	StenoT2                       // -T
	StenoS2                       // -S
	StenoD                        // -D
	StenoZ                        // -Z
)

type Stroke uint32

type StenoKey struct {
	str string
	bit uint32
}

var leftKeys = []StenoKey{
	{"#", StenoHash},
	{"S", StenoS1},
	{"T", StenoT},
	{"K", StenoK},
	{"P", StenoP},
	{"W", StenoW},
	{"H", StenoH},
	{"R", StenoR},
}

var vowelKeys = []StenoKey{
	{"A", StenoA},
	{"O", StenoO},
	{"*", StenoSplat},
	{"E", StenoE},
	{"U", StenoU},
	{"-", StenoDash}, // only used to resolve ambiguity when vowels are absent
}

var rightKeys = []StenoKey{
	{"F", StenoF},
	{"R", StenoR2},
	{"P", StenoP2},
	{"B", StenoB},
	{"L", StenoL},
	{"G", StenoG},
	{"T", StenoT2},
	{"S", StenoS2},
	{"D", StenoD},
	{"Z", StenoZ},
}

var leftBits = make(map[string]uint32)
var vowelBits = make(map[string]uint32)
var rightBits = make(map[string]uint32)
var stenoBits []map[string]uint32

func init() {
	for _, k := range leftKeys {
		leftBits[k.str] = k.bit
	}
	for _, k := range vowelKeys {
		vowelBits[k.str] = k.bit
	}
	for _, k := range rightKeys {
		rightBits[k.str] = k.bit
	}
	// All three maps of bits in an array, zones are 0=left, 1=vowel, 2=right
	stenoBits = []map[string]uint32{leftBits, vowelBits, rightBits}

}

func (s Stroke) Steno() string {
	left := gatherKeys(s, leftKeys, 0)
	vowels := gatherKeys(s, vowelKeys, len(leftKeys))
	right := gatherKeys(s, rightKeys, len(leftKeys)+len(vowelKeys))

	if vowels == "" && right != "" {
		vowels = "-"
	}
	return left + vowels + right
}

func gatherKeys(s Stroke, keys []StenoKey, offset int) string {
	var b strings.Builder
	for _, key := range keys {
		if s&Stroke(key.bit) != 0 {
			b.WriteString(key.str)
		}

	}
	return b.String()
}

func ParseSteno(steno string) Stroke {
	runes := []rune(steno)
	return gatherBits(runes, 0, 0)
}

const zoneLeft = 0
const zoneVowel = 1
const zoneRight = 2

func gatherBits(runes []rune, i int, zone int) Stroke {
	if i >= len(runes) || zone > 2 {
		return 0
	}
	key := string(runes[i])
	bits := stenoBits[zone]

	// Decide when to switch to next zone:
	switch zone {
	case zoneLeft:
		if _, ok := vowelBits[key]; ok {
			return gatherBits(runes, i, 1)
		}
	case zoneVowel:
		if _, ok := vowelBits[key]; !ok {
			return gatherBits(runes, i, 2)
		}
	}
	// For zone 2 (right), always stay in zone 2 until end.

	// Set the bit if present, then recurse.
	var bit Stroke
	if b, ok := bits[key]; ok {
		bit = Stroke(b)
	}
	return bit | gatherBits(runes, i+1, zone)
}

func (s Stroke) String() string {
	return s.Steno()
}
