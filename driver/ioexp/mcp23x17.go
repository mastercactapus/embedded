package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type MCP23X17 struct {
	rw io.ReadWriter

	InvertPins  *Register16
	PullupPins  *Register16
	InputPins   *Register16
	OutputState *Register16
	InputState  *Register16
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
	return &MCP23X17{
		rw:          rw,
		InvertPins:  NewRegister16(rw, mcp16RegIPOLA),
		PullupPins:  NewRegister16(rw, mcp16RegGPPUA),
		InputPins:   NewRegister16(rw, mcp16RegIODIRA),
		OutputState: NewRegister16(rw, mcp16RegOLATA),
		InputState:  NewRegister16(rw, mcp16RegGPIOA),
	}
}

func (MCP23X17) PinCount() int { return 16 }

func (m *MCP23X17) Flush() error {
	if err := m.PullupPins.Flush(); err != nil {
		return err
	}
	if err := m.InputPins.Flush(); err != nil {
		return err
	}
	if err := m.InvertPins.Flush(); err != nil {
		return err
	}
	if err := m.OutputState.Flush(); err != nil {
		return err
	}
	return nil
}

func (m *MCP23X17) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.InputState.Get,
		SetInputFunc: m.InputPins.Set,
		SetFunc:      m.OutputState.Set,
	}
}

func (m *MCP23X17) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      m.InputState.GetBuf,
		SetInputFunc: m.InputPins.SetBuf,
		SetFunc:      m.OutputState.SetBuf,
	}
}

func (m *MCP23X17) Refresh() error { return m.InputState.Refresh() }
