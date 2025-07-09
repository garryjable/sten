package output

import (
	"fmt"
	"log"
	"time"
	"unicode"

	uinput "gopkg.in/bendahl/uinput.v1"
)

type Output struct {
	keyboard uinput.Keyboard
}

func NewVirtualOutput() (*Output, error) {
	kb, err := uinput.CreateKeyboard("/dev/uinput", []byte("dotterel-virtual"))
	if err != nil {
		return nil, err
	}
	return &Output{keyboard: kb}, nil
}

func (v *Output) Close() error {
	return v.keyboard.Close()
}

func (v *Output) TypeString(s string) error {
	for _, r := range s {
		if err := v.TypeRune(r); err != nil {
			log.Printf("failed to type rune %q: %v", r, err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

func (v *Output) TypeRune(r rune) error {
	key, shifted, ok := runeToKey(r)
	if !ok {
		return fmt.Errorf("unsupported rune: %q", r)
	}
	if shifted {
		// _ = v.keyboard.KeyDown(uinput.KeyLeftShift)
		println("shifted")
	}
	err := v.keyboard.KeyPress(key)
	if shifted {
		println("unshifted")
		// _ = v.keyboard.KeyUp(uinput.KeyLeftShift)
	}
	return err
}

// --- keycode mapping ---

func runeToKey(r rune) (int, bool, bool) {
	switch r {
	case 'a', 'A':
		return uinput.KeyA, unicode.IsUpper(r), true
	case 'b', 'B':
		return uinput.KeyB, unicode.IsUpper(r), true
	case 'c', 'C':
		return uinput.KeyC, unicode.IsUpper(r), true
	case 'd', 'D':
		return uinput.KeyD, unicode.IsUpper(r), true
	case 'e', 'E':
		return uinput.KeyE, unicode.IsUpper(r), true
	case 'f', 'F':
		return uinput.KeyF, unicode.IsUpper(r), true
	case 'g', 'G':
		return uinput.KeyG, unicode.IsUpper(r), true
	case 'h', 'H':
		return uinput.KeyH, unicode.IsUpper(r), true
	case 'i', 'I':
		return uinput.KeyI, unicode.IsUpper(r), true
	case 'j', 'J':
		return uinput.KeyJ, unicode.IsUpper(r), true
	case 'k', 'K':
		return uinput.KeyK, unicode.IsUpper(r), true
	case 'l', 'L':
		return uinput.KeyL, unicode.IsUpper(r), true
	case 'm', 'M':
		return uinput.KeyM, unicode.IsUpper(r), true
	case 'n', 'N':
		return uinput.KeyN, unicode.IsUpper(r), true
	case 'o', 'O':
		return uinput.KeyO, unicode.IsUpper(r), true
	case 'p', 'P':
		return uinput.KeyP, unicode.IsUpper(r), true
	case 'q', 'Q':
		return uinput.KeyQ, unicode.IsUpper(r), true
	case 'r', 'R':
		return uinput.KeyR, unicode.IsUpper(r), true
	case 's', 'S':
		return uinput.KeyS, unicode.IsUpper(r), true
	case 't', 'T':
		return uinput.KeyT, unicode.IsUpper(r), true
	case 'u', 'U':
		return uinput.KeyU, unicode.IsUpper(r), true
	case 'v', 'V':
		return uinput.KeyV, unicode.IsUpper(r), true
	case 'w', 'W':
		return uinput.KeyW, unicode.IsUpper(r), true
	case 'x', 'X':
		return uinput.KeyX, unicode.IsUpper(r), true
	case 'y', 'Y':
		return uinput.KeyY, unicode.IsUpper(r), true
	case 'z', 'Z':
		return uinput.KeyZ, unicode.IsUpper(r), true
	case ' ':
		return uinput.KeySpace, false, true
	case '\n':
		return uinput.KeyEnter, false, true
	case '.':
		return uinput.KeyDot, false, true
	case ',':
		return uinput.KeyComma, false, true
	default:
		return 0, false, false
	}
}
