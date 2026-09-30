package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	asm "github.com/okneniz/assembly/asm"
)

// TestSysOpAssembleWords - the IC/DC/TLBI operation spellings and the
// eret family through the text layer; every word is clang-pinned
// (cross-checked on the host: clang -c + objdump).
func TestSysOpAssembleWords(t *testing.T) {
	cases := []struct {
		src  string
		word uint32
	}{
		{"ic iallu", 0xd508751f},
		{"ic ialluis", 0xd508711f},
		{"ic ivau, x0", 0xd50b7520},
		{"ic ivau, x5", 0xd50b7525},
		{"dc zva, x0", 0xd50b7420},
		{"dc zva, x7", 0xd50b7427},
		{"dc ivac, x3", 0xd5087623},
		{"dc isw, x3", 0xd5087643},
		{"dc cvac, x3", 0xd50b7a23},
		{"dc cvau, x3", 0xd50b7b23},
		{"dc civac, x3", 0xd50b7e23},
		{"dc cisw, x3", 0xd5087e43},
		{"dc cvadp, x3", 0xd50b7d23},
		{"tlbi alle1", 0xd50c879f},
		{"tlbi alle1is", 0xd50c839f},
		{"tlbi alle1os", 0xd50c819f},
		{"tlbi alle2", 0xd50c871f},
		{"tlbi alle2is", 0xd50c831f},
		{"tlbi alle2os", 0xd50c811f},
		{"tlbi alle3", 0xd50e871f},
		{"tlbi alle3is", 0xd50e831f},
		{"tlbi alle3os", 0xd50e811f},
		{"tlbi vmalle1", 0xd508871f},
		{"tlbi vmalle1is", 0xd508831f},
		{"tlbi vmalle1os", 0xd508811f},
		{"tlbi vmalls12e1", 0xd50c87df},
		{"tlbi vmalls12e1is", 0xd50c83df},
		{"tlbi vmalls12e1os", 0xd50c81df},
		{"tlbi vae1, x3", 0xd5088723},
		{"tlbi vae1is, x3", 0xd5088323},
		{"tlbi vae1os, x3", 0xd5088123},
		{"tlbi rvae1, x3", 0xd5088623},
		{"tlbi rvae1is, x3", 0xd5088223},
		{"tlbi rvae1os, x3", 0xd5088523},
		{"tlbi aside1, x3", 0xd5088743},
		{"tlbi aside1is, x3", 0xd5088343},
		{"tlbi aside1os, x3", 0xd5088143},
		{"tlbi vae2, x3", 0xd50c8723},
		{"tlbi vae2is, x3", 0xd50c8323},
		{"tlbi vae2os, x3", 0xd50c8123},
		{"tlbi rvae2, x3", 0xd50c8623},
		{"tlbi rvae2is, x3", 0xd50c8223},
		{"tlbi rvae2os, x3", 0xd50c8523},
		{"tlbi vae3, x3", 0xd50e8723},
		{"tlbi vae3is, x3", 0xd50e8323},
		{"tlbi vae3os, x3", 0xd50e8123},
		{"tlbi rvae3, x3", 0xd50e8623},
		{"tlbi rvae3is, x3", 0xd50e8223},
		{"tlbi rvae3os, x3", 0xd50e8523},
		{"tlbi ipas2e1, x3", 0xd50c8423},
		{"tlbi ipas2e1is, x3", 0xd50c8023},
		{"tlbi ipas2e1os, x3", 0xd50c8403},
		{"tlbi ripas2e1, x3", 0xd50c8443},
		{"tlbi ripas2e1is, x3", 0xd50c8043},
		{"tlbi ripas2e1os, x3", 0xd50c8463},
		{"eret", 0xd69f03e0},
		{"eretaa", 0xd69f0bff},
		{"eretab", 0xd69f0fff},
		{"at s1e1r, x5", 0xd5087805},
		{"at s1e1w, x5", 0xd5087825},
		{"at s1e0r, x5", 0xd5087845},
		{"at s1e0w, x5", 0xd5087865},
		{"at s1e1rp, x5", 0xd5087905},
		{"at s1e1wp, x5", 0xd5087925},
		{"at s12e1r, x5", 0xd50c7885},
		{"at s12e1w, x5", 0xd50c78a5},
		{"at s12e0r, x0", 0xd50c78c0},
		{"at s1e2r, x0", 0xd50c7800},
		{"at s1e2w, x5", 0xd50c7825},
		{"at s1e3r, x5", 0xd50e7805},
		{"at s1e3w, x5", 0xd50e7825},
		{"msr sp_el0, x0", 0xd5184100},
		{"mrs x0, sp_el0", 0xd5384100},
		{"msr sp_el1, x0", 0xd51c4100},
		{"mrs x4, sp_el2", 0xd53e4104},
	}

	for _, c := range cases {
		require.Equal(t, c.word, armAssembleOne(t, c.src+"\n", 0), "case %q", c.src)
	}
}

// TestSysOpAssembleErrors - the refused shapes: an unknown operation, a
// register on an operandless one, a missing register, 32-bit and sp
// registers, and an eret with operands.
func TestSysOpAssembleErrors(t *testing.T) {
	cases := []string{
		"ic bogus",
		"tlbi zva",
		"dc zva",
		"ic iallu, x0",
		"tlbi alle1, x0",
		"ic ivau, w5",
		"dc zva, sp",
		"tlbi vae1, x3, x4",
		"eret x0",
		"eretaa #1",
		"at bogus, x0",
		"at s1e2r",
		"at s1e2r, w5",
		"ats1e2r x0", // the fused spelling: refused by gas and clang too
	}

	for _, src := range cases {
		_, errs := asm.Assemble(src+"\n", 0, New())
		require.NotEmpty(t, errs, "case %q must not assemble", src)
	}
}
