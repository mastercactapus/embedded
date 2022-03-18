package ascii

import (
	"bytes"
	"io"
	"strings"
)

func Sprint(a ...interface{}) string {
	var buf bytes.Buffer
	Fprint(&buf, a...)
	return buf.String()
}

func Sprintln(a ...interface{}) string {
	var buf bytes.Buffer
	Fprintln(&buf, a...)
	return buf.String()
}

func Sprintf(format string, a ...interface{}) string {
	var buf bytes.Buffer
	buf.Grow(len(format))
	Fprintf(&buf, format, a...)
	return buf.String()
}

func Fprint(w io.Writer, a ...interface{}) (err error) {
	f := &fPrinter{w: w}
	f.print(a)

	return f.err
}

func Fprintln(w io.Writer, a ...interface{}) (err error) {
	f := &fPrinter{w: w}
	f.print(a)
	f.writeString("\r\n")

	return f.err
}

func Fprintf(w io.Writer, format string, a ...interface{}) (err error) {
	f := &fPrinter{w: w, format: format, args: a}
	f.printf(format, a)
	return f.err
}

type fPrinter struct {
	w io.Writer

	format string
	args   []interface{}
	state  printState

	isErrorf   bool
	wrappedErr error

	err error

	param fmtParam
}

func (f *fPrinter) print(a []interface{}) {
	f.state = f.parseFormat

	for _, a := range a {
		f.printVariable(fmtParam{code: 'v'}, a)
	}
}

func (f *fPrinter) printf(format string, a []interface{}) {
	f.state = f.parseFormat

	for _, b := range []byte(format) {
		f.state = f.state(b)
	}

	for _, a := range f.args {
		f.writeString("%!(EXTRA ")
		f.printTypeName(a)
		f.writeString("=")
		f.printVariable(fmtParam{code: 'v'}, a)
		f.writeString(")")
	}
}

type printState func(byte) printState

func (f *fPrinter) write(p []byte) {
	if f.err != nil {
		return
	}
	_, f.err = f.w.Write(p)
}

func (f *fPrinter) writeByte(p byte) { f.write([]byte{p}) }

func (f *fPrinter) writeString(s string) { f.write([]byte(s)) }

type fmtParam struct {
	sign     bool
	leftJust bool
	alt      bool
	zero     bool
	code     byte
	pad      int
}

func (f *fPrinter) parseFormat(b byte) printState {
	switch b {
	case '%':
		return f.parseParam
	default:
		f.writeByte(b)
	}
	return f.parseFormat
}

func (f *fPrinter) parseParam(b byte) printState {
	switch b {
	case '%':
		f.writeByte(b)
		return f.parseFormat
	case '0':
		f.param.zero = true
	case '#':
		f.param.alt = true
	case ' ':
		f.param.zero = false
	case '-':
		f.param.leftJust = true
	case '+':
		f.param.sign = true
	default:
		if b >= '0' && b <= '9' {
			f.param.pad = f.param.pad*10 + int(b-'0')
			break
		}

		f.param.code = b
		return f.applyParam()
	}

	return f.parseParam
}

func (f *fPrinter) printTypeName(v interface{}) {
	switch v.(type) {
	case string:
		f.writeString("string")
	case []byte:
		f.writeString("[]byte")
	case int:
		f.writeString("int")
	case int8:
		f.writeString("int8")
	case int16:
		f.writeString("int16")
	case int32:
		f.writeString("int32")
	case int64:
		f.writeString("int64")
	case uint:
		f.writeString("uint")
	case uint8:
		f.writeString("uint8")
	case uint16:
		f.writeString("uint16")
	case uint32:
		f.writeString("uint32")
	case uint64:
		f.writeString("uint64")
	case float32:
		f.writeString("float32")
	case float64:
		f.writeString("float64")
	case complex64:
		f.writeString("complex64")
	case complex128:
		f.writeString("complex128")
	case bool:
		f.writeString("bool")
	default:
		f.writeString("<unknown>")
	}
}

func (f *fPrinter) applyParam() printState {
	if len(f.args) == 0 {
		f.writeString("%!")
		f.writeByte(f.param.code)
		f.writeString("(MISSING)")
		return f.parseFormat
	}

	switch f.param.code {
	case 'c':
		f.printChar(f.param, f.args[0])
	case 'd':
		f.printInt(f.param, f.args[0])
	case 's':
		f.printStr(f.param, f.args[0])
	case 'x', 'X':
		f.printHex(f.param, f.args[0])
	default:
		f.printVariable(f.param, f.args[0])
	}

	f.args = f.args[1:]
	return f.parseFormat
}

