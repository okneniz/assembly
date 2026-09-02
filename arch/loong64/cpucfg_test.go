package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCpucfgCtor(t *testing.T) {
	// llvm-mc-verified: cpucfg $t0, $t1.
	require.Equal(
		t,
		uint32(0x00006dac),
		ctorWord(t, New().Cpucfg(lreg(t, 12), lreg(t, 13))),
	)

	in := New().Cpucfg(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Cpucfg)
	require.True(t, ok, "type = %T, want Cpucfg", in)
}

func TestCpucfgDecodeEncode(t *testing.T) {
	in := decodeCpucfg(0x00006dac, 0x90000000)

	cpucfg, ok := in.(Cpucfg)
	require.True(t, ok, "type = %T, want Cpucfg", in)
	require.Equal(t, "cpucfg $t0, $t1", cpucfg.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), cpucfg.Addr())
	require.Equal(t, uint32(0x00006dac), ctorWord(t, cpucfg))
}
