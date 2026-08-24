package arm64

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/file"
)

// TestISADirectMatch checks whether the generated armISA table
// can drive decode DIRECTLY — matching the same words the hand-written
// arm64Schemas matches — when ordered by mask specificity (most bits first).
// If the mnemonics agree across the test binary, the full rewrite (armISA drives
// decode) is viable and ordering-by-specificity suffices; mismatches pinpoint
// where curated ordering is load-bearing.
//
// Read-only experiment; does not change production decode.
func TestISADirectMatch(t *testing.T) {
	ff, err := file.Detect("../../tests/examples/hello-world/hello-world")
	if err != nil {
		t.Skipf("example binary not available: %v", err)
	}

	ts, err := ff.CodeSection()
	if err != nil {
		t.Skipf("example binary not available: %v", err)
	}

	// armISA ordered by mask specificity (popcount desc), then value — a
	// deterministic first-match order that tries more-specific encodings first.
	ordered := make([]armISAEntry, len(armISA))
	copy(ordered, armISA)
	bits := func(m uint32) int {
		n := 0
		for m != 0 {
			n += int(m & 1)
			m >>= 1
		}

		return n
	}
	for i := 1; i < len(ordered); i++ {
		for j := i; j > 0 && (bits(ordered[j].Mask) > bits(ordered[j-1].Mask) ||
			(bits(ordered[j].Mask) == bits(ordered[j-1].Mask) && ordered[j].Match < ordered[j-1].Match)); j-- {
			ordered[j], ordered[j-1] = ordered[j-1], ordered[j]
		}
	}

	// Current decode (hand-written) entry per address: mnemonic + mask + value.
	type hw struct {
		name        string
		mask, value uint32
	}
	cur, err := Parse(ts.Addr)(bytes.Buffer(ts.Data))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	hwByAddr := make(map[uint64]hw, len(cur))
	off := uint64(0)
	for _, in := range cur {
		w := ts.Data[off : off+4]
		word := binary.LittleEndian.Uint32(w)
		for _, sc := range getSchemas() {
			if (word & sc.Mask) == sc.Value {
				hwByAddr[in.Addr()] = hw{
					sc.Meta.Name,
					sc.Mask,
					sc.Value,
				}
				break
			}
		}

		off += 4
	}

	// Classify each word's armISA-direct match vs the hand-written match:
	//   same     — identical mnemonic.
	//   labeling — different mnemonic for the SAME instruction (alias): cmp/subs,
	//              lsl/ubfm, mov/orr, … Resolve each mnemonic to a canonical base
	//              via the ARM alias map; same base ⇒ alias-labeling (formatter).
	//   real     — different canonical instruction ⇒ genuine decode/ordering
	//              failure (e.g. umlal matched where ldp should be).
	//   more     — armISA decoded an instruction the hand-written table left as
	//              .word (armISA is MORE complete — not a failure).
	var same, labeling, realDiff, more, none int
	for pos := 0; pos+4 <= len(ts.Data); pos += 4 {
		w := uint32(
			ts.Data[pos],
		) | uint32(
			ts.Data[pos+1],
		)<<8 | uint32(
			ts.Data[pos+2],
		)<<16 | uint32(
			ts.Data[pos+3],
		)<<24
		addr := ts.Addr + uint64(pos)
		var eA *armISAEntry
		for i := range ordered {
			if w&ordered[i].Mask == ordered[i].Match {
				eA = &ordered[i]
				break
			}
		}

		eH, ok := hwByAddr[addr]
		if !ok {
			continue
		}

		if eA == nil {
			if eH.name == ".word" {
				continue
			}

			none++
			continue
		}

		if eA.Name == eH.name {
			same++
			continue
		}

		if eH.name == ".word" {
			more++ // armISA decoded what hand-written couldn't
			continue
		}

		if aliasBase(eA.Name) == aliasBase(eH.name) {
			labeling++ // alias / formatter-layer difference
		} else {
			realDiff++ // genuine decode/ordering failure
			if realDiff <= 25 {
				t.Logf("  REAL: armISA=%-10s hand=%-10s @0x%x w=0x%08x", eA.Name, eH.name, addr, w)
			}
		}
	}

	t.Logf(
		"armISA-direct vs hand-written: same=%d labeling(alias)=%d armISA-more=%d REAL-decode-fail=%d armISA-none=%d",
		same,
		labeling,
		more,
		realDiff,
		none,
	)
}

// aliasBase maps an ARM mnemonic to its canonical base instruction, collapsing
// aliases (cmp→subs, lsl→ubfm, mov→orr, …) so an alias and its base compare
// equal. Unknown mnemonics map to themselves.
func aliasBase(m string) string {
	// Conditional branches: armISA collapses b.<cond> into one entry named "b";
	// the hand-written table has b.eq/b.ls/…. Same instruction either way.
	if m == "b" || strings.HasPrefix(m, "b.") {
		return "b.cond"
	}

	switch m {
	case "cmp", "negs":
		return "subs"
	case "cmn":
		return "adds"
	case "neg":
		return "sub"
	case "lsl", "lsr", "ubfiz", "ubfx":
		return "ubfm"
	case "asr", "sxtb", "sxth", "sxtw", "sbfiz", "sbfx":
		return "sbfm"
	case "ror":
		return "extr"
	case "mvn":
		return "orn"
	case "mov", "movz", "movn", "orr":
		return "mov"
	case "mul":
		return "madd"
	case "mneg":
		return "msub"
	case "tst":
		return "ands"
	case "cset", "cinc":
		return "csinc"
	case "csetm", "cinv":
		return "csinv"
	case "cneg":
		return "csneg"
	}

	return m
}
