//go:build pico
// +build pico

package onewire

import (
	"device/arm"
	"device/rp"
	"machine"
)

const picoMul = 416

type ctrl struct {
	mask uint32

	a, b, c, d, e, f, g, h, i, j int
}

func NewController(p machine.Pin) Controller {
	p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	p.Low()
	b := &ctrl{mask: 1 << uint32(p)}
	b.standard()
	return b
}

func (c *ctrl) standard() {
	c.a = 6 * picoMul / 10
	c.b = 64 * picoMul / 10
	c.c = 60 * picoMul / 10
	c.d = 10 * picoMul / 10
	c.e = 9 * picoMul / 10
	c.f = 55 * picoMul / 10
	c.g = 1
	c.h = 480 * picoMul / 10
	c.i = 70 * picoMul / 10
	c.j = 410 * picoMul / 10
}

func (c *ctrl) overdrive() {
	c.a = 10 * picoMul / 100
	c.b = 75 * picoMul / 100
	c.c = 75 * picoMul / 100
	c.d = 25 * picoMul / 100
	c.e = 10 * picoMul / 100
	c.f = 70 * picoMul / 100
	c.g = 25 * picoMul / 100
	c.h = 70 * picoMul / 100
	c.i = 85 * picoMul / 100
	c.j = 40 * picoMul / 100
}

func (c *ctrl) WriteBit(v bool) {
	if v {
		rp.SIO.GPIO_OE_SET.Set(c.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &c.a,
			})
		rp.SIO.GPIO_OE_CLR.Set(c.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &c.b,
			})
	} else {
		rp.SIO.GPIO_OE_SET.Set(c.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &c.c,
			})
		rp.SIO.GPIO_OE_CLR.Set(c.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &c.d,
			})
	}
}

func (c *ctrl) ReadBit() (value bool) {
	rp.SIO.GPIO_OE_SET.Set(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.a,
		})
	rp.SIO.GPIO_OE_CLR.Set(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.e,
		})
	value = rp.SIO.GPIO_IN.HasBits(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.f,
		})
	return value
}

//go:inline
func (c *ctrl) Reset() (hasDevices bool) {
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.g,
		})
	rp.SIO.GPIO_OE_SET.Set(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.h,
		})
	rp.SIO.GPIO_OE_CLR.Set(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.i,
		})
	hasDevices = !rp.SIO.GPIO_IN.HasBits(c.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &c.j,
		})
	return hasDevices
}
