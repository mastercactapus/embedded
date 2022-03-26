package rtc

import (
	"io"
	"time"

	"github.com/mastercactapus/embedded/serial"
)

func NewDS3231(rw io.ReadWriter) *DS3231 {
	return &DS3231{rw: rw}
}

type DS3231 struct {
	rw io.ReadWriter
}

type DS3231Alarm struct {
	Seconds int
	Minutes int
	Hours   int
	IsPM    bool // Only valid if Use12Hour is true.
	Day     int

	// UseWeekday indicates Day matches the day of the week.
	UseWeekday bool

	// Use12Hour indicates the hour is in 12-hour format.
	Use12Hour bool

	IgnoreSeconds bool // A1M1
	IgnoreMinutes bool // A1M2,A2M2
	IgnoreHours   bool // A1M3,A2M3
	IgnoreDay     bool // A1M4,A2M4
}

func (rtc *DS3231) SetAlarm2(a DS3231Alarm) error {
	var buf [4]byte
	buf[0] = 0x0b
	buf[1] = dsEncMinSec(a.Minutes)
	if a.Use12Hour {
		buf[1] = dsEnc12Hour(a.Hours, a.IsPM)
	} else {
		buf[2] = dsEnc24Hour(a.Hours)
	}
	buf[3] = dsEncDate(a.Day)
	if a.UseWeekday {
		buf[3] |= 1 << 6
	}
	if a.IgnoreSeconds {
		buf[0] |= 1 << 7
	}
	if a.IgnoreMinutes {
		buf[1] |= 1 << 7
	}
	if a.IgnoreHours {
		buf[2] |= 1 << 7
	}
	if a.IgnoreDay {
		buf[3] |= 1 << 7
	}
	return serial.Tx(rtc.rw, buf[:], nil)
}

func (rtc *DS3231) SetAlarm1(a DS3231Alarm) error {
	var buf [5]byte
	buf[0] = 0x07
	buf[1] = dsEncMinSec(a.Seconds)
	buf[2] = dsEncMinSec(a.Minutes)
	if a.Use12Hour {
		buf[3] = dsEnc12Hour(a.Hours, a.IsPM)
	} else {
		buf[3] = dsEnc24Hour(a.Hours)
	}
	buf[4] = dsEncDate(a.Day)
	if a.UseWeekday {
		buf[4] |= 1 << 6
	}
	if a.IgnoreSeconds {
		buf[1] |= 1 << 7
	}
	if a.IgnoreMinutes {
		buf[2] |= 1 << 7
	}
	if a.IgnoreHours {
		buf[3] |= 1 << 7
	}
	if a.IgnoreDay {
		buf[4] |= 1 << 7
	}
	return serial.Tx(rtc.rw, buf[:], nil)
}

func (rtc *DS3231) Alarm2() (*DS3231Alarm, error) {
	var buf [3]byte
	if err := serial.Tx(rtc.rw, []byte{0x0b}, buf[:]); err != nil {
		return nil, err
	}
	return &DS3231Alarm{
		Minutes:       dsMinSec(buf[0]),
		Hours:         dsHour(buf[1]),
		IsPM:          buf[1]>>5&1 == 1,
		Use12Hour:     buf[1]>>6&1 == 1,
		Day:           dsDate(buf[2]),
		UseWeekday:    buf[2]>>6&1 == 1,
		IgnoreSeconds: true,
		IgnoreMinutes: buf[0]>>7&1 == 1,
		IgnoreHours:   buf[1]>>7&1 == 1,
		IgnoreDay:     buf[2]>>7&1 == 1,
	}, nil
}

