package arm64

import (
	"bytes"
	"fmt"
	"iter"
	mrnd "math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/link"
	"github.com/okneniz/assembly/unit"
)

// The twins: two sources linked through the writers.

func TestMachO(t *testing.T) {
	img, err := Macho([]link.Source{
		{
			File: "main.S",
			Src: `
.global _start
.text
_start:
    bl bump
    bl bump
    ldr x1, =val
    ldr w0, [x1]
    movz x16, #0x200, lsl #16
    movk x16, #0x1
    svc #0x80
`,
		},
		{
			File: "util.S",
			Src: `
.data
val:
    .word 2
.text
bump:
    ldr x1, =val
    ldr w0, [x1]
    add w0, w0, #1
    str w0, [x1]
    ret
.bss
pad:
    .zero 8
`,
		},
	}, "")

	require.NoError(t, err)
	require.NotEmpty(t, img.Bytes())
}

func TestMachOErrors(t *testing.T) {
	// a duplicate definition across the sources
	_, err := Macho([]link.Source{
		{File: "a.S", Src: ".text\nx:\n  ret\n"},
		{File: "b.S", Src: ".text\nx:\n  ret\n"},
	}, "")
	require.ErrorContains(t, err, `label "x" redefined`)

	// no entry of the program's own
	_, err = Macho([]link.Source{
		{File: "a.S", Src: ".text\nfn:\n  ret\n"},
	}, "")
	require.ErrorContains(t, err, "no entry symbol")
}

func TestELF(t *testing.T) {
	blob, err := ELF([]link.Source{
		{File: "a.S", Src: ".text\nstart:\n  movz x0, #1\n  ret\n"},
		{File: "b.S", Src: ".data\nd:\n  .word 7\n"},
	}, "", 0x10000)

	require.NoError(t, err)
	require.NotEmpty(t, blob)
	require.Equal(t, []byte{0x7f, 'E', 'L', 'F'}, blob[:4])

	// e_entry: the start symbol - base + 0
	require.Equal(t, uint64(0x10000), binaryLE(blob[0x18:0x20]))
}

func binaryLE(b []byte) uint64 {
	var v uint64
	for _, x := range slices.Backward(b) {
		v = v<<8 | uint64(x)
	}

	return v
}

// --- the link law: linking separately-assembled sources is the same
// program as assembling their concatenation ------------------------------

// linkPair - the bodies of two sources with disjoint label namespaces
// (the globals ga*/gb*, the locals .LA*/.LB*) that reference each
// other's globals: the shape the law is proved on. The skeletons are
// fixed; the bodies vary.
type linkPair struct {
	a []string
	b []string
}

// sourceA/B render a side: the global entry, the body, the local label,
// the local-referencing tail, the data quad of the other side's global,
// and a bss reserve.
func (p linkPair) sourceA() string {
	return ".text\nga1:\n" + strings.Join(p.a, "\n") +
		"\n.LAret:\n  cbz x0, .LAret\n  bl gb1\n  ret\n" +
		".data\npda:\n  .quad gb1\n.bss\nzba:\n  .zero 8\n"
}

func (p linkPair) sourceB() string {
	return ".text\ngb1:\n" + strings.Join(p.b, "\n") +
		"\n.LBret:\n  cbnz x0, .LBret\n  bl ga1\n  ret\n" +
		".data\npdb:\n  .quad ga1\n.bss\nzbb:\n  .zero 4\n"
}

// pairArb - random bodies of safe instructions (register moves, imm
// arithmetic, branches to the own labels, calls into the other source):
// 1..12 lines per side; shrink drops lines. No literal pools - see the
// law's domain note.
type pairArb struct {
	rnd *mrnd.Rand
}

func (p pairArb) Generate() iter.Seq[linkPair] {
	return iter.Seq[linkPair](func(yield func(linkPair) bool) {
		for {
			if !yield(linkPair{
				a: p.genSide(".LAret", "gb1"),
				b: p.genSide(".LBret", "ga1"),
			}) {
				return
			}
		}
	})
}

