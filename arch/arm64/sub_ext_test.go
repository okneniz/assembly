package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubExtCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"sub x0,x1,x2,uxtx",
			ctorSubExt(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 0),
			0xcb226020,
		},
		{
			"sub sp,sp,x2,uxtx#3",
			ctorSubExt(t, SP, SP, xreg(t, 2), "uxtx", 3),
			0xcb226fff,
		},
		{
			"sub w1,w2,w3,uxtw#2",
			ctorSubExt(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), "uxtw", 2),
			0x4b234841,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSubExt(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 0)
	_, ok := in.(SubExt)
	require.True(t, ok, "type = %T, want SubExt", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"sub xzr,rd",
			func() error {
				_, err := New().SubExt(XZR, xreg(t, 1), xreg(t, 2), "uxtx", 0)
				return err
			},
		},
		{
			"sub xzr,rn",
			func() error {
				_, err := New().SubExt(xreg(t, 0), XZR, xreg(t, 2), "uxtx", 0)
				return err
			},
		},
		{
			"sub xzr,rm",
			func() error {
				_, err := New().SubExt(xreg(t, 0), xreg(t, 1), XZR, "uxtx", 0)
				return err
			},
		},
		{
			"sub x,w widths",
			func() error {
				_, err := New().SubExt(xreg(t, 0), wreg(t, 1), xreg(t, 2), "uxtx", 0)
				return err
			},
		},
		{
			"sub bad ext",
			func() error {
				_, err := New().SubExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "foo", 0)
				return err
			},
		},
		{
			"sub imm3 8",
			func() error {
				_, err := New().SubExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 8)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSubExt — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSubExt(t *testing.T, rd, rn, rm Reg, ext string, imm3 uint32) Instr {
	t.Helper()
	in, err := New().SubExt(rd, rn, rm, ext, imm3)
	require.NoError(t, err)
	return in
}
