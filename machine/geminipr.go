package machine

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
)

const (
	BytesPerStroke = 6
)

// In the Gemini PR protocol, each packet consists of exactly six bytes
// and the most significant bit (MSB) of every byte is used exclusively
// to indicate whether that byte is the first byte of the packet
// (MSB=1) or one of the remaining five bytes of the packet (MSB=0). As
// such, there are really only seven bits of steno data in each packet
// byte. This is why the STENO_KEY_CHART below is visually presented as
// six rows of seven elements instead of six rows of eight elements.
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
	portName    string
	baudRate    int
	callback    StrokeCallback
	serialPort  *serial.Port
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

// NewGeminiPrMachine creates a new Gemini PR machine instance.
func NewGeminiPrMachine(portName string, baudRate int, cb StrokeCallback) *GeminiPrMachine {
	return &GeminiPrMachine{
		portName:    portName,
		baudRate:    baudRate,
		callback:    cb,
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}
}

// processPacket decodes a 6-byte Gemini PR packet into pressed keys
func (g *GeminiPrMachine) processPacket(packet []byte) []string {
	if len(packet) != BytesPerStroke {
		return nil
	}

	// Verify packet structure - first byte should have MSB=1, others MSB=0
	if packet[0]&0x80 == 0 {
		return nil // Invalid packet
	}
	for i := 1; i < BytesPerStroke; i++ {
		if packet[i]&0x80 != 0 {
			return nil // Invalid packet
		}
	}

	var pressedKeys []string

	// Process each byte and extract the 7 data bits
	for byteIndex := 0; byteIndex < BytesPerStroke; byteIndex++ {
		dataBits := packet[byteIndex] & 0x7F // Remove MSB, keep 7 data bits

		// Check each of the 7 bits
		for bitIndex := 0; bitIndex < 7; bitIndex++ {
			if dataBits&(1<<uint(6-bitIndex)) != 0 { // Check from bit 6 down to bit 0
				keyIndex := byteIndex*7 + bitIndex
				if keyIndex < len(stenoKeyChart) {
					key := stenoKeyChart[keyIndex]
					// Skip reserved keys
					if key != "res1" && key != "res2" && key != "pwr" {
						pressedKeys = append(pressedKeys, key)
					}
				}
			}
		}
	}

	return pressedKeys
}

// StartCapture begins capturing strokes from the serial port
func (g *GeminiPrMachine) StartCapture() error {
	config := &serial.Config{
		Name: g.portName,
		Baud: g.baudRate,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}

	g.serialPort = port
	go g.readLoop()

	return nil
}

// StopCapture stops capturing strokes and closes the serial port
func (g *GeminiPrMachine) StopCapture() {
	close(g.stopChan)
	<-g.stoppedChan // Wait for readLoop to finish

	if g.serialPort != nil {
		g.serialPort.Close()
		g.serialPort = nil
	}
}

// readLoop continuously reads packets from the serial port
func (g *GeminiPrMachine) readLoop() {
	defer close(g.stoppedChan)

	buffer := make([]byte, BytesPerStroke)

	for {
		select {
		case <-g.stopChan:
			return
		default:
			// Read one byte at a time to find packet start
			n, err := g.serialPort.Read(buffer[:1])
			if err != nil {
				log.Printf("Serial read error: %v", err)
				continue
			}
			if n == 0 {
				continue
			}

			// Check if this is the start of a packet (MSB = 1)
			if buffer[0]&0x80 == 0 {
				continue // Not a packet start, keep looking
			}

			// Read the remaining 5 bytes
			bytesRead := 1
			for bytesRead < BytesPerStroke {
				n, err := g.serialPort.Read(buffer[bytesRead:BytesPerStroke])
				if err != nil {
					log.Printf("Serial read error: %v", err)
					break
				}
				bytesRead += n
			}

			if bytesRead == BytesPerStroke {
				// Process the complete packet
				keys := g.processPacket(buffer)
				if len(keys) > 0 && g.callback != nil {
					g.callback(keys)
				}
			}
		}
	}
}
