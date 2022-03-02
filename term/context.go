package term

import (
	"context"

	"github.com/mastercactapus/embedded/term/ansi"
)

type (
	ctxKey int
	envKey string
)

const (
	ctxKeyCmd ctxKey = iota
)

type cmdContext struct {
	sh   *Shell
	env  *CommandEnv
	desc string
	fp   *FlagParser
}

// Printer will return the printer associated with the current context.
func Printer(ctx context.Context) *ansi.Printer {
	cmd, ok := ctx.Value(ctxKeyCmd).(*cmdContext)
	if !ok {
		return nil
	}

	return cmd.sh.p
}

// env will return the value of the environment variable with the given name.
func env(ctx context.Context, name string) string {
	val, _ := ctx.Value(envKey(name)).(string)
	return val
}

// withEnv will return a new context with the given environment variable.
func withEnv(ctx context.Context, name, val string) context.Context {
	return context.WithValue(ctx, envKey(name), val)
}

func Parse(ctx context.Context) *FlagParser {
	cmd, ok := ctx.Value(ctxKeyCmd).(*cmdContext)
	if !ok {
		return nil
	}
	if cmd.fp == nil {
		cmd.fp = NewFlagParser(cmd.env, func(name string) string {
			return env(ctx, name)
		})
	}

	return cmd.fp
}
