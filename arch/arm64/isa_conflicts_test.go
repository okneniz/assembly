package arm64

import (
	"testing"

	"github.com/okneniz/assembly/disasm"
)

// TestISAConflictFixes pins the five hand-written schemas that used to
// conflict with the XML-generated armISA table (TestISAAudit used to report
// conflict=5: cmeq/sqrshl/abs/fmax/fmin). Per-case verdict, llvm-mc-verified
// (docker assembly-tests, 2026-09-06): the XML was right everywhere; the
// hand-written constants were transcriptions of OTHER instructions:
//
//	cmeq   0x0E208C00 → cmtst (vector register: cmeq is U=1, cmtst U=0)
//	sqrshl 0x0E203C00 → cmge  (vector register; real sqrshl is 0x0E205C00)
//	abs    0x0E20C000 → a smull-shaped word with bit23 pinned 0 - a phantom
//	                    matching nothing real (real abs is 0x0E20B800)
//	fmax   0x1E600C00 → fcsel d0, d1, d2, eq
//	fmin   0x1E601C00 → fcsel d0, d1, d2, ne
//
// The words below are exactly what llvm-mc assembles these texts into.
func TestISAConflictFixes(t *testing.T) {
	for _, c := range []struct {
		word uint32
		text string
	}{
		{0x6EA28C20, "cmeq.4s v0, v1, v2"},
		{0x4EA28C20, "cmtst.4s v0, v1, v2"},
		{0x4EA25C20, "sqrshl.4s v0, v1, v2"},
		{0x4EA0B820, "abs.4s v0, v1"},
		{0x1E624820, "fmax d0, d1, d2"},
		{0x1E224820, "fmax s0, s1, s2"},
		{0x1E625820, "fmin d0, d1, d2"},
		{0x1E225820, "fmin s0, s1, s2"},
	} {
		in, derr := decodeOne(c.word)
		_ = derr
		got := in.ObjDump(disasm.DefaultViewCtx())
		if got != c.text {
			t.Errorf("%#010x: got %q, want %q", c.word, got, c.text)
		}
	}
}
