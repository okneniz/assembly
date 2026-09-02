package loong64

import (
	"encoding/binary"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/dtree"
)

// Parse - constructor of a combinator that decodes LoongArch machine code
// (fixed 32-bit width, little-endian) from a parsec buffer into []Instr.
// The word sub-combinator is built once per assembly; Try rolls the
// position back on a truncated tail (<4 bytes), Many swallows its error
// and drives the loop to the end of the buffer. Every 32-bit word reaches
// the output: unrecognized encodings become a .word instruction so the
// total line count matches objdump.
func Parse(baseAddr uint64) parsec.Combinator[byte, int, []Instr] {
	word := bytes.ReadAs[uint32](4, "loong64: word", binary.LittleEndian)
	instr := func(buf parsec.Buffer[byte, int]) (Instr, parsec.Error[int]) {
		addr := baseAddr + uint64(buf.Position())

		w, err := word(buf)
		if err != nil {
			return nil, err
		}

		return decodeOne(w, addr), nil
	}

	return parsec.Many(0, bytes.Try(instr))
}

// decodeCtor - the constructor of a table entry (the decision-tree payload).
type decodeCtor = func(word uint32, addr uint64) Instr

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

		rules = append(rules, dtree.Rule[decodeCtor]{
			Mask:    mm[1],
			Match:   mm[0],
			Payload: e.ctor,
		})
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
func decodeOne(word uint32, addr uint64) Instr {
	if ctor, ok := decodeTree.Lookup(word); ok {
		return ctor(word, addr)
	}

	return Unknown{base: newBase(addr, word)}
}
