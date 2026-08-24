package arm64

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file"
)

// TestISAMapAnalysis builds the arm64Schemas → armISA correspondence map
// and classifies every entry. The classes decide an entry's fate when
// switching the decoder to generated bits:
//
//	exact       — the schema's {Mask,Value} matches the XML entry's {Match,Mask}
//	              bit-for-bit: the substitution is safe, decode behavior
//	              does not change.
//	genLooser   — the XML entry is looser (its mask ⊆ the schema's mask,
//	              values agree): the substitution WIDENS the match set — this
//	              is the variant ambiguity that killed armISA-merge (3999
//	              mismatches). A candidate for manual review; safe only if
//	              the widened words are not stolen from later schemas.
//	genNarrower — the schema's mask ⊆ the entry's mask (the handwritten
//	              entry covers a family of iclasses, e.g. ldp+stp): the
//	              substitution NARROWS the set — words fall through to the
//	              tail/unknown. The bits stay handwritten.
//	absent      — the name is not in the XML (aliases: b.eq..b.nv, mov,
//	              cmp...). The bits stay handwritten forever (an alias has no
//	              encoding of its own).
//	consistent  — only the shared bits agree (field boundaries/sf).
//	conflict    — agrees with no entry of the same name.
func TestISAMapAnalysis(t *testing.T) {
	byName := make(map[string][]int, len(armISA))
	for i := range armISA {
		byName[armISA[i].Name] = append(byName[armISA[i].Name], i)
	}

	var exact, genLooser, looserAmb, genNarrower, absent, consistent, conflict int
	var looserExamples, narrowerExamples []string
	widenHist := map[int]int{} // widened bit → number of schemas
	widenSF := 0               // widening by exactly bit 31 (sf)
	widenField := 0            // freed ⊆ {sf} ∪ XML fields, outside exclusions

	for si := range arm64Schemas {
		s := &arm64Schemas[si]
		cands := byName[s.Meta.Name]
		if len(cands) == 0 {
			absent++
			continue
		}

		exactIdx := -1
		for _, gi := range cands {
			g := &armISA[gi]
			if g.Mask == s.Mask && g.Match == s.Value {
				exactIdx = gi
				break
			}
		}

		if exactIdx >= 0 {
			exact++
			continue
		}

		var looser []int
		narrower := false
		for _, gi := range cands {
			g := &armISA[gi]
			// g is looser: g.Mask ⊆ s.Mask, the values on g's bits agree.
			if (g.Mask&s.Mask) == g.Mask && (s.Value&g.Mask) == g.Match {
				looser = append(looser, gi)
			}

			// g is narrower: s.Mask ⊆ g.Mask, the values on s's bits agree.
			if (s.Mask&g.Mask) == s.Mask && (g.Match&s.Mask) == s.Value {
				narrower = true
			}
		}

		switch {
		case len(looser) == 1:
			genLooser++
			g := &armISA[looser[0]]
			freed := s.Mask &^ g.Mask
			widenHist[popcount(freed)]++
			if widenAcceptable(s, g) {
				widenField++
			}

			if freed == 1<<31 {
				widenSF++
			}
		case len(looser) > 1:
			looserAmb++
			// the policy takes the FIRST subsuming candidate
			// (equivalent masks give the same verdict)
			if widenAcceptable(s, &armISA[looser[0]]) {
				widenField++
			}

			if n := len(looserExamples); n < 15 {
				looserExamples = append(looserExamples, fmt.Sprintf(
					"%-10s hw=%08x/%08x — %d looser XML entries (ambiguous)",
					s.Meta.Name, s.Value, s.Mask, len(looser)))
			}
		case narrower:
			genNarrower++
			if n := len(narrowerExamples); n < 15 {
				narrowerExamples = append(narrowerExamples, fmt.Sprintf(
					"%-10s hw=%08x/%08x — covers >=%d iclasses (family)",
					s.Meta.Name, s.Value, s.Mask, len(cands)))
			}
		default:
			agree := false
			for _, gi := range cands {
				g := &armISA[gi]
				shared := g.Mask & s.Mask
				if shared != 0 && (s.Value&shared) == (g.Match&shared) {
					agree = true
					break
				}
			}

			if agree {
				consistent++
			} else {
				conflict++
			}
		}
	}

	t.Logf(
		"A1 map: exact=%d genLooser=%d (of which ambiguous %d) genNarrower=%d absent=%d consistent=%d conflict=%d | total %d",
		exact,
		genLooser,
		looserAmb,
		genNarrower,
		absent,
		consistent,
		conflict,
		len(arm64Schemas),
	)
	t.Logf("genLooser widened-bits histogram: %v", widenHist)
	t.Logf("decodeOne on generated bits (exact+widen-sf): %d schemas", overriddenCount())
	// widenField includes sf splits (the predicate accepts freed={31})
	require.Equal(
		t,
		exact+widenField,
		overriddenCount(),
		"policy and analysis diverged: exact=%d widenField=%d",
		exact,
		widenField,
	)
	t.Log("genLooser with widening >1 bit (keep-hand candidates):")
	for _, e := range looserExamples {
		t.Log("  " + e)
	}

	t.Log("genNarrower (family; bits stay handwritten):")
	for _, e := range narrowerExamples {
		t.Log("  " + e)
	}
}

// TestISATailCoverage measures the effect of the generated tail: words not
// recognized by handwritten schemas get a mnemonic and raw fields (Generic)
// instead of an anonymous .word. Plus a synthetic check: a word built from
// a tail entry must decode into something meaningful (Generic or an
// intercepting schema) — but never into Unknown.
func TestISATailCoverage(t *testing.T) {
	ff, err := file.Detect("../../tests/examples/hello-world/hello-world")
	if err != nil {
		t.Skipf("example binary not available: %v", err)
	}

	ts, err := ff.CodeSection()
	if err != nil {
		t.Skipf("example binary not available: %v", err)
	}

	var unknown, generic int
	insts, err := Parse(ts.Addr)(bytes.Buffer(ts.Data))
	require.NoError(t, err)
	for _, in := range insts {
		switch in.(type) {
		case Unknown:
			unknown++
		case Generic:
			generic++
		}
	}

	t.Logf(
		"armISA tail: %d entries; on the hello-world corpus: generic=%d, .word=%d",
		isaTailLen(),
		generic,
		unknown,
	)

	step := max(isaTailLen()/200, 1)
	var gen, shadowedBySchema int
	for k := 0; k < isaTailLen(); k += step {
		e := isaTailEntry(k)
		w := e.Match | (0x9e3779b9 &^ e.Mask) // fixed field bits
		switch decodeOne(w, 0).(type) {
		case Generic:
			gen++
		default:
			shadowedBySchema++ // encoding shadowed by a handwritten schema — priority
		}
	}

	t.Logf("synthetic (step %d): generic=%d, intercepted by schemas=%d",
		step, gen, shadowedBySchema)
	if unknown > 0 {
		t.Logf("warning: %d words are still .word — the tail did not cover them", unknown)
	}
}
