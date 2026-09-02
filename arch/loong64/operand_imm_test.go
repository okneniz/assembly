package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The role types: one happy-path and one out-of-range value each (the
// constructors are the only range gate - the scatter side trusts them).
func TestImmRoleTypes(t *testing.T) {
	ok := []struct {
		name string
		make func(int64) (interface{ Val() int64 }, error)
		v    int64
	}{
		{"si12", func(v int64) (interface{ Val() int64 }, error) { return New().Imm12(v) }, -2048},
		{"ui12", func(v int64) (interface{ Val() int64 }, error) { return New().UImm12(v) }, 4095},
		{"si14", func(v int64) (interface{ Val() int64 }, error) { return New().Imm14(v) }, 16380},
		{"si16", func(v int64) (interface{ Val() int64 }, error) { return New().Imm16(v) }, -32768},
		{
			"off16",
			func(v int64) (interface{ Val() int64 }, error) { return New().Off16(v) },
			131068,
		},
		{"si20", func(v int64) (interface{ Val() int64 }, error) { return New().Imm20(v) }, 524287},
		{"ui5", func(v int64) (interface{ Val() int64 }, error) { return New().UImm5(v) }, 31},
		{"ui6", func(v int64) (interface{ Val() int64 }, error) { return New().UImm6(v) }, 63},
		{
			"alsl shift",
			func(v int64) (interface{ Val() int64 }, error) { return New().Shift3(v) },
			4,
		},
		{"ui8", func(v int64) (interface{ Val() int64 }, error) { return New().UImm8(v) }, 255},
		{"csr", func(v int64) (interface{ Val() int64 }, error) { return New().UImm14(v) }, 16383},
		{"code", func(v int64) (interface{ Val() int64 }, error) { return New().Code15(v) }, 32767},
	}
	for _, tc := range ok {
		x, err := tc.make(tc.v)
		require.NoError(t, err, tc.name)
		require.Equal(t, tc.v, x.Val(), tc.name)
	}

	bad := []struct {
		name string
		make func(int64) (interface{ Val() int64 }, error)
		v    int64
	}{
		{
			"si12 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm12(v) },
			2048,
		},
		{
			"si12 under",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm12(v) },
			-2049,
		},
		{
			"ui12 negative",
			func(v int64) (interface{ Val() int64 }, error) { return New().UImm12(v) },
			-1,
		},
		{
			"ui12 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().UImm12(v) },
			4096,
		},
		{
			"si14 misaligned",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm14(v) },
			6,
		},
		{
			"si14 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm14(v) },
			16384,
		},
		{
			"si16 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm16(v) },
			32768,
		},
		{
			"off16 misaligned",
			func(v int64) (interface{ Val() int64 }, error) { return New().Off16(v) },
			2,
		},
		{
			"off16 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Off16(v) },
			131072,
		},
		{
			"si20 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Imm20(v) },
			524288,
		},
		{"ui5 over", func(v int64) (interface{ Val() int64 }, error) { return New().UImm5(v) }, 32},
		{"ui6 over", func(v int64) (interface{ Val() int64 }, error) { return New().UImm6(v) }, 64},
		{
			"alsl shift zero",
			func(v int64) (interface{ Val() int64 }, error) { return New().Shift3(v) },
			0,
		},
		{
			"alsl shift over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Shift3(v) },
			5,
		},
		{
			"ui8 over",
			func(v int64) (interface{ Val() int64 }, error) { return New().UImm8(v) },
			256,
		},
		{
			"csr over",
			func(v int64) (interface{ Val() int64 }, error) { return New().UImm14(v) },
			16384,
		},
		{
			"code over",
			func(v int64) (interface{ Val() int64 }, error) { return New().Code15(v) },
			32768,
		},
	}
	for _, tc := range bad {
		_, err := tc.make(tc.v)
		require.Error(t, err, tc.name)
	}
}

func TestBitsScatter(t *testing.T) {
	// The pure scatters: validated values in, bits out.
	require.Equal(t, uint32(0xff0<<10), scatterS(-16, 10, 12))
	require.Equal(t, uint32(0x1ffffe0), scatterS(-1, 5, 20))
	require.Equal(t, uint32(0xff<<10), scatterU(255, 10, 12))
	require.Equal(t, uint32(0x1f|0xfffe<<10), scatterD5k16(-2))
	require.Equal(t, uint32(0x3ff|0xfffe<<10), scatterD10k16(-2))
	require.Equal(t, uint32(0x1f|0xabcd<<10), scatterD10k16(0x1fabcd))

	// The gathers are their inverses.
	require.Equal(t, int64(-2), d5k16Imm(0x3ff|0xfffe<<10))
	require.Equal(t, int64(-2), d10k16Imm(0x3ff|0xfffe<<10))
	require.Equal(t, int64(0x1fabcd), d10k16Imm(0x1f|0xabcd<<10))
	require.Equal(t, int64(-16), sField(0xff0<<10, 10, 12))
	require.Equal(t, uint32(0xff), uField(0xff0<<10, 14, 8))
}

func TestBEncodeOutOfRange(t *testing.T) {
	in := New().B(1 << 28)

	_, err := in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "does not fit")

	in = New().B(2)
	_, err = in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "not word-aligned")
}

func TestBeqzEncodeOutOfRange(t *testing.T) {
	in := New().Beqz(lreg(t, 13), 1<<23)

	_, err := in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "does not fit")
}
