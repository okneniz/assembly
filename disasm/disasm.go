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
// The instruction carries no length: that is a property of the byte stream
// (the per-arch InstrLen), not of the semantics.
type ObjDump interface {
	// ObjDump returns the instruction's representation in objdump style:
	// mnemonic and operands ("ldr x28, [x29, #24]"), without the address and
	// machine-code columns - those are added by this package. ctx is the
	// representation context (ViewCtx).
	ObjDump(ctx ViewCtx) string
}

// ViewCtx is the instruction representation context (view layer): text
// rendering parameters: the machine-code column style and the
// instruction's own address (PC-relative forms print absolute targets
// from it). Future representation parameters are added here as methods
// (like EncOpts in riscv - encoding modes).
type ViewCtx interface {
	// Style is the machine-code column style: space-separated bytes (Mach-O)
	// or a hex word (ELF); see text.StyleFor.
	Style() text.CodeStyle

	// Addr is the address of the instruction itself: PC-relative operands
	// (branches, literal loads, adrp) render their absolute targets as
	// ctx.Addr() + offset.
	Addr() uint64
}

// viewCtx is the canonical implementation of ViewCtx.
type viewCtx struct {
	style text.CodeStyle
	addr  uint64
}

func (c viewCtx) Style() text.CodeStyle {
	return c.style
}

func (c viewCtx) Addr() uint64 {
	return c.addr
}

// DefaultViewCtx is the representation context without an address (zero):
// PC-relative operands render relative to 0.
func DefaultViewCtx() ViewCtx {
	return viewCtx{style: text.CodeBytes}
}

// ViewCtxAt is the representation context of an instruction at addr
// (byte style): PC-relative operands render absolute targets.
func ViewCtxAt(addr uint64) ViewCtx {
	return viewCtx{style: text.CodeBytes, addr: addr}
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

// ctx is the ViewCtx derived from the output parameters and the
// instruction's address.
func (o Options) ctx(addr uint64) ViewCtx {
	return viewCtx{style: o.Style, addr: addr}
}

// Line returns a single instruction's line: "<addr>:\t<code>\t<text>".
// code is exactly this instruction's bytes (the caller walks the stream
// with the per-arch InstrLen and slices them out).
func Line(addr uint64, code []byte, in ObjDump, opts Options) string {
	raw, n := instrBytes(code)
	return fmt.Sprintf(
		"%x:\t%s\t%s",
		addr,
		text.FormatCode(raw, n, opts.Style),
		in.ObjDump(opts.ctx(addr)),
	)
}

// Write disassembles the buffer code (located at address base) into the
// instructions instrs - one line per instruction (the Line format).
// size is the arch's stream rule (InstrLen): the byte length of the
// instruction at the head of the remaining code. Write is generic so that
// it can accept slices of concrete instructions ([ ]arm64.Instr,
// [ ]riscv.Instr) without manually converting them to an interface slice.
func Write[T ObjDump](w io.Writer, base uint64, code []byte, instrs []T, opts Options, size func(code []byte) int) error {
	off := 0
	for _, in := range instrs {
		if off < len(code) {
			n := size(code[off:])
			if _, err := fmt.Fprintln(w, Line(base+uint64(off), code[off:off+n], in, opts)); err != nil {
				return err
			}

			off += n
		}
	}

	return nil
}

// instrBytes is the buffer as a little-endian word (the 2 or 4 bytes of an
// instruction encoding). A truncated buffer tail is returned as is.
func instrBytes(code []byte) (uint32, int) {
	n := len(code)
	if n > 4 {
		n = 4
	}

	var raw uint32
	for i := range n {
		raw |= uint32(code[i]) << (8 * i)
	}

	return raw, n
}
