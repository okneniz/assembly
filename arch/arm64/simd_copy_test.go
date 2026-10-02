package arm64

// Decode words of the SIMD copy family (dup/ins/smov/umov), pinned
// against clang. The mov spellings are the canonical print aliases:
// INS general (mov.sz vd[idx], wn), the umov fill forms and the scalar
// dup (mov <b|h|s|d>n, vn.sz[idx]).

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimdCopyWordsDecode(t *testing.T) {
	for _, c := range []struct {
		word uint32
		text string
	}{
		{0x4e0b1e23, "mov.b v3[5], w17"},    // ins (general)
		{0x6e052e69, "ins.b v9[2], v19[5]"}, // ins (element)
		{0x0e010669, "dup.8b v9, v19[0]"},   // dup (element)
		{0x0e040669, "dup.2s v9, v19[0]"},
		{0x0e043e65, "mov w5, v19.s[0]"}, // umov, the fill alias
		{0x4e042e65, "smov x5, v19.s[0]"},
		{0x0e012e65, "smov w5, v19.b[0]"},
		{0x4e0e2e65, "smov x5, v19.h[3]"},
		{0x5e040669, "mov s9, v19.s[0]"}, // dup (scalar)
		{0x5e070669, "mov b9, v19.b[3]"},
		{0x5e060669, "mov h9, v19.h[1]"},
		{0x5e180669, "mov d9, v19.d[1]"},
	} {
		in, err := DecodeWord(c.word)
		require.NoError(t, err, "case %#x", c.word)
		require.Equal(t, c.text, in.ObjDump(nil), "case %#x", c.word)
	}
}

// TestSimdCopyUnallocated — the classes the ISA does not allocate
// (clang refuses the spelling, the word decodes as unknown): smov .s
// into a w register, the scalar dup with Q=0.
func TestSimdCopyUnallocated(t *testing.T) {
	for _, w := range []uint32{0x0e042e65} {
		in, err := DecodeWord(w)
		require.NoError(t, err, "case %#x", w)
		require.IsType(t, Unknown{}, in, "case %#x", w)
	}
}

// TestSimdCopyCtorBounds — the operand constraints, clang-pinned: the
// ins/umov width rules, the smov .s destination, the dup scalar index.
func TestSimdCopyCtorBounds(t *testing.T) {
	for _, c := range []struct {
		name string
		call func() error
	}{
		{"smov .s into w", func() error {
			_, err := New().Smov(wreg(t, 5), vregOf(t, 19), "s", 0)
			return err
		}},
		{"ins .b into x", func() error {
			_, err := New().Ins(vregOf(t, 3), 5, xreg(t, 0), "b")
			return err
		}},
		{"umov .b into x", func() error {
			_, err := New().Umov(xreg(t, 5), vregOf(t, 19), "b", 0)
			return err
		}},
		{"dup scalar index out of range", func() error {
			_, err := New().DupScalar(vregOf(t, 9), vregOf(t, 19), "s", 4)
			return err
		}},
	} {
		assertErr(t, c.name, c.call())
	}
}

func vregOf(t *testing.T, n int) VReg {
	t.Helper()
	r, err := V(n)
	require.NoError(t, err)
	return r
}
