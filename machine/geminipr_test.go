// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package machine

import (
	"io"
	"testing"
)

type MockSerialPort struct {
	packets []StrokePacket
	cursor  int
	closed  bool
}

func (m *MockSerialPort) Read(p []byte) (int, error) {
	if m.cursor >= len(m.packets) {
		return 0, io.ErrClosedPipe // simulate port closure
	}
	copy(p, m.packets[m.cursor][:])

	m.cursor++
	return 6, nil
}

func (m *MockSerialPort) Close() error {
	m.closed = true
	return nil
}

func MakeGeminiPacket(keys ...string) StrokePacket {
	var pkt StrokePacket
	pkt[0] = 0x80 // Always set MSB of first byte per protocol

	for _, key := range keys {
		if pos, ok := geminiBits[key]; ok {
			row, bit := pos[0], pos[1]
			// Key bits are bits 6..0 (MSB used for header), so we shift accordingly
			mask := byte(0x80 >> (bit + 1))
			pkt[row] |= mask
		}
	}
	return pkt
}

// Gemini PR protocol key-to-byte/bit position mapping
var geminiBits = map[string][2]int{
	// Row 0 (packet header + number bar)
	"Fn": {0, 0},
	"#1": {0, 1},
	"#2": {0, 2},
	"#3": {0, 3},
	"#4": {0, 4},
	"#5": {0, 5},
	"#6": {0, 6},

	// Row 1 (left hand)
	"S1-": {1, 0},
	"S2-": {1, 1},
	"T-":  {1, 2},
	"K-":  {1, 3},
	"P-":  {1, 4},
	"W-":  {1, 5},
	"H-":  {1, 6},

	// Row 2 (middle/center)
	"R-":   {2, 0},
	"A-":   {2, 1},
	"O-":   {2, 2},
	"*1":   {2, 3},
	"*2":   {2, 4},
	"res1": {2, 5},
	"res2": {2, 6},

	// Row 3 (stars/power/right vowels)
	"pwr": {3, 0},
	"*3":  {3, 1},
	"*4":  {3, 2},
	"-E":  {3, 3},
	"-U":  {3, 4},
	"-F":  {3, 5},
	"-R":  {3, 6},

	// Row 4 (right hand)
	"-P": {4, 0},
	"-B": {4, 1},
	"-L": {4, 2},
	"-G": {4, 3},
	"-T": {4, 4},
	"-S": {4, 5},
	"-D": {4, 6},

	// Row 5 (far right, number bar, -Z)
	"#7": {5, 0},
	"#8": {5, 1},
	"#9": {5, 2},
	"#A": {5, 3},
	"#B": {5, 4},
	"#C": {5, 5},
	"-Z": {5, 6},
}

func TestProcessPacket_InvalidFirstByte(t *testing.T) {
	packet := StrokePacket{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err := packet.toStroke()
	if err == nil {
		t.Errorf("Packet parsed incorrectly")
	}
}

func TestProcessPacket_InvalidOtherByte(t *testing.T) {
	packet := StrokePacket{0x80, 0x80, 0x00, 0x00, 0x00, 0x00}
	_, err := packet.toStroke()
	if err == nil {
		t.Errorf("Packet parsed incorrectly")
	}
}

func TestGeminiMachine(t *testing.T) {
	cases := []struct {
		name     string
		input    []StrokePacket
		expected []string
	}{
		{
			name: "Single Keys",
			input: []StrokePacket{
				MakeGeminiPacket("#1"),
				MakeGeminiPacket("#2"),
				MakeGeminiPacket("#3"),
				MakeGeminiPacket("#4"),
				MakeGeminiPacket("#5"),
				MakeGeminiPacket("#6"),
				MakeGeminiPacket("#7"),
				MakeGeminiPacket("#8"),
				MakeGeminiPacket("#9"),
				MakeGeminiPacket("S1-"),
				MakeGeminiPacket("S2-"),
				MakeGeminiPacket("T-"),
				MakeGeminiPacket("K-"),
				MakeGeminiPacket("P-"),
				MakeGeminiPacket("W-"),
				MakeGeminiPacket("H-"),
				MakeGeminiPacket("R-"),
				MakeGeminiPacket("A-"),
				MakeGeminiPacket("O-"),
				MakeGeminiPacket("*1"),
				MakeGeminiPacket("*2"),
				MakeGeminiPacket("*3"),
				MakeGeminiPacket("*4"),
				MakeGeminiPacket("-E"),
				MakeGeminiPacket("-U"),
				MakeGeminiPacket("-F"),
				MakeGeminiPacket("-R"),
				MakeGeminiPacket("-P"),
				MakeGeminiPacket("-B"),
				MakeGeminiPacket("-L"),
				MakeGeminiPacket("-G"),
				MakeGeminiPacket("-T"),
				MakeGeminiPacket("-S"),
				MakeGeminiPacket("-D"),
				MakeGeminiPacket("-Z"),
			},
			expected: []string{
				"#", "#", "#", "#", "#", "#", "#", "#", "#",
				"S", "S", "T", "K", "P", "W", "H", "R",
				"A", "O", "*", "*", "*", "*", "E", "U",
				"-F", "-R", "-P", "-B", "-L", "-G", "-T", "-S", "-D", "-Z",
			},
		},
		{
			name: "Strokes",
			input: []StrokePacket{
				MakeGeminiPacket("T-", "K-", "A-", "O-", "-P", "-L", "-D"),
				MakeGeminiPacket("S1-", "K-", "P-"),
				MakeGeminiPacket("K-", "O-", "-P", "-L"),
				MakeGeminiPacket("P-", "R-", "A-", "O-", "-E", "-T"),
				MakeGeminiPacket("-T", "-D", "-Z"),
				MakeGeminiPacket("S2-", "O-", "-E"),
			},
			expected: []string{"TKAOPLD", "SKP", "KOPL", "PRAOET", "-TDZ", "SOE"},
		},
		// More cases...
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &MockSerialPort{packets: tc.input}
			machine := NewGeminiPrMachine("mock", 9600)
			machine.port = mock // inject mock port

			done := make(chan struct{})
			go func() { machine.StartCapture(); close(done) }()

			var output []string
			for stroke := range machine.Strokes() {
				output = append(output, stroke.Steno())
			}
			machine.StopCapture()
			<-done

			if len(output) != len(tc.expected) {
				t.Fatalf("expected %d outputs, got %d", len(tc.expected), len(output))
			}
			for i, want := range tc.expected {
				if output[i] != want {
					t.Errorf("at %d: want %q, got %q", i, want, output[i])
				}
			}
			if !mock.closed {
				t.Errorf("serial port was not closed")
			}
		})
	}
}
