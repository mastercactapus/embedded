package at

import "strings"

// EscapeString escapes any backslash, newline or extra
// characters in s.
func EscapeString(s string, extra ...rune) string {
	s = strings.ReplaceAll(s, `\`, `\e`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	for _, r := range extra {
		s = strings.ReplaceAll(s, string(r), `\`+string(r))
	}
	return s
}

// UnescapeString is the inverse of EscapeString.
func UnescapeString(s string, extra ...rune) string {
	for _, r := range extra {
		s = strings.ReplaceAll(s, `\`+string(r), string(r))
	}
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\r`, "\r")
	s = strings.ReplaceAll(s, `\e`, `\`)
	return s
}
