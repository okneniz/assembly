package asm

// The "unresolved / resolved instruction" contract:
// asm/<arch> parses text into unresolved instructions (expression operands),
// the core computes addresses and calls Resolve; arch/<arch> encodes resolved
// instructions purely - the exact inverse of decode, without a resolver.

import (
	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
	"github.com/okneniz/assembly/unit"
)

// Syntax is the per-arch syntax layer (asm/<arch>): the grammar of
// unresolved instructions and .option modes. Created per assembly; the core
// calls ResetOptions before each pass (symmetry).
type Syntax interface {
	// Instruction is the full parsec grammar of an instruction (mnemonic
	// and operands); the result is an unresolved instruction.
	Instruction() (comb parsec.Combinator[rune, parsecstrings.Position, Unresolved, parsec.Stateless])

	// Comment is the parsec grammar of a comment up to end of line (the
	// characters are arch-dependent: '#' - RISC-V, '//' - both). The
	// newline is not consumed.
	Comment() parsec.Combinator[rune, parsecstrings.Position, string, parsec.Stateless]

	// Separator is the statement separator rune of the syntax: ';' - ARM
	// (a line may hold several statements, as in GAS), 0 - the newline is
	// the only separator (RISC-V, Loong).
	Separator() rune

	// ApplyOption applies a .option value (each arch interprets its own).
	ApplyOption(name string) error

	// ResetOptions resets the modes to defaults.
	ResetOptions()
}

// Unresolved is an unresolved instruction: value slots are still
// expressions. Resolve is first called under the placeholder environment of
// the layout pass (all symbols = the instruction's own address: pc-relative
// offsets are zero); the layout is then RELAXED - the walk repeats with the
// frozen symbol table of the previous iteration until the sizes stabilize
// (see walkLayout). Size decisions may therefore depend on symbol values,
// with one obligation: sizes may only GROW from the placeholder seed (the
// relaxation is monotone and terminates); anything else is an
// encoding-vs-reservation error in pass 2.
type Unresolved interface {
	// Resolve evaluates the expressions via ctx and builds the resolved
	// instruction.
	Resolve(ctx unit.Ctx) (unit.Resolved, error)
}

// addrCtx - see unit.NewCtx (the concrete Ctx of the resolve walks).

// PoolUser is an optional capability of an unresolved instruction: it needs
// a slot in the literal pool of its subsection (GAS: ldr xN, =literal).
// PoolReq returns the slot value (evaluated when the pool is written) and
// its size in bytes (4/8); the instruction receives the slot address via the
// reserved resolver name. Deduplication is by (slot, ExprKey): identical
// literals share a slot.
type PoolUser interface {
	Unresolved
	PoolReq() (*expr.Expr, int, bool)
}

// countingWriter counts written bytes (for the layout pass).
type countingWriter struct {
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.n += int64(len(p))
	return len(p), nil
}

// sizeOf determines the instruction size by a trial Resolve with the
// sizing environment (placeholder on the first layout walk, the relaxation
// chain afterwards) and a count of the written bytes. The sizes of the
// FINAL layout walk are what pass 2 encodes against (see Unresolved on the
// relaxation).
func sizeOf(in Unresolved, ctx unit.Ctx) (int, error) {
	res, err := in.Resolve(ctx)
	if err != nil {
		return 0, err
	}

	var c countingWriter
	if _, err := res.Encode(&c); err != nil {
		return 0, err
	}

	return int(c.n), nil
}
