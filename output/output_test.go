// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

import (
	"fmt"
	"testing"
)

// MockRobotgo implements RobotgoInterface for testing
type MockRobotgo struct {
	Valid    bool
	KeyDowns []string
	KeyUps   []string
	KeyTaps  []string
	TypeStrs []string
}

func (m *MockRobotgo) IsValid() bool {
	return m.Valid
}

func (m *MockRobotgo) KeyDown(key string) {
	m.KeyDowns = append(m.KeyDowns, key)
}

func (m *MockRobotgo) KeyUp(key string) {
	m.KeyUps = append(m.KeyUps, key)
}

func (m *MockRobotgo) KeyTap(key string) {
	m.KeyTaps = append(m.KeyTaps, key)
}

func (m *MockRobotgo) TypeStr(s string) {
	m.TypeStrs = append(m.TypeStrs, s)
}

func TestNewVirtualOutput(t *testing.T) {
	tests := []struct {
		name    string
		valid   bool
		wantErr bool
	}{
		{"Valid robotgo", true, false},
		{"Invalid robotgo", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRobotgo{Valid: tt.valid}
			o, err := newVirtualOutputWithRobotgo(mock)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVirtualOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && o == nil {
				t.Errorf("NewVirtualOutput() returned nil Output")
			}
		})
	}
}

// newVirtualOutputWithRobotgo allows injecting mock for testing
func newVirtualOutputWithRobotgo(rg RobotgoInterface) (*Output, error) {
	if !rg.IsValid() {
		return nil, fmt.Errorf("robotgo initialization failed")
	}
	return &Output{robotgo: rg}, nil
}

func TestClose(t *testing.T) {
	mock := &MockRobotgo{Valid: true}
	o, _ := newVirtualOutputWithRobotgo(mock)
	if err := o.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantTypeStrs []string
		wantKeyDowns []string
		wantKeyUps   []string
	}{
		{
			name:         "Lowercase string",
			input:        "hello",
			wantTypeStrs: []string{"hello"},
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
		{
			name:         "Uppercase string",
			input:        "HELLO",
			wantTypeStrs: []string{"HELLO"},
			wantKeyDowns: []string{"shift"},
			wantKeyUps:   []string{"shift"},
		},
		{
			name:         "Mixed case string",
			input:        "HeLLo",
			wantTypeStrs: []string{"H", "e", "LL", "o"},
			wantKeyDowns: []string{"shift", "shift"},
			wantKeyUps:   []string{"shift", "shift"},
		},
		{
			name:         "Special characters",
			input:        "hi, world.",
			wantTypeStrs: []string{"hi", ",", " world", "."},
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
		{
			name:         "Unsupported rune",
			input:        "hi😊",
			wantTypeStrs: []string{"hi", "😊"},
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRobotgo{Valid: true}
			o, _ := newVirtualOutputWithRobotgo(mock)
			err := o.TypeString(tt.input)
			if err != nil {
				t.Errorf("TypeString(%q) error = %v, want nil", tt.input, err)
			}
			if !equalSlices(mock.TypeStrs, tt.wantTypeStrs) {
				t.Errorf("TypeString(%q) TypeStrs = %v, want %v", tt.input, mock.TypeStrs, tt.wantTypeStrs)
			}
			if !equalSlices(mock.KeyDowns, tt.wantKeyDowns) {
				t.Errorf("TypeString(%q) KeyDowns = %v, want %v", tt.input, mock.KeyDowns, tt.wantKeyDowns)
			}
			if !equalSlices(mock.KeyUps, tt.wantKeyUps) {
				t.Errorf("TypeString(%q) KeyUps = %v, want %v", tt.input, mock.KeyUps, tt.wantKeyUps)
			}
		})
	}
}

func TestTypeRune(t *testing.T) {
	tests := []struct {
		name         string
		input        rune
		wantKeyTaps  []string
		wantTypeStrs []string
		wantKeyDowns []string
		wantKeyUps   []string
		wantErr      bool
	}{
		{
			name:         "Lowercase letter",
			input:        'a',
			wantKeyTaps:  []string{"a"},
			wantTypeStrs: nil,
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
		{
			name:         "Uppercase letter",
			input:        'A',
			wantKeyTaps:  []string{"a"},
			wantTypeStrs: nil,
			wantKeyDowns: []string{"shift"},
			wantKeyUps:   []string{"shift"},
		},
		{
			name:         "Space",
			input:        ' ',
			wantKeyTaps:  []string{"space"},
			wantTypeStrs: nil,
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
		{
			name:         "Unsupported rune",
			input:        '😊',
			wantKeyTaps:  nil,
			wantTypeStrs: []string{"😊"},
			wantKeyDowns: nil,
			wantKeyUps:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRobotgo{Valid: true}
			o, _ := newVirtualOutputWithRobotgo(mock)
			err := o.TypeRune(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeRune(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !equalSlices(mock.KeyTaps, tt.wantKeyTaps) {
				t.Errorf("TypeRune(%q) KeyTaps = %v, want %v", tt.input, mock.KeyTaps, tt.wantKeyTaps)
			}
			if !equalSlices(mock.TypeStrs, tt.wantTypeStrs) {
				t.Errorf("TypeRune(%q) TypeStrs = %v, want %v", tt.input, mock.TypeStrs, tt.wantTypeStrs)
			}
			if !equalSlices(mock.KeyDowns, tt.wantKeyDowns) {
				t.Errorf("TypeRune(%q) KeyDowns = %v, want %v", tt.input, mock.KeyDowns, tt.wantKeyDowns)
			}
			if !equalSlices(mock.KeyUps, tt.wantKeyUps) {
				t.Errorf("TypeRune(%q) KeyUps = %v, want %v", tt.input, mock.KeyUps, tt.wantKeyUps)
			}
		})
	}
}

func TestRuneToKey(t *testing.T) {
	tests := []struct {
		name      string
		input     rune
		wantKey   string
		wantShift bool
		wantOK    bool
	}{
		{"Lowercase a", 'a', "a", false, true},
		{"Uppercase A", 'A', "a", true, true},
		{"Space", ' ', "space", false, true},
		{"Enter", '\n', "enter", false, true},
		{"Dot", '.', ".", false, true},
		{"Comma", ',', ",", false, true},
		{"Unsupported rune", '😊', "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, shift, ok := runeToKey(tt.input)
			if key != tt.wantKey {
				t.Errorf("runeToKey(%q) key = %v, want %v", tt.input, key, tt.wantKey)
			}
			if shift != tt.wantShift {
				t.Errorf("runeToKey(%q) shift = %v, want %v", tt.input, shift, tt.wantShift)
			}
			if ok != tt.wantOK {
				t.Errorf("runeToKey(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
		})
	}
}

// equalSlices compares two string slices for equality
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
