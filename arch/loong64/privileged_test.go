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
	csr, err := NewUImm14(5)
	require.NoError(t, err)

	op5, err := NewUImm5(5)
	require.NoError(t, err)

	off, err := NewImm12(8)
	require.NoError(t, err)

	lvl, err := NewUImm8(1)
	require.NoError(t, err)

	code, err := NewCode15(1)
	require.NoError(t, err)

	family := []struct {
		mnem string
		ctor func() Instr
	}{
		{"csrrd", func() Instr { return NewCsrrd(lreg(t, 12), csr) }},
		{"csrwr", func() Instr { return NewCsrwr(lreg(t, 12), csr) }},
		{"csrxchg", func() Instr { return NewCsrxchg(lreg(t, 12), lreg(t, 13), csr) }},
		{"cacop", func() Instr { return NewCacop(op5, lreg(t, 13), off) }},
		{"lddir", func() Instr { return NewLddir(lreg(t, 12), lreg(t, 13), lvl) }},
		{"ldpte", func() Instr { return NewLdpte(lreg(t, 13), lvl) }},
		{"iocsrrd.b", func() Instr { return NewIocsrrdB(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.h", func() Instr { return NewIocsrrdH(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.w", func() Instr { return NewIocsrrdW(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrrd.d", func() Instr { return NewIocsrrdD(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.b", func() Instr { return NewIocsrwrB(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.h", func() Instr { return NewIocsrwrH(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.w", func() Instr { return NewIocsrwrW(lreg(t, 12), lreg(t, 13)) }},
		{"iocsrwr.d", func() Instr { return NewIocsrwrD(lreg(t, 12), lreg(t, 13)) }},
		{"tlbclr", NewTlbclr},
		{"tlbflush", NewTlbflush},
		{"tlbsrch", NewTlbsrch},
		{"tlbrd", NewTlbrd},
		{"tlbwr", NewTlbwr},
		{"tlbfill", NewTlbfill},
		{"ertn", NewErtn},
		{"idle", func() Instr { return NewIdle(code) }},
		{"invtlb", func() Instr { return NewInvtlb(op5, lreg(t, 13), lreg(t, 14)) }},
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
