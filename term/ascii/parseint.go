package ascii

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
		return 0, Errorf("term: invalid int: %s", s)
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
				return 0, Errorf("term: invalid int: %s", s)
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
				return 0, Errorf("term: invalid int: %s", s)
			}
		}
	case buf[0] == '0' && buf[1] == 'o':
		for _, c := range buf[2:] {
			switch {
			case c >= '0' && c <= '7':
				val = val<<3 + int(c-'0')
			default:
				return 0, Errorf("term: invalid int: %s", s)
			}
		}
	default:
		for _, c := range buf {
			if c < '0' || c > '9' {
				return 0, Errorf("term: invalid int: %s", s)
			}
			val = val*10 + int(c-'0')
		}
	}

	if neg {
		val = -val
	}
	return val, nil
}
