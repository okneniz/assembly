package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestExtWHCtor(t *testing.T) {
	// llvm-mc-verified: ext.w.h $t0, $t1.
	require.Equal(
		t,
		uint32(0x000059ac),
		ctorWord(t, New().ExtWH(lreg(t, 12), lreg(t, 13))),
	)

	in := New().ExtWH(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ExtWH)
	require.True(t, ok, "type = %T, want ExtWH", in)
}

func TestExtWHDecodeEncode(t *testing.T) {
	in := decodeExtWH(0x000059ac, 0x90000000)

	extwh, ok := in.(ExtWH)
	require.True(t, ok, "type = %T, want ExtWH", in)
	require.Equal(t, "ext.w.h $t0, $t1", extwh.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), extwh.Addr())
	require.Equal(t, uint32(0x000059ac), ctorWord(t, extwh))
}
