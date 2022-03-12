package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type MCP23X17 struct {
	rw io.ReadWriter

	// Invert pins will invert the polarity of the input pins.
	InvertPins  uint16
	PullupPins  uint16
	InputPins   uint16
	OutputState uint16
}

const (
	// Registers A and B are one after another
	// so we can just send two bytes to the
	// first for all these calls.
	mcpRegOLATA  = 0x14
	mcpRegIPOLA  = 0x02
	mcpRegGPPUA  = 0x0C
	mcpRegIODIRA = 0x00
	mcpRegGPIOA  = 0x12
)

// NewMCP23017 is a convenience method that returns a PinReadWriter for a MCP23017-compatible I2C device.
func NewMCP23017(bus i2c.Bus, addr uint16) *MCP23X17 {
	if addr == 0 {
		addr = 0x20
	}
	return NewMCP23X17(i2c.NewDevice(bus, addr))
}

// NewMCP23X17 is a convenience method that returns a PinReadWriter for a MCP23x17-compatible serial device.
func NewMCP23X17(rw io.ReadWriter) *MCP23X17 {
	return &MCP23X17{rw: rw}
}

func (MCP23X17) PinCount() int { return 16 }

func (m *MCP23X17) Configure() error {
	if err := m.write(mcpRegGPPUA, m.PullupPins); err != nil {
		return err
	}
	if err := m.write(mcpRegIODIRA, m.InputPins); err != nil {
		return err
	}
	if err := m.write(mcpRegIPOLA, m.InputPins); err != nil {
		return err
	}
	if err := m.write(mcpRegOLATA, m.OutputState); err != nil {
		return err
	}
	return nil
}

func (m *MCP23X17) write(reg uint8, v uint16) error {
	_, err := m.rw.Write([]byte{reg, byte(v & 0xff), byte(v >> 8)})
	return err
}

func (m *MCP23X17) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.getPin,
		SetInputFunc: m.setIODIR,
		SetFunc:      m.setOLAT,
	}
}

func (m *MCP23X17) setIODIR(n int, v bool) error {
	if v {
		m.InputPins |= 1 << uint16(n)
	} else {
		m.InputPins &= ^(1 << uint16(n))
	}
	return m.write(mcp8RegIODIR, m.OutputState)
}

func (m *MCP23X17) setOLAT(n int, v bool) error {
	if v {
		m.OutputState |= 1 << uint16(n)
	} else {
		m.OutputState &= ^(1 << uint16(n))
	}
	return m.write(mcp8RegOLAT, m.OutputState)
}

func (m *MCP23X17) getPin(n int) (bool, error) {
	b, err := m.read()
	if err != nil {
		return false, err
	}

	return b&(1<<uint(n)) != 0, nil
}

func (m *MCP23X17) read() (uint16, error) {
	var buf [2]byte
	_, err := m.rw.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | (uint16(buf[1]) << 8), nil
}
