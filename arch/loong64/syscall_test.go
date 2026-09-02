package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSyscallCtor(t *testing.T) {
	c0, err := New().Code15(0)
	require.NoError(t, err)

	// llvm-mc-verified: syscall 0.
	require.Equal(t, uint32(0x002b0000), ctorWord(t, New().Syscall(c0)))
}

func TestSyscallDecodeEncode(t *testing.T) {
	x, ok := decodeSyscall(0x002b0000, 0x90000000).(Syscall)
	require.True(t, ok, "type = %T, want Syscall", x)
	require.Equal(t, "syscall 0", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(0), x.code.val)
	require.Equal(t, uint32(0x002b0000), ctorWord(t, x))
}
