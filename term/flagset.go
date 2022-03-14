package term

import (
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
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

type (
	Flags struct {
		cmd *CmdLine
		env func(string) (string, bool)

		flags    map[string]*flagInfo
		flagList []string

		helpParams string
		examples   []flagExample
		showHelp   *bool

		parsedArgs []string
	}
)

type flagInfo struct {
	Flag
	Bool  bool
	Type  string
	Parse func(string) error

	wasSet bool
}

type flagExample struct {
	Cmdline string
	Desc    string
}

func NewFlagSet(cmd *CmdLine, env func(string) (string, bool)) *Flags {
	fs := &Flags{
		cmd:        cmd,
		env:        env,
		flags:      make(map[string]*flagInfo),
		helpParams: "[parameters ...]",
	}
	fs.showHelp = fs.Bool(Flag{Name: "help", Short: 'h', Desc: "Show this help message"})
	return fs
}

// SetHelpParameters sets the parameters for the usage output.
func (fs *Flags) SetHelpParameters(s string) { fs.helpParams = s }

func (fs *Flags) Example(cmdline, desc string) {
	fs.examples = append(fs.examples, flagExample{cmdline, desc})
}

func (fs *Flags) Parse() error {
	var boolFlags []string
	for _, f := range fs.flags {
		if !f.Bool {
			continue
		}
		if f.Name != "" {
			boolFlags = append(boolFlags, f.Name)
		}
		if f.Short != 0 {
			boolFlags = append(boolFlags, string(f.Short))
		}
	}

	vals, args, err := parseFlags(boolFlags, fs.cmd.Args)
	if err != nil {
		return usageErr{err: err}
	}
	fs.parsedArgs = args

	for name, val := range vals {
		f, ok := fs.flags[name]
		if !ok {
			return fs.UsageError("unknown flag %s", name)
		}

		err := f.Parse(val)
		if err != nil {
			return fs.UsageError("flag %s: %w", name, err)
		}
		f.wasSet = true
	}

	for _, f := range fs.flags {
		if f.wasSet {
			continue
		}

		env, ok := fs.env(f.Name)
		if ok {
			err := f.Parse(env)
			if err != nil {
				return fs.UsageError("flag %s: %w", f.Name, err)
			}
			continue
		}

		if f.Def != "" {
			err := f.Parse(f.Def)
			if err != nil {
				return fs.UsageError("flag %s: %w", f.Name, err)
			}
			continue
		}

		if !f.Req {
			continue
		}

		return fs.UsageError("flag %s is required", f.Name)
	}

	if *fs.showHelp {
		return usageErr{fs: fs}
	}

	return nil
}

func (fs *Flags) UsageError(format string, a ...interface{}) error {
	return usageErr{fs: fs, err: ansi.Errorf(format, a...)}
}

func (fs *Flags) Args() []string {
	return fs.parsedArgs
}

func (fs *Flags) Arg(n int) string {
	if n < len(fs.parsedArgs) {
		return fs.parsedArgs[n]
	}

	return ""
}

func (fs *Flags) flag(info *flagInfo) {
	if _, ok := fs.flags[info.Name]; ok {
		panic("flag --" + info.Name + " already defined")
	}
	if _, ok := fs.flags[string(info.Short)]; ok {
		panic("flag -" + string(info.Short) + " already defined")
	}

	if info.Name != "" {
		fs.flags[info.Name] = info
		fs.flagList = append(fs.flagList, info.Name)
	}
	if info.Short != 0 {
		fs.flags[string(info.Short)] = info
		if info.Name == "" {
			fs.flagList = append(fs.flagList, string(info.Short))
		}
	}
}

func (fs *Flags) Enum(f Flag, vals ...string) *string {
	var res string

	fs.flag(&flagInfo{Flag: f, Type: strings.Join(vals, "|"), Parse: func(value string) error {
		res = value
		for _, v := range vals {
			if v == value {
				return nil
			}
		}

		return fs.UsageError("flag '%s' must be one of %s", f.Name, strings.Join(vals, ", "))
	}})

	return &res
}

func (fs *Flags) Bytes(f Flag) *[]byte {
	var res []byte

	fs.flag(&flagInfo{Flag: f, Type: "byte ...", Parse: func(value string) error {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			v, err := parseUint8(part)
			if err != nil {
				return err
			}
			res = append(res, v)
		}

		return nil
	}})

	return &res
}

func (fs *Flags) String(f Flag) *string {
	var res string

	fs.flag(&flagInfo{Flag: f, Type: "string", Parse: func(value string) error {
		res = value
		return nil
	}})

	return &res
}

func (fs *Flags) Bool(f Flag) *bool {
	var res bool

	fs.flag(&flagInfo{Flag: f, Bool: true, Parse: func(value string) (err error) {
		res, err = parseBool(value)
		return err
	}})

	return &res
}

func (fs *Flags) Int(f Flag) *int {
	var res int

	fs.flag(&flagInfo{Flag: f, Type: "int", Parse: func(value string) (err error) {
		res, err = ParseInt(value)
		return err
	}})

	return &res
}

func (fs *Flags) Uint16(f Flag) *uint16 {
	var res uint16

	fs.flag(&flagInfo{Flag: f, Type: "uint16", Parse: func(value string) (err error) {
		res, err = parseUint16(value)
		return err
	}})

	return &res
}

func (fs *Flags) Byte(f Flag) *byte {
	var res byte

	fs.flag(&flagInfo{Flag: f, Type: "byte", Parse: func(value string) (err error) {
		res, err = parseUint8(value)
		return err
	}})

	return &res
}
