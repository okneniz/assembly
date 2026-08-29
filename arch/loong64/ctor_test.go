package loong64

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// ctorBytes - the 4 bytes encoded by the constructor.
func ctorBytes(t *testing.T, in Instr) []byte {
	t.Helper()

	var buf bytes.Buffer
	_, err := in.Encode(&buf, 0x90000000)
	require.NoError(t, err, "Encode %s", in.ObjDump(disasm.DefaultViewCtx()))
	require.Len(
		t,
		buf.Bytes(),
		4,
		"Encode %s: %d bytes",
		in.ObjDump(disasm.DefaultViewCtx()),
		buf.Len(),
	)

	return buf.Bytes()
}

// ctorWord - the 32-bit word of the encoded instruction (at the fixed
// helper pc; Encode is pc-independent for every form now).
func ctorWord(t *testing.T, in Instr) uint32 {
	t.Helper()

	return binary.LittleEndian.Uint32(ctorBytes(t, in))
}

// lreg - a register by number; a validation error fails the test.
func lreg(t *testing.T, n int) Reg {
	t.Helper()

	r, err := R(n)
	require.NoError(t, err)

	return r
}
