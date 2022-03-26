package bustool

import (
	"time"

	"github.com/mastercactapus/embedded/driver/rtc"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddRTC(sh *term.Shell) *term.Shell {
	rtcSh := sh.NewSubShell("rtc", "Interact with a DS3231 RTC device over I2C.", func(ra term.RunArgs) error {
		addr := ra.Uint16(term.Flag{Name: "dev", Short: 'd', Def: "0x68", Env: "DEV", Desc: "Device addresss.", Req: true})
		if err := ra.Parse(); err != nil {
			return err
		}

		bus := ra.Get("i2c").(i2c.Bus)

		ra.Set("rtc", rtc.NewDS3231(i2c.NewDevice(bus, *addr)))

		return nil
	})

	rtcSh.AddCommands(RTCCommands...)
	return rtcSh
}

var RTCCommands = []term.Command{
	{Name: "date", Desc: "Read current date.", Exec: func(ra term.RunArgs) error {
		if err := ra.Parse(); err != nil {
			return err
		}

		rtc := ra.Get("rtc").(*rtc.DS3231)

		n, err := rtc.Now()
		if err != nil {
			return err
		}

		ra.Println(n.Format(time.RFC3339))

		return nil
	}},
	{Name: "set", Desc: "Set current date.", Exec: func(ra term.RunArgs) error {
		use12 := ra.Bool(term.Flag{Name: "12", Desc: "Use 12-hour clock."})
		timeStr := ra.String(term.Flag{Name: "time", Short: 't', Desc: "Time to set (instead of system clock) in ISO format."})
		if err := ra.Parse(); err != nil {
			return err
		}

		rtc := ra.Get("rtc").(*rtc.DS3231)
		if *timeStr == "" {
			return rtc.SetTime(time.Now(), *use12)
		}

		t, err := time.Parse(time.RFC3339, *timeStr)
		if err != nil {
			return err
		}
		return rtc.SetTime(t, *use12)
	}},
}
