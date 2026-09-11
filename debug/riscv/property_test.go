package riscv

import (
	"bytes"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	parsecbytes "github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arb"
	arbriscv "github.com/okneniz/assembly/arb/riscv"
	arch "github.com/okneniz/assembly/arch/riscv"
)

// TestInstrLenMatchesDecoder - the target's length heuristic (the low
// two bits of the first byte: 11 starts a 32-bit instruction,
// anything else a compressed one) agrees with the decoder on every
// generated encoding. The breakpoint kind and the disassembly step
// must equal the real instruction length, compressed or not.
func TestInstrLenMatchesDecoder(t *testing.T) {
	rnd := arb.Rnd(42)
	tgt := NewTarget()
	decode := arch.MakeDecoder()

	// one generator of encoded instruction bytes per arb sampler; the
	// law is the same, the samplers cover the RVC-compressible forms
	// (add/addi/lw/sw compress under register and immediate constraints)
	cases := []struct {
		name string
		gen  ohsnap.Arbitrary[[]byte]
	}{
		{name: "Add", gen: encoded(arbriscv.Add(rnd))},
		{name: "Sub", gen: encoded(arbriscv.Sub(rnd))},
		{name: "Addi", gen: encoded(arbriscv.Addi(rnd))},
		{name: "Lui", gen: encoded(arbriscv.Lui(rnd))},
		{name: "Lw", gen: encoded(arbriscv.Lw(rnd))},
		{name: "Ld", gen: encoded(arbriscv.Ld(rnd))},
		{name: "Sw", gen: encoded(arbriscv.Sw(rnd))},
		{name: "Sd", gen: encoded(arbriscv.Sd(rnd))},
	}

	compressed := 0
	for _, c := range cases {
		ohsnap.Check(t, 12000, c.gen, func(b []byte) bool {
			instrs, err := decode(parsecbytes.Buffer(b))
			if err != nil || len(instrs) != 1 {
				return false
			}

			if len(b) == 2 {
				compressed++
			}

			return tgt.InstrLen(b) == instrs[0].Len() && instrs[0].Len() == len(b)
		})
	}

	// the law must not hold vacuously over 32-bit forms only: some
	// sampled encodings must actually be compressed
	if compressed == 0 {
		t.Fatal("no compressed (2-byte) encodings sampled: the RVC half of the law is untested")
	}

	t.Logf("compressed encodings sampled: %d", compressed)
}

// encoded maps an arb params sampler to the encoded bytes of its
// instruction: the params build the instruction themselves (the Instr
// method of every arb/riscv params type; an encode failure maps to nil -
// the property fails on it, encoding is part of the samplers' own
// contract).
func encoded[P interface{ Instr() arch.Instr }](
	params ohsnap.Arbitrary[P],
) ohsnap.Arbitrary[[]byte] {
	return ohsnap.Map(params, func(p P) []byte {
		var buf bytes.Buffer
		if _, err := p.Instr().Encode(&buf, arch.EncOpts{}); err != nil {
			return nil
		}

		return buf.Bytes()
	})
}
