package mem

import (
	"io"
	"time"

	"github.com/mastercactapus/embedded/i2c"
)

// NewAT24C01 creates a AT24C01-compatible Pager.
//
// You most likely want NewAT24C01A unless you are sure you need the 4-byte page size.
func NewAT24C01(bus i2c.Bus, addr uint16) *Pager {
	if addr == 0 {
		addr = 0x50
	}
	return NewPager(i2c.NewDevice(bus, addr), PagerConfig{
		PageSize:    4,
		AddressSize: 1,
		Capacity:    128,
		WriteDelay:  5 * time.Millisecond,
	})
}

// NewAT24C01A introduced an 8-byte page size (as opposed to the original 4-byte page size)
//
// Future versions (B,C,D etc...) should be compatible and use this.
func NewAT24C01A(bus i2c.Bus, addr uint16) *Pager {
	if addr == 0 {
		addr = 0x50
	}
	return NewPager(i2c.NewDevice(bus, addr), PagerConfig{
		PageSize:    8,
		AddressSize: 1,
		Capacity:    128,
		WriteDelay:  5 * time.Millisecond,
	})
}

func NewAT24C02(bus i2c.Bus, addr uint16) *Pager {
	if addr == 0 {
		addr = 0x50
	}
	return NewPager(i2c.NewDevice(bus, addr), PagerConfig{
		PageSize:    8,
		AddressSize: 1,
		Capacity:    256,
		WriteDelay:  5 * time.Millisecond,
	})
}

func newMultiDevice(bus i2c.Bus, baseAddr uint16, cfg PagerConfig, numDevices int) io.ReadWriteSeeker {
	if baseAddr == 0 {
		baseAddr = 0x50
	}

	var devs []io.ReadWriteSeeker
	for i := 0; i < numDevices; i++ {
		devs = append(devs, NewPager(i2c.NewDevice(bus, baseAddr+uint16(i)), cfg))
	}

	return Join(devs...)
}

func NewAT24C04(bus i2c.Bus, addr uint16) io.ReadWriteSeeker {
	if addr == 0 {
		addr = 0x50
	}

	return newMultiDevice(bus, addr, PagerConfig{
		PageSize:    16,
		AddressSize: 1,
		Capacity:    256,
		WriteDelay:  5 * time.Millisecond,
	}, 2)
}

func NewAT24C08(bus i2c.Bus, addr uint16) io.ReadWriteSeeker {
	if addr == 0 {
		addr = 0x50
	}

	return newMultiDevice(bus, addr, PagerConfig{
		PageSize:    16,
		AddressSize: 1,
		Capacity:    256,
		WriteDelay:  5 * time.Millisecond,
	}, 4)
}

func NewAT24C16(bus i2c.Bus, addr uint16) io.ReadWriteSeeker {
	if addr == 0 {
		addr = 0x50
	}

	return newMultiDevice(bus, addr, PagerConfig{
		PageSize:    16,
		AddressSize: 1,
		Capacity:    256,
		WriteDelay:  5 * time.Millisecond,
	}, 8)
}

func NewAT24C32(bus i2c.Bus, addr uint16) *Pager {
	if addr == 0 {
		addr = 0x50
	}
	return NewPager(i2c.NewDevice(bus, addr), PagerConfig{
		PageSize:    32,
		AddressSize: 2,
		Capacity:    4096,
		WriteDelay:  5 * time.Millisecond,
	})
}

func NewAT24C64(bus i2c.Bus, addr uint16) *Pager {
	if addr == 0 {
		addr = 0x50
	}
	return NewPager(i2c.NewDevice(bus, addr), PagerConfig{
		PageSize:    32,
		AddressSize: 2,
		Capacity:    8192,
		WriteDelay:  5 * time.Millisecond,
	})
}
