package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevb2WCtor(t *testing.T) {
	// llvm-mc-verified: revb.2w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000039ac),
		ctorWord(t, NewRevb2W(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevb2W(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revb2W)
	require.True(t, ok, "type = %T, want Revb2W", in)
}

func TestRevb2WDecodeEncode(t *testing.T) {
	in := decodeRevb2W(0x000039ac, 0x90000000)

	revb2w, ok := in.(Revb2W)
	require.True(t, ok, "type = %T, want Revb2W", in)
	require.Equal(t, "revb.2w $t0, $t1", revb2w.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revb2w.Addr())
	require.Equal(t, uint32(0x000039ac), ctorWord(t, revb2w))
}
