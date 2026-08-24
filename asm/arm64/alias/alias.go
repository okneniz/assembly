// Package alias provides ARM64 instruction aliases (mnemonics without
// their own encoding): a layer ABOVE arch/arm64, the real instructions.
// An alias is purely assembler syntax mapped 1:1 (same word, same
// length): cmp x0, #1 is the encoding of subs xzr, x0, #1; cset is
// csinc with an inverted condition; sxtb/ubfiz are SBFM/UBFM with
// computed immr/imms. Alias constructors build the arch base
// instruction structures (arch.*Of/family constructor machinery) and go
// through the same self-verify encodeARM.
//
// Wiring is done by injecting the constructors into the syntax layer
// backend: NewASMBackend() = arm64.NewWithCtors(aliasCtors). Pure
// arm64.New() remains without alias constructors (the legacy
// armAliasMap path still handles aliases: it is the fallback for
// symbolic immediates and does not depend on the constructors).
package alias

import (
	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/asm"
	arm64 "github.com/okneniz/assembly/asm/arm64"
)

// NewASMBackend returns the ARM64 Syntax with aliases on top of the
// syntax layer (asm/arm64) for asm.Assemble.
func NewASMBackend() asm.Syntax {
	return arm64.NewWithCtors(aliasCtors)
}

// Assemble is the full ARM64 assembly (syntax layer + aliases) by the
// asm core.
func Assemble(src string, base uint64) (*asm.Result, []asm.AsmError) {
	return asm.Assemble(src, base, NewASMBackend())
}

// aliasCtors maps mnemonics to alias constructors.
var aliasCtors = map[string]arch.ArmCtor{
	// add/sub family: cmp/cmn (Rd = zr), neg/negs (Rn = zr)
	"cmp": newCmp("subs"), "cmn": newCmp("adds"),
	"neg": newNeg("sub"), "negs": newNeg("subs"),
	// logical family: tst (ands Rd=zr), mvn (orn Rn=zr), mov
	"tst": newTst, "mvn": newMvn, "mov": newMov,
	// madd/msub: mul/mneg (ra = zr)
	"mul": newMul, "mneg": newMneg,
	// csel family: inverted condition
	"cset": newCset, "csetm": newCsetm, "cinc": newCinc, "cinv": newCinv, "cneg": newCneg,
	// bitfield: sxt* and *bfiz/*bfx (lsb+width -> immr/imms)
	"sxtb": newSxt("sxtb", 7), "sxth": newSxt("sxth", 15), "sxtw": newSxt("sxtw", 31),
	"ubfiz": newBf("ubfiz", true, true), "ubfx": newBf("ubfx", true, false),
	"sbfiz": newBf("sbfiz", false, true), "sbfx": newBf("sbfx", false, false),
}
