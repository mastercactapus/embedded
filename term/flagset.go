package term

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
)

type (
	Arg struct {
		Name string
		Desc string // Description of the argument
		Req  bool   // Mark as required
	}

	Flag struct {
		Name  string
		Env   string // Env var name
		Def   string // Default value
		Desc  string // Description of the flag
		Req   bool   // Mark as required
		Short rune
	}
)

type FlagSet struct {
	cmd *CmdLine
	set *getopt.Set
	env func(string) (string, bool)

	args []argInfo

	examples []flagExample
	showHelp bool
}

type argInfo struct {
	Arg
	v interface{}
}

type flagExample struct {
	Cmdline string
	Desc    string
}
type usageErr struct {
	fs  *FlagSet
	err error
}

func (e usageErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "usage requested"
}

func (e usageErr) PrintUsage(w io.Writer) {
	e.fs.set.PrintUsage(w)
}

func NewFlagSet(cmd *CmdLine, env func(string) (string, bool)) *FlagSet {
	fs := &FlagSet{cmd: cmd, env: env, set: getopt.New()}

	fs.flag(Flag{Name: "help", Short: 'h', Desc: "Show this help message"}, flagVal{&fs.showHelp}).SetFlag()
	return fs
}

// SetHelpParameters sets the parameters for the usage output.
func (fs *FlagSet) SetHelpParameters(s string) {
	fs.set.SetParameters(s)
}

func (fs *FlagSet) Example(cmdline, desc string) {
	fs.examples = append(fs.examples, flagExample{cmdline, desc})
}

func (fs *FlagSet) Parse() error {
	err := fs.set.Getopt(fs.cmd.Args, nil)
	if err != nil {
		return usageErr{fs: fs, err: err}
	}
	if fs.showHelp {
		return usageErr{fs: fs}
	}

	// TODO: check for required args
	for i, arg := range fs.args {
		switch v := arg.v.(type) {
		case *[]string:
			*v = fs.set.Args()
		case *[]byte:
			*v = []byte(fs.set.Arg(i))
		}
	}

	return nil
}

func (fs *FlagSet) UsageError(format string, a ...interface{}) error {
	return usageErr{fs: fs, err: fmt.Errorf(format, a...)}
}

func (fs *FlagSet) flag(f Flag, v getopt.Value) getopt.Option {
	if f.Env != "" {
		envVal, _ := fs.env(f.Env)
		if envVal != "" {
			f.Def = envVal
		}
	}

	if f.Def != "" {
		v.Set(f.Def, nil)
	}
	opt := fs.set.FlagLong(v, f.Name, f.Short, f.Desc)

	if f.Req && f.Def == "" {
		opt.Mandatory()
	}
	return opt
}

func (fs *FlagSet) Args() []string {
	return fs.set.Args()
}

func (fs *FlagSet) Bytes(f Flag) *[]byte {
	var v []byte
	fs.flag(f, flagVal{&v})
	return &v
}

func (fs *FlagSet) String(f Flag) *string {
	var v string
	fs.flag(f, flagVal{&v})
	return &v
}

func (fs *FlagSet) Bool(f Flag) *bool {
	var v bool
	fs.flag(f, flagVal{&v}).SetFlag()
	return &v
}

func (fs *FlagSet) Int(f Flag) *int {
	var v int
	fs.flag(f, flagVal{&v})
	return &v
}

func (fs *FlagSet) Byte(f Flag) *byte {
	var v byte
	fs.flag(f, flagVal{&v})
	return &v
}

type flagVal struct {
	value interface{}
}

func (f flagVal) String() string {
	switch v := f.value.(type) {
	case *string:
		return *v
	case *bool:
		return strconv.FormatBool(*v)
	case *int:
		return strconv.Itoa(*v)
	case *byte:
		return fmt.Sprintf("0x%02x", *v)
	case *[]byte:
		var s []string
		for _, b := range *v {
			s = append(s, fmt.Sprintf("0x%02x", b))
		}
		return strings.Join(s, ",")
	default:
		panic(fmt.Sprintf("unsupported flag type %T", v))
	}
}

func (f flagVal) Set(s string, opt getopt.Option) error {
	switch v := f.value.(type) {
	case *[]byte:
		bits := strings.Split(s, ",")
		bytes := make([]byte, 0, len(bits))
		for _, s := range bits {
			if len(s) == 0 {
				continue
			}
			if s[0] == 's' {
				bytes = append(bytes, []byte(s[1:])...)
				continue
			}
			b, err := strconv.ParseUint(s, 0, 8)
			if err != nil {
				return err
			}
			bytes = append(bytes, byte(b))
		}
		*v = bytes
	case *string:
		*v = s
	case *bool:
		if s == "" && opt != nil && opt.IsFlag() {
			*v = true
			break
		}
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*v = b
	case *int:
		p, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return err
		}
		*v = int(p)
	case *byte:
		p, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return err
		}
		*v = byte(p)
	default:
		return fmt.Errorf("unsupported flag type %T", v)
	}
	return nil
}
