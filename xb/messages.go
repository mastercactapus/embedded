package xb

const (
	ignore uint8 = iota
	reset
	hello
	setInput
	setPin
	getPin

	i2cSetup
	i2cTx

	spiSetup
	spiRead
	spiSetFill
	spiWrite
	spiReadWrite
	spiReadWriteByte
)

type SPIConfig struct {
	Mode uint8
	Baud uint32

	MISO, MOSI, SCLK uint8
}

func (cfg *SPIConfig) encode() []byte {
	var m msgEnc
	m.addByte(1, cfg.Mode)
	m.addUint32(2, cfg.Baud)
	m.addByte(3, cfg.MISO)
	m.addByte(4, cfg.MOSI)
	m.addByte(5, cfg.SCLK)
	return m.data
}

func (cfg *SPIConfig) decode(b []byte) {
	m := msgDec{b}
	cfg.Mode = m.getByte(1)
	cfg.Baud = m.getUint32(2)
	cfg.MISO = m.getByte(3)
	cfg.MOSI = m.getByte(4)
	cfg.SCLK = m.getByte(5)
}

func (cfg *I2CConfig) encode() []byte {
	var m msgEnc
	m.addByte(1, cfg.SDA)
	m.addByte(2, cfg.SCL)
	return m.data
}

func (cfg *I2CConfig) decode(data []byte) {
	var m msgDec
	m.data = data
	cfg.SDA = m.getByte(1)
	cfg.SCL = m.getByte(2)
}

type I2CConfig struct {
	SDA, SCL uint8
}

type Request struct {
	Cmd uint8 `json:"c,omitempty"`

	Pin   uint8 `json:"p,omitempty"`
	State bool  `json:"s,omitempty"`

	I2CAddr  uint16 `json:"a,omitempty"`
	DataByte byte   `json:"b,omitempty"`
	ReadN    uint16 `json:"n,omitempty"`

	Data []byte `json:"d,omitempty"`

	I2CConfig *I2CConfig `json:"i2c,omitempty"`
	SPIConfig *SPIConfig `json:"spi,omitempty"`
}

type msgEnc struct {
	data []byte
}
type msgDec struct {
	data []byte
}

func (m *msgEnc) addByte(field uint8, val byte) {
	if val == 0 {
		return
	}
	m.data = append(m.data, field, val)
}

func (m *msgDec) getByte(field uint8) (b byte) {
	if len(m.data) < 2 {
		return 0
	}
	if m.data[0] != field {
		return 0
	}
	b = m.data[1]
	m.data = m.data[2:]
	return b
}

func (m *msgEnc) addBool(field uint8, val bool) {
	if val {
		m.data = append(m.data, field)
	}
}

func (m *msgDec) getBool(field uint8) bool {
	if len(m.data) < 1 {
		return false
	}
	if m.data[0] != field {
		return false
	}
	m.data = m.data[1:]
	return true
}

func (m *msgEnc) addUint16(field uint8, val uint16) {
	if val == 0 {
		return
	}
	m.data = append(m.data, field, byte(val>>8), byte(val))
}

func (m *msgDec) getUint16(field uint8) uint16 {
	if len(m.data) < 3 {
		return 0
	}
	if m.data[0] != field {
		return 0
	}
	v := uint16(m.data[1])<<8 | uint16(m.data[2])
	m.data = m.data[3:]
	return v
}

func (m *msgEnc) addUint32(field uint8, val uint32) {
	if val == 0 {
		return
	}
	m.data = append(m.data, field, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func (m *msgDec) getUint32(field uint8) uint32 {
	if len(m.data) < 5 {
		return 0
	}
	if m.data[0] != field {
		return 0
	}
	v := uint32(m.data[1])<<24 | uint32(m.data[2])<<16 | uint32(m.data[3])<<8 | uint32(m.data[4])
	m.data = m.data[5:]
	return v
}

func (m *msgEnc) addString(field uint8, val string) {
	if val == "" {
		return
	}
	m.addData(field, []byte(val))
}

func (m *msgDec) getString(field uint8) string {
	return string(m.getData(field))
}

func (m *msgEnc) addData(field uint8, data []byte) {
	if len(data) == 0 {
		return
	}

	m.data = append(m.data, field, byte(len(data)>>8), byte(len(data)))
	m.data = append(m.data, data...)
}

func (m *msgDec) getData(field uint8) []byte {
	if len(m.data) < 3 {
		return nil
	}
	if m.data[0] != field {
		return nil
	}
	l := uint16(m.data[1])<<8 | uint16(m.data[2])
	if len(m.data) < int(l+3) {
		return nil
	}
	d := m.data[3 : l+3]
	m.data = m.data[l+3:]
	return d
}

func (req *Request) encode() []byte {
	var m msgEnc
	m.addByte('c', req.Cmd)
	m.addByte('p', req.Pin)
	m.addBool('s', req.State)
	m.addUint16('a', req.I2CAddr)
	m.addByte('b', req.DataByte)
	m.addUint16('n', req.ReadN)
	m.addData('d', req.Data)

	if req.I2CConfig != nil {
		m.addData('I', req.I2CConfig.encode())
	}
	if req.SPIConfig != nil {
		m.addData('S', req.SPIConfig.encode())
	}

	return m.data
}

func (req *Request) decode(data []byte) {
	m := msgDec{data}
	req.Cmd = m.getByte('c')
	req.Pin = m.getByte('p')
	req.State = m.getBool('s')
	req.I2CAddr = m.getUint16('a')
	req.DataByte = m.getByte('b')
	req.ReadN = m.getUint16('n')
	req.Data = m.getData('d')

	i2cData := m.getData('I')
	if len(i2cData) > 0 {
		req.I2CConfig = &I2CConfig{}
		req.I2CConfig.decode(i2cData)
	} else {
		req.I2CConfig = nil
	}

	spiData := m.getData('S')
	if len(spiData) > 0 {
		req.SPIConfig = &SPIConfig{}
		req.SPIConfig.decode(spiData)
	} else {
		req.SPIConfig = nil
	}
}

type Response struct {
	Err      string `json:"e,omitempty"`
	State    bool   `json:"s,omitempty"`
	PinCount uint8  `json:"n,omitempty"`
	DataByte byte   `json:"b,omitempty"`
	Data     []byte `json:"d,omitempty"`
}

func (resp *Response) encode() []byte {
	var m msgEnc
	m.addString('e', resp.Err)
	m.addBool('s', resp.State)
	m.addByte('p', resp.PinCount)
	m.addByte('b', resp.DataByte)
	m.addData('d', resp.Data)
	return m.data
}

func (resp *Response) decode(data []byte) {
	m := msgDec{data}
	resp.Err = m.getString('e')
	resp.State = m.getBool('s')
	resp.PinCount = m.getByte('p')
	resp.DataByte = m.getByte('b')
	resp.Data = m.getData('d')
}
