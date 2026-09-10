package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRegNames - number -> ABI name -> regBits (round-trip of register names).
func TestRegNames(t *testing.T) {
	for n := range 32 {
		r := xreg(t, n)
		require.Equal(t, rvRegNames[n], r.String(), "X(%d).String()", n)
		require.Equal(t, uint32(r.Num()), regBits(r.name()), "regBits(%s)", r.name())
	}

	require.Equal(t, "zero", Zero.String())
	require.Equal(t, "ra", Ra.String())
	require.Equal(t, "sp", Sp.String())
	_, err := X(-1)
	require.Error(t, err)
	_, err = X(32)
	require.Error(t, err)
	require.ErrorContains(t, err, "0..31")
}

// TestImmValidation - immediate and offset ranges.
func TestImmValidation(t *testing.T) {
	// out of range
	errCases := []struct {
		name string
		kind string
		v    int64
	}{
		{"New().Imm12(2048)", "imm12", 2048},
		{"New().Imm12(-2049)", "imm12", -2049},
		{"New().Imm20(-1)", "imm20", -1},
		{"New().Imm20(0x100000)", "imm20", 0x100000},
		{"New().Off(2048)", "off", 2048},
	}
	for _, c := range errCases {
		var err error
		switch c.kind {
		case "imm12":
			_, err = New().Imm12(c.v)
		case "imm20":
			_, err = New().Imm20(c.v)
		case "off":
			_, err = New().Off(c.v)
		}

		require.Error(t, err, "case %q: out of range", c.name)
	}

	// Boundary values are valid.
	okCases := []struct {
		name string
		kind string
		v    int64
	}{
		{"New().Imm12(-2048)", "imm12", -2048},
		{"New().Imm12(2047)", "imm12", 2047},
		{"New().Imm20(0)", "imm20", 0},
		{"New().Imm20(0xfffff)", "imm20", 0xfffff},
		{"New().Off(-2048)", "off", -2048},
		{"New().Off(2047)", "off", 2047},
	}
	for _, c := range okCases {
		var err error
		switch c.kind {
		case "imm12":
			_, err = New().Imm12(c.v)
		case "imm20":
			_, err = New().Imm20(c.v)
		case "off":
			_, err = New().Off(c.v)
		}

		require.NoError(t, err, "case %q: boundary", c.name)
	}
}
