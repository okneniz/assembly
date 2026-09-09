package loong64

import (
	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/dtree"
)

// Parse - constructor of a combinator that decodes LoongArch machine code
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
			return 0, parsec.NewParseError(buf.Position(), "loong64: word")
		}

		w |= uint32(b) << (8 * i)
	}

	return w, nil
}

// decodeCtor - the constructor of a table entry (the decision-tree payload).
type decodeCtor = func(word uint32) Instr

// decodeRules - the decodeTable rules in priority order; match/mask are
// authoritative - from the generated loongEncodings (loongarch-opcodes
// tables). The table mnemonics are the loongEncodings keys as-is.
func decodeRules() []dtree.Rule[decodeCtor] {
	rules := make([]dtree.Rule[decodeCtor], 0, len(decodeTable))

	for _, e := range decodeTable {
		mm, ok := loongEncodings[e.name]
		if !ok {
			continue // entries without an encoding are unreachable
		}

		rules = append(rules, dtree.NewRule(mm[1], mm[0], e.ctor))
	}

	return rules
}

// decodeTree - the decision tree over the table: a derivative of the
// registry data, built once at package load; decoding is a descent over
// bits instead of a linear scan of the table.
var decodeTree = dtree.New(decodeRules())

// decodeOne picks a 32-bit instruction with the decision tree; the
// first-match order of decodeTable (order = priority, for example
// csrrd/csrwr before the csrxchg encoding they overlap) is preserved by
// dtree. The unrecognized - Unknown.
func decodeOne(word uint32) Instr {
	if ctor, ok := decodeTree.Lookup(word); ok {
		return ctor(word)
	}

	return Unknown{base: newBase(word)}
}
