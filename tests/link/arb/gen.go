package arb

// The specimen generator: one random point of the link axes per call -
// the file count (1..6), which file holds the entry, which files define
// functions and which are data-only or empty, the reference directions
// (the entry may call forward and backward across the whole program,
// the same function twice), the label kinds (.Lloop locals of the same
// name in every looping file, numeric locals, globals), the data shapes
// (words, quad pointers to functions), the bss reserves - and the
// noise sizes that vary every layout dimension the linker owns.

import (
	"fmt"
	"iter"
	mrnd "math/rand/v2"
	"slices"
	"strings"

	ohsnap "github.com/okneniz/oh-snap"
)

// Gen walks the axes: one Program per call.
type Gen struct {
	rnd *mrnd.Rand
}

// NewGen is the generator over rnd.
func NewGen(rnd *mrnd.Rand) *Gen {
	return &Gen{rnd: rnd}
}

// Program builds one specimen in the dialect (the structure and the
// value are dialect-free; the spelling and the size-only noise follow
// it).
func (g *Gen) Program(arch Arch) Program {
	d := dialectOf(arch)
	n := 1 + g.rnd.IntN(6)
	files := make([]FileSpec, n)
	names := make([]string, 0, 2*n)

	for i := range files {
		for j := range g.rnd.IntN(3) {
			// fn* dodges the exact register names of the dialects (f1
			// is an FP register even in GAS - an operand cannot mean
			// the symbol then)
			name := fmt.Sprintf("fn%d_%d", i, j)
			// the add constant fits the signed imm12 of every dialect
			// (the arm64 unsigned imm12 takes it too)
			files[i].funcs = append(files[i].funcs, FuncSpec{
				name:  name,
				add:   uint32(1 + g.rnd.IntN(2047)),
				noise: g.noise(d),
			})
			names = append(names, name)
		}
	}

	// the entry axis first: any file, its seed, its call sequence
	// (forward and backward references, repeats - every callee defined)
	entry := g.rnd.IntN(n)
	entryLoop := 0
	if g.rnd.IntN(2) == 0 {
		entryLoop = 1 + g.rnd.IntN(12)
	}

	calls := make([]string, 0, 8)
	if len(names) > 0 {
		for range g.rnd.IntN(9) {
			calls = append(calls, names[g.rnd.IntN(len(names))])
		}
	}

	files[entry].entry = &EntrySpec{
		init:  uint32(g.rnd.IntN(65536)),
		loop:  entryLoop,
		calls: calls,
		noise: g.noise(d),
	}

	// the loop axis: the same .Lloop name in many files (the isolation
	// axis), at most one owner per file - the entry owns its file's
	// loop when it has one
	for i := range files {
		fs := &files[i]
		if len(fs.funcs) == 0 || g.rnd.IntN(3) != 0 {
			continue
		}

		if i == entry && entryLoop > 0 {
			continue
		}

		owner := &fs.funcs[g.rnd.IntN(len(fs.funcs))]
		owner.loop = 1 + g.rnd.IntN(12)
	}

	// the data and bss axes (layout only)
	for i := range files {
		for range g.rnd.IntN(4) {
			item := DataItem{value: g.rnd.Uint32()}
			if len(names) > 0 && g.rnd.IntN(3) == 0 {
				item.quad = true
				item.name = names[g.rnd.IntN(len(names))]
			}

			files[i].data = append(files[i].data, item)
		}

		if g.rnd.IntN(3) == 0 {
			files[i].bss = 1 + g.rnd.IntN(32)
		}
	}

	return Program{entry: entry, files: files}
}

// noise is 0..3 size-only lines (a nop, a numeric-local branch
// dance, .word padding) - the branch verb follows the dialect.
func (g *Gen) noise(d dialect) []string {
	out := make([]string, 0, 3)
	for range g.rnd.IntN(4) {
		switch g.rnd.IntN(3) {
		case 0:
			out = append(out, "    nop\n")
		case 1:
			out = append(out, d.jump+"    nop\n1:\n")
		default:
			out = append(out, fmt.Sprintf("    .word %d\n", g.rnd.IntN(65536)))
		}
	}

	return out
}

