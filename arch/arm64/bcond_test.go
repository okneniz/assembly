package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBcondCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"b.eq 0x1000",
			ctorBcond(t, "eq", 0x1000),
			0x54000000,
		},
		{
			"b.ne 0x1010",
			ctorBcond(t, "ne", 0x1010),
			0x54000081,
		},
		{
			"b.hs 0xffc",
			ctorBcond(t, "hs", 0xffc),
			0x54ffffe2,
		},
		{
			"b.eq 0x100ffc",
			ctorBcond(t, "eq", 0x100ffc),
			0x547fffe0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorBcond(t, "eq", 0x1000)
	_, ok := in.(Bcond)
	require.True(t, ok, "type = %T, want Bcond", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"b.cond bad cond",
			func() error {
				_, err := New().Bcond("foo", 0x1000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBcond — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBcond(t *testing.T, cond string, target int64) Instr {
	t.Helper()
	in, err := New().Bcond(cond, target)
	require.NoError(t, err)
	return in
}
