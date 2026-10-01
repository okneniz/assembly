// Package alias — arbitrary test generators (oh-snap) for the ARM64
// aliases: the text layer above arch/arm64 (cmp is subs with Rd=zr,
// cset is csinc with an inverted condition, sxtb is SBFM...), one
// generator per alias mnemonic of asm/arm64/alias. The parameters are
// the alias's own operands, String() renders the alias text, Instr()
// assembles that text through the alias layer and decodes the word —
// the generated space is exactly what the assembler accepts.
package alias

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"strconv"

	arm64 "github.com/okneniz/assembly/arch/arm64"
	parsec "github.com/okneniz/parsec"
	parsecbytes "github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arb"
	asmalias "github.com/okneniz/assembly/asm/arm64/alias"
)

// stream — the lazy source of a value constructor (arb.Stream).
func stream[T any](f func() T) iter.Seq[T] {
	return arb.Stream(f)
}

// instrOfText — the instruction of one alias text (assemble + decode);
// the bridge of the alias families to the arch instruction.
func instrOfText(src string) (arm64.Instr, error) {
	res, errs := asmalias.Assemble(src, 0)
	if len(errs) != 0 {
		return nil, fmt.Errorf("%q: %w", src, errs[0])
	}

	if len(res.Sections) == 0 || len(res.Sections[0].Data) != 4 {
		return nil, fmt.Errorf("%q: not a single word", src)
	}

	ins, err := arm64.MakeDecoder()(
		parsec.Stateless{},
		parsecbytes.Buffer(res.Sections[0].Data),
	)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", src, err)
	}

	if len(ins) != 1 {
		return nil, fmt.Errorf("%q: %d instructions", src, len(ins))
	}

	return ins[0], nil
}

// genCond14 — a condition of the cset/cinc family (al/nv are not alias
// conditions there: the inverted form must stay a real condition).
func genCond14(rnd *rand.Rand) string {
	names := arm64.CondNames()
	return names[rnd.IntN(14)]
}

// condCanonical — the shrink target of a condition: "eq".
func condCanonical() string {
	return arm64.CondNames()[0]
}

// itoa — the decimal text of an immediate.
func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}

// utoa — the decimal text of an unsigned immediate.
func utoa(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

// halved — the halving-toward-zero shrink candidates of an immediate.
func halved(v int64) []int64 {
	if v == 0 {
		return nil
	}

	var out []int64
	for d := v / 2; d > 0; d /= 2 {
		out = append(out, d)
	}

	out = append(out, 0)
	return out
}

// uhalved — the halving shrink candidates of an unsigned field.
func uhalved(v uint32) []uint32 {
	if v == 0 {
		return nil
	}

	var out []uint32
	for d := v / 2; d > 0; d /= 2 {
		out = append(out, d)
	}

	out = append(out, 0)
	return out
}

// shiftKinds — the shift names of the shifted-register forms.
func shiftKinds() []string {
	return []string{"lsl", "lsr", "asr", "ror"}
}

// addSubShiftKinds — the shifts of the add/sub shifted-register forms
// (ror belongs to the logical family only).
func addSubShiftKinds() []string {
	return []string{"lsl", "lsr", "asr"}
}

// shiftSuffix — ", lsl #4" of a shifted-register operand, or "".
func shiftSuffix(sh string, amt int64) string {
	if sh == "" {
		return ""
	}

	return ", " + sh + " #" + itoa(amt)
}
