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
	sh *Shell
	*CmdLine
	desc string
	fs   *FlagSet
	env  *Env
}

// Printer will return the printer associated with the current context.
func Printer(ctx context.Context) *ansi.Printer {
	cmd, ok := ctx.Value(ctxKeyCmd).(*cmdContext)
	if !ok {
		return nil
	}

	return cmd.sh.p
}

// getEnv will return the value of the environment variable with the given name.
func getEnv(ctx context.Context, name string) string {
	val, _ := ctx.Value(envKey(name)).(string)
	return val
}

// withEnv will return a new context with the given environment variable.
func withEnv(ctx context.Context, name, val string) context.Context {
	return context.WithValue(ctx, envKey(name), val)
}

func Flags(ctx context.Context) *FlagSet {
	cmd, ok := ctx.Value(ctxKeyCmd).(*cmdContext)
	if !ok {
		panic("Parse called but no command context")
	}

	return cmd.fs
}
