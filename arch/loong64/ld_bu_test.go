package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdBuCtor(t *testing.T) {
	// llvm-mc-verified: ld.bu $t0, $t1, 8.
	require.Equal(
		t,
		uint32(0x2a0021ac),
		ctorWord(t, NewLdBu(lreg(t, 12), lreg(t, 13), imm12v(t, 8))),
	)
}

func TestLdBuDecodeEncode(t *testing.T) {
	in := decodeLdBu(0x2a0021ac, 0x90000000)

	x, ok := in.(LdBu)
	require.True(t, ok, "type = %T, want LdBu", in)
	require.Equal(t, "ld.bu $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x2a0021ac), ctorWord(t, x))
}
