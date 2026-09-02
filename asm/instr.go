package asm

// The "unresolved / resolved instruction" contract:
// asm/<arch> parses text into unresolved instructions (expression operands),
// the core computes addresses and calls Resolve; arch/<arch> encodes resolved
// instructions purely - the exact inverse of decode, without a resolver.

import (
	"io"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// Ctx is the evaluation environment of an unresolved instruction: the
// address of the instruction itself (PC-relative encodings, '.') and the
// symbol resolver. Built by the core; it carries no arch-specific modes -
// those live in Resolved.
type Ctx interface {
	// Addr is the absolute address of the instruction itself.
	Addr() uint64

	// Resolve is the value of a symbol by name; ok=false for unknown ones.
	Resolve(name string) (uint64, bool)
}

// Syntax is the per-arch syntax layer (asm/<arch>): the grammar of
// unresolved instructions and .option modes. Created per assembly; the core
// calls ResetOptions before each pass (symmetry).
type Syntax interface {
	// Instruction is the full parsec grammar of an instruction (mnemonic
	// and operands); the result is an unresolved instruction.
	Instruction() (comb parsec.Combinator[rune, parsecstrings.Position, Unresolved])

	// Comment is the parsec grammar of a comment up to end of line (the
	// characters are arch-dependent: '#' - RISC-V, ';' - ARM, '//' - both).
	// The newline is not consumed.
	Comment() parsec.Combinator[rune, parsecstrings.Position, string]

	// ApplyOption applies a .option value (each arch interprets its own).
	ApplyOption(name string) error

	// ResetOptions resets the modes to defaults.
	ResetOptions()
}

// Unresolved is an unresolved instruction: value slots are still
// expressions. Resolve must be stable under the placeholder environment of
// the layout pass (all symbols = the instruction's own address): the size of
// the encoded result does not depend on the VALUES of symbols, only on their
// presence - an implementation obligation (riscv: symbol targets are not
// compressed).
type Unresolved interface {
	// Resolve evaluates the expressions via ctx and builds the resolved
	// instruction.
	Resolve(ctx Ctx) (Resolved, error)
}

// Resolved is a resolved instruction: all operands are numbers. Encoding is
// pure, no environment needed (arch modes are values closed over by the
// implementation, see the riscv EncOpts approach).
type Resolved interface {
	// Encode writes the encoding bytes.
	Encode(w io.Writer) (int64, error)
}

// addrCtx is the core's concrete Ctx.
type addrCtx struct {
	addr    uint64
	resolve func(string) (uint64, bool)
}

func newCtx(addr uint64, resolve func(string) (uint64, bool)) Ctx {
	return addrCtx{addr: addr, resolve: resolve}
}

func (c addrCtx) Addr() uint64 {
	return c.addr
}

func (c addrCtx) Resolve(name string) (uint64, bool) {
	if c.resolve == nil {
		return 0, false
	}

	return c.resolve(name)
}

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

// sizeOf determines the instruction size by a trial Resolve with a
// placeholder environment and a count of the written bytes. Deterministic
// across passes - the same call in code generation gives the same length
// (see Unresolved on placeholder stability).
func sizeOf(in Unresolved, ctx Ctx) (int, error) {
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
