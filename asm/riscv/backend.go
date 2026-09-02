// Package riscv is the RISC-V (RV64GC) assembler syntax layer: the
// instruction grammar, .option modes, and unevaluated instructions
// (operands with expression slots). It sits ABOVE arch/riscv
// (concrete instructions: decode/encode) and the asm core (sections,
// symbols, passes).
// Pseudo-instructions
// (nop/li/mv/ret/call/la/...) are a layer above this package (pseudo).
// A mirror of Parse: the operand grammar accepts both GNU-as syntax
// and ObjDump()'s own output (round-trip).
package riscv

import (
	"fmt"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// New returns a RISC-V Syntax for asm.Assemble. An instance is
// intended for a single assembly: ApplyOption mutates the (.option)
// modes; create separate instances for parallel assemblies.
func New() *Backend {
	return &Backend{}
}

// Backend implements asm.Syntax; it owns the (.option) encoding modes:
// noRVC disables auto-compression, optStack holds push/pop snapshots.
type Backend struct {
	noRVC    bool
	optStack []bool
}

// ApplyOption handles .option values: norvc/rvc toggle auto-compression,
// push/pop are a stack of mode snapshots; the rest (pic/nopic/relax/
// norelax/csr-check/arch) are outside the model (an absolute assembler
// without relaxation, a fixed RV64GC profile) and do not affect the
// bytes - they are ignored.
func (b *Backend) ApplyOption(name string) error {
	switch name {
	case "norvc":
		b.noRVC = true
	case "rvc":
		b.noRVC = false
	case "push":
		b.optStack = append(b.optStack, b.noRVC)
	case "pop":
		if n := len(b.optStack); n > 0 {
			b.noRVC = b.optStack[n-1]
			b.optStack = b.optStack[:n-1]
		}
	}

	return nil
}

// ResetOptions resets the modes to the defaults (the core calls it
// before every pass: .option is applied symmetrically from one
// initial state).
func (b *Backend) ResetOptions() {
	b.noRVC = false
	b.optStack = nil
}

// Instruction is the grammar "mnemonic operands" (comma-separated
// operands).
func (b *Backend) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		skipSpaces(buf)
		name, err := cMnemonic(buf)
		if err != nil {
			return nil, err
		}

		// mnemonic boundary: followed by a space/comma/end of line
		if r, ok := peekRune(buf); ok && r != ' ' && r != '\t' && r != ',' && r != '\n' {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown mnemonic %q", name))
		}

		ops, err := ParseOps(buf)
		if err != nil {
			return nil, err
		}

		return instr{mnem: name, ops: ops, be: b}, nil
	}
}

// Comment parses '#' and '//' to the end of the line.
func (b *Backend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	body := parsecstrings.Many(4, expr.CNotNL)
	hash := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.Eq("comment", '#')), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	slash := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.String("comment", "//")), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	return parsecstrings.Choice("comment", parsecstrings.Try(slash), parsecstrings.Try(hash))
}

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr, the
// symbolic flag is absorbed into EncOpts.NoRVC - see resolve.go).
type instr struct {
	mnem string
	ops  []Op
	be   *Backend
}

// Resolve evaluates the expressions and builds the evaluated
// instruction; the "slot was symbolic" flag + the current .option
// norvc -> NoRVC (symbolic targets are not compressed - so that pass
// sizes converge).
func (in instr) Resolve(ctx asm.Ctx) (asm.Resolved, error) {
	return in.resolve(ctx, in.be != nil && in.be.noRVC)
}

// ResolveForm evaluates the "mnemonic + operands" form in the ctx
// environment with the current .option modes (pseudo expansion: each
// form starts from its own address, pass a derived ctx).
func (b *Backend) ResolveForm(mnem string, ops []Op, ctx asm.Ctx) (asm.Resolved, error) {
	return instr{mnem: mnem, ops: ops, be: b}.resolve(ctx, b.noRVC)
}
