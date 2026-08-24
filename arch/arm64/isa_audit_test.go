package arm64

import (
	"fmt"
	"testing"
)

// TestISAAudit (P1) checks every hand-written arm64Schemas entry against the
// XML-generated armISA table. A hand-written encoding (Mask/Value) must be a
// SUB-ENCODING of some generated entry of the same mnemonic: i.e. the generated
// mask bits are all pinned by the hand-written mask, to the same values.
//
// Outcome categories:
//   - ok:      hand-written is a faithful refinement of an XML encoding.
//   - absent:  mnemonic not in the XML (typically alias mnemonics: mov, cmp, …).
//   - conflict: hand-written encoding disagrees with the XML (real discrepancy).
//
// A consistency diagnostic of the handwritten schemas against the XML
// table: not a gate, only a report. The schema→entry map — isa_map.go.
func TestISAAudit(t *testing.T) {
	var ok, absent, conflict, consistent int
	var conflicts []string

	for _, s := range arm64Schemas {
		name := s.Meta.Name
		var gen []armISAEntry
		for _, e := range armISA {
			if e.Name == name {
				gen = append(gen, e)
			}
		}

		if len(gen) == 0 {
			absent++
			continue
		}

		subsumed := false
		agreeShared := false
		for _, g := range gen {
			// g.Mask ⊆ s.Mask and the shared bits agree ⟹ hw ⊆ generated.
			if (g.Mask&s.Mask) == g.Mask && (s.Value&g.Mask) == g.Match {
				subsumed = true
				break
			}

			// otherwise: do they at least agree on the bits both pin?
			shared := g.Mask & s.Mask
			if shared != 0 && (s.Value&shared) == (g.Match&shared) {
				agreeShared = true
			}
		}

		switch {
		case subsumed:
			ok++
		case agreeShared:
			consistent++ // boundary/field difference (e.g. sf pinned by hw, field in XML)
		default:
			conflict++
			if len(conflicts) < 20 {
				conflicts = append(conflicts, fmt.Sprintf(
					"%-8s hw=0x%08x/0x%08x — disagrees with all %d XML entries",
					name,
					s.Value,
					s.Mask,
					len(gen),
				))
			}
		}
	}

	t.Logf("ISA audit: %d hand-written — ok=%d consistent(boundary)=%d absent=%d conflict=%d",
		len(arm64Schemas), ok, consistent, absent, conflict)
	for _, c := range conflicts {
		t.Log(c)
	}
}
