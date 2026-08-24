package arm64

// Build points and self-verify for the syntax layer (asm/arm64): the
// constructor dispatcher by mnemonic, legacy candidates/handlers (by the
// schemas' decode table) and single-word decoding for text cross-checking.

import "fmt"

// BuildInstr builds a computed instruction by mnemonic and operands
// (the armCtors dispatcher). Errors — in the form of encoding paths.
func BuildInstr(mnem string, ops []VOp) (Instr, error) {
	ctor, ok := armCtors[mnem]
	if !ok {
		return nil, fmt.Errorf("unknown instruction %q", mnem)
	}

	st, err := ctor(ops)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mnem, err)
	}

	return st, nil
}

// LegacyCandidates — legacy candidates for a mnemonic (the decode-table
// order and the alias map; encodeARM tries them when the constructors fail).
func LegacyCandidates(mnem string) []*Schema {
	return candidatesFor(mnem)
}

// BuildLegacy — legacy encoding of a candidate from computed operands:
// format-family handlers + word assembly via inverse transforms.
// pc — the instruction address (branches).
func BuildLegacy(s *Schema, mnem string, ops []VOp, pc uint64) (uint32, error) {
	fields, err := armFieldsFor(s, resolvedInstr{mnem: mnem, ops: ops}, ctx{Addr: pc})
	if err != nil {
		return 0, err
	}

	return packFields(s, fields)
}

// DecodeWord — decode a single word (self-verify of encodeARM: the decoder
// must reproduce the source text).
func DecodeWord(w uint32, addr uint64) Instr {
	return decodeOne(w, addr)
}

// Schemas — the ordered decode table (first-match-wins; for testing
// tools).
func Schemas() []Schema {
	return getSchemas()
}

// IsKeywordSym — a name operand from the keyword tables (conditions,
// system registers, prfm hints): it stays a name instead of being computed
// by the resolver.
func IsKeywordSym(sym string) bool {
	return isKeyword(sym)
}
