//go:build pico
// +build pico

package i2c

import (
	"device/arm"
	"device/rp"
	"machine"
)

// TODO: revisit clock timing vs. baud
type ctrl struct {
	sdaMask, sclMask uint32

	half, qtr int
}

func NewController(sda, scl machine.Pin) Controller {
	sda.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	sda.Low()
	scl.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	scl.Low()

	b := &ctrl{
		sdaMask: 1 << uint32(sda),
		sclMask: 1 << uint32(scl),
	}
	b.SetBaudRate(100e3)

	return b
}

func (c *ctrl) SetBaudRate(br uint32) {
	c.half = int(416 * 1000000 / br / 20)
	c.qtr = int(416 * 1000000 / br / 40)
}

//go:inline
func wait(cyc *int) {
	arm.AsmFull(`
	ldr {}, {cyc}
	1:
	subs {}, #1
	bne 1b
`,
		map[string]interface{}{
			"cyc": cyc,
		})
}

// TODO: set timeout based on baud
func (c *ctrl) clockUp() {
	rp.SIO.GPIO_OE_CLR.Set(c.sclMask)
	for !rp.SIO.GPIO_IN.HasBits(c.sclMask) {
		// clock stretching
	}
}

// Start will send a start condition on the bus.
func (c *ctrl) Start() {
	c.clockUp()
	wait(&c.half)
	rp.SIO.GPIO_OE_SET.Set(c.sdaMask)
	wait(&c.half)
	rp.SIO.GPIO_OE_SET.Set(c.sclMask)
	wait(&c.half)
}

func (c *ctrl) WriteBit(v bool) {
	if v {
		rp.SIO.GPIO_OE_CLR.Set(c.sdaMask)
	} else {
		rp.SIO.GPIO_OE_SET.Set(c.sdaMask)
	}
	wait(&c.half)
	c.clockUp()
	wait(&c.half)
	rp.SIO.GPIO_OE_SET.Set(c.sclMask)
	wait(&c.half)
}

func (c *ctrl) ReadBit() (value bool) {
	rp.SIO.GPIO_OE_CLR.Set(c.sdaMask)
	wait(&c.half)
	c.clockUp()
	wait(&c.half)
	value = rp.SIO.GPIO_IN.HasBits(c.sdaMask)
	if !value {
		// keep it low
		rp.SIO.GPIO_OE_SET.Set(c.sdaMask)
	}
	wait(&c.qtr)
	rp.SIO.GPIO_OE_SET.Set(c.sclMask)
	wait(&c.qtr)

	if !value {
		rp.SIO.GPIO_OE_CLR.Set(c.sdaMask)
	}
	return value
}

// Stop will send a stop condition on the bus.
func (c *ctrl) Stop() {
	rp.SIO.GPIO_OE_SET.Set(c.sdaMask)
	wait(&c.half)
	c.clockUp()
	wait(&c.half)
	rp.SIO.GPIO_OE_CLR.Set(c.sdaMask)
	wait(&c.half)
}

func I2C0() *I2C { return New(NewController(machine.I2C0_SDA_PIN, machine.I2C0_SCL_PIN)) }
