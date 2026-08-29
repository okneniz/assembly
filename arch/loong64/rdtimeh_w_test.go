package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRdtimehWCtor(t *testing.T) {
	// llvm-mc-verified: rdtimeh.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000065ac),
		ctorWord(t, NewRdtimehW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRdtimehW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RdtimehW)
	require.True(t, ok, "type = %T, want RdtimehW", in)
}

func TestRdtimehWDecodeEncode(t *testing.T) {
	in := decodeRdtimehW(0x000065ac, 0x90000000)

	rdtimehw, ok := in.(RdtimehW)
	require.True(t, ok, "type = %T, want RdtimehW", in)
	require.Equal(t, "rdtimeh.w $t0, $t1", rdtimehw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), rdtimehw.Addr())
	require.Equal(t, uint32(0x000065ac), ctorWord(t, rdtimehw))
}
