package i2c

func (i2c *I2C) WriteRegister(addr, reg byte, p []byte) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte(addr << 1); err != nil {
		return 0, err
	}

	if err := i2c.WriteByte(reg); err != nil {
		return 0, err
	}

	return i2c.Write(p)
}

func (i2c *I2C) ReadRegister(addr, reg byte, p []byte) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte((addr << 1) | 1); err != nil {
		return 0, err
	}

	if err := i2c.WriteByte(reg); err != nil {
		return 0, err
	}

	return i2c.Read(p)
}

func (i2c *I2C) Tx(addr byte, r, w []byte) error {
	i2c.Start()
	defer i2c.Stop()

	if len(w) > 0 {
		if err := i2c.WriteByte(addr << 1); err != nil {
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
		if err := i2c.WriteByte((addr << 1) | 1); err != nil {
			return err
		}

		if _, err := i2c.Read(r); err != nil {
			return err
		}
	}

	return nil
}
