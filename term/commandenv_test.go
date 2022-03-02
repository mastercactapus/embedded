package term_test

import (
	"testing"

	"github.com/mastercactapus/embedded/term"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommandEnv(t *testing.T) {
	env, err := term.ParseCommandEnv("A=B foo -bar -baz=BAR -- foo -okay M=F")
	require.NoError(t, err)

	assert.Equal(t, "foo", env.Name)

	assert.Contains(t, env.Flags, "-bar")
	assert.Contains(t, env.Flags, "-baz=BAR")
	assert.Len(t, env.Flags, 2)

	assert.Contains(t, env.LocalEnv, "A=B")
	assert.Len(t, env.LocalEnv, 1)

	assert.Contains(t, env.Args, "foo")
	assert.Contains(t, env.Args, "-okay")
	assert.Contains(t, env.Args, "M=F")
	assert.Len(t, env.Args, 3)

	env, err = term.ParseCommandEnv("foo -h")
	require.NoError(t, err)

	assert.Equal(t, "foo", env.Name)
	assert.Contains(t, env.Flags, "-h")
	assert.Len(t, env.Flags, 1)
	assert.Len(t, env.Args, 0)

	env, err = term.ParseCommandEnv("export FOO=bar")
	require.NoError(t, err)

	assert.Equal(t, "export", env.Name)
	assert.Contains(t, env.Args, "FOO=bar")
	assert.Len(t, env.Args, 1)

	env, err = term.ParseCommandEnv("tx -d=0 -f")
	require.NoError(t, err)

	assert.Equal(t, "tx", env.Name)
	assert.Contains(t, env.Flags, "-d=0")
	assert.Contains(t, env.Flags, "-f")
	assert.Len(t, env.Flags, 2)
	assert.Len(t, env.Args, 0)
}
