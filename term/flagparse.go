package term

import (
	"strings"

	"github.com/mastercactapus/embedded/term/ascii"
)

type flagParse struct {
	boolFlags []string

	args  []string
	flags map[string]string
	err   error

	cur   flagValue
	state flagParseFunc
}

func (p *flagParse) isBool() bool {
	for _, b := range p.boolFlags {
		if b == p.cur.name {
			return true
		}
	}
	return false
}

type flagValue struct {
	name  string
	value string
}

func parseFlags(boolFlags, args []string) (map[string]string, []string, error) {
	p := flagParse{boolFlags: boolFlags, state: flagParseNext, flags: make(map[string]string)}
	for _, arg := range args {
		p.state = p.state(&p, arg)
		if p.err != nil {
			return nil, nil, p.err
		}
	}
	return p.flags, p.args, nil
}

type flagParseFunc func(*flagParse, string) flagParseFunc

func eqIndex(s string) int {
	for i, v := range s {
		if v == '=' {
			return i
		}
	}
	return -1
}

func (p *flagParse) setFlag() {
	if _, ok := p.flags[p.cur.name]; ok {
		if len(p.cur.name) == 1 {
			p.err = ascii.Errorf("flag -%s already set", p.cur.name)
		} else {
			p.err = ascii.Errorf("flag --%s already set", p.cur.name)
		}
		return
	}

	p.flags[p.cur.name] = p.cur.value
}

func flagParseNext(p *flagParse, arg string) flagParseFunc {
	switch {
	case arg == "--":
		return flagParseArgs
	case strings.HasPrefix(arg, "--"):
		eq := eqIndex(arg)
		if eq > -1 {
			p.cur.name = arg[2:eq]
			p.cur.value = arg[eq+1:]
			break
		}
		p.cur.name = arg[2:]
		if p.isBool() {
			p.cur.value = "true"
			p.setFlag()
			break
		}
		return flagParseValue
	case arg[0] == '-':
		return flagParseShort(p, arg[1:])
	}

	return flagParseNext
}

func flagParseShort(p *flagParse, arg string) flagParseFunc {
	p.cur.name = arg[:1]
	if len(arg) > 1 && arg[1] == '=' {
		p.cur.value = arg[2:]
		p.setFlag()
		return flagParseNext
	}
	if len(arg) == 1 && !p.isBool() {
		return flagParseValue
	}

	p.cur.value = "true"
	p.setFlag()
	if len(arg) > 1 {
		return flagParseShort(p, arg[1:])
	}

	return flagParseNext
}

func flagParseValue(p *flagParse, arg string) flagParseFunc {
	p.cur.value = arg
	p.setFlag()
	return flagParseNext
}

func flagParseArgs(p *flagParse, arg string) flagParseFunc {
	p.args = append(p.args, arg)
	return flagParseNext
}
