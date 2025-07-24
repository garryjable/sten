// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package machine

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sten/stroke"
	"time"

	"github.com/tarm/serial"
)

type SerialPort interface {
	Read(p []byte) (int, error)
	Close() error
}

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

// Can Be overridden to customize layouts
var GeminiDefaults = map[string]string{
	"#1":  "#-",
	"#2":  "#-",
	"#3":  "#-",
	"#4":  "#-",
	"#5":  "#-",
	"#6":  "#-",
	"#7":  "#-",
	"#8":  "#-",
	"#9":  "#-",
	"#A":  "#-",
	"#B":  "#-",
	"#C":  "#-",
	"S1-": "S-",
	"S2-": "S-",
	"T-":  "T-",
	"K-":  "K-",
	"P-":  "P-",
	"W-":  "W-",
	"H-":  "H-",
	"R-":  "R-",
	"A-":  "A",
	"O-":  "O",
	"*1":  "*",
	"*2":  "*",
	"*3":  "*",
	"*4":  "*",
	"-E":  "E",
	"-U":  "U",
	"-F":  "-F",
	"-R":  "-R",
	"-P":  "-P",
	"-B":  "-B",
	"-L":  "-L",
	"-G":  "-G",
	"-T":  "-T",
	"-S":  "-S",
	"-D":  "-D",
	"-Z":  "-Z",
	// Optionally map Fn, pwr, res1, etc., as needed.
}

// GeminiPrMachine represents a Gemini PR stenotype machine.
type GeminiPrMachine struct {
	portName   string
	baudRate   int
	port       SerialPort
	strokeChan chan stroke.Stroke
}

// NewGeminiPrMachine creates a new Gemini PR machine instance.
func NewGeminiPrMachine(portName string, baudRate int) *GeminiPrMachine {
	return &GeminiPrMachine{
		portName:   portName,
		baudRate:   baudRate,
		strokeChan: make(chan stroke.Stroke, 64),
	}
}

// StartCapture opens the serial port and starts reading strokes.
func (m *GeminiPrMachine) StartCapture() error {
	if m.port == nil {
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
	}
	go m.readLoop()
	return nil
}

// StopCapture stops reading and closes the serial port.
func (m *GeminiPrMachine) StopCapture() {
	if m.port != nil {
		m.port.Close()
		m.port = nil
	}
}

// reads packets and sends output via the callback
func (m *GeminiPrMachine) readLoop() {
	defer close(m.strokeChan)

	packet := StrokePacket{}
	for {
		_, err := m.port.Read(packet[:])
		if err != nil {
			// Only print unexpected errors
			if err == io.EOF {
				continue
			}
			log.Printf("serial read error: %v", err)
			break
		}
		stroke, err := packet.toStroke()
		if err == nil {
			m.strokeChan <- stroke
		}
	}
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

func (m *GeminiPrMachine) Strokes() chan stroke.Stroke {
	return m.strokeChan
}

func (packet *StrokePacket) toStroke() (stroke.Stroke, error) {
	if !packet.isValid() {
		return 0, errors.New("Invalid Stroke Packet")
	}
	var keys []string
	for row, b := range packet {
		for bit := 1; bit <= 7; bit++ {
			mask := byte(0x80 >> bit)
			if b&mask != 0 {
				if key, ok := GeminiDefaults[keyChart[row][bit-1]]; ok {
					keys = append(keys, key)
				}
			}
		}
	}
	stenoKeys := stroke.JoinKeys(keys)
	return stroke.ParseSteno(stenoKeys), nil
}
