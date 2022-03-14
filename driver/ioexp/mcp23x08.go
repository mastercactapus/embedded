package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type MCP23X08 struct {
	rw io.ReadWriter

	// Invert pins will invert the polarity of the input pins.
	InvertPins  uint8
	PullupPins  uint8
	InputPins   uint8
	OutputState uint8

	lastRead uint8
}

const (
	mcp8RegOLAT  = 0x0a
	mcp8RegIPOL  = 0x01
	mcp8RegGPPU  = 0x06
	mcp8RegIODIR = 0x00
	mcp8RegGPIO  = 0x09
)

// NewMCP23008 is a convenience method that returns a PinReadWriter for a MCP23008-compatible I2C device.
func NewMCP23008(bus i2c.Bus, addr uint16) *MCP23X08 {
	if addr == 0 {
		addr = 0x20
	}
	return NewMCP23X08(i2c.NewDevice(bus, addr))
}

// NewMCP23X08 is a convenience method that returns a PinReadWriter for a MCP23x08-compatible serial device.
func NewMCP23X08(rw io.ReadWriter) *MCP23X08 {
	return &MCP23X08{rw: rw}
}

func (MCP23X08) PinCount() int { return 8 }

func (m *MCP23X08) Flush() error {
	if err := m.write(mcp8RegGPPU, m.PullupPins); err != nil {
		return err
	}
	if err := m.write(mcp8RegIODIR, m.InputPins); err != nil {
		return err
	}
	if err := m.write(mcp8RegIPOL, m.InvertPins); err != nil {
		return err
	}
	if err := m.write(mcp8RegOLAT, m.OutputState); err != nil {
		return err
	}
	return nil
}

func (m *MCP23X08) write(reg, val uint8) error {
	_, err := m.rw.Write([]byte{reg, val})
	return err
}

func (m *MCP23X08) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.getPin,
		SetInputFunc: m.setIODIR,
		SetFunc:      m.setOLAT,
	}
}

func (m *MCP23X08) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.getPinBuf,
		SetInputFunc: m.setIODIRBuf,
		SetFunc:      m.setOLATBuf,
	}
}

func (m *MCP23X08) setIODIRBuf(n int, v bool) error {
	if v {
		m.InputPins |= 1 << uint8(n)
	} else {
		m.InputPins &= ^(1 << uint8(n))
	}
	return nil
}

func (m *MCP23X08) setIODIR(n int, v bool) error {
	m.setIODIRBuf(n, v)
	return m.write(mcp8RegIODIR, m.InputPins)
}

func (m *MCP23X08) setOLATBuf(n int, v bool) error {
	if v {
		m.OutputState |= 1 << uint8(n)
	} else {
		m.OutputState &= ^(1 << uint8(n))
	}
	return nil
}

func (m *MCP23X08) setOLAT(n int, v bool) error {
	m.setOLATBuf(n, v)
	return m.write(mcp8RegOLAT, m.OutputState)
}

func (m *MCP23X08) getPinBuf(n int) (bool, error) {
	return m.lastRead&(1<<uint(n)) != 0, nil
}

func (m *MCP23X08) getPin(n int) (bool, error) {
	if err := m.Refresh(); err != nil {
		return false, err
	}
	return m.getPinBuf(n)
}

func (m *MCP23X08) Refresh() error {
	if br, ok := m.rw.(io.ByteReader); ok {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		m.lastRead = b
		return nil
	}

	var buf [1]byte
	_, err := m.rw.Read(buf[:])
	if err != nil {
		return err
	}
	m.lastRead = buf[0]
	return nil
}
