package machine

import (
	"errors"
	"testing"
	"time"

	"github.com/tarm/serial"
)

func TestProcessPacket_ValidPacket(t *testing.T) {
	called := false
	var received []string

	m := NewGeminiPrMachine("", 0, func(keys []string) {
		called = true
		received = keys
	})

	// packet := []byte{0x80, 0x10, 0x00, 0x00, 0x00, 0x00} // T- key only
	packet := []byte{0x80, 0x18, 0x20, 0x00, 0x01, 0x00}

	err := m.processPacket(packet)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("callback was not called")
	}
	if len(received) == 0 || received[0] != "T-" {
		t.Errorf("expected T-, got %v", received)
	}
}

func TestProcessPacket_InvalidFirstByte(t *testing.T) {
	m := NewGeminiPrMachine("", 0, nil)
	packet := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := m.processPacket(packet)
	if err == nil || err.Error() != "first byte MSB not set" {
		t.Errorf("expected first byte MSB error, got %v", err)
	}
}

func TestProcessPacket_InvalidOtherByte(t *testing.T) {
	m := NewGeminiPrMachine("", 0, nil)
	packet := []byte{0x80, 0x80, 0x00, 0x00, 0x00, 0x00}
	err := m.processPacket(packet)
	if err == nil || !errors.Is(err, err) {
		t.Errorf("expected byte 1 MSB set error, got %v", err)
	}
}

func TestStartStopCapture(t *testing.T) {
	m := NewGeminiPrMachine("/dev/null", 9600, nil)
	err := m.StartCapture()
	if err == nil {
		t.Error("expected failure on opening /dev/null as serial port")
	}
}

func TestNewGeminiPrMachineDefaults(t *testing.T) {
	m := NewGeminiPrMachine("test", 1234, nil)
	if m == nil {
		t.Fatal("expected machine instance, got nil")
	}
	if m.baudRate != 1234 || m.portName != "test" {
		t.Errorf("unexpected config: %+v", m)
	}
	if m.callback != nil {
		t.Errorf("expected nil callback")
	}
}

func TestReadLoopStops(t *testing.T) {
	m := NewGeminiPrMachine("", 0, nil)
	fakePort := &serial.Port{}
	m.serialPort = fakePort
	go m.StopCapture()
	m.readLoop() // should return quickly
}

type MockSerialPort struct {
	readDelay time.Duration
	error     error
}

func (m *MockSerialPort) Read(p []byte) (int, error) {
	time.Sleep(m.readDelay)
	return 0, m.error
}

func (m *MockSerialPort) Close() error {
	return nil
}
