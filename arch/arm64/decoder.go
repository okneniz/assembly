package arm64

import (
	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/dtree"
)

// decodeCtor - constructor of a registry entry (decision tree payload).
type decodeCtor = func(word uint32) Instr

// schemaRules - rules of the curated schemas in priority order. The match
// bits come from the generated armISA when the entry is mapped
// (schemaISAEntry, step A): exact matches and sf splits; the remaining
// schemas (families, aliases, soft XML bits) are matched against
// hand-written Mask/Value. An entry without a ctor is a barrier (payload
// nil): the word goes to the isaTail tail, bypassing the lower schemas.
func schemaRules() []dtree.Rule[decodeCtor] {
	schemas := getSchemas()
	rules := make([]dtree.Rule[decodeCtor], 0, len(schemas))

	for i, s := range schemas {
		match, mask := s.Value, s.Mask
		if e := schemaISAEntry(i); e != nil {
			match, mask = e.Match, e.Mask
		}

		rules = append(rules, dtree.NewRule(mask, match, s.ctor))
	}

	return rules
}

// tailRules - rules of the generated isaTail tail (step B): armISA entries
// with names not covered by the schemas, ordered by specificity. The word
// gets a mnemonic and raw fields (Generic) instead of a generic .word. The
// tail goes AFTER the schemas - the curated decode order does not change.
func tailRules() []dtree.Rule[decodeCtor] {
	var rules []dtree.Rule[decodeCtor]

	for k := 0; isaTailEntry(k) != nil; k++ {
		e := isaTailEntry(k)
		rules = append(rules, dtree.NewRule(e.Mask, e.Match, func(w uint32) Instr {
			return decodeGeneric(e, w)
		}))
	}

	return rules
}

// Decision trees over the registry are derivatives of the registry data
// (like the schemas themselves): built once at package load, decoding is a
// bit-by-bit descent instead of a linear scan over ~3.7k entries.
var (
	schemaTree = dtree.New(schemaRules())
	tailTree   = dtree.New(tailRules())
)

// decodeOne decodes a single 32-bit word with the decision trees; the
// registry's first-match is preserved by dtree. A schema match with a ctor
// builds its structure; a missing match or a barrier (schema without a
// ctor) sends the word to the isaTail tail; the unrecognized becomes
// Unknown (.word) so the total line count matches objdump.
func decodeOne(word uint32) Instr {
	if ctor, ok := schemaTree.Lookup(word); ok && ctor != nil {
		return ctor(word)
	}

	if ctor, ok := tailTree.Lookup(word); ok {
		return ctor(word)
	}

	return decodeUnknown(word)
}

// Parse — constructor of a combinator that decodes ARM64 machine code
// (fixed 32-bit width, little-endian) from a parsec buffer into []Instr.
// The word reader is a plain function: unlike bytes.ReadAs (which rebuilds
// its Count combinator on every invocation) it allocates nothing per word.
// Try rolls the position back on a truncated tail (<4 bytes), Many
// swallows its error and drives the loop to the end of the buffer.
// Every 32-bit word reaches the output: unrecognized encodings become a
// .word instruction so the total line count matches objdump. The
// instructions are position-independent: addresses live in the view
// context (disasm), not in the structures.
func MakeDecoder() parsec.Combinator[byte, int, []Instr] {
	instr := func(buf parsec.Buffer[byte, int]) (Instr, parsec.Error[int]) {
		w, err := decodeWordLE(buf)
		if err != nil {
			return nil, err
		}

		return decodeOne(w), nil
	}

	return parsec.Many(0, bytes.Try(instr))
}

// decodeWordLE reads a 4-byte little-endian word; a truncated tail is an
// error (the caller's Try rolls the position back).
func decodeWordLE(buf parsec.Buffer[byte, int]) (uint32, parsec.Error[int]) {
	var w uint32
	for i := range 4 {
		b, err := buf.Read(true)
		if err != nil {
			return 0, parsec.NewParseError(buf.Position(), "arm64: word")
		}

		w |= uint32(b) << (8 * i)
	}

	return w, nil
}
