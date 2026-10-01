// Package arb - the generator of link specimens: N assembly sources
// (the axes: which file defines what, the direction of every reference,
// the local/global/numeric label kinds, the data/bss shapes) rendered
// TOGETHER with the value the linked program must exit with - the same
// seed produces the sources and their passport, so a property checks
// the program's meaning, not its bytes.
//
// The mini-ABI: w0 is the accumulator (in and out of every call), each
// function adds its constant (and its loop count) into it, the entry
// seeds w0, calls a sequence of functions and exits. Everything else -
// nops, branch dances, .word padding, data words, quad pointers, bss
// reserves - changes sizes and layout only, never the value.
package arb

import (
	"fmt"
	"strings"
)

// Program is one generated link specimen.
type Program struct {
	// ExitLinux selects the exit idiom of the entry epilogue: macOS
	// (svc #0x80, x16 = 1<<16|1) or Linux (svc #0, x8 = 93). The exit
	// value is the same w0 either way.
	ExitLinux bool

	entry int        // the file holding _start
	files []FileSpec // in link order
}

// FileSpec is one source file of the specimen: its functions, the entry
// (one file of the program), its data items and its bss reserve.
type FileSpec struct {
	funcs []FuncSpec
	entry *EntrySpec
	data  []DataItem
	bss   int
}

// FuncSpec is one global function: it adds Add into w0 (plus one w0++
// per loop round) and returns.
type FuncSpec struct {
	name  string
	add   uint32
	loop  int // 0: none; 1..12 rounds (at most one loop per file: .Lloop)
	noise []string
}

// EntrySpec is the _start of the program: it seeds w0 with Init, runs
// its own loop rounds, calls the sequence and exits with w0.
type EntrySpec struct {
	init  uint32
	loop  int
	calls []string
	noise []string
}

// DataItem is one word of the data stream: a plain constant or an
// absolute pointer to a function (never dereferenced - the PIE gap;
// it rides along for the layout and the merge).
type DataItem struct {
	quad  bool // .quad name instead of .word const
	name  string
	value uint32
}

// Source is one rendered file of the program.
type Source struct {
	Name string
	Src  string
}

// Sources renders every file of the program.
func (p Program) Sources() []Source {
	out := make([]Source, len(p.files))
	for i := range p.files {
		out[i] = Source{
			Name: fmt.Sprintf("t%d.s", i),
			Src:  p.files[i].render(i, p.ExitLinux),
		}
	}

	return out
}

// Value is the w0 the program exits with (32-bit wrapping arithmetic,
// as in the registers).
func (p Program) Value() uint32 {
	e := p.files[p.entry].entry
	v := e.init + uint32(e.loop)
	for _, name := range e.calls {
		f := p.find(name)
		v += f.add + uint32(f.loop)
	}

	return v
}

// Exit is the process exit code of the program (the kernel keeps the
// low byte of w0).
func (p Program) Exit() int {
	return int(p.Value() & 0xFF)
}

// String renders the specimen for a failure report: the expected exit,
// the expected data size, every source.
func (p Program) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "exit=%d datamem=%d\n", p.Exit(), p.DataMem())
	for _, s := range p.Sources() {
		fmt.Fprintf(&b, "--- %s\n%s", s.Name, s.Src)
	}

	return b.String()
}

// DataMem is the expected memory size of the linked data stream: the
// data of every file (the quad pointers 8-aligned inside their file,
// as rendered - .align pads count) plus the bss of every file (the
// single tail).
func (p Program) DataMem() int {
	n := 0
	for i := range p.files {
		off := 0
		for _, d := range p.files[i].data {
			if d.quad {
				off += (8 - off%8) % 8
				off += 8
				continue
			}

			off += 4
		}

		n += off + p.files[i].bss
	}

	return n
}

// find locates a function by its global name.
func (p Program) find(name string) FuncSpec {
	for i := range p.files {
		for _, f := range p.files[i].funcs {
			if f.name == name {
				return f
			}
		}
	}

	return FuncSpec{name: name}
}

