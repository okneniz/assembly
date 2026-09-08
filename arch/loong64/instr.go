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
// numbers, the environment is just the address), into text
// (disasm.ObjDump, for the disassembler). There is NO common type with
// other architectures; each instruction is its own structure
// (AddW, AddiW, Beq, ...).
type Instr interface {
	disasm.ObjDump
	Addr() uint64

	// Encode encodes the computed instruction: pc is the absolute address
	// (PC-relative forms). There is no resolver - an exact inverse of
	// decode, and no modes: LA64 words are fixed 32-bit, no compression.
	Encode(w io.Writer, pc uint64) (int64, error)
}

// base - the bookkeeping fields of every instruction: address, raw word,
// and length. This is not instruction semantics but a bookkeeping record
// for Addr/Len.
type base struct {
	addr   uint64
	raw    uint32
	length int
}

func newBase(addr uint64, raw uint32) base {
	return base{
		addr:   addr,
		raw:    raw,
		length: 4,
	}
}

func (b base) Addr() uint64 {
	return b.addr
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

func (i Unknown) Encode(w io.Writer, _ uint64) (int64, error) {
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
