// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package machine

import "sten/stroke"

type Machine interface {
	StartCapture() error
	StopCapture()
	Strokes() <-chan stroke.Stroke
}
