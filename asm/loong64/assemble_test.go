package loong64

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	asm "github.com/okneniz/assembly/asm"
)

func assembleOne(t *testing.T, src string, addr uint64) []byte {
	t.Helper()

	res, errs := asm.Assemble(src, addr, New())
	require.Empty(t, errs, "assemble %q: %v", src, errs)
	require.NotEmpty(t, res.Sections, "assemble %q: no sections", src)

	return res.Sections[0].Data
}

func TestAssembleWords(t *testing.T) {
	for _, tc := range []struct {
		src  string
		word uint32
	}{
		{"add.w $t0, $t1, $t2", 0x001039ac},
		{"add.w $r12, $r13, $r14", 0x001039ac},
		{"add.d $t0, $t1, $t2", 0x0010b9ac},
		{"addi.w $t0, $t1, -16", 0x02bfc1ac},
		{"addi.w $t0, $t1, 6*3-34", 0x02bfc1ac},
		{"andi $t0, $t1, 0xf0f", 0x037c3dac},
		{"lu12i.w $t0, 0x12345", 0x142468ac},
		{"slli.d $t0, $t1, 3", 0x00410dac},
		{"ld.d $t0, $t1, 8", 0x28c021ac},
		{"st.w $t0, $t1, -8", 0x29bfe1ac},
		{"csrrd $t0, 5", 0x0400140c},
		{"tlbsrch", 0x06482800},
		{"dbar 0", 0x38720000},
		{"alsl.w $t0, $t1, $t2, 3", 0x000539ac},
		{"bstrins.w $t0, $t1, 5, 3", 0x00650dac},
	} {
		b := assembleOne(t, tc.src, 0)
		require.Len(t, b, 4, tc.src)
		require.Equal(t, tc.word, binary.LittleEndian.Uint32(b), tc.src)
	}
}

func TestAssembleComments(t *testing.T) {
	b := assembleOne(t, "add.w $t0, $t1, $t2 # hash comment\n", 0)
	require.Equal(t, uint32(0x001039ac), binary.LittleEndian.Uint32(b))

	b = assembleOne(t, "// leading comment\nadd.w $t0, $t1, $t2", 0)
	require.Len(t, b, 4)
}

func TestAssembleLabelsAndBranches(t *testing.T) {
	// Forward and backward local labels resolve to absolute targets; the
	// offsets are computed against each instruction's own pc.
	b := assembleOne(t, "beqz $t1, 1f\n1: b 1b\n", 0x90000000)
	require.Len(t, b, 8, "two 4-byte instructions")

	// beqz at 0x90000000 targets 0x90000004: +1 word (offs16 = 1).
	require.Equal(t, uint32(0x400005a0), binary.LittleEndian.Uint32(b))
	// b at 0x90000004 jumps back to itself (1b = its own line): 0 words.
	require.Equal(t, uint32(0x50000000), binary.LittleEndian.Uint32(b[4:]))

	// A forward (+1 word) and a backward (-1 word) branch.
	b = assembleOne(t, "1: b 1f\nb 1b\n1: b 1b\n", 0)
	require.Equal(t, uint32(0x50000800), binary.LittleEndian.Uint32(b))
	require.Equal(t, uint32(0x53ffffff), binary.LittleEndian.Uint32(b[4:]))
	require.Equal(t, uint32(0x50000000), binary.LittleEndian.Uint32(b[8:]))
}

func TestAssembleExprSymbol(t *testing.T) {
	// The current-position symbol '.' in an expression.
	b := assembleOne(t, "b .+8\n", 0x90000000)
	require.Equal(t, uint32(0x50000800), binary.LittleEndian.Uint32(b))

	// Symbol arithmetic over a label: the offset to the next word.
	b = assembleOne(t, "1: jirl $t0, $t1, 1f-.\n1: andi $zero, $zero, 0\n", 0x90000000)
	require.Equal(t, uint32(0x4c0005ac), binary.LittleEndian.Uint32(b))
}

func TestAssembleErrors(t *testing.T) {
	for _, src := range []string{
		"add.w $t0, $t1",        // wrong count
		"add.w $t0, 5, $t2",     // number as a register
		"addi.w $t0, $t1, $t2",  // register as an immediate
		"addi.w $t0, $t1, 2048", // si12 out of range
		"no.such.mnem $t0, $t1", // unknown mnemonic
		"add.w t0, t1, t2",      // missing $ prefix
		"addi.w $t0, $t1, data", // unknown symbol
	} {
		_, errs := asm.Assemble(src, 0, New())
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}
}

func TestBackendApplyOption(t *testing.T) {
	be := New()
	require.NoError(t, be.ApplyOption("anything"))
	be.ResetOptions()
}

func TestOpAccessors(t *testing.T) {
	reg := OpReg("$t0")
	require.Equal(t, "$t0", reg.Reg())
	require.True(t, reg.IsReg())
	require.Nil(t, reg.Expr())

	e := OpNum(8)
	require.False(t, e.IsReg())
	require.NotNil(t, e.Expr())

	name, err := WantReg(reg)
	require.NoError(t, err)
	require.Equal(t, "$t0", name)

	_, err = WantReg(e)
	require.ErrorContains(t, err, "want register")

	_, err = WantReg(OpReg("$x0"))
	require.ErrorContains(t, err, "unknown register")

	got, err := WantExpr(e)
	require.NoError(t, err)
	require.NotNil(t, got)

	_, err = WantExpr(reg)
	require.ErrorContains(t, err, "want immediate")
}

func TestGrammarEdges(t *testing.T) {
	// A mnemonic immediately followed by more identifier characters is
	// unknown (the boundary check).
	_, errs := asm.Assemble("add.wX $t0, $t1, $t2", 0, New())
	require.NotEmpty(t, errs)

	// A trailing comma / a lone comma are operand errors.
	for _, src := range []string{"add.w $t0, $t1, $t2,", "add.w ,"} {
		_, errs := asm.Assemble(src, 0, New())
		require.NotEmpty(t, errs, "assemble %q: want errors", src)
	}

	// ResetOptions/ApplyOption are no-ops but must be callable.
	be := New()
	be.ResetOptions()
	require.NoError(t, be.ApplyOption("anything"))

	// ResolveForm evaluates a form directly (the pseudo machinery's
	// seam).
	be = New()
	res, err := be.ResolveForm(
		"add.w",
		[]Op{OpReg("$t0"), OpReg("$t1"), OpReg("$zero")},
		testCtx{addr: 0x90000000},
	)
	require.NoError(t, err)

	var buf bytes.Buffer
	n, err := res.Encode(&buf)
	require.NoError(t, err)
	require.Equal(t, int64(4), n)
	// add.w at pc 0x90000000: pc-independent.
	require.Equal(t, uint32(0x001001ac), binary.LittleEndian.Uint32(buf.Bytes()))
}

// testCtx is a fixed evaluation environment for ResolveForm.
type testCtx struct {
	addr uint64
}

func (c testCtx) Addr() uint64 {
	return c.addr
}

func (c testCtx) Resolve(string) (uint64, bool) {
	return 0, false
}
