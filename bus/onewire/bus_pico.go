//go:build pico
// +build pico

package onewire

import (
	"device/arm"
	"device/rp"
	"machine"
)

const picoMul = 416

type bus struct {
	mask uint32

	a, b, c, d, e, f, g, h, i, j int
}

func NewBus(p machine.Pin) Bus {
	p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	p.Low()
	b := &bus{mask: 1 << uint32(p)}
	b.standard()
	return b
}

func (b *bus) standard() {
	b.a = 6 * picoMul / 10
	b.b = 64 * picoMul / 10
	b.c = 60 * picoMul / 10
	b.d = 10 * picoMul / 10
	b.e = 9 * picoMul / 10
	b.f = 55 * picoMul / 10
	b.g = 1
	b.h = 480 * picoMul / 10
	b.i = 70 * picoMul / 10
	b.j = 410 * picoMul / 10
}

func (b *bus) overdrive() {
	b.a = 10 * picoMul / 100
	b.b = 75 * picoMul / 100
	b.c = 75 * picoMul / 100
	b.d = 25 * picoMul / 100
	b.e = 10 * picoMul / 100
	b.f = 70 * picoMul / 100
	b.g = 25 * picoMul / 100
	b.h = 70 * picoMul / 100
	b.i = 85 * picoMul / 100
	b.j = 40 * picoMul / 100
}

func (b *bus) WriteBit(v bool) {
	if v {
		rp.SIO.GPIO_OE_SET.Set(b.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &b.a,
			})
		rp.SIO.GPIO_OE_CLR.Set(b.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &b.b,
			})
	} else {
		rp.SIO.GPIO_OE_SET.Set(b.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &b.c,
			})
		rp.SIO.GPIO_OE_CLR.Set(b.mask)
		arm.AsmFull(`
			ldr {}, {cyc}
			1:
			subs {}, #1
			bne 1b
		`,
			map[string]interface{}{
				"cyc": &b.d,
			})
	}
}

func (b *bus) ReadBit() (value bool) {
	rp.SIO.GPIO_OE_SET.Set(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.a,
		})
	rp.SIO.GPIO_OE_CLR.Set(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.e,
		})
	value = rp.SIO.GPIO_IN.HasBits(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.f,
		})
	return value
}

//go:inline
func (b *bus) Reset() (hasDevices bool) {
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.g,
		})
	rp.SIO.GPIO_OE_SET.Set(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.h,
		})
	rp.SIO.GPIO_OE_CLR.Set(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.i,
		})
	hasDevices = !rp.SIO.GPIO_IN.HasBits(b.mask)
	arm.AsmFull(`
		ldr {}, {cyc}
		1:
		subs {}, #1
		bne 1b
	`,
		map[string]interface{}{
			"cyc": &b.j,
		})
	return hasDevices
}
