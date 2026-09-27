package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBarrierBuild - the words of the typed ctors over the domains
// (the CRm nibble at 11:8; llvm-pinned).
func TestBarrierBuild(t *testing.T) {
	var b Builder

	cases := []struct {
		name string
		in   func() (Instr, error)
		word uint32
	}{
		{"dmb sy", func() (Instr, error) { return b.Dmb(Sy) }, 0xd5033fbf},
		{"dmb st", func() (Instr, error) { return b.Dmb(St) }, 0xd5033ebf},
		{"dmb ish", func() (Instr, error) { return b.Dmb(Ish) }, 0xd5033bbf},
		{"dmb ishst", func() (Instr, error) { return b.Dmb(Ishst) }, 0xd5033abf},
		{"dmb nsh", func() (Instr, error) { return b.Dmb(Nsh) }, 0xd5033bbf - 0x400}, // CRm 7
		{"dmb oshst", func() (Instr, error) { return b.Dmb(Oshst) }, 0xd50332bf},
		{"dsb sy", func() (Instr, error) { return b.Dsb(Sy) }, 0xd5033f9f},
		{"dsb ish", func() (Instr, error) { return b.Dsb(Ish) }, 0xd5033b9f},
		{"isb", func() (Instr, error) { return b.Isb() }, 0xd5033fdf},
	}
	for _, c := range cases {
		in, err := c.in()
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	// a cast domain outside the eight is refused (the constructor, not
	// the type, is the guard - BarrierDomain is a nibble-shaped number)
	_, err := b.Dmb(BarrierDomain(0x0))
	assertErr(t, "dmb unassigned domain", err)
	_, err = b.Dsb(BarrierDomain(0x5))
	assertErr(t, "dsb unassigned domain", err)
}

// TestBarrierDecode - the CRm nibble back to the option; an unassigned
// nibble is not the instruction.
func TestBarrierDecode(t *testing.T) {
	dmb := decodeBarrierOf("dmb")

	cases := []struct {
		word   uint32
		domain BarrierDomain
	}{
		{0xd5033fbf, Sy},
		{0xd5033ebf, St},
		{0xd50332bf, Oshst},
	}
	for _, c := range cases {
		in, err := dmb(c.word)
		require.NoError(t, err, "case %#x", c.word)
		b, ok := in.(Barrier)
		require.True(t, ok, "type = %T, want Barrier", in)
		require.Equal(t, c.domain, b.domain, "case %#x", c.word)
		require.Equal(t, uint32(0xd50330bf), buildWord(t, in)&^0xf00, "case %#x", c.word)
	}

	_, err := dmb(0xd50331bf) // CRm 1 - unassigned
	require.Error(t, err)

}
