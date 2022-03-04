package term_test

import (
	"testing"

	"github.com/mastercactapus/embedded/term"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCmdLine(t *testing.T) {
	cmd, err := term.ParseCmdLine("cmd")
	require.NoError(t, err)
	assert.Equal(t, "cmd", cmd.Args[0])
	assert.Len(t, cmd.Args, 1)
	assert.Len(t, cmd.Env, 0)

	cmd, err = term.ParseCmdLine("foo=bar `` \"\" test")
	require.NoError(t, err)
	assert.Len(t, cmd.Args, 1)
	assert.Len(t, cmd.Env, 1)
	assert.Equal(t, "test", cmd.Args[0])
	assert.Equal(t, "foo=bar", cmd.Env[0])

	cmd, err = term.ParseCmdLine(`foo "test\x77" `)
	require.NoError(t, err)
	assert.Len(t, cmd.Args, 2)
	assert.Equal(t, "foo", cmd.Args[0])
	assert.Equal(t, "test\x77", cmd.Args[1])
}
