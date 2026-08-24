// Command assembly-diff compares the output of the local disassembler with
// objdump over the code section of a given binary. A tool for iteratively
// improving instruction coverage: shows the match percentage and the first
// divergences.
//
// The architecture is determined automatically from the binary (cputype /
// e_machine); the -arch flag overrides it.
//
// Usage:
//
//	assembly-diff <binary>
//	assembly-diff <binary> -n 30        # show up to 30 divergences
//	assembly-diff <binary> -arch arm64  # override the architecture
//
// Requires objdump on PATH (Apple's llvm-based /usr/bin/objdump supports
// --macho). Its output is parsed by the objdump package (shared with the
// differential test).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
	"github.com/okneniz/assembly/text"
)

func main() {
	n := flag.Int("n", 15, "max mismatches to print")
	archFlag := flag.String("arch", "", "override architecture (arm64|riscv64)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: assembly-diff <binary> [-n N] [-arch ARCH]")
		os.Exit(2)
	}

	path := flag.Arg(0)

	ff, err := file.Detect(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "detect:", err)
		os.Exit(1)
	}

	sec, err := ff.CodeSection()
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(1)
	}

	kind := ff.ArchKind()
	if *archFlag != "" {
		k, ok := parseArchFlag(*archFlag)
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown -arch %q (want arm64 or riscv64)\n", *archFlag)
			os.Exit(2)
		}

		kind = k
	}

	// ours, keyed by address. Architecture dispatches directly to the arch
	// package's Parse (no shared decoder/instruction type); lines are composed
	// by disasm from the instruction's own ObjDump text.
	style := text.StyleFor(ff.Name())
	opts := disasm.NewOptions(style)
	ours := map[uint64]string{}
	switch kind {
	case file.ArchARM64:
		insts, err := arm64.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		if err != nil {
			fmt.Fprintln(os.Stderr, "decode:", err)
			os.Exit(1)
		}

		for _, in := range insts {
			a := in.Addr()
			ours[a] = objdump.Normalize(disasm.Line(a, sec.Data[a-sec.Addr:], in, opts))
		}
	case file.ArchRISCV64:
		insts, err := riscv.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		if err != nil {
			fmt.Fprintln(os.Stderr, "decode:", err)
			os.Exit(1)
		}

		for _, in := range insts {
			ours[in.Addr()] = objdump.Normalize(
				disasm.Line(in.Addr(), sec.Data[in.Addr()-sec.Addr:], in, opts),
			)
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported architecture (kind %d)\n", kind)
		os.Exit(1)
	}

	out, err := objdump.Run(context.Background(), objdump.Args(ff.Name(), kind, path))
	if err != nil {
		fmt.Fprintln(os.Stderr, "objdump:", err)
		os.Exit(1)
	}

	obj := objdump.ParseByAddr(string(out))
	matched, mismatched := 0, 0
	var mism []string
	for addr, objLine := range obj {
		ol, ok := ours[addr]
		if !ok {
			continue
		}

		if objdump.StripComments(ol) == objdump.StripComments(objLine) {
			matched++
		} else {
			mismatched++
			if len(mism) < *n {
				mism = append(
					mism,
					fmt.Sprintf("  %#x\n    our: %s\n    obj: %s", addr, ol, objLine),
				)
			}
		}
	}

	pct := 0.0
	if total := len(obj); total > 0 {
		pct = float64(matched) * 100 / float64(total)
	}

	fmt.Printf("total=%d  matched=%d  mismatched=%d\n", len(obj), matched, mismatched)
	fmt.Printf("match: %.2f%%\n", pct)
	fmt.Println("first mismatches:")
	if len(mism) == 0 {
		fmt.Println("  (none)")
	}

	for _, m := range mism {
		fmt.Println(m)
	}
}

// parseArchFlag parses the string value of the -arch flag into an ArchKind.
func parseArchFlag(s string) (file.ArchKind, bool) {
	switch s {
	case "arm64", "aarch64":
		return file.ArchARM64, true
	case "riscv64", "riscv":
		return file.ArchRISCV64, true
	}

	return file.ArchUnknown, false
}
