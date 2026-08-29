package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestExtWBCtor(t *testing.T) {
	// llvm-mc-verified: ext.w.b $t0, $t1.
	require.Equal(
		t,
		uint32(0x00005dac),
		ctorWord(t, NewExtWB(lreg(t, 12), lreg(t, 13))),
	)

	in := NewExtWB(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ExtWB)
	require.True(t, ok, "type = %T, want ExtWB", in)
}

func TestExtWBDecodeEncode(t *testing.T) {
	in := decodeExtWB(0x00005dac, 0x90000000)

	extwb, ok := in.(ExtWB)
	require.True(t, ok, "type = %T, want ExtWB", in)
	require.Equal(t, "ext.w.b $t0, $t1", extwb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), extwb.Addr())
	require.Equal(t, uint32(0x00005dac), ctorWord(t, extwb))
}
