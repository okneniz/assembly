package riscv

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Name - the architecture as consumers see it (test/diff).
const Name = "RISCV64"

// Instr - a RISC-V instruction: it represents itself - encodes itself into
// bytes (Encode, the computed form: operand values are already numbers,
// PC-relative ones are offsets; o is only the modes - there is no address
// environment) and into text (disasm.ObjDump, for the disassembler;
// absolute targets are computed from the address in the view context).
// There is NO common type with other architectures; each instruction is
// its own structure (Lui, Addi, Beq, ...).
type Instr interface {
	disasm.ObjDump

	// Encode encodes the computed instruction; o is the modes (NoRVC).
	// There is no resolver - an exact inverse of decode.
	Encode(w io.Writer, o EncOpts) (int64, error)
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

// imm - the value of an immediate operand: a concrete number. Symbolic
// slots are evaluated by the syntax layer (asm/riscv) BEFORE the structure
// is built - a computed instruction contains no holes.
type imm struct {
	val int64
}

// immNum - a concrete value (decoding and construction).
func immNum(v int64) imm {
	return imm{val: v}
}

// text - the %#x output (for ObjDump of decoded instructions).
func (m imm) text() string {
	return fmt.Sprintf("%#x", m.val)
}

// regBits - the register number by ABI/xN name (for encoding). An invalid
// name is already rejected by the parser-side constructor (wantReg).
func regBits(name string) uint32 {
	n, err := rvRegNumOf(name, false)
	if err != nil {
		return 0
	}

	return n
}

// fregBits - the floating register number (f operands of flw/fld/fsw/fsd).
func fregBits(name string) uint32 {
	n, err := rvRegNumOf(name, true)
	if err != nil {
		return 0
	}

	return n
}

// writeWord/writeHalf - writing an encoded instruction (LE).
func writeWord(w io.Writer, word uint32) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint32(nil, word))
	return int64(n), err
}

func writeHalf(w io.Writer, half uint16) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint16(nil, half))
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

func (i Unknown) Encode(w io.Writer, _ EncOpts) (int64, error) {
	if i.length == 2 {
		return writeHalf(w, uint16(i.raw))
	}

	return writeWord(w, i.raw)
}
