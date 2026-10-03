package linktests

// The monolithic law, one subtest per architecture: linking the sources
// of a specimen SEPARATELY-assembled is the same program as assembling
// their concatenation - the same text bytes, the same data bytes, the
// same memory size, the same global addresses (the .L locals differ by
// design: a link namespaces them behind the file key, a monolith keeps
// them bare - the laws compare the shared namespace only).
//
// The domain is quad-free specimens: a quad rides behind an .align, and
// alignment is SECTION-scoped - each file's data is its own section in
// a link (the pads are per file), one merged section in a monolith (the
// pads are global). Both are honest layouts, just not the same one -
// exactly as per-object sections differ from a merged output section in
// every real linker. The quads stay covered by the resolve, image and
// exec laws.

import (
	"bytes"
	"maps"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/link"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
	"github.com/okneniz/assembly/unit"
)

func TestLinkMatchesMonolithic(t *testing.T) {
	for _, ac := range arches() {
		t.Run(ac.name, func(t *testing.T) {
			arb := linkarb.NewArchProgramArb(seedRnd(t), ac.arch)
			ohsnap.Check(t, 2000, arb, func(a linkarb.ArchProgram) bool {
				if a.Prog.HasQuad() || a.Prog.LoopingFiles() > 1 {
					return true // outside the domain, see the notes above
				}

				srcs := sources(a)

				linked := link.Link(ac.deps(), srcs, "", flatPlace)
				if len(linked.Errs) != 0 {
					t.Logf("link errors: %v", linked.Errs)
					return false
				}

				var raw []string
				for i := range srcs {
					raw = append(raw, srcs[i].Src)
				}

				u := unit.New()
				if errs := ac.assembleUnit(u, "mono.s", strings.Join(raw, "\n")); len(errs) > 0 {
					return false // outside the property domain (a broken source)
				}

				mono := u.Resolve(flatPlace)
				if len(mono.Errs) != len(linked.Errs) {
					return false
				}

				lt, lerr := linked.EncodeText()
				mt, merr := mono.EncodeText()
				if lerr != nil || merr != nil {
					return false
				}

				ld, lderr := linked.EncodeData()
				md, mderr := mono.EncodeData()
				if lderr != nil || mderr != nil {
					return false
				}

				if !bytes.Equal(lt, mt) {
					t.Logf("text bytes differ: link %d, mono %d", len(lt), len(mt))
					return false
				}

				if !bytes.Equal(ld, md) {
					t.Logf("data bytes differ: link %d, mono %d", len(ld), len(md))
					return false
				}

				if linked.DataMem != mono.DataMem {
					t.Logf("datamem %d, mono %d", linked.DataMem, mono.DataMem)
					return false
				}

				if !maps.Equal(globalsOf(linked.Syms), globalsOf(mono.Syms)) {
					t.Logf("global symbols differ: link %v, mono %v",
						globalsOf(linked.Syms), globalsOf(mono.Syms))
					return false
				}

				return true
			})
		})
	}
}

// globalsOf keeps the shared namespace only: the .L locals of a link
// sit behind their "file:" keys, the ones of a monolith sit bare - both
// drop out.
func globalsOf(syms map[string]uint64) map[string]uint64 {
	out := make(map[string]uint64, len(syms))
	for name, addr := range syms {
		if strings.Contains(name, ":") || strings.HasPrefix(name, ".L") {
			continue
		}

		out[name] = addr
	}

	return out
}
