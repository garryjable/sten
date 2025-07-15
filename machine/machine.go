// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package machine

type Machine interface {
	StartCapture() error
	StopCapture()
	SetCallback(cb StrokeCallback)
}
