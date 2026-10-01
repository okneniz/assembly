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

// Program builds one specimen.
func (g *Gen) Program() Program {
	n := 1 + g.rnd.IntN(6)
	files := make([]FileSpec, n)
	names := make([]string, 0, 2*n)

	for i := range files {
		for j := range g.rnd.IntN(3) {
			name := fmt.Sprintf("f%d_%d", i, j)
			files[i].funcs = append(files[i].funcs, FuncSpec{
				name:  name,
				add:   uint32(g.rnd.IntN(4096)),
				noise: g.noise(),
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
		noise: g.noise(),
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

// noise is 0..3 size-only lines (nops, a numeric-local branch dance,
// .word padding).
func (g *Gen) noise() []string {
	out := make([]string, 0, 3)
	for range g.rnd.IntN(4) {
		switch g.rnd.IntN(3) {
		case 0:
			out = append(out, "    nop\n")
		case 1:
			out = append(out, "    b 1f\n    nop\n1:\n")
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

// ProgramArb adapts the generator to oh-snap: every check draws a fresh
// specimen; shrink drops noise lines, data items and calls, and counts
// loop rounds down - a failing specimen shrinks to a minimal lying
// program.
type ProgramArb struct {
	gen *Gen
}

// NewProgramArb is the arbitrary over the generator.
func NewProgramArb(rnd *mrnd.Rand) ohsnap.Arbitrary[Program] {
	return ProgramArb{gen: NewGen(rnd)}
}

func (a ProgramArb) Generate() iter.Seq[Program] {
	return iter.Seq[Program](func(yield func(Program) bool) {
		for {
			if !yield(a.gen.Program()) {
				return
			}
		}
	})
}

func (a ProgramArb) Shrink(p Program) iter.Seq[Program] {
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

	return slices.Values(out)
}
