package at

import "strings"

type Response struct {
	Data []string

	OK bool
}

// SetValue sets the value of a response parameter.
func (r *Response) SetValue(name, value string) {
	name = strings.ToUpper(name)
	r.Data = append(r.Data, EscapeString(name, '=')+"="+EscapeString(value))
}

// Value returns the value of a response parameter.
//
// If multiple values are present, only the first one is returned.
func (r Response) Value(name string) string {
	name = strings.ToUpper(name)
	prefix := EscapeString(name, '=') + "="
	for _, line := range r.Data {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		return UnescapeString(strings.TrimPrefix(line, prefix))
	}
	return ""
}

// Values is like Value, but returns all values.
func (r Response) Values(name string) []string {
	name = strings.ToUpper(name)
	prefix := EscapeString(name, '=') + "="
	values := make([]string, 0)
	for _, line := range r.Data {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		values = append(values, UnescapeString(strings.TrimPrefix(line, prefix)))
	}
	return values
}
