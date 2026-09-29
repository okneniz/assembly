package alias

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/unit"
)

// TestMacroVentry is the traps.S shape: ".macro ventry label" (the
// comma-less seL4 spelling), the \label substitution inside an
// identifier (trap_\label), .align in the body, and the expansions
// landing back to back in the statement stream.
func TestMacroVentry(t *testing.T) {
	src := `.macro ventry label
    .align 7
    b trap_\label
.endm
ventry 0
ventry 1
trap_0: ret
trap_1: ret
`
	res, errs := Assemble(src, 0x1000)
	require.Empty(t, errs, "errs: %v", errs)

	d := res.Sections[0].Data
	require.Len(t, d, 140, "b, pad to 128, b, two rets")

	// ventry 0: no pad @0, b trap_0 @0 -> trap_0 @0x1084 (off 132)
	require.Equal(t, uint32(0x14000021), binary.LittleEndian.Uint32(d[0:]))
	// ventry 1: pad to 128, b trap_1 @128 -> trap_1 @0x1088 (off 136)
	require.Equal(t, uint32(0x14000002), binary.LittleEndian.Uint32(d[128:]))
	for i := 4; i < 128; i++ {
		require.Zerof(t, d[i], "the alignment padding at %d", i)
	}

	require.Equal(t, uint32(0xd65f03c0), binary.LittleEndian.Uint32(d[132:]))
	require.Equal(t, uint64(0x1084), res.Symbols["trap_0"])
	require.Equal(t, uint64(0x1088), res.Symbols["trap_1"])

	// the unit mode assembles the same macro source to the same bytes
	u := unit.New()
	require.Empty(t, AssembleUnit(u, "traps.S", src))
	fixed := u.Resolve(flatPlace(0x1000, 0x80000000))
	require.Empty(t, fixed.Errs)
	text, err := fixed.EncodeText()
	require.NoError(t, err)
	require.Equal(t, d, text)
}

// TestMacroParams is the argument discipline: positional mapping,
// defaults, the empty string for a missing argument (as in GAS), the
// no-comma header, a macro invoking a macro, and the comment/separator
// cut of the invocation line.
func TestMacroParams(t *testing.T) {
	src := `.macro save, reg, amt=16 // the frame of one register
    str \reg, [sp, #-\amt]!
.endm
.macro save2 amt
    save x1, \amt
    save x2 ; save x3
.endm
save2 8
save x9 // a trailing comment
`
	res, errs := Assemble(src, 0x1000)
	require.Empty(t, errs, "errs: %v", errs)

	// save2 #8: save x1,#8; save x2,#16 (default); save x3,#16
	d := res.Sections[0].Data
	require.Len(t, d, 16)

	// str x1, [sp, #-8]!: pre-index, imm9 = -8
	require.Equal(t, uint32(0xF81F8FE1), binary.LittleEndian.Uint32(d[0:]), "str x1, #8")
	// str x2, [sp, #-16]!: the default amt; imm9=-16 (llvm-checked pins)
	require.Equal(t, uint32(0xF81F0FE2), binary.LittleEndian.Uint32(d[4:]), "str x2, default 16")
	require.Equal(t, uint32(0xF81F0FE3), binary.LittleEndian.Uint32(d[8:]), "str x3, default 16")
	require.Equal(t, uint32(0xF81F0FE9), binary.LittleEndian.Uint32(d[12:]), "str x9, #16")
}

// TestMacroErrors is the refusal list: too many arguments, an unknown
// \param, a missing .endm, a nested definition.
func TestMacroErrors(t *testing.T) {
	for _, c := range []struct {
		src  string
		want string
	}{
		{
			src:  ".macro m, a\nnop\n.endm\nm 1, 2\n",
			want: "arguments for",
		},
		{
			src:  ".macro m, a\nb \\x\n.endm\nm 1\n",
			want: "unknown parameter",
		},
		{
			src:  ".macro m, a\nnop\n",
			want: ".endm missing",
		},
		{
			src:  ".macro m, a\n.macro n\nnop\n.endm\n.endm\nm 1\n",
			want: "nested .macro",
		},
	} {
		_, errs := Assemble(c.src, 0)
		require.NotEmpty(t, errs, "case %q", c.src)
		require.Contains(
			t,
			strings.Join(errStrings(errs), "; "),
			c.want,
			"case %q",
			c.src,
		)
	}
}
