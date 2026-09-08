package loong64

import (
	"encoding/binary"
	"io"
	"strconv"

	"github.com/okneniz/assembly/disasm"
)

// Name - the architecture as consumers see it (test/diff).
const Name = "LOONG64"

// Instr - a LoongArch instruction: it represents itself - encodes itself
// into bytes (Encode, the computed form: operand values are already
// numbers, PC-relative ones are offsets - there is no environment), into
// text (disasm.ObjDump, for the disassembler; absolute targets are
// computed from the address in the view context). There is NO common
// type with other architectures; each instruction is its own structure
// (AddW, AddiW, Beq, ...).
type Instr interface {
	disasm.ObjDump

	// Encode encodes the computed instruction. There is no resolver and
	// no environment - an exact inverse of decode, and no modes: LA64
	// words are fixed 32-bit, no compression.
	Encode(w io.Writer) (int64, error)
}

// base - the bookkeeping fields of every instruction: raw word and
// length. This is not instruction semantics but a bookkeeping record
// for Len.
type base struct {
	raw    uint32
	length int
}

func newBase(raw uint32) base {
	return base{
		raw:    raw,
		length: 4,
	}
}

func (b base) Len() int {
	return b.length
}

// writeWord - writing an encoded instruction (LE).
func writeWord(w io.Writer, word uint32) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint32(nil, word))
	return int64(n), err
}

// Unknown - a word matching no encoding (the ".word" fallback):
// displayed as <unknown> (in angle brackets - objdiff comparison
// treats it as an annotation), Encode re-emits the raw word.
type Unknown struct {
	base
}

func newUnknown(base_ base) Unknown {
	return Unknown{base: base_}
}

func (i Unknown) ObjDump(_ disasm.ViewCtx) string {
	return "<unknown>"
}

func (i Unknown) Encode(w io.Writer) (int64, error) {
	return writeWord(w, i.raw)
}

// imm - the value of an immediate operand: a concrete number. Symbolic
// slots are evaluated by the syntax layer (asm/loong64) BEFORE the
// structure is built - a computed instruction contains no holes. Branch
// and jump targets are absolute addresses.
type imm struct {
	val int64
}

// immNum - a concrete value (decoding and construction).
func immNum(v int64) imm {
	return imm{val: v}
}

// text - the %d output (for ObjDump of decoded instructions: LoongArch
// objdump prints immediates and branch targets in decimal).
func (m imm) text() string {
	return strconv.FormatInt(m.val, 10)
}
