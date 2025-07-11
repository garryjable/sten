// output/dev.go
//go:build !prod
// +build !prod

package output

import (
	"fmt"

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

func NewVirtualOutput() (*DevOutput, error) {
	if !robotgo.IsValid() {
		return nil, fmt.Errorf("robotgo initialization failed")
	}
	return &DevOutput{}, nil
}

func (d *DevOutput) Close() error {
	return nil
}

// TypeString types a full Unicode string
func (d *DevOutput) TypeString(s string) error {
	d.Type(s)
	return nil
}
