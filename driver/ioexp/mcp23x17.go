package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/serial"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type MCP23X17 struct {
	rw io.ReadWriter
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

var (
	_ PinReadWriter = (*MCP23X17)(nil)
	_ InputSetter   = (*MCP23X17)(nil)
	_ InvertSetter  = (*MCP23X17)(nil)
	_ PullupSetter  = (*MCP23X17)(nil)
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
	return &MCP23X17{rw}
}

func (MCP23X17) PinCount() int { return 16 }

// SetInvertPins allows inverting the polarity of input pin values.
// > If a bit is set, the corresponding GPIO register bit will
// > reflect the inverted value on the pin.
func (mcp *MCP23X17) SetInvertPins(p Valuer) error { return mcp.setPins(mcpRegIPOLA, p) }
func (mcp *MCP23X17) WritePins(p Valuer) error     { return mcp.setPins(mcpRegOLATA, p) }
func (mcp *MCP23X17) SetInputPins(p Valuer) error  { return mcp.setPins(mcpRegIODIRA, p) }
func (mcp *MCP23X17) SetPullupPins(p Valuer) error { return mcp.setPins(mcpRegGPPUA, p) }
func (mcp *MCP23X17) ReadPins() (PinState, error)  { return mcp.getPins(mcpRegGPIOA) }

func (mcp *MCP23X17) SetInvertPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcpRegIPOLA, pins, mask)
}

func (mcp *MCP23X17) WritePinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcpRegOLATA, pins, mask)
}

func (mcp *MCP23X17) SetInputPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcpRegIODIRA, pins, mask)
}

func (mcp *MCP23X17) SetPullupPinsMask(pins, mask Valuer) error {
	return mcp.maskPins(mcpRegGPPUA, pins, mask)
}

func (mcp *MCP23X17) maskPins(reg uint8, pins, mask Valuer) error {
	p, err := mcp.getPins(reg)
	if err != nil {
		return err
	}
	return mcp.setPins(reg, ApplyPinsMaskN(p, pins, mask, mcp.PinCount()))
}

func (mcp *MCP23X17) setPins(reg byte, p Valuer) error {
	v := PinUint16(p)
	_, err := mcp.rw.Write([]byte{reg, byte(v & 0xff), byte(v >> 8)})
	return err
}

func (mcp *MCP23X17) getPins(reg byte) (PinState, error) {
	var buf [2]byte
	err := serial.Tx(mcp.rw, []byte{reg}, buf[:])
	if err != nil {
		return nil, err
	}
	p := Pin16(uint16(buf[0]) | uint16(buf[1])<<8)
	return &p, nil
}
