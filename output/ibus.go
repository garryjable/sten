package output

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type IBUSOutput struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func NewIBUSOutput() (*IBUSOutput, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	obj := conn.Object("org.freedesktop.IBus", "/org/freedesktop/IBus/Panel")

	return &IBUSOutput{
		conn: conn,
		obj:  obj,
	}, nil
}

func (i *IBUSOutput) SendText(text string) error {
	call := i.obj.Call("org.freedesktop.IBus.Panel.TextCommitted", 0, text)
	if call.Err != nil {
		return fmt.Errorf("failed to send text over IBUS: %w", call.Err)
	}
	return nil
}

func (i *IBUSOutput) Close() error {
	return i.conn.Close()
}
