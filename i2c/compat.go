package i2c

import (
	"tinygo.org/x/drivers"
)

var _ drivers.I2C = (*I2C)(nil)

func (i2c *I2C) WriteRegister(addr, reg byte, p []byte) error {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte(addr << 1); err != nil {
		return err
	}

	if err := i2c.WriteByte(reg); err != nil {
		return err
	}

	_, err := i2c.Write(p)
	return err
}

func (i2c *I2C) ReadRegister(addr, reg byte, p []byte) error {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte((addr << 1) | 1); err != nil {
		return err
	}

	if err := i2c.WriteByte(reg); err != nil {
		return err
	}

	_, err := i2c.Read(p)
	return err
}

func (i2c *I2C) Tx(addr uint16, w, r []byte) error {
	i2c.Start()
	defer i2c.Stop()

	if len(w) > 0 {
		if err := i2c.writeAddress(addr, modeWrite); err != nil {
			return err
		}

		if _, err := i2c.Write(w); err != nil {
			return err
		}
	}

	if len(r) > 0 {
		if len(w) > 0 {
			// repeated start
			i2c.Start()
		}
		if err := i2c.writeAddress(addr, modeRead); err != nil {
			return err
		}

		if _, err := i2c.Read(r); err != nil {
			return err
		}
	}

	return nil
}
