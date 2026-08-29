package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestScrelDCtor(t *testing.T) {
	// llvm-mc-verified: screl.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x38578dac),
		ctorWord(t, NewScrelD(lreg(t, 12), lreg(t, 13))),
	)

	in := NewScrelD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ScrelD)
	require.True(t, ok, "type = %T, want ScrelD", in)
}

func TestScrelDDecodeEncode(t *testing.T) {
	in := decodeScrelD(0x38578dac, 0x90000000)

	screld, ok := in.(ScrelD)
	require.True(t, ok, "type = %T, want ScrelD", in)
	require.Equal(t, "screl.d $t0, $t1", screld.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), screld.Addr())
	require.Equal(t, 4, screld.Len())
	require.Equal(t, uint32(0x38578dac), ctorWord(t, screld))
}
