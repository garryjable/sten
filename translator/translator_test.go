// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package translator

import (
	"fmt"
	"sten/output"
	"sten/stroke"
	"testing"
)

// MockDictionary is a basic stub for dictionary.Dictionary
type MockDictionary struct {
	entries map[string]string
}

func (m *MockDictionary) Lookup(outline fmt.Stringer) (string, bool) {
	val, ok := m.entries[outline.String()]
	return val, ok
}

func TestTranslator(t *testing.T) {
	type testCase struct {
		name       string
		dict       map[string]string
		outlineCap int
		strokes    []string
		expected   []output.Output
	}

	cases := []testCase{
		{
			name: "Single",
			dict: map[string]string{
				"STPH":      "hello",
				"STPH/STPH": "Beyond ouline cap",
				"*":         "=undo",
			},
			strokes: []string{
				"STPH",
				"STPH",
				"*",
				"TPHOEPB", // unknown translation
			},
			outlineCap: 1,
			expected: []output.Output{
				{"hello ", ""},
				{"hello ", ""},
				{"", "hello "},
				{"TPHOEPB ", ""},
			},
		},
		{
			name: "Multi",
			dict: map[string]string{
				"U":                 "you",
				"R":                 "are",
				"EUPB":              "in",
				"TE":                "the",
				"EUPB/TE/HREB/TWAL": "intellectual",
				"HREB/TWAL":         "translations should be greedy",
				"*":                 "=undo",
			},
			strokes: []string{
				"U",
				"R",
				"EUPB",
				"TE",
				"HREB",
				"TWAL",
				"*",
				"TWAL",
			},
			outlineCap: 4,
			expected: []output.Output{
				{"you ", ""},
				{"are ", ""},
				{"in ", ""},
				{"the ", ""},
				{"HREB ", ""},
				{"intellectual ", "in the HREB "},
				{"in the HREB ", "intellectual "},
				{"intellectual ", "in the HREB "},
			},
		},
		{
			name: "Recursive Outlines",
			dict: map[string]string{
				"STKPHEPL":           "dismember",
				"STKPHEPL/PWER":      "dismember",
				"STKPHEPL/PWER/-PLT": "dismemberment",
				"*":                  "=undo",
			},
			strokes: []string{
				"STKPHEPL",
				"PWER",
				"-PLT",
				"*",
				"*",
				"*",
			},
			outlineCap: 3,
			expected: []output.Output{
				{"dismember ", ""},
				{"dismember ", "dismember "},
				{"dismemberment ", "dismember "},
				{"dismember ", "dismemberment "}, // Should not rewrite dismember twice
				{"dismember ", "dismember "},
				{"", "dismember "},
			},
		},
		{
			name: "Prefer Old Multi",
			dict: map[string]string{
				"U":                   "you",
				"R":                   "are",
				"EUPB":                "in",
				"TE":                  "the",
				"U/R/EUPB/TE":         "you're into",
				"EUPB/TE/HREB/TWAL":   "translations can't replace multi with older starting stroke",
				"EUPB/TE/HREB/TWAL/E": "translations can't replace multi with older starting stroke",
			},
			strokes: []string{
				"U",
				"R",
				"EUPB",
				"TE",
				"HREB",
				"TWAL",
				"E",
			},
			outlineCap: 5,
			expected: []output.Output{
				{"you ", ""},
				{"are ", ""},
				{"in ", ""},
				{"you're into ", "you are in "},
				{"HREB ", ""},
				{"TWAL ", ""},
				{"E ", ""},
			},
		},
		{
			name: "Suffixes",
			dict: map[string]string{
				"KOPL":         "come",
				"KOPL/PHRAOET": "complete",
				"#P*EU":        "{^}.py",
				"*":            "=undo",
			},
			strokes: []string{
				"KOPL",
				"PHRAOET",
				"#P*EU",
				"*",
				"*",
				"*",
			},
			outlineCap: 2,
			expected: []output.Output{
				{"come ", ""},
				{"complete ", "come "},
				{".py ", " "},
				{" ", ".py "},
				{"come ", "complete "},
				{"", "come "},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := make(chan stroke.Stroke, len(tc.strokes))
			tr := NewTranslator(&MockDictionary{tc.dict}, tc.outlineCap, in)
			go tr.Run()
			for _, steno := range tc.strokes {
				in <- stroke.ParseSteno(steno)
			}
			close(in)
			i := 0
			fmt.Println("testing output")
			for out := range tr.Out() {
				fmt.Println(out)
				if i >= len(tc.expected) {
					t.Fatalf("got more outputs than expected: %+v", out)
				}
				if out != tc.expected[i] {
					t.Errorf("at %d: expected %+v, got %+v", i, tc.expected[i], out)
				}
				i++
			}
			if i != len(tc.expected) {
				t.Fatalf("expected %d outputs, got %d", len(tc.expected), i)
			}
		})
	}
}
