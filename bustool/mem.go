package bustool

import (
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"io"
	"math/rand"

	"github.com/mastercactapus/embedded/driver/mem"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
)

type memDevice interface {
	io.ReadWriteSeeker
	io.ReaderAt
	io.WriterAt
}

func AddMem(sh *term.Shell) *term.Shell {
	memSh := sh.NewSubShell("mem", "Interact with an AT24Cxx-compatible EEPROM device over I2C.", func(r term.RunArgs) error {
		addr := r.Uint16(term.Flag{Name: "dev", Short: 'd', Def: "0x50", Env: "DEV", Desc: "Device addresss.", Req: true})
		size := r.Int(term.Flag{Name: "size", Short: 'm', Def: "1", Desc: "Memory size in kbits.", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		var dev memDevice
		bus := r.Get("i2c").(i2c.Bus)

		switch *size {
		case 1:
			dev = mem.NewAT24C01A(bus, *addr)
		case 2:
			dev = mem.NewAT24C02(bus, *addr)
		case 4:
			dev = mem.NewAT24C04(bus, *addr).(memDevice)
		case 8:
			dev = mem.NewAT24C08(bus, *addr).(memDevice)
		case 16:
			dev = mem.NewAT24C16(bus, *addr).(memDevice)
		case 32:
			dev = mem.NewAT24C32(bus, *addr)
		case 64:
			dev = mem.NewAT24C64(bus, *addr)
		default:
			return r.UsageError("unsupported memory size %d", *size)
		}
		r.Set("mem", dev)

		return nil
	})

	memSh.AddCommands(MemCommands...)
	return memSh
}

func size(s io.Seeker) int {
	s.Seek(0, io.SeekEnd)
	size, _ := s.Seek(0, io.SeekCurrent)
	return int(size)
}

var MemCommands = []term.Command{
	{Name: "r", Desc: "Read device data.", Exec: func(r term.RunArgs) error {
		start := r.Int(term.Flag{Short: 's', Def: "0", Desc: "Position to start from.", Req: true})
		count := r.Int(term.Flag{Short: 'n', Def: "0", Desc: "Number of bytes to read, if zero read to end."})
		if err := r.Parse(); err != nil {
			return err
		}

		mem := r.Get("mem").(memDevice)

		_, err := mem.Seek(int64(*start), 0)
		if err != nil {
			return err
		}

		wc := hex.Dumper(r)
		if *count == 0 {
			io.Copy(wc, mem)
		} else {
			io.CopyN(wc, mem, int64(*count))
		}

		return wc.Close()
	}},
	{Name: "w", Desc: "Write device data.", Exec: func(r term.RunArgs) error {
		r.SetHelpParameters("[data]")
		start := r.Int(term.Flag{Short: 's', Def: "0", Desc: "Position to start from.", Req: true})
		binData := r.Bytes(term.Flag{Name: "data", Short: 'b', Desc: "Write bytes (comma separated) before arg data."})
		if err := r.Parse(); err != nil {
			return err
		}

		mem := r.Get("mem").(memDevice)

		data := append(*binData, []byte(r.Arg(0))...)
		_, err := mem.WriteAt(data, int64(*start))
		if err != nil {
			return err
		}

		return nil
	}},
	{Name: "crc", Desc: "Calculate CRC value of bytes.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		crc := crc32.ChecksumIEEE([]byte(r.Arg(0)))
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, crc)
		r.Printf("CRC: 0x%02x 0x%02x 0x%02x 0x%02x\n", data[0], data[1], data[2], data[3])
		return nil
	}},
	{Name: "format", Desc: "Clear all data.", Exec: func(r term.RunArgs) error {
		start := r.Int(term.Flag{Short: 'p', Def: "0", Desc: "Position to start from.", Req: true})
		count := r.Int(term.Flag{Short: 'n', Def: "0", Desc: "Number of bytes to wipe, if zero clear to end."})
		value := r.Byte(term.Flag{Short: 'v', Def: "0xff", Desc: "Value to write."})
		// TODO: group
		rnd := r.Bool(term.Flag{Name: "random", Desc: "Fill with random data."})
		seq := r.Bool(term.Flag{Name: "sequential", Desc: "Fill with sequential data."})
		if err := r.Parse(); err != nil {
			return err
		}

		mem := r.Get("mem").(memDevice)

		if *count == 0 {
			*count = size(mem) - *start
		}
		if *count <= 0 {
			return nil
		}

		data := make([]byte, *count)
		for i := range data {
			switch {
			case *seq:
				data[i] = byte(i)
			case *rnd:
				data[i] = byte(rand.Intn(256))
			default:
				data[i] = *value
			}
		}
		_, err := mem.WriteAt(data, int64(*start))
		if err != nil {
			return err
		}

		return nil
	}},
}
