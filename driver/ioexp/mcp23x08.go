package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/serial"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type MCP23X08 struct {
	rw io.ReadWriter
}

const (
	mcp8RegOLAT  = 0x0a
	mcp8RegIPOL  = 0x01
	mcp8RegGPPU  = 0x06
	mcp8RegIODIR = 0x00
	mcp8RegGPIO  = 0x09
)

var (
	_ PinReadWriter = (*MCP23X08)(nil)
	_ InputSetter   = (*MCP23X08)(nil)
	_ InvertSetter  = (*MCP23X08)(nil)
	_ PullupSetter  = (*MCP23X08)(nil)
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
	return &MCP23X08{rw}
}

func (MCP23X08) PinCount() int { return 8 }

// SetInvertPins allows inverting the polarity of input pin values.
// > If a bit is set, the corresponding GPIO register bit will
// > reflect the inverted value on the pin.
func (mcp *MCP23X08) SetInvertPins(p Valuer) error { return mcp.setPins(mcpRegIPOLA, p) }
func (mcp *MCP23X08) WritePins(p Valuer) error     { return mcp.setPins(mcpRegOLATA, p) }
func (mcp *MCP23X08) SetInputPins(p Valuer) error  { return mcp.setPins(mcpRegIODIRA, p) }
func (mcp *MCP23X08) SetPullupPins(p Valuer) error { return mcp.setPins(mcpRegGPPUA, p) }
func (mcp *MCP23X08) ReadPins() (PinState, error)  { return mcp.getPins(mcpRegGPIOA) }

func (mcp *MCP23X08) SetInvertPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcp8RegIPOL, pins, mask)
}

func (mcp *MCP23X08) WritePinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcp8RegOLAT, pins, mask)
}

func (mcp *MCP23X08) SetInputPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcp8RegIODIR, pins, mask)
}

func (mcp *MCP23X08) SetPullupPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcp8RegGPPU, pins, mask)
}

func (mcp *MCP23X08) maskPins(reg uint8, pins, mask Valuer) error {
	p, err := mcp.getPins(reg)
	if err != nil {
		return err
	}
	return mcp.setPins(reg, ApplyPinsMaskN(p, pins, mask, mcp.PinCount()))
}

func (mcp *MCP23X08) setPins(reg byte, p Valuer) error {
	_, err := mcp.rw.Write([]byte{reg, PinByte(p)})
	return err
}

func (mcp *MCP23X08) getPins(reg byte) (PinState, error) {
	var buf [1]byte
	err := serial.Tx(mcp.rw, []byte{reg}, buf[:])
	if err != nil {
		return nil, err
	}
	p := Pin8(buf[0])
	return &p, nil
}
