package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// TestSysOpBuild - every operation spelling of the table builds and
// encodes its word (the Rt forms join x3); the wrong-class spellings
// and register shapes are refused by the constructor.
func TestSysOpBuild(t *testing.T) {
	for name, rec := range sysOps {
		rt := ""
		want := rec.enc
		if rec.hasRt {
			rt = "x3"
			want |= 3
		}

		in, err := SysOpOf(rec.mnem, name, rt)
		require.NoError(t, err, "case %s %s", rec.mnem, name)
		require.Equal(t, want, buildWord(t, in), "case %s %s", rec.mnem, name)
	}
}

// TestSysOpBuildErrors - the constructor guards: an unknown spelling, a
// spelling of another class, a 32-bit or sp register, a missing and an
// extra register operand.
func TestSysOpBuildErrors(t *testing.T) {
	cases := []struct {
		mnem string
		op   string
		rt   string
		why  string
	}{
		{"ic", "bogus", "", "unknown op"},
		{"tlbi", "zva", "x0", "a dc spelling under tlbi"},
		{"ic", "ivau", "w3", "a 32-bit rt"},
		{"ic", "ivau", "sp", "sp as rt"},
		{"ic", "ivau", "", "a missing rt"},
		{"tlbi", "alle1", "x0", "an extra rt"},
		{"dc", "zva", "x31", "a bad register name"},
	}

	for _, c := range cases {
		_, err := SysOpOf(c.mnem, c.op, c.rt)
		require.Error(t, err, "case %q", c.why)
	}
}

// TestSysOpDecode - every table word decodes back to its spelling (the
// Rt-bearing ones at x3); a word of the class the table does not know
// keeps the honest raw-field print.
func TestSysOpDecode(t *testing.T) {
	for name, rec := range sysOps {
		w := rec.enc
		want := rec.mnem + " " + name
		if rec.hasRt {
			w |= 3
			want += ", x3"
		}

		in, err := DecodeWord(w)
		require.NoError(t, err, "case %s", name)
		require.Equal(t, want, instrText(t, in), "case %s", name)
	}

	cases := []struct {
		word uint32
		text string
	}{
		// ic iallu with Rt != 31: not the operandless word, not any
		// other spelling - the class print with raw fields
		{0xd508752f, "dc op1=0x0 CRm=0x5 op2=0x1 Rt=0xf"},
		// a tlbi CRm/op2 pair no operation owns
		{0xd50c859f, "tlbi op1=0x4 CRm=0x5 op2=0x4 Rt=0x1f"},
	}

	for _, c := range cases {
		in, err := DecodeWord(c.word)
		require.NoError(t, err, "case %q", c.text)
		require.Equal(t, c.text, instrText(t, in))
	}
}

// TestSysOpSpellingUniqueness - the words of the table are pairwise
// distinct: the decode's reverse lookup relies on it.
func TestSysOpSpellingUniqueness(t *testing.T) {
	seen := make(map[uint32]string, len(sysOps))
	for name, rec := range sysOps {
		prev, dup := seen[rec.enc]
		require.False(t, dup, "words of %s and %s collide", name, prev)
		seen[rec.enc] = name
	}
}

// instrText - the instruction's own ObjDump text at address zero.
func instrText(t *testing.T, in Instr) string {
	t.Helper()

	return in.ObjDump(disasm.ViewCtxAt(0))
}
