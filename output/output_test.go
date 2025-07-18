// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

import (
	"testing"
)

// TestDevOutputType manually verifies that typing a string does not panic.
// Since robotgo actually types to the real system, this test only confirms no crash.
func TestDevOutputType(t *testing.T) {
	out := NewDevOutputService()
	// You can manually observe the output if needed
	out.Type("hello world")

	t.Log("Typed 'hello world ' to the system keyboard (manually verify if needed)")
}
