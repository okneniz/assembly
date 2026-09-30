package asm

// SourceUnit - a parsed .S source of the unit mode, ready to deposit:
// Deposit lays its text and data sections into a unit, DepositBss its
// zero-fill ones. The split is the linker's - every source of a link
// deposits its sections first and its reserves after, so the bss of the
// whole program aggregates into the unit's single tail; a single-source
// caller takes AssembleUnit instead and never sees the two phases.

import (
	"github.com/okneniz/assembly/unit"
)

// SourceUnit is the interpreted source between parsing and depositing:
// its sections, labels, pools and sets are final, nothing resolves until
// the receiving unit's phase boundary.
type SourceUnit struct {
	src *unitSource
}

// ParseSourceUnit parses and interprets a .S source (the unit-mode
// sibling of Assemble): the statements are walked, the records sized
// under the placeholder environment, the labels, pools and sets
// recorded. file names the origin in the unit's line map and isolates
// the .L locals of the source when several of them link into one unit.
func ParseSourceUnit(file, src string, be Syntax) (*SourceUnit, []AsmError) {
	be.ResetOptions()
	stmts := parseSource([]rune(src), be)

	us, errs := buildUnitSource(stmts, be)
	if len(errs) > 0 {
		return nil, errs
	}

	us.file = file
	return &SourceUnit{src: us}, nil
}

// Deposit lays the text and data sections of the source into the unit,
// in first-appearance order: the named labels join the unit's namespace
// through the source's name scope, the record runs between them deposit
// as one deferred record each.
func (s *SourceUnit) Deposit(u *unit.Unit) {
	depositUnit(u, s.src, false)
}

// DepositBss lays the NOBITS sections of the source into the unit's
// data stream as zero-fill reserves - after every source's Deposit, the
// bss of the whole program is the unit's one tail.
func (s *SourceUnit) DepositBss(u *unit.Unit) {
	depositUnit(u, s.src, true)
}
