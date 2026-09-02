package loong64

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/text"
)

// Name - the architecture as consumers see it (test/diff/server).
const Name = "LOONG64"

// Instr - a LoongArch instruction: it represents itself - encodes itself
// into bytes (Encode, the computed form: operand values are already
// numbers, the environment is just the address), into text
// (disasm.ObjDump, for the disassembler) and into JSON (json.Marshaler,
// for the server). There is NO common type with other architectures; each
// instruction is its own structure (AddW, AddiW, Beq, ...).
type Instr interface {
	disasm.ObjDump
	json.Marshaler
	Addr() uint64

	// Encode encodes the computed instruction: pc is the absolute address
	// (PC-relative forms). There is no resolver - an exact inverse of
	// decode, and no modes: LA64 words are fixed 32-bit, no compression.
	Encode(w io.Writer, pc uint64) (int64, error)
}

// base - the bookkeeping fields of every instruction: address, raw word,
// and length. This is not instruction semantics but a bookkeeping record
// for Addr/Len/MarshalJSON.
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

// marshalDTO - the shared JSON skeleton for the server (server response
// keys; full is the full ObjDump text).
func (b base) marshalDTO(mnem, full, group string, fields map[string]any) ([]byte, error) {
	type dto struct {
		Addr        string         `json:"addr"`
		Bytes       string         `json:"bytes"`
		Raw         string         `json:"raw"`
		Mnemonic    string         `json:"mnemonic"`
		Operands    string         `json:"operands"`
		Group       string         `json:"group"`
		Tags        []string       `json:"tags"`
		Description string         `json:"description"`
		Aliases     []string       `json:"aliases"`
		DocURL      string         `json:"docurl"`
		Fields      map[string]any `json:"fields"`
	}
	if fields == nil {
		fields = map[string]any{}
	}

	return json.Marshal(dto{
		Addr:     fmt.Sprintf("0x%x", b.addr),
		Bytes:    text.FormatCode(b.raw, b.length, text.CodeWord),
		Raw:      fmt.Sprintf("0x%x", b.raw),
		Mnemonic: mnem,
		Operands: full,
		Group:    group,
		Fields:   fields,
	})
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

func (i Unknown) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(".word", i.ObjDump(disasm.DefaultViewCtx()), "", nil)
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
