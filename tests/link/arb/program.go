// Package arb - the generator of link specimens: N assembly sources
// (the axes: which file defines what, the direction of every reference,
// the local/global/numeric label kinds, the data/bss shapes) rendered
// TOGETHER with the value the linked program must exit with - the same
// seed produces the sources and their passport, so a property checks
// the program's meaning, not its bytes.
//
// The mini-ABI: one accumulator register per dialect (w0/a0/$a0) is the
// in-out argument of every call, each function adds its constant (and
// its loop count) into it, the entry seeds the accumulator, calls a
// sequence of functions and exits. Everything else - nops, branch
// dances, .word padding, data words, quad pointers, bss reserves -
// changes sizes and layout only, never the value.
//
// The specimen structure and its value are dialect-free; the spelling
// is one of three (Arch): arm64, riscv64, LoongArch64 - the accumulator
// and the exit idiom follow the dialect, the arithmetic does not.
package arb

import (
	"fmt"
	"strings"
)

// Arch is the instruction dialect a program renders in.
type Arch int

const (
	// Arm64 - macOS (svc #0x80) and Linux (svc #0) exits.
	Arm64 Arch = iota
	// Riscv - Linux (li a7, 93; ecall).
	Riscv
	// Loong64 - Linux ($ registers, syscall 0).
	Loong64
)

// dialect is the spelling of the mini-ABI verbs of one arch.
type dialect struct {
	init   func(v uint32) string // seed the accumulator
	loop   func(n int) string    // the whole countdown block (acc += n)
	call   func(name string) string
	add    func(c uint32) string // acc += c
	ret    string
	exit   string // the Linux exit idiom
	exitMX string // the macOS one (arm64 alone)
	jump   string // the noise branch over one line
}

func dialectOf(a Arch) dialect {
	switch a {
	case Riscv:
		return dialect{
			init: func(v uint32) string { return fmt.Sprintf("    li a0, %d\n", v) },
			loop: func(n int) string {
				return fmt.Sprintf(
					"    li t0, %d\n.Lloop:\n"+
						"    addi a0, a0, 1\n"+ //nolint:dupword // the instruction spells its operands
						"    addi t0, t0, -1\n"+ //nolint:dupword // the instruction spells its operands
						"    bne t0, zero, .Lloop\n",
					n,
				)
			},
			call: func(name string) string { return fmt.Sprintf("    call %s\n", name) },
			add: func(c uint32) string {
				//nolint:dupword // the instruction spells its operands
				return fmt.Sprintf("    addi a0, a0, %d\n", c)
			},
			ret:  "    ret\n",
			exit: "    li a7, 93\n    ecall\n",
			jump: "    j 1f\n",
		}
	case Loong64:
		return dialect{
			init: func(v uint32) string { return fmt.Sprintf("    li.d $a0, %d\n", v) },
			loop: func(n int) string {
				return fmt.Sprintf(
					"    addi.w $t0, $zero, %d\n.Lloop:\n"+
						"    addi.d $a0, $a0, 1\n"+
						"    addi.w $t0, $t0, -1\n"+
						"    bne $t0, $zero, .Lloop\n",
					n,
				)
			},
			call: func(name string) string { return fmt.Sprintf("    bl %s\n", name) },
			add:  func(c uint32) string { return fmt.Sprintf("    addi.d $a0, $a0, %d\n", c) },
			ret:  "    jirl $zero, $ra, 0\n",
			exit: "    addi.d $a7, $zero, 93\n    syscall 0\n",
			jump: "    b 1f\n",
		}
	default:
		return dialect{
			init: func(v uint32) string { return fmt.Sprintf("    movz w0, #%d\n", v) },
			loop: func(n int) string {
				return fmt.Sprintf(
					"    movz w9, #%d\n.Lloop:\n"+
						"    add w0, w0, #1\n"+ //nolint:dupword // the instruction spells its operands
						"    subs w9, w9, #1\n"+ //nolint:dupword // the instruction spells its operands
						"    b.ne .Lloop\n",
					n,
				)
			},
			call: func(name string) string { return fmt.Sprintf("    bl %s\n", name) },
			add: func(c uint32) string {
				//nolint:dupword // the instruction spells its operands
				return fmt.Sprintf("    add w0, w0, #%d\n", c)
			},
			ret:    "    ret\n",
			exit:   "    movz x8, #93\n    svc #0\n",
			exitMX: "    movz x16, #0x200, lsl #16\n    movk x16, #0x1\n    svc #0x80\n",
			jump:   "    b 1f\n",
		}
	}
}

