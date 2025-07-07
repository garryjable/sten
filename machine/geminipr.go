// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package machine

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

const (
	BytesPerStroke = 6
)

var (
	// Gemini PR key chart - 6 rows * 7 bits
	stenoKeyChart = []string{
		"Fn", "#1", "#2", "#3", "#4", "#5", "#6",
		"S1-", "S2-", "T-", "K-", "P-", "W-", "H-",
		"R-", "A-", "O-", "*1", "*2", "res1", "res2",
		"pwr", "*3", "*4", "-E", "-U", "-F", "-R",
		"-P", "-B", "-L", "-G", "-T", "-S", "-D",
		"#7", "#8", "#9", "#A", "#B", "#C", "-Z",
	}
)

// StrokeCallback is the function type called when a stroke is decoded.
type StrokeCallback func([]string)

// GeminiPrMachine represents a Gemini PR stenotype machine.
type GeminiPrMachine struct {
	portName string
	baudRate int
	callback StrokeCallback
	port     *serial.Port
	stopping chan struct{} //
	stopped  chan struct{}
}

// NewGeminiPrMachine creates a new Gemini PR machine instance.
func NewGeminiPrMachine(portName string, baudRate int, cb StrokeCallback) *GeminiPrMachine {
	return &GeminiPrMachine{
		portName: portName,
		baudRate: baudRate,
		callback: cb,
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

// StartCapture opens the serial port and starts reading strokes.
func (m *GeminiPrMachine) StartCapture() error {
	c := &serial.Config{
		Name:        m.portName,
		Baud:        m.baudRate,
		ReadTimeout: time.Second * 2,
	}

	port, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}
	m.port = port

	go m.readLoop()
	return nil
}

// StopCapture stops reading and closes the serial port.
func (m *GeminiPrMachine) StopCapture() {
	close(m.stopping)
	<-m.stopped
	if m.port != nil {
		m.port.Close()
		m.port = nil
	}
}

func (m *GeminiPrMachine) readLoop() {
	defer close(m.stopped)

	packet := [6]byte{}

	for {
		select {
		case <-m.stopping:
			return
		default:
			n, err := m.port.Read(packet[:])
			if err != nil {
				// Only print unexpected errors
				if !errors.Is(err, os.ErrDeadlineExceeded) && err != io.EOF {
					log.Printf("serial read error: %v", err)
				}
				continue
			}
			if n != BytesPerStroke {
				continue
			}
			if err := m.processPacket(packet); err != nil {
				log.Printf("invalid packet: %v", err)
			}
		}
	}
}

// processPacket validates and decodes a Gemini PR packet.
func (m *GeminiPrMachine) processPacket(packet [BytesPerStroke]byte) error {
	// Validate packet: first byte MSB must be 1, others must be 0
	if packet[0]&0x80 == 0 {
		return errors.New("first byte MSB not set")
	}
	for i := 1; i < len(packet); i++ {
		if packet[i]&0x80 != 0 {
			return fmt.Errorf("byte %d MSB set", i)
		}
	}

	stenoKeys := []string{}

	for i, b := range packet {
		for bit := 1; bit <= 7; bit++ {
			mask := byte(0x80 >> bit)
			if b&mask != 0 {
				index := i*7 + (bit - 1)
				if index < len(stenoKeyChart) {
					stenoKeys = append(stenoKeys, stenoKeyChart[index])
				}
			}
		}
	}

	// Notify callback with decoded keys
	if m.callback != nil && len(stenoKeys) > 0 {
		m.callback(stenoKeys)
	}
	return nil
}
