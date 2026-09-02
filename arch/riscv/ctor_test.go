package riscv

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// ctorBytes - the bytes encoded by the constructor (2 or 4: RVC compression).
func ctorBytes(t *testing.T, in Instr) []byte {
	t.Helper()
	var buf bytes.Buffer
	_, err := in.Encode(&buf, 0x1000, EncOpts{})
	require.NoError(t, err, "Encode %s", in.ObjDump(disasm.DefaultViewCtx()))
	require.Contains(
		t,
		[]int{2, 4},
		buf.Len(),
		"Encode %s: %d bytes",
		in.ObjDump(disasm.DefaultViewCtx()),
		buf.Len(),
	)
	return buf.Bytes()
}

// ctorWord - the 32-bit word (for non-compressible forms).
func ctorWord(t *testing.T, in Instr) uint32 {
	t.Helper()
	b := ctorBytes(t, in)
	require.Len(t, b, 4, "%s: compressed to %d bytes", in.ObjDump(disasm.DefaultViewCtx()), len(b))
	return binary.LittleEndian.Uint32(b)
}

// xreg - a register by number; a validation error fails the test.
func xreg(t *testing.T, n int) Reg {
	t.Helper()
	r, err := X(n)
	require.NoError(t, err)
	return r
}

// imm12 - a signed immediate -2048..2047; a validation error fails the test.
func imm12(t *testing.T, v int64) Imm12 {
	t.Helper()
	i, err := New().Imm12(v)
	require.NoError(t, err)
	return i
}

// imm20 - a U-type immediate 0..0xfffff; a validation error fails the test.
func imm20(t *testing.T, v int64) Imm20 {
	t.Helper()
	i, err := New().Imm20(v)
	require.NoError(t, err)
	return i
}

// off - a load/store byte offset -2048..2047; a validation error fails the test.
func off(t *testing.T, v int64) Off {
	t.Helper()
	o, err := New().Off(v)
	require.NoError(t, err)
	return o
}
