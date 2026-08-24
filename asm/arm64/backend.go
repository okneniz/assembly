// Package arm64 is the ARM64 assembler syntax layer: the instruction
// grammar, unevaluated operands (expression slots), the backend, and
// encoding with self-verify (encodeARM: encode + decoding with the
// arch's own decoder + comparing normalized texts). It sits ABOVE
// arch/arm64 (concrete instructions: decode/encode) and the asm core
// (sections, symbols, passes).
// Aliases (the cmp/neg/mov/cset family) are a layer above this
// package (alias).
package arm64

import (
	"encoding/binary"
	"io"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// New returns an ARM64 Syntax for asm.Assemble (real instructions).
// An instance is intended for a single assembly.
func New() *Backend {
	return &Backend{}
}

// Backend implements asm.Syntax. It owns the injected alias ctors; the
// model has no .option modes for ARM64 (ApplyOption is a no-op).
type Backend struct {
	extraCtors map[string]arch.ArmCtor
}

// NewWithCtors returns an asm.Syntax with extra constructors (aliases,
// asm/arm64/alias): they run before armCtors and the legacy path and go
// through the same self-verify encodeARM.
func NewWithCtors(extra map[string]arch.ArmCtor) *Backend {
	return &Backend{extraCtors: extra}
}

// ApplyOption is a no-op: the model has no .option modes for ARM64.
func (b *Backend) ApplyOption(string) error {
	return nil
}

// ResetOptions is a no-op (see ApplyOption).
func (b *Backend) ResetOptions() {}

// armAsmInstr is an unevaluated instruction (mnemonic + operands); it
// implements Unresolved: resolve+encodeARM with self-verify
// (encode.go).
type armAsmInstr struct {
	mnem  string
	ops   []armOp
	ctors map[string]arch.ArmCtor
}

// Resolve evaluates the expressions and encodes with self-verify.
func (in armAsmInstr) Resolve(c asm.Ctx) (asm.Resolved, error) {
	word, err := encodeARM(in, ctx{Addr: c.Addr(), Resolve: c.Resolve, extra: in.ctors})
	if err != nil {
		return nil, err
	}

	return encodedWord{w: word}, nil
}

// encodedWord is an evaluated instruction as a ready-made word.
type encodedWord struct {
	w uint32
}

func (e encodedWord) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint32(nil, e.w))
	return int64(n), err
}

func newArmAsmInstr(mnem string, ctors map[string]arch.ArmCtor) armAsmInstr {
	return armAsmInstr{
		mnem:  mnem,
		ctors: ctors,
	}
}

// PoolReq reports the literal pool request of the wrapper form: ldr
// with a "=expression" operand (mirrors newLdrPool, without building
// a per-instr structure).
func (in armAsmInstr) PoolReq() (*expr.Expr, int, bool) {
	if in.mnem != "ldr" || len(in.ops) != 2 || in.ops[1].kind != armOpLit {
		return nil, 0, false
	}

	if in.ops[0].reg == "" {
		return nil, 0, false
	}

	if in.ops[0].reg[0] == 'w' {
		return in.ops[1].expr, 4, true
	}

	return in.ops[1].expr, 8, true
}
