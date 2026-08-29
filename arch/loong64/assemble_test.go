package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// One assembled word per shape family: the operand-kind/count validation
// and the delegation to the role constructors, via BuildInstr. The
// per-instruction encodings are covered by the per-file tests.
func TestBuildInstrShapes(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }
	n := OpNum

	for _, tc := range []struct {
		mnem string
		ops  []Op
		word uint32
	}{
		{"add.w", []Op{r(12), r(13), r(14)}, 0x001039ac},
		{"amadd.w", []Op{r(12), r(13), r(14)}, 0x386135cc},
		{"clz.d", []Op{r(12), r(13)}, 0x000025ac},
		{"asrtle.d", []Op{r(13), r(14)}, 0x000139a0},
		{"addi.w", []Op{r(12), r(13), n(-16)}, 0x02bfc1ac},
		{"andi", []Op{r(12), r(13), n(3855)}, 0x037c3dac},
		{"addu16i.d", []Op{r(12), r(13), n(-1)}, 0x13fffdac},
		{"lu12i.w", []Op{r(12), n(74565)}, 0x142468ac},
		{"slli.w", []Op{r(12), r(13), n(3)}, 0x00408dac},
		{"slli.d", []Op{r(12), r(13), n(3)}, 0x00410dac},
		{"bytepick.w", []Op{r(12), r(13), r(14), n(3)}, 0x0009b9ac},
		{"bytepick.d", []Op{r(12), r(13), r(14), n(3)}, 0x000db9ac},
		{"bstrins.w", []Op{r(12), r(13), n(5), n(3)}, 0x00650dac},
		{"bstrins.d", []Op{r(12), r(13), n(63), n(0)}, 0x00bf01ac},
		{"alsl.w", []Op{r(12), r(13), r(14), n(3)}, 0x000539ac},
		{"beq", []Op{r(13), r(12), n(8)}, 0x580009ac},
		{"beqz", []Op{r(13), n(8)}, 0x400009a0},
		{"b", []Op{n(8)}, 0x50000800},
		{"jirl", []Op{r(12), r(13), n(4)}, 0x4c0005ac},
		{"break", []Op{n(1)}, 0x002a0001},
		{"ldptr.w", []Op{r(12), r(13), n(8)}, 0x240009ac},
		{"preld", []Op{n(5), r(13), n(8)}, 0x2ac021a5},
		{"preldx", []Op{n(5), r(13), r(14)}, 0x382c39a5},
		{"lddir", []Op{r(12), r(13), n(1)}, 0x064005ac},
		{"ldpte", []Op{r(13), n(1)}, 0x064405a0},
		{"csrrd", []Op{r(12), n(5)}, 0x0400140c},
		{"csrxchg", []Op{r(12), r(13), n(5)}, 0x040015ac},
		{"tlbsrch", nil, 0x06482800},
	} {
		in, err := BuildInstr(tc.mnem, tc.ops)
		require.NoError(t, err, tc.mnem)
		require.Equal(t, tc.word, ctorWord(t, in), tc.mnem)
	}
}

func TestBuildInstrErrors(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }
	n := OpNum

	for _, tc := range []struct {
		mnem  string
		ops   []Op
		wantE string
	}{
		{"no.such", nil, "unknown instruction"},
		{"add.w", []Op{r(1), r(2)}, "want 3 operands"},
		{"add.w", []Op{r(1), r(2), n(3)}, "want a register"},
		{"addi.w", []Op{r(1), n(2), n(3)}, "want a register"},
		{"addi.w", []Op{r(1), r(2), r(3)}, "want an immediate"},
		{"addi.w", []Op{r(1), r(2), n(2048)}, "outside"},
		{"andi", []Op{r(1), r(2), n(-1)}, "outside"},
		{"bstrins.w", []Op{r(1), r(2), n(3), n(5)}, "msb is below lsb"},
		{"alsl.w", []Op{r(1), r(2), r(3), n(5)}, "outside"},
		{"csrrd", []Op{r(1), n(16384)}, "outside"},
		{"tlbsrch", []Op{n(1)}, "want 0 operands"},
	} {
		_, err := BuildInstr(tc.mnem, tc.ops)
		require.ErrorContains(t, err, tc.wantE, tc.mnem)
	}
}

func TestRegNumOf(t *testing.T) {
	for name, num := range map[string]uint32{
		"$zero": 0, "$ra": 1, "$sp": 3, "$a0": 4, "$a7": 11,
		"$t0": 12, "$t8": 20, "$r21": 21, "$fp": 22, "$s0": 23, "$s8": 31,
		"$r0": 0, "$r12": 12, "$r31": 31,
	} {
		v, err := RegNumOf(name)
		require.NoError(t, err, name)
		require.Equal(t, num, v, name)
	}

	for _, name := range []string{"t0", "$x0", "$r32", "$r-1", "$r", ""} {
		_, err := RegNumOf(name)
		require.Error(t, err, name)
	}
}

func TestMnemonicsAndRegNames(t *testing.T) {
	require.Len(t, Mnemonics(), 248)

	names := RegNames()
	require.Equal(t, "$zero", names[0])
	require.Equal(t, "$t0", names[12])
	require.Equal(t, "$r21", names[21])
	require.Equal(t, "$s8", names[31])
}
