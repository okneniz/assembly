package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestOrnCtor(t *testing.T) {
	// llvm-mc-verified: orn $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001639ac),
		ctorWord(t, NewOrn(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewOrn(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Orn)
	require.True(t, ok, "type = %T, want Orn", in)
}

func TestOrnDecodeEncode(t *testing.T) {
	in := decodeOne(0x001639ac, 0x90000000)

	x, ok := in.(Orn)
	require.True(t, ok, "type = %T, want Orn", in)
	require.Equal(t, "orn $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001639ac), ctorWord(t, x))
}
