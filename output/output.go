// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

type OutputService interface {
	Run()
}

type Output struct {
	Write string
	Undo  string
}

func NewOutput(write, undo string) Output {
	return Output{
		Write: write,
		Undo:  undo,
	}
}
