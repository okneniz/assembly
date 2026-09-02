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
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"New().Imm12(2048)",
			func() error {
				_, err := New().Imm12(2048)
				return err
			},
		},
		{
			"New().Imm12(-2049)",
			func() error {
				_, err := New().Imm12(-2049)
				return err
			},
		},
		{
			"New().Imm20(-1)",
			func() error {
				_, err := New().Imm20(-1)
				return err
			},
		},
		{
			"New().Imm20(0x100000)",
			func() error {
				_, err := New().Imm20(0x100000)
				return err
			},
		},
		{
			"New().Off(2048)",
			func() error {
				_, err := New().Off(2048)
				return err
			},
		},
	} {
		err := c.call()
		require.Error(t, err, "case %q: out of range", c.name)
	}

	// Boundary values are valid.
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"New().Imm12(-2048)",
			func() error {
				_, err := New().Imm12(-2048)
				return err
			},
		},
		{
			"New().Imm12(2047)",
			func() error {
				_, err := New().Imm12(2047)
				return err
			},
		},
		{
			"New().Imm20(0)",
			func() error {
				_, err := New().Imm20(0)
				return err
			},
		},
		{
			"New().Imm20(0xfffff)",
			func() error {
				_, err := New().Imm20(0xfffff)
				return err
			},
		},
		{
			"New().Off(-2048)",
			func() error {
				_, err := New().Off(-2048)
				return err
			},
		},
		{
			"New().Off(2047)",
			func() error {
				_, err := New().Off(2047)
				return err
			},
		},
	} {
		err := c.call()
		require.NoError(t, err, "case %q", c.name)
	}

	require.Equal(t, "-0x4", imm12(t, -4).String())
}
