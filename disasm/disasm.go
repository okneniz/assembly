// Package disasm is the post-processing of decoded instructions into
// disassembly text. The package knows nothing about instruction structure: any
// architecture whose instructions implement ObjDump is disassembled by it as
// is. The line format is objdump style, "<addr>:\t<code>\t<mnemonic+operands>".
package disasm

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/text"
)

// ObjDump is the instruction interface for disassembly: an instruction turns
// itself into text (just as, for the assembler, it turns itself into bytes).
type ObjDump interface {
	// ObjDump returns the instruction's representation in objdump style:
	// mnemonic and operands ("ldr x28, [x29, #24]"), without the address and
	// machine-code columns - those are added by this package. ctx is the
	// representation context (ViewCtx).
	ObjDump(ctx ViewCtx) string

	// Len returns the length of the instruction in bytes (2/4): the
	// machine-code column is sliced and the addresses of subsequent
	// instructions are computed from it.
	Len() int
}

// ViewCtx is the instruction representation context (view layer): text
// rendering parameters. Today there is a single parameter - the machine-code
// column style; future representation parameters are added here as methods
// (like EncOpts in riscv - encoding modes).
type ViewCtx interface {
	// Style is the machine-code column style: space-separated bytes (Mach-O)
	// or a hex word (ELF); see text.StyleFor.
	Style() text.CodeStyle
}

// viewCtx is the canonical implementation of ViewCtx.
type viewCtx struct {
	style text.CodeStyle
}

func (c viewCtx) Style() text.CodeStyle {
	return c.style
}

// DefaultViewCtx is the canonical representation context (byte style).
func DefaultViewCtx() ViewCtx {
	return viewCtx{style: text.CodeBytes}
}

// Options is the output parameters (adapted into a ViewCtx in Line/Write).
type Options struct {
	// Style is the machine-code column style: space-separated bytes (Mach-O)
	// or a hex word (ELF); see text.StyleFor.
	Style text.CodeStyle
}

func NewOptions(style text.CodeStyle) Options {
	return Options{Style: style}
}

// ctx is the ViewCtx derived from the output parameters.
func (o Options) ctx() ViewCtx {
	return viewCtx{style: o.Style}
}

// Line returns a single instruction's line: "<addr>:\t<code>\t<text>".
// code is the buffer starting at this instruction (the first in.Len() bytes are used).
func Line(addr uint64, code []byte, in ObjDump, opts Options) string {
	raw, n := instrBytes(code, in.Len())
	return fmt.Sprintf(
		"%x:\t%s\t%s",
		addr,
		text.FormatCode(raw, n, opts.Style),
		in.ObjDump(opts.ctx()),
	)
}

// Write disassembles the buffer code (located at address base) into the
// instructions instrs - one line per instruction (the Line format). It is
// generic so that it can accept slices of concrete instructions
// ([ ]arm64.Instr, [ ]riscv.Instr) without manually converting them to an
// interface slice.
func Write[T ObjDump](w io.Writer, base uint64, code []byte, instrs []T, opts Options) error {
	off := 0
	for _, in := range instrs {
		if off < len(code) {
			if _, err := fmt.Fprintln(w, Line(base+uint64(off), code[off:], in, opts)); err != nil {
				return err
			}
		}

		off += in.Len()
	}

	return nil
}

// instrBytes is the first n bytes of the buffer as a little-endian word
// (n=2 -> a uint32 of two bytes). A truncated buffer tail is returned as is
// (shorter than n).
func instrBytes(code []byte, n int) (uint32, int) {
	if n != 2 {
		n = 4
	}

	if len(code) < n {
		n = len(code)
	}

	var raw uint32
	for i := range n {
		raw |= uint32(code[i]) << (8 * i)
	}

	return raw, n
}
