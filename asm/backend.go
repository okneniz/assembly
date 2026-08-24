// Package asm is the arch-neutral core of the assembler (text → machine
// code), the mirror of decode for the opposite direction. It owns everything
// that does not depend on the architecture: the line grammar (labels,
// directives, string literals - on the base rune atoms of the asm/expr
// package), the symbol table, sections, directives, and two-pass assembly.
// GAS expressions are the adjacent package asm/expr: AST, a grammar with
// precedences, and an evaluator.
// Arch specifics are only the parsing grammar (Syntax); encoding is done by
// the instruction itself (Unresolved → Resolved): the core does not know
// instruction types, makes no type-assertions, and knows no file formats.
//
// The syntax is a GNU-as compatible subset: labels, directives, expressions
// (| ^ & << >> + - * / %), plus its own assembly output (objdump style with
// absolute branch addresses) - this gives the round-trip oracle
// assemble(ObjDump(instr)) == instr.
//
// v1 limitations: no relocations and no .o writing (la/call/tail resolve to
// absolute addresses), no .macro, no relaxation.
package asm

import (
	"fmt"
	"strconv"
	"strings"

	parsecstrings "github.com/okneniz/parsec/strings"
)

// poolPrefix is the reserved prefix of literal pool slot names;
// \x00 is not part of identifiers - there are no clashes with user symbols
// (as with '.', also a special resolver name).
const poolPrefix = "\x00pool/"

// poolName is the literal pool slot name: deterministic from the slot size
// and the canonical expression key (deduplication of identical literals).
func poolName(slot int, key string) string {
	return poolPrefix + strconv.Itoa(slot) + "/" + key
}

// isPoolName reports whether the name belongs to the literal pool.
func isPoolName(name string) bool {
	return strings.HasPrefix(name, poolPrefix)
}

// PoolSelf is the reserved resolver name: the address of the literal pool
// slot of the instruction's OWN literal (GAS: ldr xN, =literal). The core
// substitutes it for the instruction that requested a slot (PoolUser) - the
// slot naming scheme does not cross the layer boundary.
const PoolSelf = "\x00pool/self"

// placeholderResolve is the layout pass resolver: every symbol equals the
// address of the instruction itself, so all relative offsets are zero and
// pass the range checks. Instructions must build their decisions (compression,
// pseudo-expansion) not on symbol values but on their presence - then both
// passes converge (see RISC-V: symbol targets are not compressed).
func placeholderResolve(addr uint64) func(string) (uint64, bool) {
	return func(string) (uint64, bool) {
		return addr, true
	}
}

// Section is an output section of assembly. NOBITS (.bss): Data=nil,
// Size is the memory reserve; for regular ones Size=len(Data).
type Section struct {
	Name   string
	Addr   uint64 // absolute address (the first section is base, then consecutive)
	Data   []byte // nil for NOBITS
	Size   int    // size in memory
	Nobits bool   // .bss: no file data
}

func NewSection(name string, addr uint64, data []byte) Section {
	return Section{
		Name: name,
		Addr: addr,
		Data: data,
		Size: len(data),
	}
}

// NewNobitsSection is a section without file data (.bss): reserve is the
// size in memory.
func NewNobitsSection(name string, addr uint64, reserve int) Section {
	return Section{
		Name:   name,
		Addr:   addr,
		Size:   reserve,
		Nobits: true,
	}
}

// Result is the assembly result: sections in order of appearance and the
// symbol table.
type Result struct {
	Sections []Section
	Symbols  map[string]uint64
	Globals  []string // .global/.globl
}

func NewResult(symbols map[string]uint64, globals []string) *Result {
	return &Result{
		Symbols: symbols,
		Globals: globals,
	}
}

// AsmError is an assembly error with a position in the source (1-based).
type AsmError struct {
	Line uint
	Col  uint
	Msg  string
}

func NewAsmError(line uint, col uint, msg string) AsmError {
	return AsmError{
		Line: line,
		Col:  col,
		Msg:  msg,
	}
}

func (e AsmError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Col, e.Msg)
}

// posErr builds an AsmError from a parsec position (0-based → 1-based) and
// a message.
func posErr(pos parsecstrings.Position, msg string) AsmError {
	return NewAsmError(pos.Line()+1, pos.Column()+1, msg)
}
