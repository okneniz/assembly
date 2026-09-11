// Package load is the shared input path of the debugging frontends:
// a .s source (assembled in-process - the symbols, the line map, and
// the source text come along) or a binary with its optional -sym
// sidecar become the executor image plus the debug views over it.
// The arch table (the target, the assembler, the ELF machine/flags,
// the default base) is Setup.
package load

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/debug/session"
	"github.com/okneniz/assembly/file"
)

// Build turns the input into the executor image and the debug views.
// An assembly failure carries every error as a "file:line:col: msg"
// line of the returned error. The image is an ELF, or the raw section
// concatenation for the arches whose loader takes a blob (Raw).
func Build(setup Setup, in Input) (Result, error) {
	base := setup.Base
	if in.Base != "" {
		b, err := strconv.ParseUint(in.Base, 0, 64)
		if err != nil {
			return Result{}, fmt.Errorf("assembly/load: bad base %q: %w", in.Base, err)
		}

		base = b
	}

	if in.BinPath != "" {
		return buildBinary(in.BinPath, in.SymPath)
	}

	return buildSource(setup, in.SrcPath, base)
}

// buildBinary loads a prebuilt image: the bytes run as-is, the symbols
// come from the sidecar (no line map, no source text).
func buildBinary(binPath, symPath string) (Result, error) {
	img, err := os.ReadFile(binPath)
	if err != nil {
		return Result{}, err
	}

	syms := map[string]uint64{}
	if symPath != "" {
		syms, err = readSymbols(symPath)
		if err != nil {
			return Result{}, err
		}
	}

	return Result{Img: img, Syms: syms, SrcName: binPath}, nil
}

// buildSource assembles the source in-process: the entry is the
// start/_start label (the base otherwise), the line map points back
// into the source file.
func buildSource(setup Setup, srcPath string, base uint64) (Result, error) {
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return Result{}, err
	}

	res, errs := setup.Assemble(string(src), base)
	if len(errs) > 0 {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, fmt.Sprintf("%s:%d:%d: %s", srcPath, e.Line, e.Col, e.Msg))
		}

		return Result{}, errors.New(strings.Join(msgs, "\n"))
	}

	entry := base
	for _, name := range []string{"start", "_start"} {
		if a, ok := res.Symbols[name]; ok {
			entry = a
			break
		}
	}

	img := []byte{}
	if setup.Raw {
		for _, sec := range res.Sections {
			img = append(img, sec.Data...)
		}
	} else {
		img, err = file.WriteELF(
			setup.Machine,
			setup.Flags,
			base,
			entry,
			fileSections(res.Sections),
		)
		if err != nil {
			return Result{}, err
		}
	}

	out := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		out = append(out, session.NewLine(srcPath, int(e.Line), e.Addr, e.Size))
	}

	return Result{
		Img:      img,
		Syms:     res.Symbols,
		Lines:    out,
		SrcLines: strings.Split(string(src), "\n"),
		SrcName:  srcPath,
	}, nil
}

// readSymbols parses a -sym sidecar: one "0x<addr> <name>" per line.
func readSymbols(path string) (map[string]uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	syms := map[string]uint64{}
	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		addr, err := strconv.ParseUint(fields[0], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("assembly/load: sym file %q: bad line %q", path, line)
		}

		syms[fields[1]] = addr
	}

	return syms, nil
}

// fileSections converts the asm sections into file-package sections
// for the ELF emitter (the same view the assembly CLI builds).
func fileSections(secs []asm.Section) []file.Section {
	out := make([]file.Section, len(secs))
	for i, s := range secs {
		data := s.Data
		if s.Nobits {
			data = make([]byte, s.Size)
		}

		out[i] = *file.NewSection(s.Name, "", s.Addr, 0, uint64(s.Size), data)
	}

	return out
}
