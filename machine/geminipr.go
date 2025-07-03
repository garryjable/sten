package machine

import "github.com/tarm/serial"

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
