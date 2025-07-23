// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

import (
	"github.com/go-vgo/robotgo"
)

type DevOutputService struct {
	output chan Output
}

func NewDevOutputService(output chan Output) *DevOutputService {
	s := &DevOutputService{
		output: output, // Buffered channel for performance
	}
	return s

}

func (s *DevOutputService) Run() {
	for out := range s.output {
		switch out.Type {
		case Writing:
			robotgo.TypeStr(out.Text)
		case Undoing:
			for range []rune(out.Text) {
				robotgo.KeyTap("backspace")
			}
		}
	}
}
