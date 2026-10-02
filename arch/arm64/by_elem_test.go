package arm64

// Decode words of the by-element family, pinned against clang. The
// fcmla lane-index layouts (.4h — bit 21, .4s — bit 11, .8h — both)
// and the long families' result arrangements (.4s reads .h sources,
// .2d reads .s) are the shapes the property round trips cannot
// distinguish from a self-consistent mis-encode.

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByElemWordsDecode(t *testing.T) {
	for _, c := range []struct {
		word uint32
		text string
	}{
		{0x0f758269, "mul.4h v9, v19, v5[3]"},
		{0x2fb40bdf, "mla.2s v31, v30, v20[3]"},
		{0x4f69c907, "sqdmulh.8h v7, v8, v9[6]"},
		{0x4fd41a69, "fmla.2d v9, v19, v20[1]"},
		{0x6faf9820, "fmulx.4s v0, v1, v15[3]"},
		{0x6f851a69, "fcmla.4s v9, v19, v5[1], #0"},   // the index in bit 11
		{0x2f543269, "fcmla.4h v9, v19, v20[0], #90"}, // 5-bit Vm
		{0x6f655a69, "fcmla.8h v9, v19, v5[3], #180"}, // the index in bits 11+21
		{0x4f45a269, "smull2.4s v9, v19, v5[0]"},      // .8h sources
		{0x0f75aa69, "smull.4s v9, v19, v5[7]"},       // .4h sources, 3-bit index
		{0x4fb4a269, "smull2.2d v9, v19, v20[1]"},     // .4s sources, 5-bit Vm
		{0x4f522820, "smlal2.4s v0, v1, v2[5]"},
		{0x2f99a883, "umull.2d v3, v4, v25[2]"},
		{0x0fa57a69, "sqdmlsl.2d v9, v19, v5[3]"},
	} {
		in, err := DecodeWord(c.word)
		require.NoError(t, err, "case %#x", c.word)
		require.Equal(t, c.text, in.ObjDump(nil), "case %#x", c.word)
	}
}

// TestByElemCtorBounds — the operand constraints, clang-pinned: the
// fcmla arrangements (.2s/.2d are unallocated) and index bounds, the
// long families' source widths (the .8h-result class is unallocated),
// the 4-bit Vm of the .h sources.
func TestByElemCtorBounds(t *testing.T) {
	for _, c := range []struct {
		name string
		call func() error
	}{
		{"fcmla .2s", func() error {
			_, err := New().FcmlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "2s", 0, 0)
			return err
		}},
		{"fcmla .2d", func() error {
			_, err := New().FcmlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "2d", 0, 0)
			return err
		}},
		{"fcmla .4s lane 2", func() error {
			_, err := New().FcmlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "4s", 2, 0)
			return err
		}},
		{"fcmla .4h lane 2", func() error {
			_, err := New().FcmlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "4h", 2, 0)
			return err
		}},
		{"smull .8h result", func() error {
			_, err := New().SmullElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "8h", true, 0)
			return err
		}},
		{"mla .h source above v15", func() error {
			_, err := New().MlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 20), "4h", 0)
			return err
		}},
		{"mla .h lane 8", func() error {
			_, err := New().MlaElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "4h", 8)
			return err
		}},
		{"smull .2d lane 4", func() error {
			_, err := New().SmullElem(vregOf(t, 0), vregOf(t, 1), vregOf(t, 2), "2d", true, 4)
			return err
		}},
	} {
		assertErr(t, c.name, c.call())
	}
}
