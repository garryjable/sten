// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

type OutputService interface {
	Run()
}

type OutputType int

const (
	Writing OutputType = iota
	Undoing
)

type Output struct {
	Type OutputType
	Text string
}

func (t OutputType) String() string {
	if t == Writing {
		return "Write"
	} else if t == Undoing {
		return "Undo"
	}
	return "Unknown"
}

func NewWrite(text string) Output {
	return Output{Type: Writing, Text: text}
}

func NewUndo(text string) Output {
	return Output{Type: Undoing, Text: text}
}
