package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPrfmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"prfm pldl1keep,[x1]",
			ctorPrfm(t, xreg(t, 1)),
			0xf9800020,
		},
		{
			"prfm pldl1keep,[sp]",
			ctorPrfm(t, SP),
			0xf98003e0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorPrfm(t, xreg(t, 1))
	_, ok := in.(Prfm)
	require.True(t, ok, "type = %T, want Prfm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"prfm w1 base",
			func() error {
				_, err := New().Prfm(wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorPrfm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorPrfm(t *testing.T, rn Reg) Instr {
	t.Helper()
	in, err := New().Prfm(rn)
	require.NoError(t, err)
	return in
}
