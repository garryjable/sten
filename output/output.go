// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

type Output interface {
	Type(text string)
	Undo(text string)
}

type CommandType int

const (
	TypeCommand CommandType = iota
	UndoCommand
)

type OutputCommand struct {
	Type CommandType
	Text string
}
