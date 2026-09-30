package asm

// The fragment parser on the mock backend: sizes under the placeholder
// environment, the numeric local table (Nb/Nf by statement order, the
// redefinable GAS discipline), the fragment contract rejections, and the
// resolve error carrying the fragment line. The arm64 fragment oracle
// (fragment bytes == full assembly bytes) lives in asm/arm64/alias.

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/unit"
)

func TestFragmentSize(t *testing.T) {
	// data-only: the source and the expected byte size
	tests := []struct {
		src  string
		size int
	}{
		{"pad 4", 4},
		{"pad 4\npad 2", 6},
		{"", 0},
		{"# a comment\n\npad 4\n# another", 4},
		{"1:\n2:\npad 4\n3:", 4}, // labels carry no bytes
		{"pad 4\n1:\npad 4\n1:\npad 4", 12},
	}

	for _, tt := range tests {
		f, errs := ParseFragment(tt.src, mockBackend{})
		require.Empty(t, errs, tt.src)
		require.Equal(t, tt.size, f.Size(), tt.src)
	}
}

func TestFragmentNumericLocals(t *testing.T) {
	// "pad 4 / 1: / pad 4 / 2: / pad 4 / 1: / pad 4": the pads sit at
	// offsets 0, 4, 8, 12; the defs of "1" at statements 1 and 5
	// (offsets 4 and 12), the def of "2" at statement 3 (offset 8)
	src := "pad 4\n1:\npad 4\n2:\npad 4\n1:\npad 4"
	f, errs := ParseFragment(src, mockBackend{})
	require.Empty(t, errs)
	require.Equal(t, 16, f.Size())

	// b: at the statement itself or earlier; f: strictly later
	tests := []struct {
		name string
		stmt int
		off  int
		ok   bool
	}{
		{"1b", 3, 4, true},  // the first 1: is behind
		{"1b", 6, 12, true}, // the redefinition wins for later statements
		{"1f", 0, 4, true},  // the first 1: is ahead
		{"1f", 2, 12, true}, // skips to the redefinition
		{"2b", 2, 0, false}, // 2: is not at or behind statement 2
		{"2b", 3, 8, true},  // the label of the statement itself
		{"2f", 4, 0, false}, // 2: is not strictly ahead anymore
		{"3b", 6, 0, false}, // no def of 3 at all
	}

	for _, tt := range tests {
		v, ok := f.localAddr(tt.name, tt.stmt, 0)
		require.Equal(t, tt.ok, ok, tt.name)
		if ok {
			require.Equal(t, uint64(tt.off), v, tt.name)
		}
	}
}

func TestFragmentRejects(t *testing.T) {
	// data-only: the source and the message fragment of the rejection
	tests := []struct {
		src string
		msg string
	}{
		{"name:", "named label"},
		{"name: pad 4", "named label"},
		{".word 5", "directive"},
		{".text\npad 4", "directive"},
		{"pool 4", "literal pools"},
		{"garbage", "instruction"},
	}

	for _, tt := range tests {
		f, errs := ParseFragment(tt.src, mockBackend{})
		require.Nil(t, f, tt.src)
		require.NotEmpty(t, errs, tt.src)
		require.ErrorContains(t, errs[0], tt.msg, tt.src)
	}
}

func TestFragmentToleratesIgnoredDirectives(t *testing.T) {
	// the compiler-output class the file path ignores (.arch_extension of
	// the kernel's inline asm among it) produces nothing to layout - the
	// fragment skips it; everything else still rejects (TestFragmentRejects)
	src := ".arch_extension fp\npad 4\n.arch_extension nofp\npad 4\n"
	f, errs := ParseFragment(src, mockBackend{})
	require.Empty(t, errs, src)
	require.NotNil(t, f, src)
	require.Equal(t, 8, f.Size(), src)
}

func TestFragmentRejectsReportLines(t *testing.T) {
	_, errs := ParseFragment("pad 4\nname:\npad 4", mockBackend{})
	require.Len(t, errs, 1)
	require.Equal(t, uint(2), errs[0].Line)
}

func TestFragmentFailsAtParse(t *testing.T) {
	// "fail" breaks under the placeholder environment already - the size
	// walk resolves every instruction, so resolve failures without
	// symbols surface at parse (the ones needing the host - a bad local,
	// an out-of-range host offset - surface at resolve, see the alias
	// oracle tests)
	_, errs := ParseFragment("pad 4\nfail", mockBackend{})
	require.Len(t, errs, 1)
	require.Equal(t, uint(2), errs[0].Line)
}

func TestFragmentResolveBytes(t *testing.T) {
	// the fragment resolves to its own bytes at its own addresses: each
	// pad lands at base + its offset
	f, errs := ParseFragment("pad 4\npad 2", mockBackend{})
	require.Empty(t, errs)

	rs, err := f.Resolve(unit.NewCtx(0x1000, nil))
	require.NoError(t, err)
	require.Len(t, rs, 2)

	var buf bytes.Buffer
	for _, r := range rs {
		_, err := r.Encode(&buf)
		require.NoError(t, err)
	}

	require.Equal(t, append(
		bytes.Repeat([]byte{0xAB}, 4),
		bytes.Repeat([]byte{0xAB}, 2)...,
	), buf.Bytes())
}