func (f *fPrinter) printChar(p fmtParam, v interface{}) {
	printInt := func() {
		if !p.leftJust {
			for i := 10; i < p.pad; i++ {
				f.writeByte(' ')
			}
		}
		// TODO: actually encode
		f.writeString(`\U`)
		f.printHex(fmtParam{code: 'x', zero: true, pad: 8}, v)
		if p.leftJust {
			for i := 10; i < p.pad; i++ {
				f.writeByte(' ')
			}
		}
	}

	var b byte
	switch v := v.(type) {
	case int:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case int8:
		b = byte(v)
	case int16:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case int32:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case int64:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case uint:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case uint8:
		b = byte(v)
	case uint16:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case uint32:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case uint64:
		if v > 127 {
			printInt()
			return
		}
		b = byte(v)
	case string:
	default:
		f.printVariable(p, v)
		return
	}

	if p.pad < 2 {
		f.writeByte(b)
		return
	}

	if p.leftJust {
		f.writeByte(b)
	}
	for i := 0; i < p.pad-1; i++ {
		f.writeByte(' ')
	}
	if !p.leftJust {
		f.writeByte(b)
	}
}

func (f *fPrinter) printVariable(p fmtParam, v interface{}) printState {
	if f.isErrorf && f.wrappedErr == nil && p.code == 'w' {
		err, ok := v.(error)
		if ok {
			f.wrappedErr = err
			f.writeString(err.Error())
			return f.parseFormat
		}
	}
	if p.code != 'v' {
		f.writeString("%!")
		f.writeByte(p.code)
		f.writeString("(")
		f.printTypeName(v)
		f.writeString("=")
	}
	switch v := v.(type) {
	case bool:
		if v {
			f.writeString("true")
		} else {
			f.writeString("false")
		}
	case string:
		f.printStr(p, v)
	case []string:
		f.printStr(p, "{"+strings.Join(v, " ")+"}")
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
		f.printInt(p, v)
	case []byte:
		f.writeByte('[')
		for i, b := range v {
			if i > 0 {
				f.writeByte(' ')
			}
			f.printInt(fmtParam{}, b)
		}
		f.writeByte(']')
	default:
		f.printStr(p, "?")
	}
	if p.code != 'v' {
		f.writeString(")")
	}
	return f.parseFormat
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
		f.printVariable(p, v)
		return
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
	var b uint64
	var neg bool
	switch v := v.(type) {
	case uint:
		b = uint64(v)
	case uint8:
		b = uint64(v)
	case uint16:
		b = uint64(v)
	case uint32:
		b = uint64(v)
	case uint64:
		b = v
	case int:
		if v < 0 {
			neg = true
			b = uint64(-v)
		} else {
			b = uint64(v)
		}
	case int8:
		if v < 0 {
			neg = true
			b = uint64(-v)
		} else {
			b = uint64(v)
		}
	case int16:
		if v < 0 {
			neg = true
			b = uint64(-v)
		} else {
			b = uint64(v)
		}
	case int32:
		if v < 0 {
			neg = true
			b = uint64(-v)
		} else {
			b = uint64(v)
		}
	case int64:
		if v < 0 {
			neg = true
			b = uint64(-v)
		} else {
			b = uint64(v)
		}
	default:
		f.printVariable(p, v)
		return
	}
	var buf bytes.Buffer
	for b > 0 {
		buf.WriteByte(byte(b%10 + '0'))
		b /= 10
	}
	if neg || p.sign {
		p.pad--
	}

	if !p.leftJust {
		if p.zero {
			if neg {
				f.writeByte('-')
			} else if p.sign {
				f.writeByte('+')
			}
		}
		for i := 0; i < p.pad-buf.Len(); i++ {
			if p.zero {
				f.writeByte('0')
			} else {
				f.writeByte(' ')
			}
		}
		if !p.zero {
			if neg {
				f.writeByte('-')
			} else if p.sign {
				f.writeByte('+')
			}
		}
	} else {
		if neg {
			f.writeByte('-')
		} else if p.sign {
			f.writeByte('+')
		}
	}
	// write bytes in reverse order
	for i := buf.Len() - 1; i >= 0; i-- {
		f.writeByte(buf.Bytes()[i])
	}

	if p.leftJust {
		for i := p.pad - buf.Len(); i != 0; i-- {
			f.writeByte(' ')
		}
	}
}

const maxInt = 32 << (^uint(0) >> 63)

func (f *fPrinter) printHex(p fmtParam, v interface{}) {
	if p.alt {
		if p.code == 'X' {
			f.write([]byte("0X"))
		} else {
			f.write([]byte("0x"))
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
		f.printVariable(p, v)
		return
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
		if p.zero {
			f.writeByte('0')
		} else {
			f.writeByte(' ')
		}
	}

	for i, b := range b {
		if i > 0 || b>>4 != 0 {
			f.writeByte(digits[b>>4])
		}
		f.writeByte(digits[b&0xf])
	}
}

const digits = "0123456789abcdef"
