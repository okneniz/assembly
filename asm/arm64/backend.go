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
	"fmt"
	"io"
	"strconv"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	arch "github.com/okneniz/assembly/arch/arm64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// New returns an ARM64 Syntax for asm.Assemble (real instructions).
// An instance is intended for a single assembly.
func New() *Backend {
	return NewWithCtors(nil)
}

// Backend implements asm.Syntax. It owns the injected alias ctors; the
// model has no .option modes for ARM64 (ApplyOption is a no-op). The
// grammar is built once in NewWithCtors and captured here -
// Instruction/Comment return the stored values.
type Backend struct {
	extraCtors       map[string]arch.ArmCtor
	g                *armGrammar
	parseInstruction parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved]
	parseComment     parsec.Combinator[rune, parsecstrings.Position, string]
}

// NewWithCtors returns an asm.Syntax with extra constructors (aliases,
// asm/arm64/alias): they run before armCtors and the legacy path and go
// through the same self-verify encodeARM.
func NewWithCtors(extra map[string]arch.ArmCtor) *Backend {
	b := &Backend{extraCtors: extra}
	b.g = makeGrammar()
	b.parseInstruction = b.makeInstructionParser()
	b.parseComment = makeCommentParser()
	return b
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
	word, err := encodeARM(in, newCtx(c.Addr(), c.Resolve, in.ctors))
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

// Instruction is the grammar "mnemonic comma-separated operands"; built
// once in NewWithCtors (makeInstructionParser).
func (b *Backend) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return b.parseInstruction
}

// Comment parses ';' and '//' to the end of the line ('#' is NOT a
// comment - it is the imm prefix!); built once in NewWithCtors.
func (b *Backend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	return b.parseComment
}

// newComment builds the comment combinator (';' and '//' to the end of
// the line).
func makeCommentParser() parsec.Combinator[rune, parsecstrings.Position, string] {
	body := parsecstrings.Many(4, expr.MakeNotNewlineParser())
	semi := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.Eq("comment", ';')), body),
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
	return parsecstrings.Choice("comment", parsecstrings.Try(slash), parsecstrings.Try(semi))
}

// makeInstructionParser builds the instruction combinator: the mnemonic
// (longest-match + boundary check) and the operand list.
func (b *Backend) makeInstructionParser() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	g := b.g

	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		expr.SkipSpaces(buf)
		mnem, err := g.parseMnemonic(buf)
		if err != nil {
			return nil, err
		}

		if r, ok := expr.PeekRune(buf); ok && r != ' ' && r != '\t' && r != ',' && r != '\n' {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown mnemonic %q", mnem))
		}

		in := newArmAsmInstr(mnem, b.extraCtors)

		expr.SkipSpaces(buf)
		if expr.AtEOL(buf) {
			return in, nil
		}

		op, err := g.parseOperand(buf)
		if err != nil {
			return nil, err
		}

		if rerr := g.parseRegIndex(buf, &op); rerr != nil {
			return nil, rerr
		}

		in.ops = append(in.ops, op)
		for {
			save := buf.Position()
			expr.SkipSpaces(buf)
			if _, err := g.parseComma(buf); err != nil {
				if rerr := expr.Rewind(buf, save); rerr != nil {
					return nil, rerr
				}

				break
			}

			expr.SkipSpaces(buf)
			op, err := g.parseOperand(buf)
			if err != nil {
				return nil, err
			}

			if rerr := g.parseRegIndex(buf, &op); rerr != nil {
				return nil, rerr
			}

			in.ops = append(in.ops, op)
		}

		return in, nil
	}
}

// parseRegIndex parses the vector register index suffix: "v8[0]". The
// suffix is optional: its absence is not an error, only a failed
// rewind is.
func (g *armGrammar) parseRegIndex(
	buf parsec.Buffer[rune, parsecstrings.Position],
	op *armOp,
) parsec.Error[parsecstrings.Position] {
	if op.kind != armOpReg {
		return nil
	}

	save := buf.Position()
	if _, err := g.parseLBrack(buf); err != nil {
		return expr.Rewind(buf, save)
	}

	var digits []rune
	for {
		d, derr := g.parseDigit(buf)
		if derr != nil {
			break
		}

		digits = append(digits, d)
	}

	if len(digits) == 0 {
		return expr.Rewind(buf, save)
	}

	if _, err := g.parseRBrack(buf); err != nil {
		return expr.Rewind(buf, save)
	}

	idx, err := strconv.ParseInt(string(digits), 10, 32)
	if err != nil {
		return parsec.NewParseError(buf.Position(), "register index: "+err.Error())
	}

	op.expr = expr.Num(idx)
	return nil
}
