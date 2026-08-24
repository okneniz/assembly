package riscv

// The RISC-V instruction grammar: a mnemonic (longest-match against
// the arch decoding table) and operands - memory off(reg), a register
// (ABI or xN/fN), an expression (asm/expr). It builds unevaluated
// operands (Op): value slots are expressions, evaluated at resolution.

import (
	"fmt"

	"github.com/okneniz/parsec/common"
	parsecstrings "github.com/okneniz/parsec/strings"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm/expr"
)

// asmMnemonics is all accepted mnemonics (the arch decoding table;
// pseudo adds its own layer). It is built by a function in a var
// initialization: cMnemonic (MapStrings builds a trie at creation)
// depends on asmMnemonics, while init() would run after the vars.
var asmMnemonics = buildAsmMnemonics()

func buildAsmMnemonics() map[string]string {
	m := map[string]string{}
	for _, name := range arch.Mnemonics() {
		m[name] = name
	}

	return m
}

// cMnemonic is the mnemonic table (longest-match).
var cMnemonic = parsecstrings.MapStrings("mnemonic", asmMnemonics)

// asmReg is a register from the source: an ABI name or xN/fN.
type asmReg struct {
	num uint32
	fp  bool
}

// asmRegNum maps a name -> register (int + FP, both notations). It is
// built in a var initialization: cRegOperand (MapStrings builds a
// trie at creation) depends on it, while init() would run after the
// vars.
var asmRegNum = buildAsmRegNum()

func buildAsmRegNum() map[string]asmReg {
	m := map[string]asmReg{}
	for i, n := range arch.RegNames() {
		m[n] = asmReg{num: uint32(i), fp: false}
		m[fmt.Sprintf("x%d", i)] = asmReg{num: uint32(i), fp: false}
	}

	for i, n := range arch.FpRegNames() {
		m[n] = asmReg{num: uint32(i), fp: true}
		m[fmt.Sprintf("f%d", i)] = asmReg{num: uint32(i), fp: true}
	}

	return m
}

// cOperand parses an operand: memory | register | expression.
var cOperand = parsecstrings.Choice("operand",
	parsecstrings.Try(cMemOperand),
	parsecstrings.Try(cRegOperand),
	parsecstrings.Try(func() common.Combinator[rune, parsecstrings.Position, Op] {
		return func(buf common.Buffer[rune, parsecstrings.Position]) (Op, common.Error[parsecstrings.Position]) {
			e, err := expr.CExpr()(buf)
			if err != nil {
				return Op{}, err
			}

			return OpExpr(e), nil
		}
	}()),
)

// cRegOperand parses a register name (ABI or xN/fN).
var cRegOperand = parsecstrings.Cast(
	parsecstrings.MapStrings("register", asmRegNum),
	func(r asmReg) (Op, error) {
		return OpReg(regName(r)), nil
	},
)

// regName restores the canonical name (the arch name tables are used
// to match the decoder output).
func regName(r asmReg) string {
	if r.fp {
		return arch.FpRegNames()[r.num]
	}

	return arch.RegNames()[r.num]
}

// cMemOperand — [expr] '(' reg ')': "0x8(sp)", "(a0)".
var cMemOperand = func() common.Combinator[rune, parsecstrings.Position, Op] {
	lparen := parsecstrings.Try(parsecstrings.Eq("'('", '('))
	rparen := parsecstrings.Try(parsecstrings.Eq("')'", ')'))
	return func(buf common.Buffer[rune, parsecstrings.Position]) (Op, common.Error[parsecstrings.Position]) {
		var off *expr.Expr
		if r, ok := expr.PeekRune(buf); ok && r == '(' {
			if _, err := lparen(buf); err != nil {
				return Op{}, err
			}
		} else {
			e, err := expr.CExpr()(buf)
			if err != nil {
				return Op{}, err
			}

			off = e
			expr.SkipSpaces(buf)
			if _, err := lparen(buf); err != nil {
				return Op{}, err
			}
		}

		r, err := cRegOperand(buf)
		if err != nil {
			return Op{}, err
		}

		if _, err := rparen(buf); err != nil {
			return Op{}, err
		}

		return OpMemExpr(r.Reg(), off), nil
	}
}()

// skipSpaces consumes spaces (except the newline).
func skipSpaces(buf common.Buffer[rune, parsecstrings.Position]) {
	expr.SkipSpaces(buf)
}

// peekRune returns the next rune without consuming it; ok=false at
// EOF.
func peekRune(buf common.Buffer[rune, parsecstrings.Position]) (rune, bool) {
	return expr.PeekRune(buf)
}

// ParseOps parses the operand list after the mnemonic (to the end of
// the line, comma-separated) - the assembler's operand grammar.
func ParseOps(
	buf common.Buffer[rune, parsecstrings.Position],
) ([]Op, common.Error[parsecstrings.Position]) {
	var ops []Op
	expr.SkipSpaces(buf)
	if expr.AtEOL(buf) {
		return ops, nil
	}

	op, err := cOperand(buf)
	if err != nil {
		return nil, err
	}

	ops = append(ops, op)
	for {
		save := buf.Position()
		expr.SkipSpaces(buf)
		if _, err := expr.CComma(buf); err != nil {
			if rerr := expr.Rewind(buf, save); rerr != nil {
				return nil, rerr
			}

			break
		}

		expr.SkipSpaces(buf)
		op, err := cOperand(buf)
		if err != nil {
			return nil, err
		}

		ops = append(ops, op)
	}

	return ops, nil
}