// clone is a deep copy of the specimen (every shrink candidate edits
// its own tree).
func (p Program) clone() Program {
	c := p
	c.files = make([]FileSpec, len(p.files))
	for i := range p.files {
		c.files[i] = p.files[i]
		c.files[i].funcs = slices.Clone(p.files[i].funcs)
		c.files[i].data = slices.Clone(p.files[i].data)

		if e := p.files[i].entry; e != nil {
			ec := EntrySpec{
				init:  e.init,
				loop:  e.loop,
				calls: slices.Clone(e.calls),
				noise: slices.Clone(e.noise),
			}
			c.files[i].entry = &ec
		}
	}

	return c
}

// ArchProgram binds a specimen to its instruction dialect: the value
// oh-snap holds and prints on failure carries the sources of THAT
// dialect (structure, exit and data bookkeeping are dialect-free).
type ArchProgram struct {
	Arch Arch
	Prog Program
}

// String renders the bound specimen for a failure report: the dialect,
// the expected exit, the expected data size, every source.
func (a ArchProgram) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "arch=%d exit=%d datamem=%d\n", a.Arch, a.Prog.Exit(), a.Prog.DataMem())
	for _, s := range a.Prog.Sources(a.Arch) {
		fmt.Fprintf(&b, "--- %s\n%s", s.Name, s.Src)
	}

	return b.String()
}

// archArb adapts the generator to oh-snap, one dialect at a time.
type archArb struct {
	gen  *Gen
	arch Arch
}

// NewArchProgramArb is the arbitrary of specimens in one dialect.
func NewArchProgramArb(rnd *mrnd.Rand, arch Arch) ohsnap.Arbitrary[ArchProgram] {
	return archArb{gen: NewGen(rnd), arch: arch}
}

func (a archArb) Generate() iter.Seq[ArchProgram] {
	return iter.Seq[ArchProgram](func(yield func(ArchProgram) bool) {
		for {
			if !yield(ArchProgram{Arch: a.arch, Prog: a.gen.Program(a.arch)}) {
				return
			}
		}
	})
}

func (a archArb) Shrink(v ArchProgram) iter.Seq[ArchProgram] {
	candidates := a.gen.shrink(v.Prog)
	out := make([]ArchProgram, len(candidates))
	for i, p := range candidates {
		out[i] = ArchProgram{Arch: v.Arch, Prog: p}
	}

	return slices.Values(out)
}

// shrink builds the candidates of one specimen: drop a noise line, drop
// a data item, drop a call, count a loop down - a failing specimen
// shrinks to a minimal lying program.
func (g *Gen) shrink(p Program) []Program {
	var out []Program

	add := func(edit func(*Program)) {
		c := p.clone()
		edit(&c)
		out = append(out, c)
	}

	for i := range p.files {
		fs := p.files[i]

		for j := range fs.funcs {
			f := fs.funcs[j]
			for k := range f.noise {
				add(func(c *Program) {
					n := &c.files[i].funcs[j].noise
					*n = slices.Delete(slices.Clone(*n), k, k+1)
				})
			}

			if f.loop > 1 {
				add(func(c *Program) { c.files[i].funcs[j].loop-- })
			}
		}

		if e := fs.entry; e != nil {
			for k := range e.noise {
				add(func(c *Program) {
					n := &c.files[i].entry.noise
					*n = slices.Delete(slices.Clone(*n), k, k+1)
				})
			}

			if e.loop > 1 {
				add(func(c *Program) { c.files[i].entry.loop-- })
			}

			for k := range e.calls {
				add(func(c *Program) {
					n := &c.files[i].entry.calls
					*n = slices.Delete(slices.Clone(*n), k, k+1)
				})
			}
		}

		for k := range fs.data {
			add(func(c *Program) {
				n := &c.files[i].data
				*n = slices.Delete(slices.Clone(*n), k, k+1)
			})
		}
	}

	return out
}
