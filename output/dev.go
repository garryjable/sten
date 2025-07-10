// output/dev.go
//go:build !prod
// +build !prod

package output

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/go-vgo/robotgo"
)

type DevOutput struct{}

func (d *DevOutput) Type(text string) {
	robotgo.TypeStr(text)
}

func (d *DevOutput) Backspace(n int) {
	for i := 0; i < n; i++ {
		robotgo.KeyTap("backspace")
	}
}

// NewVirtualOutput initializes a robotgo-based output
func NewVirtualOutput() (*DevOutput, error) {
	if !robotgo.IsValid() {
		return nil, fmt.Errorf("robotgo initialization failed")
	}
	return &DevOutput{}, nil
}

// Close is a no-op for robotgo (no resources to close)
func (o *DevOutput) Close() error {
	return nil
}

// TypeString types a string using robotgo
func (o *DevOutput) TypeString(s string) error {
	for _, r := range s {
		if err := o.TypeRune(r); err != nil {
			log.Printf("failed to type rune %q: %v", r, err)
		}
	}
	return nil
}

// TypeRune types a single rune, handling shift for uppercase
func (o *DevOutput) TypeRune(r rune) error {
	key, shifted, ok := runeToKey(r)
	if !ok {
		return fmt.Errorf("unsupported rune: %q", r)
	}
	if shifted {
		robotgo.KeyDown("shift")
	}
	robotgo.KeyTap(key)
	if shifted {
		robotgo.KeyUp("shift")
	}
	return nil
}

func runeToKey(r rune) (string, bool, bool) {
	// Handle Latin letters (a-z, A-Z)
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return strings.ToLower(string(r)), unicode.IsUpper(r), true
	}
	// Handle digits
	switch r {
	case '0':
		return "0", false, true
	case '1':
		return "1", false, true
	case '2':
		return "2", false, true
	case '3':
		return "3", false, true
	case '4':
		return "4", false, true
	case '5':
		return "5", false, true
	case '6':
		return "6", false, true
	case '7':
		return "7", false, true
	case '8':
		return "8", false, true
	case '9':
		return "9", false, true
	}
	// Handle common punctuation and special keys
	switch r {
	case ' ':
		return "space", false, true
	case '\n':
		return "enter", false, true
	case '.':
		return ".", false, true
	case ',':
		return ",", false, true
	case '!':
		return "1", true, true
	case '@':
		return "2", true, true
	case '#':
		return "3", true, true
	case '$':
		return "4", true, true
	case '%':
		return "5", true, true
	case '^':
		return "6", true, true
	case '&':
		return "7", true, true
	case '*':
		return "8", true, true
	case '(':
		return "9", true, true
	case ')':
		return "0", true, true
	case '-':
		return "-", false, true
	case '_':
		return "-", true, true
	case '=':
		return "=", false, true
	case '+':
		return "=", true, true
	case ';':
		return ";", false, true
	case ':':
		return ";", true, true
	case '\'':
		return "'", false, true
	case '"':
		return "'", true, true
	case '/':
		return "/", false, true
	case '?':
		return "/", true, true
	case '[':
		return "[", false, true
	case ']':
		return "]", false, true
	case '{':
		return "[", true, true
	case '}':
		return "]", true, true
	case '\\':
		return "\\", false, true
	case '|':
		return "\\", true, true
	}
	// For all other Unicode runes, return the rune as a string for TypeStr
	return string(r), false, false
}
