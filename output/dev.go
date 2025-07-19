// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

import (
	"github.com/go-vgo/robotgo"
)

type DevOutputService struct {
	cmds chan OutputCommand
}

func NewDevOutputService() *DevOutputService {
	s := &DevOutputService{
		cmds: make(chan OutputCommand, 32), // Buffered channel for performance
	}
	go s.loop()
	return s

}

func (out *DevOutputService) loop() {
	for cmd := range out.cmds {
		switch cmd.Type {
		case TypeCommand:
			robotgo.TypeStr(cmd.Text + " ")
		case UndoCommand:
			for range []rune(cmd.Text + " ") {
				robotgo.KeyTap("backspace")
			}
		}
	}
}

func (s *DevOutputService) Type(text string) {
	s.cmds <- OutputCommand{Type: TypeCommand, Text: text}
}

func (s *DevOutputService) Undo(text string) {
	s.cmds <- OutputCommand{Type: UndoCommand, Text: text}
}
