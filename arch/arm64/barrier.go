package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Barrier - dmb/dsb/isb with a shareability domain: the CRm field at
// bits 11:8 over the base word. isb takes Sy only (the bare "isb" of
// the text layer is its spelling), so Isb() carries no domain at all.

// BarrierDomain - the barrier shareability domain (the CRm nibble).
type BarrierDomain uint8

const (
	Oshst BarrierDomain = 0b0010 // outer shareable, stores only
	Osh   BarrierDomain = 0b0011 // outer shareable
	Nshst BarrierDomain = 0b0110 // non-shareable, stores only
	Nsh   BarrierDomain = 0b0111 // non-shareable
	Ishst BarrierDomain = 0b1010 // inner shareable, stores only
	Ish   BarrierDomain = 0b1011 // inner shareable
	St    BarrierDomain = 0b1110 // full system, stores only
	Sy    BarrierDomain = 0b1111 // full system, all accesses
)

// barrierNames - the BarrierDomain back to the option spelling (decode side;
// the map is total over the eight domains).
var barrierNames = map[BarrierDomain]string{
	Oshst: "oshst", Osh: "osh", Nshst: "nshst", Nsh: "nsh",
	Ishst: "ishst", Ish: "ish", St: "st", Sy: "sy",
}

// barrierBase - the fixed part of each mnemonic's word.
var barrierBase = map[string]uint32{
	"dmb": 0xD50330BF,
	"dsb": 0xD503309F,
	"isb": 0xD50330DF,
}

// Barrier - dmb/dsb/isb with a shareability domain.
type Barrier struct {
	base

	name   string
	domain BarrierDomain
}

// newBarrier - the Barrier constructor (the decoder calls it with the
// nibble it read; the typed ctors below are the user surface).
func newBarrier(b base, name string, domain BarrierDomain) (Barrier, error) {
	if _, ok := barrierBase[name]; !ok {
		return Barrier{}, unknownBarrier(name)
	}

	if _, ok := barrierNames[domain]; !ok {
		return Barrier{}, unknownDomain(name, domain)
	}

	return Barrier{base: b, name: name, domain: domain}, nil
}

func (i Barrier) ObjDump(_ disasm.ViewCtx) string {
	if i.name == "isb" {
		return i.name // llvm canon: the bare spelling is the print form
	}

	return i.name + " " + barrierNames[i.domain]
}

func (i Barrier) Encode(w io.Writer) (int64, error) {
	return writeWord(w, barrierBase[i.name]|uint32(i.domain)<<8)
}

// domainOf - the BarrierDomain of a text-layer option spelling (the assembler
// vocabulary; the caller reports the unknown option).
func domainOf(option string) (BarrierDomain, bool) {
	for d, name := range barrierNames {
		if name == option {
			return d, true
		}
	}

	return 0, false
}
