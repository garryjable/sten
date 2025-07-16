// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package machine

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sten/stroke"
	"time"

	"github.com/tarm/serial"
)

// Standard stenotype interface for a Gemini PR machine.
//
//     KEYS_LAYOUT =
//         #1 #2  #3 #4 #5 #6 #7 #8 #9 #A #B #C
//         Fn S1- T- P- H- *1 *3 -F -P -L -T -D
//            S2- K- W- R- *2 *4 -R -B -G -S -Z
//                   A- O-       -E -U
//         pwr
//         res1
//         res2

// In the Gemini PR protocol, each packet consists of exactly six bytes
// and the most significant bit (MSB) of every byte is used exclusively
// to indicate whether that byte is the first byte of the packet
// (MSB=1) or one of the remaining five bytes of the packet (MSB=0). As
// such, there are really only seven bits of steno data in each packet
// byte. This is why the keyChart below is represented as
// six rows of seven elements instead of six rows of eight elements.

type StrokePacket [6]byte
type GeminiBoard [6][7]string

var keyChart = GeminiBoard{
	{"Fn", "#1", "#2", "#3", "#4", "#5", "#6"},
	{"S1-", "S2-", "T-", "K-", "P-", "W-", "H-"},
	{"R-", "A-", "O-", "*1", "*2", "res1", "res2"},
	{"pwr", "*3", "*4", "-E", "-U", "-F", "-R"},
	{"-P", "-B", "-L", "-G", "-T", "-S", "-D"},
	{"#7", "#8", "#9", "#A", "#B", "#C", "-Z"},
}

// StrokeCallback is the function type called when a stroke is decoded.
type StrokeCallback func(*stroke.Stroke)

// GeminiPrMachine represents a Gemini PR stenotype machine.
type GeminiPrMachine struct {
	portName string
	baudRate int
	callback StrokeCallback
	port     *serial.Port
	stopping chan struct{} //
	stopped  chan struct{}
}

func (m *GeminiPrMachine) SetCallback(cb StrokeCallback) {
	m.callback = cb
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

// reads packets and sends output via the callback
func (m *GeminiPrMachine) readLoop() {
	defer close(m.stopped)

	packet := StrokePacket{}

	for {
		select {
		case <-m.stopping:
			return
		default:
			_, err := m.port.Read(packet[:])
			if err != nil {
				// Only print unexpected errors
				if errors.Is(err, os.ErrDeadlineExceeded) || err == io.EOF {
					continue
				}
				log.Printf("serial read error: %v", err)
			}
			stroke, err := packet.toStroke()
			if err == nil {
				m.callback(stroke)
			}
		}
	}
}

// toStroke decodes a Gemini PR packet into a chord of key presses
func (packet StrokePacket) toStroke() (*stroke.Stroke, error) {
	if !packet.isValid() {
		return &stroke.Stroke{}, errors.New("Invalid Stroke Packet")
	}
	stroke := make(stroke.Stroke, 0, 42) // max keys
	for row, b := range packet {
		for bit := 1; bit <= 7; bit++ {
			mask := byte(0x80 >> bit)
			if b&mask != 0 {
				stroke = append(stroke, keyChart[row][bit-1])
			}
		}
	}
	return &stroke, nil
}

// Validate packet: first byte MSB must be 1, others must be 0
func (p StrokePacket) isValid() bool {
	if p[0]&0x80 == 0 {
		return false
	}
	for _, b := range p[1:] {
		if b&0x80 != 0 {
			return false
		}
	}
	return true
}
