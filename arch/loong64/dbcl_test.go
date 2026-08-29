package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDbclCtor(t *testing.T) {
	// llvm-mc-verified: dbcl 1.
	v, err := NewCode15(1)
	require.NoError(t, err)
	require.Equal(t, uint32(0x002a8001), ctorWord(t, NewDbcl(v)))
}

func TestDbclDecodeEncode(t *testing.T) {
	in := decodeDbcl(0x002a8001, 0x90000000)
	require.Equal(t, "dbcl 1", in.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x002a8001), ctorWord(t, in))
}