func (p pairArb) Shrink(v linkPair) iter.Seq[linkPair] {
	shrunk := func(lines []string, replace func(linkPair, []string) linkPair) []linkPair {
		out := make([]linkPair, 0, len(lines))
		for i := range lines {
			dropped := slices.Delete(slices.Clone(lines), i, i+1)
			out = append(out, replace(v, dropped))
		}

		return out
	}

	out := append(
		shrunk(v.a, func(p linkPair, lines []string) linkPair {
			p.a = lines
			return p
		}),
		shrunk(v.b, func(p linkPair, lines []string) linkPair {
			p.b = lines
			return p
		})...,
	)

	return slices.Values(out)
}

// genSide - one side's body: n random lines from the safe pool.
func (p pairArb) genSide(local, other string) []string {
	pools := []func() string{
		func() string { return fmt.Sprintf("  movz x0, #%d", p.rnd.IntN(4096)) },
		func() string { return fmt.Sprintf("  add x2, x2, #%d", p.rnd.IntN(4096)) },
		func() string { return fmt.Sprintf("  sub x3, x3, #%d", p.rnd.IntN(4096)) },
		func() string { return "  nop" },
		func() string { return "  b " + local },
		func() string { return "  bl " + other },
		func() string { return fmt.Sprintf("  .word %d", p.rnd.IntN(65536)) },
	}

	n := 1 + p.rnd.IntN(12)
	out := make([]string, n)
	for i := range out {
		out[i] = pools[p.rnd.IntN(len(pools))]()
	}

	return out
}

// TestLinkMatchesMonolithic - the law of the stage: for label-disjoint
// sources, linking a.S and b.S separately-assembled is the SAME program
// as assembling their concatenation - the same text bytes, the same data
// bytes, the same memory size, the same global addresses. The domain is
// pool-free sources: a literal pool belongs to ITS source's section tail
// (as in every real linker - the pool of an object travels with it), so
// a pooled source is a different program, not a broken link.
func TestLinkMatchesMonolithic(t *testing.T) {
	place := func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return 0x1000, 0x8000
	}

	ohsnap.Check(t, 2000, pairArb{rnd: seedRnd(t)}, func(p linkPair) bool {
		linked := link.Link(
			link.Deps{Parse: alias.ParseSourceUnit},
			[]link.Source{
				{File: "a.S", Src: p.sourceA()},
				{File: "b.S", Src: p.sourceB()},
			},
			"",
			place,
		)

		u := unit.New()
		if errs := alias.AssembleUnit(u, "mono.S", p.sourceA()+"\n"+p.sourceB()); len(errs) > 0 {
			return false // outside the property domain (a broken source)
		}

		mono := u.Resolve(place)

		if len(linked.Errs) != len(mono.Errs) {
			return false
		}

		lt, lerr := linked.EncodeText()
		mt, merr := mono.EncodeText()
		if lerr != nil || merr != nil {
			return false
		}

		ld, lerr := linked.EncodeData()
		md, merr := mono.EncodeData()
		if lerr != nil || merr != nil {
			return false
		}

		for _, name := range []string{"ga1", "gb1", "pda", "pdb", "zba", "zbb"} {
			if linked.Syms[name] != mono.Syms[name] {
				return false
			}
		}

		return bytes.Equal(lt, mt) && bytes.Equal(ld, md) && linked.DataMem == mono.DataMem
	})
}

// seedRnd - the deterministic seed of the property suite (ASSEMBLY_SEED,
// default 42), logged so a failure can be reproduced.
func seedRnd(t *testing.T) *mrnd.Rand {
	t.Helper()

	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	t.Logf("seed: %d (ASSEMBLY_SEED)", seed)
	return arb.Rnd(seed)
}
