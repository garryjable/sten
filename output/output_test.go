package output

import (
	"errors"
	"testing"

	uinput "gopkg.in/bendahl/uinput.v1"
)

// mockKeyboard is a mock that satisfies the uinput.Keyboard interface
type mockKeyboard struct {
	pressed []int
	fail    bool
}

func (m *mockKeyboard) KeyPress(code int) error {
	if m.fail {
		return errors.New("mock fail")
	}
	m.pressed = append(m.pressed, code)
	return nil
}

func (m *mockKeyboard) Close() error { return nil }

func TestTypeRune_Success(t *testing.T) {
	mock := &mockKeyboard{}
	out := &Output{keyboard: mock}

	err := out.TypeRune('a')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mock.pressed) != 1 || mock.pressed[0] != runeToKey('a') {
		t.Errorf("expected KeyA pressed, got %v", mock.pressed)
	}
}

func TestTypeRune_Unsupported(t *testing.T) {
	mock := &mockKeyboard{}
	out := &Output{keyboard: mock}

	err := out.TypeRune('@') // unsupported
	if err != nil {
		t.Errorf("expected no error for unsupported key, got %v", err)
	}
	if len(mock.pressed) != 0 {
		t.Errorf("expected no keypress, got %v", mock.pressed)
	}
}

func TestTypeRune_Failure(t *testing.T) {
	mock := &mockKeyboard{fail: true}
	out := &Output{keyboard: mock}

	err := out.TypeRune('a')
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTypeString_Mixed(t *testing.T) {
	mock := &mockKeyboard{}
	out := &Output{keyboard: mock}

	err := out.TypeString("ab ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{uinput.KeyA, uinput.KeyB, uinput.KeySpace}
	if len(mock.pressed) != len(expected) {
		t.Errorf("expected %d key presses, got %d", len(expected), len(mock.pressed))
	}
	for i, k := range expected {
		if mock.pressed[i] != k {
			t.Errorf("expected key %d to be %v, got %v", i, k, mock.pressed[i])
		}
	}
}
