package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddDCtor(t *testing.T) {
	// llvm-mc-verified: amadd.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3861b5cc),
		ctorWord(t, NewAmaddD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmaddD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddD)
	require.True(t, ok, "type = %T, want AmaddD", in)
}

func TestAmaddDDecodeEncode(t *testing.T) {
	in := decodeAmaddD(0x3861b5cc, 0x90000000)

	amaddd, ok := in.(AmaddD)
	require.True(t, ok, "type = %T, want AmaddD", in)
	require.Equal(t, "amadd.d $t0, $t1, $t2", amaddd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amaddd.Addr())
	require.Equal(t, 4, amaddd.Len())
	require.Equal(t, uint32(0x3861b5cc), ctorWord(t, amaddd))
}
