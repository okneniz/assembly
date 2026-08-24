package arm64

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/text"
)

// Name — the architecture as consumers see it (test/diff/server).
const Name = "ARM64"

// Instr — an ARM64 instruction: represents itself — encodes itself into
// bytes (Encode: computed form, operand values are already numbers, the
// environment — only the address), into text (disasm.ObjDump), into JSON
// (for the server) and knows its address (the key of diff tools). There is
// NO common type with other architectures; every instruction is its own
// struct (Ldr, AddImm, Csel, ...).
type Instr interface {
	disasm.ObjDump
	json.Marshaler
	Addr() uint64

	// Encode encodes the computed instruction: pc — the absolute address
	// (PC-relative forms). There is no resolver — an exact inverse of decode.
	Encode(w io.Writer, pc uint64) (int64, error)
}

// base — the bookkeeping fields of every instruction: address, raw word
// and length (32-bit ones are always 4; bookkeeping for Addr/Len/MarshalJSON).
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

// marshal — the common JSON skeleton for the server (server response keys;
// full — the full ObjDump text).
func (b base) marshal(mnem, full, group string, fields map[string]any) ([]byte, error) {
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
		Bytes:    text.FormatCode(b.raw, b.length, text.CodeBytes),
		Raw:      fmt.Sprintf("0x%x", b.raw),
		Mnemonic: mnem,
		Operands: full,
		Group:    group,
		Fields:   fields,
	})
}

// armRegName — the x/w register name by number and width.
func armRegName(n uint32, is64 bool) string {
	if is64 {
		return regNameX(n)
	}

	return regNameW(n)
}

// writeWord — write an encoded instruction (LE).
func writeWord(w io.Writer, word uint32) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint32(nil, word))
	return int64(n), err
}

// sfMatch — the encoding base by the first register's width.
func sfMatch(rd string, matchX, matchW uint32) (uint32, error) {
	if rd[0] == 'x' {
		return matchX, nil
	}

	if matchW == 0 {
		return 0, errors.New("32-bit form not supported")
	}

	return matchW, nil
}

// regNums2/3 — register numbers for encoding.
func regNums2(a, b string) (uint32, uint32, error) {
	na, err := armRegNum(a)
	if err != nil {
		return 0, 0, err
	}

	nb, err := armRegNum(b)
	if err != nil {
		return 0, 0, err
	}

	return na, nb, nil
}

func regNums3(a, b, c string) (uint32, uint32, uint32, error) {
	na, nb, err := regNums2(a, b)
	if err != nil {
		return 0, 0, 0, err
	}

	nc, err := armRegNum(c)
	if err != nil {
		return 0, 0, 0, err
	}

	return na, nb, nc, nil
}

// condNum — the condition number by name (inverse of condNames).
func condNum(c string) (uint32, error) {
	for i, n := range condNames {
		if n == c {
			return uint32(i), nil
		}
	}

	return 0, fmt.Errorf("unknown condition %q", c)
}

// imm — an immediate operand value: a concrete number. Symbolic slots are
// computed before the struct is built (resolve in encodeARM) — the computed
// instruction contains no holes.
type imm struct {
	val int64
}

// immNum — a concrete value (decoding and construction).
func immNum(v int64) imm {
	return imm{val: v}
}

// textHex — 0x%x output (for ObjDump of decoded instructions).
func (m imm) textHex() string {
	return fmt.Sprintf("0x%x", m.val)
}

// Schema — a decode-table entry: mask/value for selecting the encoding plus
// the formatting config (Fields/Formatter/FullFormat/Meta). An ARM-specific
// type (RISC-V uses its own decodeTable).
type Schema struct {
	Mask       uint32
	Value      uint32
	Fields     []Field
	Meta       Meta
	Formatter  string
	FullFormat bool
	// ctor — the constructor of the per-instruction struct: the decoding
	// contract (all table entries carry it).
	ctor func(word uint32, addr uint64) Instr
}

func NewSchema(
	mask uint32,
	value uint32,
	fields []Field,
	meta Meta,
	formatter string,
	fullFormat bool,
	ctor func(word uint32, addr uint64) Instr,
) Schema {
	return Schema{
		Mask:       mask,
		Value:      value,
		Fields:     fields,
		Meta:       meta,
		Formatter:  formatter,
		FullFormat: fullFormat,
		ctor:       ctor,
	}
}

// Unknown — a word matching no encoding (the ".word" fallback): displayed
// as .word #0x...; Encode re-emits the raw word.
type Unknown struct {
	base

	word uint32
}

func decodeUnknown(w uint32, addr uint64) Instr {
	return Unknown{
		base: newBase(addr, w),
		word: w,
	}
}

func (i Unknown) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(".word #0x%x", i.word)
}

func (i Unknown) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, i.word)
}

func (i Unknown) MarshalJSON() ([]byte, error) {
	return i.marshal(".word", i.ObjDump(disasm.DefaultViewCtx()), "", nil)
}

// shiftNumByName — the shift-kind number (inverse of shiftNames).
func shiftNumByName(name string) (uint32, error) {
	for i, n := range shiftNames {
		if n == name {
			return uint32(i), nil
		}
	}

	return 0, fmt.Errorf("unknown shift %q", name)
}

// regBitsX — the x register number for encoding (adr/adrp/mrs/msr).
func regBitsX(name string) uint32 {
	n, err := armRegNum(name)
	if err != nil {
		return 0
	}

	return n
}

// invSysRegChecked — the 15-bit system register key by name.
func invSysRegChecked(name string) uint32 {
	k, err := invSysReg(name)
	if err != nil {
		return 0
	}

	return k
}
