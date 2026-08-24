package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubImmCtor(t *testing.T) {
	got := ctorWord(
		t,
		ctorSubImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 0x42), NoSh12),
	)
	require.Equal(t, uint32(0xd1010820), got, "sub")
	got = ctorWord(t, ctorSubImm(t, wreg(t, 0), WSP, imm12(t, 1), NoSh12))
	require.Equal(t, uint32(0x510007e0), got, "sub w0,wsp")
	in := ctorSubImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 1), NoSh12)
	_, ok := in.(SubImm)
	require.True(t, ok, "type = %T, want SubImm", in)
	_, err := NewSubImm(XZR, xreg(t, 1), imm12(t, 1), NoSh12)
	assertErr(t, "sub imm zr", err)
}
