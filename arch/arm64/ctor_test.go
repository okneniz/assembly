package arm64

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// ctorWord — the word encoded by a constructor.
func ctorWord(t *testing.T, in Instr) uint32 {
	t.Helper()
	var buf bytes.Buffer
	_, err := in.Encode(&buf, 0x1000)
	require.NoError(t, err, "Encode %s", in.ObjDump(disasm.DefaultViewCtx()))
	require.Len(t, buf.Bytes(), 4, "Encode %s", in.ObjDump(disasm.DefaultViewCtx()))
	return binary.LittleEndian.Uint32(buf.Bytes())
}

// assertErr — the call must return an error (a contextual constraint).
func assertErr(t *testing.T, name string, err error) {
	t.Helper()
	require.Error(t, err, "%s: error expected", name)
}

// xreg/wreg/imm12/imm16/imm6 — valid operand values for table literals
// (an error is impossible by construction).
func xreg(t *testing.T, n int) Reg {
	t.Helper()
	r, err := X(n)
	require.NoError(t, err)
	return r
}

func wreg(t *testing.T, n int) Reg {
	t.Helper()
	r, err := W(n)
	require.NoError(t, err)
	return r
}

func imm12(t *testing.T, v int64) Imm12 {
	t.Helper()
	iv, err := New().Imm12(v)
	require.NoError(t, err)
	return iv
}

func imm16(t *testing.T, v int64) Imm16 {
	t.Helper()
	iv, err := New().Imm16(v)
	require.NoError(t, err)
	return iv
}

func imm6(t *testing.T, v int64) Imm6 {
	t.Helper()
	iv, err := New().Imm6(v)
	require.NoError(t, err)
	return iv
}

// ctorAddImm and its siblings — instruction constructor wrappers for
// table literals: valid operands, an error is impossible by construction.
func ctorAddImm(t *testing.T, rd, rn Reg, imm Imm12, sh Sh12) Instr {
	t.Helper()
	in, err := New().AddImm(rd, rn, imm, sh)
	require.NoError(t, err)
	return in
}

func ctorAddShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().AddShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}

func ctorSubImm(t *testing.T, rd, rn Reg, imm Imm12, sh Sh12) Instr {
	t.Helper()
	in, err := New().SubImm(rd, rn, imm, sh)
	require.NoError(t, err)
	return in
}

func ctorSubShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().SubShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}

func ctorMovz(t *testing.T, rd Reg, imm Imm16, hw Hw) Instr {
	t.Helper()
	in, err := New().Movz(rd, imm, hw)
	require.NoError(t, err)
	return in
}

func ctorMovk(t *testing.T, rd Reg, imm Imm16, hw Hw) Instr {
	t.Helper()
	in, err := New().Movk(rd, imm, hw)
	require.NoError(t, err)
	return in
}

func ctorLdr(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldr(rt, rn, off)
	require.NoError(t, err)
	return in
}

func ctorStr(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Str(rt, rn, off)
	require.NoError(t, err)
	return in
}

func ctorRet(t *testing.T, rn Reg) Instr {
	t.Helper()
	in, err := New().Ret(rn)
	require.NoError(t, err)
	return in
}
