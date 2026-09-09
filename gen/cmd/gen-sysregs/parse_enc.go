package main

// parsec grammar for the values of the v attribute of the <enc .../> element
// from the ARM System Register XML. Formats (see the previous parseEncV,
// adapted from golang.org/x/arch/arm64/arm64gen):
//
//	"0b110"       — a plain binary literal
//	"0b010:n[3]"  — binary prefix + one bit of the array index n
//	"0b1:n[1:0]"  — binary prefix + a bit slice of the index n
//	"n[3:0]"      — a bit slice of the index n without a prefix
//	"n[2]"        — one bit of the index n without a prefix
//
// The grammar is built by its constructors (no package vars): a
// constructor call allocates the combinators once.

import (
	"fmt"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// encValue — a parsed value: a binary prefix and an optional n[hi:lo] part
// parameterized by the register's array index.
type encValue struct {
	prefix uint64 // bits left of ':' (0 for a plain slice)
	param  bool   // has an n[...] part
	shift  uint   // lo: low bit of the slice
	width  uint   // width of the slice in bits
}

func newEncValue(prefix uint64, param bool, shift uint, width uint) encValue {
	return encValue{
		prefix: prefix,
		param:  param,
		shift:  shift,
		width:  width,
	}
}

// resolve substitutes the array index n and returns the encoding field value.
func (e encValue) resolve(instance int) uint64 {
	if !e.param {
		return e.prefix
	}

	return e.prefix<<e.width | uint64((uint(instance)>>e.shift)&((1<<e.width)-1))
}

// slice — a bit slice of an index: [hi] or [hi:lo].
type slice struct {
	hi, lo uint
}

func newSlice(hi uint, lo uint) slice {
	return slice{
		hi: hi,
		lo: lo,
	}
}

// newEncValueC — the full value of the v attribute (one of the three
// formats above).
func makeEncValueParser() parsec.Combinator[rune, strings.Position, encValue] {
	colon := strings.Try(strings.Eq("':'", ':'))
	n := strings.Try(strings.Eq("'n'", 'n'))

	// decUint — a decimal number (bit index) → uint.
	decUint := func(what string) parsec.Combinator[rune, strings.Position, uint] {
		digit := strings.Try(strings.Digit("decimal digit"))
		return strings.Cast(
			strings.Some(2, what, digit),
			func(rs []rune) (uint, error) {
				v, err := strconv.ParseUint(string(rs), 10, 32)
				if err != nil {
					return 0, fmt.Errorf("%s: %w", what, err)
				}

				return uint(v), nil
			},
		)
	}

	// cSlice — the contents of square brackets: "3" (hi==lo) or "3:0".
	sliceC := strings.Squares(strings.Choice("expected bit index or range",
		strings.Try(func() parsec.Combinator[rune, strings.Position, slice] {
			hi := decUint("high bit index")
			return func(buf parsec.Buffer[rune, strings.Position]) (slice, parsec.Error[strings.Position]) {
				h, err := hi(buf)
				if err != nil {
					return slice{}, err
				}

				if _, err := colon(buf); err != nil {
					return slice{}, err
				}

				l, err := decUint("low bit index")(buf)
				if err != nil {
					return slice{}, err
				}

				return newSlice(h, l), nil
			}
		}()),
		strings.Cast(decUint("bit index"), func(d uint) (slice, error) {
			return newSlice(d, d), nil
		}),
	))

	// binLit — "0b" and a non-empty sequence of binary digits → uint64.
	binDigit := strings.Try(strings.OneOf("binary digit", '0', '1'))
	binLit := strings.Cast(
		strings.Skip(strings.Try(strings.String("expected 0b", "0b")),
			strings.Some(8, "binary number", binDigit)),
		func(rs []rune) (uint64, error) {
			v, err := strconv.ParseUint(string(rs), 2, 64)
			if err != nil {
				return 0, fmt.Errorf("binary number: %w", err)
			}

			return v, nil
		},
	)

	return strings.Choice(
		"malformed enc value",
		strings.Try(
			func(buf parsec.Buffer[rune, strings.Position]) (encValue, parsec.Error[strings.Position]) {
				// "0b010:n[3]" / "0b1:n[1:0]"
				p, err := binLit(buf)
				if err != nil {
					return encValue{}, err
				}

				if _, err := colon(buf); err != nil {
					return encValue{}, err
				}

				if _, err := n(buf); err != nil {
					return encValue{}, err
				}

				sl, err := sliceC(buf)
				if err != nil {
					return encValue{}, err
				}

				v, cerr := toEncValue(p, sl)
				if cerr != nil {
					return encValue{}, parsec.NewParseError(buf.Position(), cerr.Error())
				}

				return v, nil
			},
		),
		strings.Try(
			func(buf parsec.Buffer[rune, strings.Position]) (encValue, parsec.Error[strings.Position]) {
				// "n[3:0]" / "n[2]"
				if _, err := n(buf); err != nil {
					return encValue{}, err
				}

				sl, err := sliceC(buf)
				if err != nil {
					return encValue{}, err
				}

				v, cerr := toEncValue(0, sl)
				if cerr != nil {
					return encValue{}, parsec.NewParseError(buf.Position(), cerr.Error())
				}

				return v, nil
			},
		),
		strings.Try(
			func(buf parsec.Buffer[rune, strings.Position]) (encValue, parsec.Error[strings.Position]) {
				// "0b110"
				p, err := binLit(buf)
				if err != nil {
					return encValue{}, err
				}

				return newEncValue(p, false, 0, 0), nil
			},
		),
	)
}

// toEncValue assembles an encValue from a prefix and a slice.
func toEncValue(prefix uint64, sl slice) (encValue, error) {
	if sl.hi < sl.lo {
		return encValue{}, fmt.Errorf("bad bit range [%d:%d]", sl.hi, sl.lo)
	}

	return newEncValue(prefix, true, sl.lo, sl.hi-sl.lo+1), nil
}

// parseEncValue parses the value of the v attribute; surrounding whitespace
// is allowed, but trailing junk after the value is an error (as in the
// previous strconv approach).
func parseEncValue(v string) (encValue, error) {
	cEncValue := makeEncValueParser()
	space := strings.Try(strings.Space("whitespace"))
	eof := strings.EOF()

	body := func(buf parsec.Buffer[rune, strings.Position]) (encValue, parsec.Error[strings.Position]) {
		e, err := cEncValue(buf)
		if err != nil {
			return encValue{}, err
		}

		// eof fails exactly when the buffer is not at the end
		if _, eofErr := eof(buf); eofErr != nil {
			return encValue{}, parsec.NewParseError(
				buf.Position(),
				"unexpected trailing characters",
			)
		}

		return e, nil
	}
	val, perr := strings.ParseString(v, strings.Padded(space, body))
	if perr != nil {
		return encValue{}, fmt.Errorf("enc value %q: %w (at %s)", v, perr, perr.Position())
	}

	return val, nil
}
