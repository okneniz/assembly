package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPrivilegedJSONEncodeError - the marshal and write-error paths of
// every privileged instruction (the rest is covered by the per-file
// tests).
func TestPrivilegedJSONEncodeError(t *testing.T) {
	csr, err := New().UImm14(5)
	require.NoError(t, err)

	op5, err := New().UImm5(5)
	require.NoError(t, err)

	off, err := New().Imm12(8)
	require.NoError(t, err)

	lvl, err := New().UImm8(1)
	require.NoError(t, err)

	code, err := New().Code15(1)
	require.NoError(t, err)

	family := []struct {
		mnem string
		ctor func() Instr
	}{
		{"csrrd", func() Instr { return New().Csrrd(lreg(t, 12), csr) }},
		{"csrwr", func() Instr { return New().Csrwr(lreg(t, 12), csr) }},
		{"csrxchg", func() Instr { return New().Csrxchg(lreg(t, 12), lreg(t, 13), csr) }},
		{"cacop", func() Instr { return New().Cacop(op5, lreg(t, 13), off) }},
		{"lddir", func() Instr { return New().Lddir(lreg(t, 12), lreg(t, 13), lvl) }},
		{"ldpte", func() Instr { return New().Ldpte(lreg(t, 13), lvl) }},
		{"iocsrrd.b", func() Instr { return New().IocsrrdB(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.h", func() Instr { return New().IocsrrdH(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.w", func() Instr { return New().IocsrrdW(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.d", func() Instr { return New().IocsrrdD(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.b", func() Instr { return New().IocsrwrB(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.h", func() Instr { return New().IocsrwrH(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.w", func() Instr { return New().IocsrwrW(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.d", func() Instr { return New().IocsrwrD(lreg(t, 12), lreg(t, 13)) }},
		{"tlbclr", New().Tlbclr},
		{"tlbflush", New().Tlbflush},
		{"tlbsrch", New().Tlbsrch},
		{"tlbrd", New().Tlbrd},
		{"tlbwr", New().Tlbwr},
		{"tlbfill", New().Tlbfill},
		{"ertn", New().Ertn},
		{"idle", func() Instr { return New().Idle(code) }},
		{"invtlb", func() Instr { return New().Invtlb(op5, lreg(t, 13), lreg(t, 14)) }},
	}

	for _, f := range family {
		in := f.ctor()

		b, err := in.MarshalJSON()
		require.NoError(t, err, f.mnem)

		var dto map[string]any
		require.NoError(t, json.Unmarshal(b, &dto), f.mnem)
		require.Equal(t, f.mnem, dto["mnemonic"], f.mnem)
		require.Equal(t, "Privileged", dto["group"], f.mnem)

		_, err = in.Encode(errWriter{}, 0)
		require.ErrorContains(t, err, "write failed", f.mnem)
	}
}
