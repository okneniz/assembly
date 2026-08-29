package loong64

// The LoongArch instruction grammar: a mnemonic (longest-match against
// the arch decoding table) and operands - a $-register or an expression
// (asm/expr). It builds unevaluated operands (Op): value slots are
// expressions, evaluated at resolution.

import (
	"fmt"

	"github.com/okneniz/parsec/common"
	parsecstrings "github.com/okneniz/parsec/strings"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/asm/expr"
)

// asmMnemonics is all accepted mnemonics (the arch decoding table; the
// pseudo layer adds its own on top). It is built by a function in a var
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

// asmRegCanon maps every accepted register spelling ($t0, $r12) to the
// canonical $-name (the decoder's print form). It is built in a var
// initialization: cRegOperand (MapStrings builds a trie at creation)
// depends on it, while init() would run after the vars.
var asmRegCanon = buildAsmRegCanon()

func buildAsmRegCanon() map[string]string {
	m := map[string]string{}
	for i, name := range arch.RegNames() {
		m[name] = name
		m[fmt.Sprintf("$r%d", i)] = name
	}

	return m
}

// cRegOperand parses a register name ($zero, $t0, $r21, ...).
var cRegOperand = parsecstrings.Cast(
	parsecstrings.MapStrings("register", asmRegCanon),
	func(name string) (Op, error) {
		return OpReg(name), nil
	},
)

// cOperand parses an operand: a register or an expression.
var cOperand = parsecstrings.Choice("operand",
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

// skipSpaces consumes spaces (except the newline).
func skipSpaces(buf common.Buffer[rune, parsecstrings.Position]) {
	expr.SkipSpaces(buf)
}

// peekRune returns the next rune without consuming it; ok=false at EOF.
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
