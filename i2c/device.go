package i2c

type Bus interface {
	Tx(addr uint16, w, r []byte) error
}

type Device struct {
	bus  Bus
	addr uint16
}
type WriterTo interface {
	WriteTo([]byte, uint16) (int, error)
}
type ReaderFrom interface {
	ReadFrom([]byte, uint16) (int, error)
}
type (
	ByteWriterTo   interface{ WriteByteTo(byte, uint16) error }
	ByteReaderFrom interface{ ReadByteFrom(uint16) (byte, error) }
)

// NewDevice returns a new Device with the given bus and address.
func NewDevice(bus Bus, addr uint16) *Device {
	return &Device{bus: bus, addr: addr}
}

func (d *Device) Write(w []byte) (int, error) {
	if wt := d.bus.(WriterTo); wt != nil {
		return wt.WriteTo(w, d.addr)
	}

	err := d.bus.Tx(d.addr, w, nil)
	if err != nil {
		return 0, err
	}

	return len(w), nil
}

func (d *Device) Read(r []byte) (int, error) {
	if rf := d.bus.(ReaderFrom); rf != nil {
		return rf.ReadFrom(r, d.addr)
	}

	err := d.bus.Tx(d.addr, nil, r)
	if err != nil {
		return 0, err
	}
	return len(r), nil
}

func (d *Device) WriteByte(b byte) error {
	if w := d.bus.(ByteWriterTo); w != nil {
		return w.WriteByteTo(b, d.addr)
	}

	_, err := d.Write([]byte{b})
	return err
}

func (d *Device) ReadByte() (byte, error) {
	if r := d.bus.(ByteReaderFrom); r != nil {
		return r.ReadByteFrom(d.addr)
	}

	var b [1]byte
	_, err := d.Read(b[:])
	return b[0], err
}