// Program is one generated link specimen.
type Program struct {
	// ExitLinux selects the exit idiom of the entry epilogue on arm64
	// (macOS svc #0x80 vs Linux svc #0); the other dialects are
	// Linux-only and ignore the flag. The exit value never changes.
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

// FuncSpec is one global function: it adds Add into the accumulator
// (plus one per loop round) and returns.
type FuncSpec struct {
	name  string
	add   uint32
	loop  int // 0: none; 1..12 rounds (at most one loop per file: .Lloop)
	noise []string
}

// EntrySpec is the _start of the program: it seeds the accumulator,
// runs its own loop rounds, calls the sequence and exits.
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

// Exit is the process exit code of the program (the kernel keeps the
// low byte of the accumulator).
func (p Program) Exit() int {
	return int(p.Value() & 0xFF)
}

// render writes one file: .global for the entry and every function
// (the GAS link interface - a bare label is a LOCAL object symbol,
// invisible to a real linker), then the functions in order (the entry
// first among them), then the data items and the bss reserve.
func (s FileSpec) render(idx int, arch Arch, exitLinux bool) string {
	d := dialectOf(arch)
	var b strings.Builder

	if s.entry != nil {
		b.WriteString(".global _start\n")
	}

	for i := range s.funcs {
		fmt.Fprintf(&b, ".global %s\n", s.funcs[i].name)
	}

	b.WriteString(".text\n")

	if s.entry != nil {
		s.entry.render(&b, d, exitLinux)
	}

	for i := range s.funcs {
		s.funcs[i].render(&b, d)
	}

	if len(s.data) > 0 {
		b.WriteString(".data\n")
		for i, dm := range s.data {
			fmt.Fprintf(&b, "d%d_%d:\n", idx, i)
			if dm.quad {
				// the pointer 8-aligned (as any real source has it -
				// ld's chained fixups refuse an unaligned pointer)
				b.WriteString("    .align 3\n")
				fmt.Fprintf(&b, "    .quad %s\n", dm.name)
			} else {
				fmt.Fprintf(&b, "    .word %d\n", dm.value)
			}
		}
	}

	if s.bss > 0 {
		fmt.Fprintf(&b, ".bss\nz%d:\n    .zero %d\n", idx, s.bss)
	}

	return b.String()
}

// render writes the entry: the seed, the loop, the calls, the
// epilogue. The noise rides AFTER the exit - off the execution path
// (size and resolution coverage without executing it).
func (e EntrySpec) render(b *strings.Builder, d dialect, exitLinux bool) {
	b.WriteString("_start:\n")
	b.WriteString(d.init(e.init))
	renderLoopInto(b, d, e.loop)

	for _, name := range e.calls {
		b.WriteString(d.call(name))
	}

	if d.exitMX != "" && !exitLinux {
		b.WriteString(d.exitMX)
	} else {
		b.WriteString(d.exit)
	}

	renderNoiseInto(b, e.noise)
}

// render writes one function: the loop, the add, the return, then the
// noise past the ret (never executed: the fall-through ends at ret).
func (f FuncSpec) render(b *strings.Builder, d dialect) {
	fmt.Fprintf(b, "%s:\n", f.name)
	renderLoopInto(b, d, f.loop)
	b.WriteString(d.add(f.add))
	b.WriteString(d.ret)
	renderNoiseInto(b, f.noise)
}

// renderLoopInto writes the countdown loop (acc += loop) - the local
// label is the same .Lloop in every file of every program: the
// isolation axis.
func renderLoopInto(b *strings.Builder, d dialect, rounds int) {
	if rounds <= 0 {
		return
	}

	b.WriteString(d.loop(rounds))
}

// renderNoiseInto writes the size-only lines (never touched by
// execution: every function ends in ret/exit before the next one
// begins).
func renderNoiseInto(b *strings.Builder, noise []string) {
	for _, n := range noise {
		b.WriteString(n)
	}
}

// HasQuad - the specimen carries a .quad pointer (an .align rides with
// it; the monolithic law's domain excludes it, see its note).
func (p Program) HasQuad() bool {
	for i := range p.files {
		for _, d := range p.files[i].data {
			if d.quad {
				return true
			}
		}
	}

	return false
}

// LoopingFiles counts the files owning a .Lloop (the label is shared
// across files BY DESIGN - a monolith of several looping files would
// define it twice, so the monolithic law's domain keeps at most one).
func (p Program) LoopingFiles() int {
	n := 0
	for i := range p.files {
		fs := &p.files[i]
		if fs.entry != nil && fs.entry.loop > 0 {
			n++
			continue
		}

		for _, f := range fs.funcs {
			if f.loop > 0 {
				n++
				break
			}
		}
	}

	return n
}

// Sources renders every file of the program in the dialect.
func (p Program) Sources(arch Arch) []Source {
	out := make([]Source, len(p.files))
	for i := range p.files {
		out[i] = Source{
			Name: fmt.Sprintf("t%d.s", i),
			Src:  p.files[i].render(i, arch, p.ExitLinux),
		}
	}

	return out
}

// Value is the accumulator value the program exits with (32-bit
// wrapping arithmetic, as in the registers).
func (p Program) Value() uint32 {
	e := p.files[p.entry].entry
	v := e.init + uint32(e.loop)
	for _, name := range e.calls {
		f := p.find(name)
		v += f.add + uint32(f.loop)
	}

	return v
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
