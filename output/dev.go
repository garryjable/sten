// output/dev.go
//go:build !prod
// +build !prod

package output

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type DevOutput struct{}

func NewVirtualOutput() (*DevOutput, error) {
	if !robotgo.IsValid() {
		return nil, fmt.Errorf("robotgo initialization failed")
	}
	return &DevOutput{}, nil
}

// TypeString types a full Unicode string
func (d *DevOutput) Type(s string) {
	robotgo.TypeStr(s + " ")
}
