package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRegNameRoundTrip — a name built from a number is parsed back into the
// same number (register name generators are the basic parameter layer).
func TestRegNameRoundTrip(t *testing.T) {
	for n := range uint32(31) {
		for _, tc := range []struct {
			name string
			gen  func(uint32) string
		}{
			{
				"regNameX",
				regNameX,
			},
			{
				"regNameW",
				regNameW,
			},
			{
				"regNameXSP",
				regNameXSP,
			},
		} {
			name := tc.gen(n)
			got, err := armRegNum(name)
			require.NoError(t, err, "%s(%d) = %q", tc.name, n, name)
			require.Equal(t, n, got, "%s(%d) = %q", tc.name, n, name)
		}
	}

	// 31st: X - xzr, XSP - sp; both parse to 31.
	for _, name := range []string{regNameX(31), regNameXSP(31), "wzr", "wsp", "sp"} {
		got, err := armRegNum(name)
		require.NoError(t, err, "armRegNum(%q)", name)
		require.Equal(t, uint32(31), got, "armRegNum(%q)", name)
	}
}

// TestRegValidation — X/W out of range return an error.
func TestRegValidation(t *testing.T) {
	for _, bad := range []int{-1, 31, 100} {
		for _, f := range []func(int) (Reg, error){X, W} {
			_, err := f(bad)
			require.Error(t, err, "X/W(%d)", bad)
		}
	}

	// Range bounds are valid.
	for _, ok := range []int{0, 30} {
		_, err := X(ok)
		require.NoError(t, err, "X(%d)", ok)
		_, err = W(ok)
		require.NoError(t, err, "W(%d)", ok)
	}
}

// TestRegNames — String() of all classes.
func TestRegNames(t *testing.T) {
	cases := []struct {
		r    Reg
		want string
		is64 bool
		bits uint32
	}{
		{
			xreg(t, 0),
			"x0",
			true,
			0,
		},
		{
			xreg(t, 30),
			"x30",
			true,
			30,
		},
		{
			wreg(t, 7),
			"w7",
			false,
			7,
		},
		{
			XZR,
			"xzr",
			true,
			31,
		},
		{
			WZR,
			"wzr",
			false,
			31,
		},
		{
			SP,
			"sp",
			true,
			31,
		},
		{
			WSP,
			"wsp",
			false,
			31,
		},
	}
	for _, c := range cases {
		require.Equal(t, c.want, c.r.String(), "%s: String()", c.want)
		require.Equal(t, c.is64, c.r.Is64(), "%s: Is64", c.want)
		require.Equal(t, c.bits, c.r.bits(), "%s: bits", c.want)
		_, err := armRegNum(c.r.name())
		require.NoError(t, err, "%s: armRegNum", c.want)
	}
}

// TestImmValidation — immediate constructors validate the range.
func TestImmValidation(t *testing.T) {
	cases := []struct {
		name string
		kind string
		v    int64
	}{
		{"New().Imm12(-1)", "imm12", -1},
		{"New().Imm12(4096)", "imm12", 4096},
		{"New().Imm16(65536)", "imm16", 65536},
		{"New().Imm6(64)", "imm6", 64},
	}
	for _, c := range cases {
		var err error
		switch c.kind {
		case "imm12":
			_, err = New().Imm12(c.v)
		case "imm16":
			_, err = New().Imm16(c.v)
		case "imm6":
			_, err = New().Imm6(c.v)
		}

		assertErr(t, c.name, err)
	}
}

// TestEnumString — string representations of enums.
func TestEnumString(t *testing.T) {
	require.Equal(t, "lsl", LSL.String(), "Shift.String()")
	require.Equal(t, "ror", ROR.String(), "Shift.String()")
	require.Equal(t, "lsl #32", Hw2.String(), "Hw2.String()")
	require.Equal(t, "lsl #12", LSL12.String(), "Sh12.String()")
	require.Equal(t, "", NoSh12.String(), "Sh12.String()")
	// Shift.String() - the inverse of shiftNumByName.
	for _, s := range []Shift{LSL, LSR, ASR, ROR} {
		n, err := shiftNumByName(s.String())
		require.NoError(t, err, "shiftNumByName(%s)", s)
		require.Equal(t, s, Shift(n), "shiftNumByName(%s)", s)
	}
}
