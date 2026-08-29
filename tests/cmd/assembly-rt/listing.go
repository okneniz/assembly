package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	parsecbytes "github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

// verdict - the outcome of the per-line check "the text assembles into
// exactly the same bytes".
type verdict int

const (
	verOK verdict = iota
	verAsmErr
	verEquiv
	verMismatch
)

// lineStats - counters of building the listing of one iteration.
type lineStats struct {
	instrLines  int
	fillerLines int
	reasons     map[string]int
	samples     []string
}

func newLineStats(reasons map[string]int) lineStats {
	return lineStats{reasons: reasons}
}

// addFiller records a filler with its reason and (the first few) samples.
func (s *lineStats) addFiller(reason string, addr uint64, line string) {
	s.fillerLines++
	s.reasons[reason]++
	if len(s.samples) < 8 {
		s.samples = append(s.samples, fmt.Sprintf("%#x: %q: %s", addr, line, reason))
	}
}

// buildListing builds a "clean" listing (without address and machine-code
// columns): the instruction line, if it assembles into exactly the original
// bytes, otherwise a .word/.half filler with the same bytes. The tail not
// covered by the decoder becomes .byte fillers. The total width of the lines
// equals len(data), so fully assembling the listing at base preserves the
// addresses of all instructions.
func buildListing(
	data []byte,
	base uint64,
	arch file.ArchKind,
	assemble func(string, uint64) (*asm.Result, []asm.AsmError),
) (string, lineStats) {
	var sb strings.Builder
	st := newLineStats(map[string]int{})
	off := 0

	switch arch {
	case file.ArchRISCV64:
		textOf := firstInstrTextRISCV
		insts, err := riscv.Parse(base)(parsecbytes.Buffer(data))
		if err != nil {
			panic(err) // unreachable: the buffer is in memory, backtracking always stays in bounds
		}

		for _, in := range insts {
			n := in.Len()
			line := cleanLine(in.ObjDump(disasm.DefaultViewCtx()))
			if line == "" || line == "<unknown>" {
				emitRaw(&sb, data[off:off+n])
				st.addFiller("unknown-riscv", in.Addr(), line)
			} else {
				switch classifyLine(assemble, line, in.Addr(), data[off:off+n], textOf) {
				case verOK:
					fmt.Fprintln(&sb, line)
					st.instrLines++
				case verAsmErr:
					emitRaw(&sb, data[off:off+n])
					st.addFiller("asm-error", in.Addr(), line)
				case verEquiv:
					emitRaw(&sb, data[off:off+n])
					st.addFiller("equiv-encoding", in.Addr(), line)
				case verMismatch:
					emitRaw(&sb, data[off:off+n])
					st.addFiller("hard-mismatch", in.Addr(), line)
				}
			}

			off += n
		}
	case file.ArchLOONGARCH64:
		textOf := firstInstrTextLOONG
		insts, err := loong64.Parse(base)(parsecbytes.Buffer(data))
		if err != nil {
			panic(err) // unreachable: the buffer is in memory, backtracking always stays in bounds
		}

		for _, in := range insts {
			line := cleanLine(in.ObjDump(disasm.DefaultViewCtx()))
			if line == "" || line == "<unknown>" {
				emitRaw(&sb, data[off:off+4])
				st.addFiller("unknown-loong64", in.Addr(), line)
			} else {
				switch classifyLine(assemble, line, in.Addr(), data[off:off+4], textOf) {
				case verOK:
					fmt.Fprintln(&sb, line)
					st.instrLines++
				case verAsmErr:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("asm-error", in.Addr(), line)
				case verEquiv:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("equiv-encoding", in.Addr(), line)
				case verMismatch:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("hard-mismatch", in.Addr(), line)
				}
			}

			off += 4
		}
	default:
		textOf := firstInstrTextARM
		insts, err := arm64.Parse(base)(parsecbytes.Buffer(data))
		if err != nil {
			panic(err) // unreachable: the buffer is in memory, backtracking always stays in bounds
		}

		for _, in := range insts {
			line := cleanLine(in.ObjDump(disasm.DefaultViewCtx()))
			if line == "" || isDecodeOnly(in) {
				emitRaw(&sb, data[off:off+4])
				st.addFiller("decode-only-arm64", in.Addr(), line)
			} else {
				switch classifyLine(assemble, line, in.Addr(), data[off:off+4], textOf) {
				case verOK:
					fmt.Fprintln(&sb, line)
					st.instrLines++
				case verAsmErr:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("asm-error", in.Addr(), line)
				case verEquiv:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("equiv-encoding", in.Addr(), line)
				case verMismatch:
					emitRaw(&sb, data[off:off+4])
					st.addFiller("hard-mismatch", in.Addr(), line)
				}
			}

			off += 4
		}
	}

	for ; off < len(data); off++ {
		fmt.Fprintf(&sb, ".byte 0x%02x\n", data[off])
		st.fillerLines++
		st.reasons["tail-byte"]++
	}

	return sb.String(), st
}

