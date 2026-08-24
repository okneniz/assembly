package arm64

import (
	"sort"
	"sync"
)

// A map from decode schemas to generated armISA entries.
// The handwritten Mask/Value in schemas.go STAY — they are the curation key:
// by them the map finds its XML entry and proves consistency (an inconsistent
// schema simply does not get mapped).
//
// Entry selection policy (conservative; classes — see TestISAMapAnalysis):
//
//	exact    — the schema's {Mask,Value} matches the entry's {Match,Mask} bit-for-bit;
//	widen-sf — the entry is looser by exactly bit 31 (sf): handwritten 64/32-bit
//	           forms of one instruction (add/add w-form, ...) refer to one
//	           iclass; the ctor must read sf from the word — the diff gates
//	           (TestDisasmVsObjdump, round-trip) watch over this.
//
// NOT mapped (stay on handwritten bits):
//   - widening >1 bit — the XML's soft should-be bits (Ra=11111 in umulh,
//     the domain in dmb, size fields): the substitution would widen the match
//     to invalid words;
//   - families (ldp also covers the XML stlp classes) and aliases (XML has none);
//   - conflicts (5 of them, handled per-case).
var (
	isaBitsOnce sync.Once
	isaBitsIdx  []int // schema index → armISA index; -1 = handwritten bits
	isaTail     []int // armISA indexes: entries not mapped to schema bits (tail B)
)

// isaTailEntry returns the tail entry at position (0..isaTailLen-1) or nil.
func isaTailEntry(k int) *armISAEntry {
	isaBitsOnce.Do(buildISAOverride)
	if k < 0 || k >= len(isaTail) {
		return nil
	}

	return &armISA[isaTail[k]]
}

// isaTailLen — the size of the generated tail (diagnostics for tests).
func isaTailLen() int {
	isaBitsOnce.Do(buildISAOverride)
	return len(isaTail)
}

// schemaISAEntry returns the generated armISA entry whose {Match, Mask} are
// used for matching the i-th schema, or nil (handwritten bits).
func schemaISAEntry(i int) *armISAEntry {
	isaBitsOnce.Do(buildISAOverride)
	if j := isaBitsIdx[i]; j >= 0 {
		return &armISA[j]
	}

	return nil
}

// overriddenCount — how many schemas are decoded with generated bits
// (diagnostics for tests).
func overriddenCount() int {
	isaBitsOnce.Do(buildISAOverride)
	n := 0
	for _, j := range isaBitsIdx {
		if j >= 0 {
			n++
		}
	}

	return n
}

// widenExcluded — families for which mask widening breaks output: the ctor
// does not read the freed bits from the word (caught by a gate: +36
// fcvt/fcvtzs mismatches — w instead of x, swapped operand types).
var widenExcluded = map[string]bool{
	// fcvt/fcvtzs/fcvtzu REMOVED: operand types are read from the word
	// (fp2.go, fpTypeBits).
	// The N bit (bit22) is an OPCODE distinguishing ubfm/sbfm, not an operand:
	// freeing it let the ubfm ctor accept SBFM encodings (+1 hard in round-trip)
	"ubfm": true, "sbfm": true, "bfm": true,
}

// widenAcceptable — whether schema s's bits may be replaced by entry g's
// generated bits: every freed bit must lie inside a declared FIELD of the
// entry (the generator declares as fields only named blocks with soft
// bracketed c-values — operands; bare 0/1 are opcode pins) or be the sf bit.
func widenAcceptable(s *Schema, g *armISAEntry) bool {
	if widenExcluded[s.Meta.Name] {
		return false
	}

	freed := s.Mask &^ g.Mask
	for b := range uint(32) {
		if freed>>b&1 == 0 || b == 31 {
			continue
		}

		inField := false
		for _, f := range g.Fields {
			if b >= f.Offset && b < f.Offset+f.Width {
				inField = true
				break
			}
		}

		if !inField {
			return false
		}
	}

	return true
}

func buildISAOverride() {
	byName := make(map[string][]int, len(armISA))
	for i := range armISA {
		byName[armISA[i].Name] = append(byName[armISA[i].Name], i)
	}

	isaBitsIdx = make([]int, len(arm64Schemas))
	for si := range arm64Schemas {
		isaBitsIdx[si] = -1
		s := &arm64Schemas[si]
		cands := byName[s.Meta.Name]
		// Pass 1: exact match (priority — deterministic choice).
		for _, gi := range cands {
			g := &armISA[gi]
			if g.Mask == s.Mask && g.Match == s.Value {
				isaBitsIdx[si] = gi
				break
			}
		}

		if isaBitsIdx[si] >= 0 {
			continue
		}

		// Pass 2: the entry is looser, but every freed bit must lie inside a
		// declared FIELD of the XML entry (the generator declares as fields
		// only named blocks with soft bracketed c-values — Ra in umulh; bare
		// 0/1 blocks stay opcode pins). The sf bit is an empirically safe
		// class. The gates are watching.
		//
		// Exceptions — ctors that HARD-CODE the type bits (22:23) from an
		// argument instead of reading them from the word (an sf-convention
		// violation): the fcvt family.
		if widenExcluded[s.Meta.Name] {
			continue
		}

		for _, gi := range cands {
			g := &armISA[gi]
			if (g.Mask&s.Mask) != g.Mask || (s.Value&g.Mask) != g.Match {
				continue
			}

			if isaBitsIdx[si] >= 0 && armISA[isaBitsIdx[si]].Mask != g.Mask {
				continue
			}

			if widenAcceptable(s, g) {
				isaBitsIdx[si] = gi
			}
		}
	}

	// Tail B: armISA entries NOT mapped to schemas (isaBitsIdx targets are
	// covered by schemas by construction and are redundant in the tail). A
	// name-based exception ("the name exists among the schemas") proved too
	// coarse: entries of families with soft bits (dup/dmb/umulh/...) left
	// words unnamed even though their encoding is known from XML. Order — by
	// specificity (mask popcount desc., then match asc.): deterministic,
	// more specific ones tried first. The tail is matched AFTER all schemas,
	// so the curated schema order is not broken by construction (all 1292
	// ordering-overlaps from the A2 failure are excluded by the tail's
	// placement alone).
	mapped := make(map[int]struct{}, len(arm64Schemas))
	for _, gi := range isaBitsIdx {
		if gi >= 0 {
			mapped[gi] = struct{}{}
		}
	}

	for i := range armISA {
		if _, redundant := mapped[i]; !redundant {
			isaTail = append(isaTail, i)
		}
	}

	sort.Slice(isaTail, func(a, b int) bool {
		ga, gb := &armISA[isaTail[a]], &armISA[isaTail[b]]
		na, nb := popcount(ga.Mask), popcount(gb.Mask)
		if na != nb {
			return na > nb
		}

		return ga.Match < gb.Match
	})
}

// popcount — the number of set bits (for sorting the tail by specificity).
func popcount(m uint32) int {
	n := 0
	for m != 0 {
		n += int(m & 1)
		m >>= 1
	}

	return n
}
