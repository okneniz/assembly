package arm64

import "strings"

// PSTATE fields of the MSR (immediate) space: the <pstatefield> operand
// of the real MSR_imm form - the alias spelling (msr daifset, #1) lives
// in the text layer; the field selects op1/op2 baked into the word base,
// the immediate rides CRm [11:8]. The svcr* fields are out: smstart/
// smstop own those words.

// pstateField - one PSTATE field: the word base (op1/op2 inside, CRm
// zeroed) and whether the immediate is a single bit.
type pstateField struct {
	base uint32
	imm1 bool
}

// pstateFields - the ps-field spellings of the MSR (immediate) decode
// case, lower-case keys (the lookup is case-insensitive, as the sysreg
// one).
var pstateFields = map[string]pstateField{
	"spsel":   {base: 0xD50040BF, imm1: false},
	"uao":     {base: 0xD500407F, imm1: false},
	"pan":     {base: 0xD500409F, imm1: false},
	"allint":  {base: 0xD501401F, imm1: true},
	"pm":      {base: 0xD501421F, imm1: true},
	"ssbs":    {base: 0xD503403F, imm1: false},
	"dit":     {base: 0xD503405F, imm1: false},
	"tco":     {base: 0xD503409F, imm1: false},
	"daifset": {base: 0xD50340DF, imm1: false},
	"daifclr": {base: 0xD50340FF, imm1: false},
}

// PstateBase - the word base of the PSTATE field spelling and whether
// its immediate is a single bit (allint/pm; the rest take the CRm
// nibble 0..15).
func PstateBase(sym string) (base uint32, imm1 bool, ok bool) {
	f, ok := pstateLookup(sym)
	return f.base, f.imm1, ok
}

// pstateLookup - the field spelling, case-insensitive.
func pstateLookup(sym string) (pstateField, bool) {
	f, ok := pstateFields[strings.ToLower(sym)]
	return f, ok
}
