package at

import "strings"

type Cmd struct {
	// FullName is always upper case.
	//
	// Includes the AT prefix, if any, as well as
	// the "?" for query commands and the "=" for
	// set commands.
	//
	// Parameters are only valid for set commands.
	FullName string

	// Command parameter(s).
	Params []string
}

// IsSet returns true if the command is a set command.
func (c Cmd) IsSet() bool {
	return strings.ContainsRune(c.FullName, '=')
}

// IsQuery returns true if the command is a query command.
func (c Cmd) IsQuery() bool {
	return !c.IsSet() && strings.HasSuffix(c.FullName, "?")
}

// Name returns the command name, without the AT prefix,
// but with the "?" for query commands and the "=" for
// set commands.
func (c Cmd) Name() string {
	n := strings.TrimPrefix(c.FullName, "AT")
	n = strings.TrimPrefix(n, "+")
	if c.IsSet() {
		n, _, _ = strings.Cut(n, "=")
	} else {
		n = strings.TrimSuffix(n, "?")
	}

	return UnescapeString(n, '=', '?')
}
