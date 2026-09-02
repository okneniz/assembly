package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdHuCtor(t *testing.T) {
	// llvm-mc-verified: ld.hu $t0, $t1, 8.
	require.Equal(
		t,
		uint32(0x2a4021ac),
		ctorWord(t, New().LdHu(lreg(t, 12), lreg(t, 13), imm12v(t, 8))),
	)
}

func TestLdHuDecodeEncode(t *testing.T) {
	in := decodeLdHu(0x2a4021ac, 0x90000000)

	x, ok := in.(LdHu)
	require.True(t, ok, "type = %T, want LdHu", in)
	require.Equal(t, "ld.hu $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x2a4021ac), ctorWord(t, x))
}
