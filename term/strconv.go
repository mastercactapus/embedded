package term

import "github.com/mastercactapus/embedded/term/ascii"

func fmtBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func fmtUint16(v uint16) string {
	var buf [6]byte
	buf[0] = '0'
	buf[1] = 'x'
	buf[2] = byte(v>>12) + '0'
	buf[3] = byte(v>>8) + '0'
	buf[4] = byte(v>>4) + '0'
	buf[5] = byte(v&0xf) + '0'
	return string(buf[:])
}

func fmtByte(v byte) string {
	var buf [4]byte
	buf[0] = '0'
	buf[1] = 'x'
	buf[2] = byte(v>>4) + '0'
	buf[3] = byte(v&0xf) + '0'
	return string(buf[:])
}

func parseBool(s string) (bool, error) {
	switch s {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	default:
		return false, ascii.Errorf("term: invalid bool: %q", s)
	}
}

func ParseInt(s string) (val int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	buf := []byte(s)
	var neg bool
	if buf[0] == '-' {
		neg = true
		buf = buf[1:]
	}

	if len(s) == 1 {
		if buf[0] >= '0' && buf[0] <= '9' {
			val = int(buf[0] - '0')
			if neg {
				val = -val
			}
			return val, nil
		}
		return 0, ascii.Errorf("term: invalid int: %s", s)
	}
	switch {
	case buf[0] == '0' && buf[1] == 'x':
		for _, c := range buf[2:] {
			switch {
			case c >= '0' && c <= '9':
				val = val<<4 + int(c-'0')
			case c >= 'a' && c <= 'f':
				val = val<<4 + int(c-'a'+10)
			case c >= 'A' && c <= 'F':
				val = val<<4 + int(c-'A'+10)
			default:
				return 0, ascii.Errorf("term: invalid int: %s", s)
			}
		}
	case buf[0] == '0' && buf[1] == 'b':
		for _, c := range buf[2:] {
			switch c {
			case '0':
				val = val << 1
			case '1':
				val = val<<1 + 1
			default:
				return 0, ascii.Errorf("term: invalid int: %s", s)
			}
		}
	case buf[0] == '0' && buf[1] == 'o':
		for _, c := range buf[2:] {
			switch {
			case c >= '0' && c <= '7':
				val = val<<3 + int(c-'0')
			default:
				return 0, ascii.Errorf("term: invalid int: %s", s)
			}
		}
	default:
		for _, c := range buf {
			if c < '0' || c > '9' {
				return 0, ascii.Errorf("term: invalid int: %s", s)
			}
			val = val*10 + int(c-'0')
		}
	}

	if neg {
		val = -val
	}
	return val, nil
}

func parseUint16(s string) (val uint16, err error) {
	if len(s) == 0 {
		return 0, nil
	}

	if len(s) == 1 {
		if s[0] >= '0' && s[0] <= '9' {
			return uint16(s[0] - '0'), nil
		}
		return 0, ascii.Errorf("term: invalid uint16: %s", s)
	}
	if s[0] == '0' && s[1] == 'x' {
		if len(s) > 8 {
			return 0, ascii.Errorf("term: invalid uint16: %s", s)
		}
		for _, c := range s[2:] {
			switch {
			case c >= '0' && c <= '9':
				val = val<<4 + uint16(c-'0')
			case c >= 'a' && c <= 'f':
				val = val<<4 + uint16(c-'a'+10)
			case c >= 'A' && c <= 'F':
				val = val<<4 + uint16(c-'A'+10)
			default:
				return 0, ascii.Errorf("term: invalid uint16: %s", s)
			}
		}
		return val, nil
	}
	if s[0] == '0' && s[1] == 'b' {
		if len(s) > 16 {
			return 0, ascii.Errorf("term: invalid uint16: %s", s)
		}
		for _, c := range s[2:] {
			switch c {
			case '0':
				val = val << 1
			case '1':
				val = val<<1 + 1
			default:
				return 0, ascii.Errorf("term: invalid uint16: %s", s)
			}
		}
		return val, nil
	}
	if s[0] == '-' {
		return 0, ascii.Errorf("term: invalid uint16: %s", s)
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, ascii.Errorf("term: invalid uint16: %s", s)
		}
		val = val*10 + uint16(c-'0')
	}
	return val, nil
}

func parseUint8(s string) (val uint8, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	if s[0] == 's' {
		return uint8(s[1]), nil
	}
	if len(s) == 1 {
		if s[0] >= '0' && s[0] <= '9' {
			return uint8(s[0] - '0'), nil
		}
		return 0, ascii.Errorf("term: invalid uint8: %s", s)
	}
	if s[0] == '0' && s[1] == 'x' {
		if len(s) > 4 {
			return 0, ascii.Errorf("term: invalid uint8: %s", s)
		}
		for _, c := range s[2:] {
			switch {
			case c >= '0' && c <= '9':
				val = val<<4 + uint8(c-'0')
			case c >= 'a' && c <= 'f':
				val = val<<4 + uint8(c-'a'+10)
			case c >= 'A' && c <= 'F':
				val = val<<4 + uint8(c-'A'+10)
			default:
				return 0, ascii.Errorf("term: invalid uint8: %s", s)
			}
		}
		return val, nil
	}
	if s[0] == '0' && s[1] == 'b' {
		if len(s) > 8 {
			return 0, ascii.Errorf("term: invalid uint8: %s", s)
		}
		for _, c := range s[2:] {
			switch c {
			case '0':
				val = val << 1
			case '1':
				val = val<<1 + 1
			default:
				return 0, ascii.Errorf("term: invalid uint8: %s", s)
			}
		}
		return val, nil
	}
	if s[0] == '-' {
		return 0, ascii.Errorf("term: invalid uint8: %s", s)
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, ascii.Errorf("term: invalid uint8: %s", s)
		}
		val = val*10 + uint8(c-'0')
	}
	return val, nil
}

func itoa(i int) string {
	var buf [7]byte
	if i >= 1000000 || i <= -1000000 {
		panic("term: itoa only supports 6-digit numbers")
	}
	if i == 0 {
		return "0"
	}

	var n int
	if i < 0 {
		n = 1
		i = -i
		buf[0] = '-'
	}
	for ; i > 0; i /= 10 {
		buf[n] = byte(i%10) + '0'
		n++
	}

	// reverse buf
	i, j := 0, n-1
	if buf[0] == '-' {
		i = 1
	}
	for ; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf[:n])
}
