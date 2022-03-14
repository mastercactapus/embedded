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

	lastRead uint16
}

const (
	// Registers A and B are one after another
	// so we can just send two bytes to the
	// first for all these calls.
	mcp16RegOLATA  = 0x14
	mcp16RegIPOLA  = 0x02
	mcp16RegGPPUA  = 0x0C
	mcp16RegIODIRA = 0x00
	mcp16RegGPIOA  = 0x12
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

func (m *MCP23X17) Flush() error {
	if err := m.write(mcp16RegGPPUA, m.PullupPins); err != nil {
		return err
	}
	if err := m.write(mcp16RegIODIRA, m.InputPins); err != nil {
		return err
	}
	if err := m.write(mcp16RegIPOLA, m.InvertPins); err != nil {
		return err
	}
	if err := m.write(mcp16RegOLATA, m.OutputState); err != nil {
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

func (m *MCP23X17) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.getPinBuf,
		SetInputFunc: m.setIODIRBuf,
		SetFunc:      m.setOLATBuf,
	}
}

func (m *MCP23X17) setIODIRBuf(n int, v bool) error {
	if v {
		m.InputPins |= 1 << uint16(n)
	} else {
		m.InputPins &= ^(1 << uint16(n))
	}
	return nil
}

func (m *MCP23X17) setIODIR(n int, v bool) error {
	if err := m.setIODIRBuf(n, v); err != nil {
		return err
	}
	return m.write(mcp16RegIODIRA, m.InputPins)
}

func (m *MCP23X17) setOLATBuf(n int, v bool) error {
	if v {
		m.OutputState |= 1 << uint16(n)
	} else {
		m.OutputState &= ^(1 << uint16(n))
	}
	return nil
}

func (m *MCP23X17) setOLAT(n int, v bool) error {
	if err := m.setOLATBuf(n, v); err != nil {
		return err
	}
	return m.write(mcp16RegOLATA, m.OutputState)
}

func (m *MCP23X17) getPinBuf(n int) (bool, error) {
	return m.lastRead&(1<<uint(n)) != 0, nil
}

func (m *MCP23X17) getPin(n int) (bool, error) {
	if err := m.Refresh(); err != nil {
		return false, err
	}
	return m.getPinBuf(n)
}

func (m *MCP23X17) Refresh() error {
	var buf [2]byte
	_, err := m.rw.Read(buf[:])
	if err != nil {
		return err
	}
	m.lastRead = uint16(buf[0]) | uint16(buf[1])<<8
	return nil
}
