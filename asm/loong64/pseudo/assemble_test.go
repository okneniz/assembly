package pseudo

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	asm "github.com/okneniz/assembly/asm"
)

func assembleOne(t *testing.T, src string, addr uint64) []byte {
	t.Helper()

	res, errs := Assemble(src, addr)
	require.Empty(t, errs, "assemble %q: %v", src, errs)
	require.NotEmpty(t, res.Sections, "assemble %q: no sections", src)

	return res.Sections[0].Data
}

func word(t *testing.T, b []byte, i int) uint32 {
	t.Helper()

	require.Greater(t, len(b), i*4, "word %d", i)

	return binary.LittleEndian.Uint32(b[i*4:])
}

func TestPseudoSingleWord(t *testing.T) {
	for _, tc := range []struct {
		src  string
		word uint32
	}{
		// llvm-mc-verified expansions.
		{"nop", 0x03400000},
		{"move $t0, $t1", 0x001501ac},
		{"not $t0, $t1", 0x001401ac},
		{"ret", 0x4c000020},
		{"jr $t1", 0x4c0001a0},
		{"bltz $t1, 8", 0x600009a0},
		{"bgez $t1, 8", 0x640009a0},
		{"bgtz $t0, 8", 0x6000080c},
		{"blez $t0, 8", 0x6400080c},
	} {
		b := assembleOne(t, tc.src, 0)
		require.Len(t, b, 4, tc.src)
		require.Equal(t, tc.word, binary.LittleEndian.Uint32(b), tc.src)
	}
}

func TestPseudoCallTail(t *testing.T) {
	b := assembleOne(t, "call 8\n", 0)
	require.Equal(t, uint32(0x54000800), word(t, b, 0))

	b = assembleOne(t, "tail 8\n", 0)
	require.Equal(t, uint32(0x50000800), word(t, b, 0))

	// A label target: call to the next instruction is +1 word.
	b = assembleOne(t, "call 1f\n1: nop\n", 0x90000000)
	require.Equal(t, uint32(0x54000400), word(t, b, 0))
}

func TestPseudoLiLadders(t *testing.T) {
	for _, tc := range []struct {
		src   string
		words []uint32
	}{
		// llvm-mc-verified chains.
		{"li.w $t0, -16", []uint32{0x02bfc00c}},
		{"li.d $t0, -16", []uint32{0x02bfc00c}},
		{"li.d $t0, 4095", []uint32{0x03bffc0c}},
		{"li.d $t0, 0x12345000", []uint32{0x142468ac}},
		{"li.w $t0, 0x12345678", []uint32{0x142468ac, 0x0399e18c}},
		{"li.d $t0, 0x12345678", []uint32{0x142468ac, 0x0399e18c}},
		{"li.d $t0, 0x123456789abcdef0", []uint32{0x153579ac, 0x03bbc18c, 0x168acf0c, 0x03048d8c}},
		{"li.d $t0, -414771", []uint32{0x15fff34c, 0x03af358c}},
	} {
		b := assembleOne(t, tc.src, 0)
		require.Len(t, b, 4*len(tc.words), tc.src)
		for i, w := range tc.words {
			require.Equal(t, w, word(t, b, i), "%s word %d", tc.src, i)
		}
	}
}

func TestPseudoLiSymbolic(t *testing.T) {
	// A symbolic operand takes the full 64-bit chain regardless of the
	// eventual value (an address, in practice - as GAS).
	b := assembleOne(t, "li.w $t0, 1f\n1: nop\n", 0x90000000)
	require.Len(t, b, 20, "li.w chain + nop")

	b = assembleOne(t, "li.d $t0, 1f\n1: nop\n", 0x90000000)
	require.Len(t, b, 20, "li.d chain + nop")
}

func TestPseudoLa(t *testing.T) {
	// la $t0, .+0x1024 at pc 0x90000000: page-aligned split
	// (hi=1, lo=36) - the pair verified with llvm-mc.
	b := assembleOne(t, "la $t0, .+0x1024\n", 0x90000000)
	require.Len(t, b, 8)
	require.Equal(t, uint32(0x1a00002c), word(t, b, 0))
	require.Equal(t, uint32(0x02c0918c), word(t, b, 1))

	// A negative low half: the sext12 carry goes into hi.
	b = assembleOne(t, "la $t0, .+0xf80\n", 0x90000000)
	// page 0x90000000, D = 0xf80: lo = -128, hi = 1.
	require.Equal(t, uint32(0x1a00002c), word(t, b, 0))
	require.Equal(t, uint32(0x02fe018c), word(t, b, 1))
}

