package arm64

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/okneniz/assembly/disasm"
)

// TestSimdDupElem - the vector-source Advanced SIMD copy forms and the
// indexed GPR-side forms, pinned against llvm-mc words (docker
// assembly-tests, 2026-09-07). These used to fall through to the tail
// (dup (element) / the scalar mov.d alias printed as bare Generic) or
// even to .word - decodeSimdCopy's imm5 guard wrongly required a
// one-hot imm5, rejecting every ins/smov/umov with a nonzero lane.
func TestSimdDupElem(t *testing.T) {
	for _, c := range []struct {
		word uint32
		text string
	}{
		// DupElem: dup (element)
		{0x4E040420, "dup.4s v0, v1[0]"},  // llvm: dup v0.4s, v1.s[0]
		{0x4E180420, "dup.2d v0, v1[1]"},  // llvm: dup v0.2d, v1.d[1]
		{0x4E070420, "dup.16b v0, v1[3]"}, // llvm: dup v0.16b, v1.b[3]
		// DupElem: ins (element)
		{0x6E180420, "ins.d v0[1], v1[0]"}, // llvm: ins v0.d[1], v1.d[0]
		// DupElem: the scalar dup alias
		{0x5E080420, "mov.d v0, v1"}, // llvm: mov.d v0, v1
		{0x5E040420, "mov.s v0, v1"}, // llvm: mov.s v0, v1
		// SimdCopy: indexed GPR-side forms (used to be .word)
		{0x4E0C2C20, "smov x0, v1.s[1]"},
		{0x0E072C20, "smov w0, v1.b[3]"},
		{0x0E073C20, "umov w0, v1.b[3]"},
		{0x4E071C20, "mov.b v0[3], w1"}, // llvm: ins v0.b[3], w1
		// the mov print alias covers the filling sizes only: s into w
		// (Q=0) and d into x (Q=1); b/h stay umov (llvm-mc words)
		{0x0E0C3C20, "mov w0, v1.s[1]"},  // input umov w0, v1.s[1]
		{0x4E183C20, "mov x0, v1.d[1]"},  // input umov/mov x0, v1.d[1]
		{0x0E1E3C20, "umov w0, v1.h[7]"}, // no alias below esize 32
	} {
		in := decodeOne(c.word)
		if got := in.ObjDump(disasm.DefaultViewCtx()); got != c.text {
			t.Errorf("%#010x: got %q, want %q", c.word, got, c.text)
		}

		var buf bytes.Buffer
		if _, err := in.Encode(&buf); err != nil || buf.Len() != 4 {
			t.Errorf("%#010x: Encode: %v (%d bytes)", c.word, err, buf.Len())
			continue
		}

		if back := binary.LittleEndian.Uint32(buf.Bytes()); back != c.word {
			t.Errorf("%#010x: encode round-trip gave %#010x", c.word, back)
		}
	}
}

// TestSimdDupElemInvalid - the imm5 forms the schema mask cannot express
// (unencodable sizes) must stay at the .word fallback, not decode as
// garbage instructions.
func TestSimdDupElemInvalid(t *testing.T) {
	for _, w := range []uint32{
		0x4E100420, // ctz(imm5)=4: no such size
		0x5E300420, // scalar imm5 not one-hot
		0x0E030C20, // dup (general) imm5 not one-hot
	} {
		if _, ok := decodeOne(w).(Unknown); !ok {
			t.Errorf("%#010x: expected the .word fallback, got %T", w, decodeOne(w))
		}
	}
}
