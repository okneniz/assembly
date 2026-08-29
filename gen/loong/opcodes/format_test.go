package opcodes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFormat(t *testing.T) {
	slots, err := ParseFormat("DJK")
	require.NoError(t, err)
	require.Equal(t, NewSlots([]rune{'D', 'J', 'K'}, nil), slots)

	slots, err = ParseFormat("EMPTY")
	require.NoError(t, err)
	require.Equal(t, NewSlots(nil, nil), slots)

	// A split immediate: [4:0] is the high part, [25:10] the low one.
	slots, err = ParseFormat("JSd5k16")
	require.NoError(t, err)
	require.Equal(t, []rune{'J'}, slots.Regs)
	require.Len(t, slots.Imms, 1)
	require.True(t, slots.Imms[0].Signed)
	require.Equal(t, []Field{NewField(0, 5), NewField(10, 16)}, slots.Imms[0].Segments)

	// Two immediates (bstrins.w): Uk5 then Um5.
	slots, err = ParseFormat("DJUk5Um5")
	require.NoError(t, err)
	require.Len(t, slots.Imms, 2)
	require.False(t, slots.Imms[0].Signed)
	require.Equal(t, NewField(10, 5), slots.Imms[0].Segments[0])
	require.Equal(t, NewField(16, 5), slots.Imms[1].Segments[0])

	// Immediate-only formats (b/bl).
	slots, err = ParseFormat("Sd10k16")
	require.NoError(t, err)
	require.Empty(t, slots.Regs)
	require.Len(t, slots.Imms, 1)
}

func TestParseFormatErrors(t *testing.T) {
	_, err := ParseFormat("")
	require.ErrorContains(t, err, "no operand slots")

	_, err = ParseFormat("ZJK")
	require.ErrorContains(t, err, "no operand slots")

	// Unrecognized tail after the slots.
	_, err = ParseFormat("DJSk12x9")
	require.ErrorContains(t, err, "unparsed tail")

	// A segment without its width: the immediate fails as a whole (Try)
	// and the tail stays unparsed.
	_, err = ParseFormat("DJSk")
	require.ErrorContains(t, err, "unparsed tail")

	// A register letter after the first immediate.
	_, err = ParseFormat("Sd5k16D")
	require.ErrorContains(t, err, "register after immediate")
}

func TestCastWidth(t *testing.T) {
	v, err := castWidth([]rune("12"))
	require.NoError(t, err)
	require.Equal(t, 12, v)

	// Not reachable through the grammar (Some caps at 2 digits); the guard
	// stays for direct callers.
	_, err = castWidth([]rune("999"))
	require.ErrorContains(t, err, "segment width")
}

func TestMask(t *testing.T) {
	for _, tc := range []struct {
		format string
		mask   uint32
	}{
		{"EMPTY", 0xffffffff},
		{"DJK", 0xffff8000},
		{"DJ", 0xfffffc00},
		{"JK", 0xffff801f},
		{"DJSk12", 0xffc00000},
		{"DJUk12", 0xffc00000},
		{"DJSk14", 0xff000000},
		{"DJSk16", 0xfc000000},
		{"DSj20", 0xfe000000},
		{"Sd10k16", 0xfc000000},
		{"JSd5k16", 0xfc000000},
		{"Ud15", 0xffff8000},
		{"DJKUa2", 0xfffe0000},
		{"DJKUa3", 0xfffc0000},
		{"DJUk5Um5", 0xffe08000},
		{"DJUk6Um6", 0xffc00000},
		{"DJUk8", 0xfffc0000},
		{"JUd5Sk12", 0xffc00000},
	} {
		slots, err := ParseFormat(tc.format)
		require.NoError(t, err, tc.format)

		mask, err := slots.Mask()
		require.NoError(t, err, tc.format)
		require.Equal(t, tc.mask, mask, tc.format)
	}
}

func TestMaskErrors(t *testing.T) {
	// Overlapping slots: the D register twice.
	slots := NewSlots([]rune{'D', 'D'}, nil)
	_, err := slots.Mask()
	require.ErrorContains(t, err, "overlapping slots")

	// A segment reaching past bit 31: [32:18].
	slots = NewSlots(nil, []Imm{NewImm(true, []Field{NewField(18, 15)})})
	_, err = slots.Mask()
	require.ErrorContains(t, err, "out of range")

	// A zero-width segment.
	slots = NewSlots(nil, []Imm{NewImm(true, []Field{NewField(10, 0)})})
	_, err = slots.Mask()
	require.ErrorContains(t, err, "out of range")

	// A negative offset (not producible by the grammar either; NewField is
	// a public constructor, so the guard is reachable).
	slots = NewSlots(nil, []Imm{NewImm(true, []Field{NewField(-1, 5)})})
	_, err = slots.Mask()
	require.ErrorContains(t, err, "out of range")
}

func TestFields(t *testing.T) {
	slots, err := ParseFormat("DJKUa2")
	require.NoError(t, err)

	require.Equal(
		t,
		[]Field{NewField(0, 5), NewField(5, 5), NewField(10, 5), NewField(15, 2)},
		slots.Fields(),
	)

	empty, err := ParseFormat("EMPTY")
	require.NoError(t, err)
	require.Empty(t, empty.Fields())
}
