package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoongEncodingsShape - a drift guard over the generated table: the
// scalar integer scope is ~250 encodings, and the spot values below are
// llvm-mc-verified (alsl.w encodes its ui3 operand shifted by one;
// csrrd/csrwr are the rj-fixed specializations of csrxchg).
func TestLoongEncodingsShape(t *testing.T) {
	require.Greater(t, len(loongEncodings), 240)

	for name, want := range map[string][2]uint32{
		"addi.w":  {0x02800000, 0xffc00000},
		"alsl.w":  {0x00040000, 0xfffe0000},
		"csrrd":   {0x04000000, 0xff0003e0},
		"csrwr":   {0x04000020, 0xff0003e0},
		"csrxchg": {0x04000000, 0xff000000},
		"lu52i.d": {0x03000000, 0xffc00000},
	} {
		require.Equal(t, want, loongEncodings[name], name)
	}
}

// TestDecodeOneUnknown - a word matching no table entry falls back to
// Unknown; decodeRules carries the joined entries (one per table row).
func TestDecodeOneUnknown(t *testing.T) {
	require.Len(t, decodeRules(), len(decodeTable))

	in := decodeOne(0xffffffff, 0x90000000)
	_, ok := in.(Unknown)
	require.True(t, ok, "type = %T, want Unknown", in)
	require.Equal(t, uint64(0x90000000), in.Addr())
}

// TestDecodeRulesSkipsUnlinked - a table row without a loongEncodings
// entry yields no rule (such rows are unreachable by construction; the
// guard keeps a typo'd mnemonic from matching everything). The row is
// appended for the test only - decodeTree stays built from the real table.
func TestDecodeRulesSkipsUnlinked(t *testing.T) {
	decodeTable = append(decodeTable, newDecodeEntry("no.such.mnemonic", decodeAddW))
	t.Cleanup(func() {
		decodeTable = decodeTable[:len(decodeTable)-1]
	})

	require.Len(t, decodeRules(), len(decodeTable)-1)
}
