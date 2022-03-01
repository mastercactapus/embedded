package i2c

import "log"

type DeviceID uint32

const reservedAddr = 0x78

func (i2c *I2C) Ping(addr byte) error {
	i2c.Start()
	defer i2c.Stop()
	log.Println("ye")
	return i2c.WriteByte((addr << 1) | 1)
}

func (i2c *I2C) PingW(addr byte) error {
	i2c.Start()
	defer i2c.Stop()

	return i2c.WriteByte(addr << 1)
}

func (i2c *I2C) DeviceID(addr byte) (DeviceID, error) {
	i2c.Start()
	defer i2c.Stop()

	// ignore nack
	i2c.WriteByte(reservedAddr)

	if err := i2c.WriteByte(addr << 1); err != nil {
		return 0, err
	}

	i2c.Start()

	i2c.WriteByte(reservedAddr | 1)

	var buf [3]byte
	_, err := i2c.Read(buf[:])
	if err != nil {
		return 0, err
	}

	var id DeviceID
	id |= DeviceID(buf[0]) << 16
	id |= DeviceID(buf[1]) << 8
	id |= DeviceID(buf[2])
	return id, err
}
