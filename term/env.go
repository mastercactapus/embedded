package term

import "sort"

// Env manages environment variables.
type Env struct {
	parent *Env
	values map[string]envVal
}

type envVal struct {
	Value string
	Set   bool
}

// NewEnv creates a new environment.
func NewEnv() *Env {
	return &Env{values: make(map[string]envVal)}
}

// SetParent will set the parent environment.
func (e *Env) SetParent(parent *Env) { e.parent = parent }

// Get will return the value of the given key.
func (e *Env) Get(key string) (string, bool) {
	if v, ok := e.values[key]; ok {
		return v.Value, v.Set
	}

	if e.parent != nil {
		return e.parent.Get(key)
	}

	return "", false
}

// Set will set the value of the given key.
func (e *Env) Set(key, value string) { e.values[key] = envVal{Set: true, Value: value} }

// Unset will remove the given key from the environment.
func (e *Env) Unset(key string) { e.values[key] = envVal{Set: false} }

// List will return a list of all keys for the current environment.
func (e *Env) List() []string {
	var keys []string
	for k := range e.values {
		keys = append(keys, k)
	}

	if e.parent != nil {
		keys = append(keys, e.parent.List()...)
	}

	// sort and remove duplicates
	sort.Strings(keys)
	var last string
	uniq := keys[:0]
	for _, k := range keys {
		if k == last {
			continue
		}
		last = k
		if v, ok := e.values[k]; ok && !v.Set {
			continue
		}
		uniq = append(uniq, k)
	}

	return keys
}
