package arm64

// Build points and self-verify for the syntax layer (asm/arm64): the
// constructor dispatcher by mnemonic, legacy candidates/handlers (by the
// schemas' decode table) and single-word decoding for text cross-checking.

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
