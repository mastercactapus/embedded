package ansi

import (
	"io"
)

func Errorf(format string, a ...interface{}) (err error) {
	return nil
}

func Fprintf(w io.Writer, format string, a ...interface{}) (err error) {
	return nil
}

type fPrinter struct {
	w   io.Writer
	err error
}

func (f *fPrinter) write(p []byte) {
}

func (f *fPrinter) writeByte(p byte) {
}

type fmtParam struct {
	sign     bool
	leftJust bool
	alt      bool
	space    bool
	zero     bool
	code     byte
	pad      int
}

func typeName(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case []byte:
		return "[]byte"
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64:
		return "uint"
	default:
		return "unknown"
	}
}

func (f *fPrinter) printStr(p fmtParam, v interface{}) {
	var b []byte
	switch v := v.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	case byte:
		b = []byte{v}
	default:
		panic("unsupported type for %s: " + typeName(v))
	}
	if p.pad == 0 || len(b) >= p.pad {
		f.write(b)
		return
	}

	if p.leftJust {
		f.write(b)
	}
	for i := 0; i < p.pad-len(b); i++ {
		f.writeByte(' ')
	}
	if !p.leftJust {
		f.write(b)
	}
}

func (f *fPrinter) printInt(p fmtParam, v interface{}) {
}

const maxInt = 32 << (^uint(0) >> 63)

func (f *fPrinter) printHex(p fmtParam, v interface{}) {
	if p.alt {
		if p.code == 'X' {
			f.write([]byte("0X"))
		} else {
			f.write([]byte("0X"))
		}
	}

	var b []byte
	switch v := v.(type) {
	case []byte:
		b = v
		if p.pad == 0 {
			p.pad = len(v) * 2
		}
	case string:
		b = []byte(v)
		if p.pad == 0 {
			p.pad = len(v) * 2
		}
	case int:
		if maxInt == 32 {
			b = []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		} else {
			b = []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		}
	case int8:
		b = []byte{byte(v)}
	case int16:
		b = []byte{byte(v >> 8), byte(v)}
	case int32:
		b = []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	case int64:
		b = []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	case uint:
		if maxInt == 32 {
			b = []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		} else {
			b = []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		}
	case uint8:
		b = []byte{byte(v)}
	case uint16:
		b = []byte{byte(v >> 8), byte(v)}
	case uint32:
		b = []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	case uint64:
		b = []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	default:
		panic("unsupported type for %x: " + typeName(v))
	}

	for i, c := range b {
		if c == 0 {
			continue
		}
		b = b[i:]
		break
	}

	hLen := len(b) * 2
	if len(b) > 0 && b[0] < 0x10 {
		hLen--
	}

	for i := 0; i < p.pad-hLen; i++ {
		if p.space {
			f.writeByte(' ')
		} else {
			f.writeByte('0')
		}
	}

	for i, b := range b {
		if i > 0 || b>>4 != 0 {
			f.writeByte(b>>4 + '0')
		}
		f.writeByte(b&0xf + '0')
	}
}
