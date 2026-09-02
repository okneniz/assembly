package opcodes

import (
	"fmt"
	"strconv"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// Field is a contiguous bit field of the encoding: bits [Lo+Width-1 : Lo].
type Field struct {
	Lo    int
	Width int
}

// NewField - a bit field at Lo, Width bits wide.
func NewField(lo, width int) Field {
	return Field{
		Lo:    lo,
		Width: width,
	}
}

// Imm is one immediate slot: a signedness plus one or more segments whose
// fields concatenate MSB-first (Sd5k16: [4:0] is imm[20:16], [25:10] is
// imm[15:0]).
type Imm struct {
	Signed   bool
	Segments []Field
}

// NewImm - an immediate slot from its signedness and segments.
func NewImm(signed bool, segments []Field) Imm {
	return Imm{
		Signed:   signed,
		Segments: segments,
	}
}

// Slots is the decoded operand layout of a format string.
type Slots struct {
	Regs []rune // register letters, format order ('D', 'J', ...)
	Imms []Imm  // immediate slots, format order
}

// NewSlots - a layout from register letters and immediate slots.
func NewSlots(regs []rune, imms []Imm) Slots {
	return Slots{
		Regs: regs,
		Imms: imms,
	}
}

// regSlots - the fixed positions of the integer register slots.
var regSlots = map[rune]Field{
	'D': NewField(0, 5),
	'J': NewField(5, 5),
	'K': NewField(10, 5),
	'A': NewField(15, 5),
}

// immStart - the starting bit of an immediate segment index.
var immStart = map[rune]int{
	'd': 0,
	'j': 5,
	'k': 10,
	'a': 15,
	'm': 16,
	'n': 18,
}

// --- format notation grammar. ---

var (
	cRegSlot   = strings.Try(strings.OneOf("register slot", 'D', 'J', 'K', 'A'))
	cImmSign   = strings.Try(strings.OneOf("immediate signedness", 'S', 'U'))
	cSegIndex  = strings.Try(strings.OneOf("segment index", 'd', 'j', 'k', 'a', 'm', 'n'))
	cSegDigits = strings.Cast(
		strings.Some(2, "segment width", strings.Try(strings.Digit("decimal digit"))),
		castWidth,
	)

	// cSegment - one index-width piece: 'k' + 12 -> {10, 12}.
	cSegment = func() parsec.Combinator[rune, strings.Position, Field] {
		return func(buf parsec.Buffer[rune, strings.Position]) (Field, parsec.Error[strings.Position]) {
			index, err := cSegIndex(buf)
			if err != nil {
				return Field{}, err
			}

			width, err := cSegDigits(buf)
			if err != nil {
				return Field{}, err
			}

			return NewField(immStart[index], width), nil
		}
	}()

	// cImm - one immediate slot: the signedness letter plus one or more
	// segments.
	cImm = strings.Try(func() parsec.Combinator[rune, strings.Position, Imm] {
		return func(buf parsec.Buffer[rune, strings.Position]) (Imm, parsec.Error[strings.Position]) {
			sign, err := cImmSign(buf)
			if err != nil {
				return Imm{}, err
			}

			segments, err := strings.Some(4, "expected segments", cSegment)(buf)
			if err != nil {
				return Imm{}, err
			}

			return NewImm(sign == 'S', segments), nil
		}
	}())
)

// castWidth - the segment width in bits (1..2 digits; the bounds are
// validated by Mask).
func castWidth(ds []rune) (int, error) {
	v, err := strconv.ParseUint(string(ds), 10, 8)
	if err != nil {
		return 0, fmt.Errorf("segment width: %w", err)
	}

	return int(v), nil
}

// ParseFormat decodes a canonical operand format string ("DJSk12") into
// its register slots and immediate segments; "EMPTY" - no operands.
// Registers precede immediates (the upstream canonical order); a register
// after the first immediate is a notation violation.
func ParseFormat(format string) (Slots, error) {
	if format == "EMPTY" {
		return NewSlots(nil, nil), nil
	}

	buf := strings.Buffer([]rune(format))

	slots := Slots{}
	for {
		if r, err := cRegSlot(buf); err == nil {
			if len(slots.Imms) > 0 {
				return Slots{}, fmt.Errorf("format %q: register after immediate", format)
			}

			slots.Regs = append(slots.Regs, r)

			continue
		}

		if imm, err := cImm(buf); err == nil {
			slots.Imms = append(slots.Imms, imm)

			continue
		}

		break
	}

	if len(slots.Regs)+len(slots.Imms) == 0 {
		return Slots{}, fmt.Errorf("format %q: no operand slots", format)
	}

	if !buf.IsEOF() {
		return Slots{}, fmt.Errorf("format %q: unparsed tail", format)
	}

	return slots, nil
}

// Fields - every bit field of the layout: the register slots and the
// immediate segments.
func (s Slots) Fields() []Field {
	fs := make([]Field, 0, len(s.Regs))
	for _, r := range s.Regs {
		fs = append(fs, regSlots[r])
	}

	for _, imm := range s.Imms {
		fs = append(fs, imm.Segments...)
	}

	return fs
}

// Mask - the encoding mask of the layout: 1 for the fixed bits, 0 for the
// operand fields. Fails on out-of-range or overlapping fields - the
// notation cannot express an overlap, so one means corrupt data.
func (s Slots) Mask() (uint32, error) {
	covered := uint32(0)

	for _, f := range s.Fields() {
		if f.Lo < 0 || f.Width <= 0 || f.Lo+f.Width > 32 {
			return 0, fmt.Errorf("field [%d:+%d]: out of range", f.Lo, f.Width)
		}

		bits := (uint32(1)<<f.Width - 1) << f.Lo
		if covered&bits != 0 {
			return 0, fmt.Errorf("field [%d:+%d]: overlapping slots", f.Lo, f.Width)
		}

		covered |= bits
	}

	return ^covered, nil
}
