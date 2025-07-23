// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package output

//// TestDevOutputType manually verifies that typing a string does not panic.
//// Since robotgo actually types to the real system, this test only confirms no crash.
//func TestDevOutputType(t *testing.T) {
//	in := make(chan Output, 2)
//	out := NewDevOutputService(in)
//
//	// Run the output processing in a goroutine
//	go out.Run()
//
//	// Send a test command
//	in <- Output{Type: Writing, Text: "hello world"}
//	close(in)
//
//	// Optionally: if your service writes to the real keyboard, you probably don't want to actually do that during tests.
//	// So just make sure it doesn't panic and you can check logs/output.
//	t.Log("Typed 'hello world' (manually verify if you want)")
//}