func TestPseudoMixedWithReal(t *testing.T) {
	b := assembleOne(t, "nop\nmove $t0, $t1\nadd.w $t0, $t0, $t2\n", 0)
	require.Len(t, b, 12)
	require.Equal(t, uint32(0x03400000), word(t, b, 0))
	require.Equal(t, uint32(0x001501ac), word(t, b, 1))
	require.Equal(t, uint32(0x0010398c), word(t, b, 2))
}

func TestPseudoErrors(t *testing.T) {
	for _, src := range []string{
		"move $t0", // want rd, rj
		"li.w $t0", // want rd, imm
		"la $t0",   // want rd, sym
		"bltz $t1", // want rj, target
	} {
		_, errs := Assemble(src, 0)
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}
}

var _ = asm.Assemble

func TestPseudoLiLadderRungs(t *testing.T) {
	// The 3-word chain: bits above 32 with a zero low half.
	b := assembleOne(t, "li.d $t0, 0x500000000", 0)
	require.Len(t, b, 12)
	require.Equal(t, uint32(0x1400000c), word(t, b, 0)) // lu12i.w $t0, 0
	require.Equal(t, uint32(0x160000ac), word(t, b, 1)) // lu32i.d $t0, 5
	require.Equal(t, uint32(0x0300018c), word(t, b, 2)) // lu52i.d $t0, $t0, 0

	// The negative page rung of li32: -4096 is a single lu12i.w.
	b = assembleOne(t, "li.w $t0, -4096", 0)
	require.Equal(t, uint32(0x15ffffec), word(t, b, 0))
}

func TestPseudoLiOutOfRange(t *testing.T) {
	// li.w cannot hold values beyond 32 bits.
	_, errs := Assemble("li.w $t0, 0x123456789", 0)
	require.NotEmpty(t, errs)
}

func TestPseudoLaOutOfRange(t *testing.T) {
	// The pcalau12i+addi.d pair spans +-2 GiB; 4 GiB does not fit.
	_, errs := Assemble("la $t0, .+0x100000000", 0)
	require.NotEmpty(t, errs)
}

func TestPseudoExpandErrors(t *testing.T) {
	for _, src := range []string{
		"jr",             // want rj
		"jr $t0, $t1",    // want rj
		"call",           // want target
		"call 1, 2",      // want target
		"tail",           // want target
		"bltz $t1",       // want rj, target
		"li.w $t0, 1, 2", // want rd, imm
	} {
		_, errs := Assemble(src, 0)
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}
}

func TestPseudoEdges(t *testing.T) {
	// .option directives reach the no-op ApplyOption through the core.
	b := assembleOne(t, ".option norvc\nnop\n", 0)
	require.Len(t, b, 4)

	// A pseudo-mnemonic run into identifier characters falls back to
	// the inner grammar (and fails there).
	_, errs := Assemble("nopx $t0", 0)
	require.NotEmpty(t, errs)

	// Numeric unary/binary expressions stay numeric; symbolic
	// arithmetic takes the fixed chain.
	b = assembleOne(t, "li.d $t0, -8+4\n", 0)
	require.Len(t, b, 4)
	require.Equal(t, uint32(0x02bff00c), word(t, b, 0))

	b = assembleOne(t, "li.d $t0, 1f+4\n1: nop\n", 0x90000000)
	require.Len(t, b, 20)

	// la operand-kind errors.
	for _, src := range []string{"la 5, sym", "la $t0, $t1"} {
		_, errs := Assemble(src, 0)
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}
}

// TestPseudoRemainingBranches - the remaining error paths: expansion
// operand errors of every kind, the dead-defensive ones via direct
// calls.
func TestPseudoRemainingBranches(t *testing.T) {
	for _, src := range []string{
		"move", // wantRR
		"move $t0, $t1, $t2",
		"not", // wantRR (not)
		"not $t0",
		"li.w", // wantRegExpr
		"li.w 5, 8",
		"li.w $t0, $t1",
		"la $t0, $t1", // WantExpr (already)
	} {
		_, errs := Assemble(src, 0)
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}

	// The unknown-pseudo guard is unreachable through the grammar (the
	// trie only produces known names); call it directly.
	_, err := expandPseudo("nosuch", nil, testCtxP{})
	require.ErrorContains(t, err, "unknown pseudo-instruction")

	// exprNumeric: nil and the unreachable fallthrough.
	require.True(t, exprNumeric(nil))
}

// testCtxP is an empty evaluation environment for direct calls.
type testCtxP struct{}

func (testCtxP) Addr() uint64 {
	return 0
}

func (testCtxP) Resolve(string) (uint64, bool) {
	return 0, false
}
