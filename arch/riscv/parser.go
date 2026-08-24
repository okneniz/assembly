package riscv

import (
	"encoding/binary"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"

	"github.com/okneniz/assembly/dtree"
)

// Parse - builds a combinator decoding variable-length RISC-V machine code
// (16-bit compressed RVC or 32-bit) from a parsec buffer
// into []Instr. The halfword sub-combinator is assembled once per build;
// Try rolls the position back on a truncated tail (a halfword shorter than
// 2 bytes or the cut-off high halfword of a 32-bit instruction), Many
// swallows its error and drives the loop to the end of the buffer.
func Parse(baseAddr uint64) common.Combinator[byte, int, []Instr] {
	half := bytes.ReadAs[uint16](2, "riscv: halfword", binary.LittleEndian)
	instr := func(buf common.Buffer[byte, int]) (Instr, common.Error[int]) {
		addr := baseAddr + uint64(buf.Position())

		lo, err := half(buf)
		if err != nil {
			return nil, err
		}

		if lo&0x3 != 0x3 {
			return compressedInstruction(uint32(lo), addr), nil
		}

		hi, err := half(buf)
		if err != nil {
			return nil, err
		}

		return decodeOne(uint32(lo)|uint32(hi)<<16, addr), nil
	}

	return common.Many(0, bytes.Try(instr))
}

// decodeCtor - the constructor of a table entry (the decision-tree payload).
type decodeCtor = func(word uint32, addr uint64) Instr

// decodeRules - the decodeTable rules in priority order; match/mask
// are authoritative - from the generated riscvEncodings (Spike encoding.h).
func decodeRules() []dtree.Rule[decodeCtor] {
	rules := make([]dtree.Rule[decodeCtor], 0, len(decodeTable))

	for _, e := range decodeTable {
		mm, ok := riscvEncodings[encName(e.name)]
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
// first-match order of decodeTable (order = priority, for example OP-IMM
// shifts before OP-IMM arithmetic) is preserved by dtree. The
// unrecognized - Unknown.
func decodeOne(word uint32, addr uint64) Instr {
	if ctor, ok := decodeTree.Lookup(word); ok {
		return ctor(word, addr)
	}

	return Unknown{base: newBase(addr, word)}
}
