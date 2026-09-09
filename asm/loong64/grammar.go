package loong64

// The LoongArch instruction grammar: a mnemonic (longest-match against
// the arch decoding table) and operands - a $-register or an expression
// (asm/expr). It builds unevaluated operands (Op): value slots are
// expressions, evaluated at resolution.
//
// The grammar is a struct built once per Backend (newGrammar): the
// combinators are values captured by the closures.

import (
	"fmt"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	arch "github.com/okneniz/assembly/arch/loong64"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

type operand = parsec.Combinator[rune, parsecstrings.Position, Op]

// buildAsmMnemonics is all accepted mnemonics (the arch decoding table;
// the pseudo layer adds its own on top).
func buildAsmMnemonics() map[string]string {
	m := map[string]string{}
	for _, name := range arch.Mnemonics() {
		m[name] = name
	}

	return m
}

// buildAsmRegCanon maps every accepted register spelling ($t0, $r12) to
// the canonical $-name (the decoder's print form).
func buildAsmRegCanon() map[string]string {
	m := map[string]string{}
	for i, name := range arch.RegNames() {
		m[name] = name
		m[fmt.Sprintf("$r%d", i)] = name
	}

	return m
}

// grammar is the instruction grammar of the syntax layer: the mnemonic
// trie (longest-match), the register operand, the full operand
// alternative, the shared expression ladder and the comma atom.
type grammar struct {
	parseMnemonic   parsec.Combinator[rune, parsecstrings.Position, string]
	parseRegOperand operand
	parseOperand    operand
	parseComma      parsec.Combinator[rune, parsecstrings.Position, rune]
	parseExpr       parsec.Combinator[rune, parsecstrings.Position, *expr.Expr]
}

// newGrammar builds the whole grammar once; the ready combinators are
// captured by the backend and reused for every line.
func makeGrammar() *grammar {
	g := &grammar{
		parseMnemonic: parsecstrings.MapStrings("mnemonic", buildAsmMnemonics()),
		parseComma:    expr.MakeCommaParser(),
		parseExpr:     expr.MakeExprParser(),
	}

	g.parseRegOperand = parsecstrings.Cast(
		parsecstrings.MapStrings("register", buildAsmRegCanon()),
		func(name string) (Op, error) {
			return OpReg(name), nil
		},
	)

	exprOperand := func(buf parsec.Buffer[rune, parsecstrings.Position]) (Op, parsec.Error[parsecstrings.Position]) {
		e, err := g.parseExpr(buf)
		if err != nil {
			return Op{}, err
		}

		return OpExpr(e), nil
	}

	g.parseOperand = parsecstrings.Choice("operand",
		parsecstrings.Try(g.parseRegOperand),
		parsecstrings.Try(exprOperand),
	)

	return g
}

// skipSpaces consumes spaces (except the newline).
func skipSpaces(buf parsec.Buffer[rune, parsecstrings.Position]) {
	expr.SkipSpaces(buf)
}

// peekRune returns the next rune without consuming it; ok=false at EOF.
func peekRune(buf parsec.Buffer[rune, parsecstrings.Position]) (rune, bool) {
	return expr.PeekRune(buf)
}

// ParseOps parses the operand list after the mnemonic (to the end of
// the line, comma-separated) - the assembler's operand grammar.
func (b *Backend) ParseOps(
	buf parsec.Buffer[rune, parsecstrings.Position],
) ([]Op, parsec.Error[parsecstrings.Position]) {
	g := b.g

	var ops []Op

	expr.SkipSpaces(buf)
	if expr.AtEOL(buf) {
		return ops, nil
	}

	op, err := g.parseOperand(buf)
	if err != nil {
		return nil, err
	}

	ops = append(ops, op)
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

		ops = append(ops, op)
	}

	return ops, nil
}

// Instruction is the grammar "mnemonic operands" (comma-separated
// operands); built once in New (b.newInstruction).
func (b *Backend) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return b.parseInstruction
}

// newInstruction builds the instruction combinator: mnemonic
// (longest-match + boundary check) and the operand list.
func (b *Backend) makeInstructionParser() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		skipSpaces(buf)
		name, err := b.g.parseMnemonic(buf)
		if err != nil {
			return nil, err
		}

		// mnemonic boundary: followed by a space/comma/end of line
		if r, ok := peekRune(buf); ok && r != ' ' && r != '\t' && r != ',' && r != '\n' {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown mnemonic %q", name))
		}

		ops, err := b.ParseOps(buf)
		if err != nil {
			return nil, err
		}

		return newInstr(name, ops), nil
	}
}

// Comment parses '#' and '//' to the end of the line; built once in New.
func (b *Backend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	return b.parseComment
}

// newComment builds the comment combinator ('#' and '//' to the end of
// the line).
func makeCommentParser() parsec.Combinator[rune, parsecstrings.Position, string] {
	body := parsecstrings.Many(4, expr.MakeNotNewlineParser())
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