// iterate - one iteration: a listing from the current bytes, a full assembly
// of the listing with a single driver call (a large-scale test of two-pass
// operation), and checks of the hashes and of listing stability between
// iterations.
func iterate(
	cur []byte,
	base uint64,
	arch file.ArchKind,
	assemble func(string, uint64) (*asm.Result, []asm.AsmError),
	n int,
	prevListingSHA string,
) *iterResult {
	it := &iterResult{}

	listing, st := buildListing(cur, base, arch, assemble)
	it.rep = iterReport{
		N:           n,
		InstrLines:  st.instrLines,
		FillerLines: st.fillerLines,
		Reasons:     st.reasons,
		Samples:     st.samples,
		ListingSHA:  shaHex([]byte(listing)),
	}

	res, errs := assemble(listing, base)
	it.rep.AsmErrs = len(errs)
	for _, e := range errs {
		if len(it.rep.ErrSamples) < 8 {
			it.rep.ErrSamples = append(it.rep.ErrSamples,
				fmt.Sprintf("line %d:%d: %s", e.Line, e.Col, e.Msg))
		}
	}

	got := textSection(res)
	if got == nil {
		it.rep.Fail = "assembler returned no .text section"
		return it
	}

	it.rep.TextSHA = shaHex(got)
	it.text = got

	switch {
	case len(got) != len(cur):
		it.rep.Fail = fmt.Sprintf("layout: assembled %d bytes, original %d", len(got), len(cur))
	case it.rep.TextSHA != shaHex(cur):
		off := firstDiff(cur, got)
		it.rep.Fail = fmt.Sprintf(
			"bytes diverged at offset %#x: want %02x got %02x", off, cur[off], got[off])
	case prevListingSHA != "" && prevListingSHA != it.rep.ListingSHA:
		it.rep.Fail = "listing not stable across iterations"
	}

	return it
}

// classifyLine assembles the line separately (at its own address) and
// classifies the result: verOK - the bytes matched; verAsmErr - the text does
// not assemble; verEquiv - the bytes differ but decode into the same text
// (an equivalent encoding: RVC compression, canonical forms); verMismatch -
// a genuine divergence.
func classifyLine(
	assemble func(string, uint64) (*asm.Result, []asm.AsmError),
	line string,
	addr uint64,
	want []byte,
	textOf func([]byte, uint64) string,
) verdict {
	res, errs := assemble(line, addr)
	if len(errs) > 0 {
		return verAsmErr
	}

	got := textSection(res)
	if got == nil {
		return verMismatch
	}

	if bytes.Equal(got, want) {
		return verOK
	}

	if textOf(got, addr) == line {
		return verEquiv
	}

	return verMismatch
}

// firstInstrTextARM - the normalized text of the first instruction of a
// buffer (arm64).
func firstInstrTextARM(b []byte, addr uint64) string {
	insts, err := arm64.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return cleanLine(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// firstInstrTextRISCV - the normalized text of the first instruction of a
// buffer (riscv64, including compressed forms).
func firstInstrTextRISCV(b []byte, addr uint64) string {
	insts, err := riscv.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return cleanLine(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// firstInstrTextLOONG - the first instruction's text of a loong64 word
// buffer.
func firstInstrTextLOONG(b []byte, addr uint64) string {
	insts, err := loong64.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return cleanLine(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// isDecodeOnly - arm64 instructions without assembly text (the Generic tail
// of armISA and the Unknown fallback).
func isDecodeOnly(in arm64.Instr) bool {
	if _, ok := in.(arm64.Generic); ok {
		return true
	}

	_, ok := in.(arm64.Unknown)
	return ok
}

// textSection extracts .text from the assembly result (the driver's default
// section).
func textSection(res *asm.Result) []byte {
	for _, s := range res.Sections {
		if s.Name == ".text" {
			return s.Data
		}
	}

	return nil
}

// emitRaw writes a filler with the same bytes: .half for 2-byte RVC words,
// .word for full instructions, .byte for other scraps.
func emitRaw(sb *strings.Builder, raw []byte) {
	switch len(raw) {
	case 2:
		fmt.Fprintf(sb, ".half 0x%04x\n", binary.LittleEndian.Uint16(raw))
	case 4:
		fmt.Fprintf(sb, ".word 0x%08x\n", binary.LittleEndian.Uint32(raw))
	default:
		for _, b := range raw {
			fmt.Fprintf(sb, ".byte 0x%02x\n", b)
		}
	}
}

// cleanLine normalizes the instruction text and strips objdump annotations.
func cleanLine(s string) string {
	return objdump.StripComments(objdump.Normalize(s))
}

// firstDiff - the offset of the first divergence between two buffers of
// equal length.
func firstDiff(a, b []byte) int {
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}

	return len(a)
}