func (rtc *DS3231) Alarm1() (*DS3231Alarm, error) {
	var buf [4]byte
	if err := serial.Tx(rtc.rw, []byte{0x07}, buf[:]); err != nil {
		return nil, err
	}
	return &DS3231Alarm{
		Seconds:       dsMinSec(buf[0]),
		Minutes:       dsMinSec(buf[1]),
		Hours:         dsHour(buf[2]),
		IsPM:          buf[2]>>5&1 == 1,
		Use12Hour:     buf[2]>>6&1 == 1,
		Day:           dsDate(buf[3]),
		UseWeekday:    buf[3]>>6&1 == 1,
		IgnoreSeconds: buf[0]>>7&1 == 1,
		IgnoreMinutes: buf[1]>>7&1 == 1,
		IgnoreHours:   buf[2]>>7&1 == 1,
		IgnoreDay:     buf[3]>>7&1 == 1,
	}, nil
}

func (rtc *DS3231) SetTime(t time.Time, use12Hour bool) error {
	s := time.Now()
	err := rtc._SetTime(t.UTC(), use12Hour)
	if err != nil {
		return err
	}
	dur := time.Since(s)
	return rtc._SetTime(t.UTC().Add(dur*2), use12Hour)
}

func (rtc *DS3231) _SetTime(t time.Time, use12Hour bool) error {
	var buf [8]byte
	buf[0] = 0x00
	buf[1] = dsEncMinSec(t.Second())
	buf[2] = dsEncMinSec(t.Minute())
	if use12Hour {
		buf[3] = dsEnc12HourFrom24(t.Hour())
	} else {
		buf[3] = dsEnc24Hour(t.Hour())
	}
	buf[4] = byte(t.Weekday()+1) & 0b111
	buf[5] = dsEncDate(t.Day())
	buf[6] = byte(t.Month()%10) | byte(t.Month()/10)<<4&1
	yr := t.Year()
	buf[7] = byte(yr%10) | (byte(yr/10%10) << 4)
	if yr >= 2100 {
		buf[6] |= 1 << 7
	}

	return serial.Tx(rtc.rw, buf[:], nil)
}

func dsEncDate(d int) byte {
	return byte(d%10) | (byte(d/10)&0b11)<<4
}

func dsDate(b byte) int {
	return int(b&0xf) + int(b>>4&0b11)*10
}

func dsEncMinSec(m int) byte {
	return byte(m%10) | (byte(m/10)&0b111)<<4
}

func dsMinSec(b byte) int {
	return int(b&0xf) + int(b>>4&0b111)*10
}

func dsEnc12HourFrom24(h int) (b byte) {
	switch {
	case h == 0:
		return dsEnc12Hour(12, false)
	case h == 12:
		return dsEnc12Hour(12, true)
	case h > 12:
		return dsEnc12Hour(h-12, true)
	}

	return dsEnc12Hour(h, false)
}

func dsEnc12Hour(h int, isPM bool) (b byte) {
	b = 1 << 6
	if isPM {
		b |= 1 << 5
	}
	return b | byte(h%10) | (byte(h/10)<<4)&1
}

func dsEnc24Hour(h int) (b byte) {
	b = byte(h % 10)
	if h > 10 {
		b |= 1 << 4
	}
	if h > 20 {
		b |= 1 << 5
	}
	return b
}

func dsHour(b byte) (hr int) {
	hr = int(b & 0xF)
	switch {
	case (b>>6&1) == 1 && b>>5&1 == 1 && hr < 12:
		hr += 12
	case (b>>6&1) == 0 && b>>5&1 == 1:
		hr += 20
	case (b>>6&1) == 0 && b>>4&1 == 1:
		hr += 10
	}
	return hr
}

func (rtc *DS3231) Now() (time.Time, error) {
	var buf [7]byte
	s := time.Now()
	if err := serial.Tx(rtc.rw, []byte{0}, buf[:]); err != nil {
		return time.Time{}, err
	}
	dur := time.Since(s)

	mon := int(buf[5]&0xF) + int(buf[5]>>4&1)*10
	year := int(buf[6]&0xF) + int(buf[6]>>4)*10 + 2000 + int(buf[5]>>7)*100

	return time.Date(
		year,
		time.Month(mon),
		dsDate(buf[4]),
		dsHour(buf[2]),
		dsMinSec(buf[1]),
		dsMinSec(buf[0]),
		0,
		time.UTC,
	).Add(dur), nil
}