// render writes one file: .global for the entry and every function
// (the GAS link interface - a bare label is a LOCAL object symbol,
// invisible to a real linker), then the functions in order (the entry
// first among them), then the data items and the bss reserve.
func (s FileSpec) render(idx int, exitLinux bool) string {
	var b strings.Builder

	if s.entry != nil {
		b.WriteString(".global _start\n")
	}

	for i := range s.funcs {
		fmt.Fprintf(&b, ".global %s\n", s.funcs[i].name)
	}

	b.WriteString(".text\n")

	if s.entry != nil {
		s.entry.render(&b, exitLinux)
	}

	for i := range s.funcs {
		s.funcs[i].render(&b)
	}

	if len(s.data) > 0 {
		b.WriteString(".data\n")
		for i, d := range s.data {
			fmt.Fprintf(&b, "d%d_%d:\n", idx, i)
			if d.quad {
				// the pointer 8-aligned (as any real source has it -
				// ld's chained fixups refuse an unaligned pointer)
				b.WriteString("    .align 3\n")
				fmt.Fprintf(&b, "    .quad %s\n", d.name)
			} else {
				fmt.Fprintf(&b, "    .word %d\n", d.value)
			}
		}
	}

	if s.bss > 0 {
		fmt.Fprintf(&b, ".bss\nz%d:\n    .zero %d\n", idx, s.bss)
	}

	return b.String()
}

// render writes the entry: the seed, the loop, the calls, the
// epilogue. The noise rides AFTER the svc - off the execution path
// (size and resolution coverage without executing it).
func (e EntrySpec) render(b *strings.Builder, exitLinux bool) {
	b.WriteString("_start:\n")
	fmt.Fprintf(b, "    movz w0, #%d\n", e.init)
	e.renderLoop(b)

	for _, name := range e.calls {
		fmt.Fprintf(b, "    bl %s\n", name)
	}

	if exitLinux {
		b.WriteString("    movz x8, #93\n    svc #0\n")
	} else {
		b.WriteString(
			"    movz x16, #0x200, lsl #16\n    movk x16, #0x1\n    svc #0x80\n",
		)
	}

	e.renderNoise(b)
}

// render writes one function: the loop, the add, the return, then the
// noise past the ret (never executed: the fall-through ends at ret).
func (f FuncSpec) render(b *strings.Builder) {
	fmt.Fprintf(b, "%s:\n", f.name)
	f.renderLoop(b)
	//nolint:dupword // the instruction spells its operands
	fmt.Fprintf(b, "    add w0, w0, #%d\n", f.add)
	b.WriteString("    ret\n")
	f.renderNoise(b)
}

// renderLoop writes the countdown loop (w0 += loop) - the local label
// is the same .Lloop in every file of every program: the isolation
// axis.
func renderLoopInto(b *strings.Builder, rounds int) {
	fmt.Fprintf(b, "    movz w9, #%d\n", rounds)
	b.WriteString(".Lloop:\n")
	b.WriteString("    add w0, w0, #1\n")  //nolint:dupword // the instruction spells its operands
	b.WriteString("    subs w9, w9, #1\n") //nolint:dupword // the instruction spells its operands
	b.WriteString("    b.ne .Lloop\n")
}

func (e EntrySpec) renderLoop(b *strings.Builder) {
	if e.loop > 0 {
		renderLoopInto(b, e.loop)
	}
}

func (f FuncSpec) renderLoop(b *strings.Builder) {
	if f.loop > 0 {
		renderLoopInto(b, f.loop)
	}
}

// renderNoise writes the size-only lines (never touched by execution:
// every function ends in ret/svc before the next one begins).
func renderNoiseInto(b *strings.Builder, noise []string) {
	for _, n := range noise {
		b.WriteString(n)
	}
}

func (e EntrySpec) renderNoise(b *strings.Builder) { renderNoiseInto(b, e.noise) }
func (f FuncSpec) renderNoise(b *strings.Builder)  { renderNoiseInto(b, f.noise) }
