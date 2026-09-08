package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCbnzBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cbnz x2,0x1008",
			buildCbnz(t, xreg(t, 2), 0x1008),
			0xb5000042,
		},
		{
			"cbnz w0,0x1014",
			buildCbnz(t, wreg(t, 0), 0x1014),
			0x350000a0,
		},
		{
			"cbnz wzr,0x1000",
			buildCbnz(t, WZR, 0x1000),
			0x3500001f,
		},
		{
			"cbnz x0,0x100ffc",
			buildCbnz(t, xreg(t, 0), 0x100ffc),
			0xb57fffe0,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildCbnz(t, xreg(t, 2), 0x1008)
	_, ok := in.(Cbnz)
	require.True(t, ok, "type = %T, want Cbnz", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"cbnz sp,rt",
			func() error {
				_, err := New().Cbnz(SP, 0x1000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildCbnz — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildCbnz(t *testing.T, rt Reg, target int64) Instr {
	t.Helper()
	in, err := New().Cbnz(rt, target)
	require.NoError(t, err)
	return in
}
