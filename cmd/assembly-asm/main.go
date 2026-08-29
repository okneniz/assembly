// Command assembly-asm assembles a source file (a GNU-as compatible subset
// + assembly's own output — round-trip) into machine code.
//
// The architecture is set with the -arch flag (arm64|riscv64); the section
// base address with -base (hex or dec). Output: a binary file (-o), a hex
// dump to stdout (--hex) and a symbol table (--sym). Errors are printed with
// line positions.
//
// Usage:
//
//	assembly-asm -arch riscv64 -base 0x1000 source.s -o prog.bin --sym symbols.txt
//	assembly-asm -arch arm64 --hex source.s
//	assembly-asm -arch riscv64 -base 0x1000 < roundtrip.s --hex
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/text"
)

func main() {
	archFlag := flag.String("arch", "", "architecture: arm64|riscv64|loong64 (required)")
	baseFlag := flag.String("base", "0", "base address of the .text section (hex or dec)")
	outPath := flag.String("o", "", "output binary file (default: stdout ignored with --hex)")
	symPath := flag.String("sym", "", "write symbol table to this file")
	hexDump := flag.Bool("hex", false, "print hex dump to stdout")
	disasmFlag := flag.Bool(
		"disasm",
		false,
		"disassemble a raw binary (with -arch/-base) instead of assembling",
	)
	formatFlag := flag.String(
		"format",
		"raw",
		"output format: raw (default) or elf — wrap sections into a minimal static ELF64 executable (runs natively on Linux arm64/riscv64 or under qemu-user)",
	)
	flag.Usage = func() {
		fmt.Fprintln(
			os.Stderr,
			"usage: assembly-asm -arch ARCH [-base ADDR] [-o FILE] [--sym FILE] [--hex] [SOURCE]",
		)
		fmt.Fprintln(os.Stderr, "  SOURCE defaults to stdin")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *archFlag == "" {
		flag.Usage()
		os.Exit(2)
	}

	var assemble func(string, uint64) (*asm.Result, []asm.AsmError)
	switch strings.ToLower(*archFlag) {
	case "arm64", "aarch64":
		assemble = alias.Assemble
	case "riscv64", "riscv":
		assemble = pseudo.Assemble
	case "loong64", "loongarch64":
		assemble = lpseudo.Assemble
	default:
		fmt.Fprintf(os.Stderr, "unknown -arch %q (want arm64, riscv64 or loong64)\n", *archFlag)
		os.Exit(2)
	}

	base, err := parseAddr(*baseFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad -base %q: %v\n", *baseFlag, err)
		os.Exit(2)
	}

	// The ELF is loaded at p_vaddr: a zero base is unusable (mmap_min_addr),
	// so for the elf format the default is the standard base 0x10000.
	if *formatFlag == "elf" && base == 0 {
		base = 0x10000
	}

	src, err := readSource(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "read:", err)
		os.Exit(1)
	}

	if *disasmFlag {
		// our own ELF format: the code starts at file offset 0x1000
		if len(src) > 4 && src[0] == 0x7f && string(src[1:4]) == "ELF" {
			src = src[0x1000:]
		}

		if err := runDisasm(src, strings.ToLower(*archFlag), base); err != nil {
			fmt.Fprintln(os.Stderr, "disasm:", err)
			os.Exit(1)
		}

		return
	}

	res, errs := assemble(string(src), base)
	if len(errs) > 0 {
		exit := 0
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", sourceName(flag.Arg(0)), e.Line, e.Col, e.Msg)
			exit = 1
		}

		if exit == 1 && len(res.Sections) == 0 {
			os.Exit(1)
		}
	}

	if *outPath != "" && len(errs) == 0 {
		var blob []byte
		if *formatFlag == "elf" {
			machine, flags, merr := elfMachine(strings.ToLower(*archFlag))
			if merr != nil {
				fmt.Fprintln(os.Stderr, merr)
				os.Exit(2)
			}

			entry := base
			for _, name := range []string{"start", "_start"} {
				if a, ok := res.Symbols[name]; ok {
					entry = a
					break
				}
			}

			elfBlob, werr := file.WriteELF(machine, flags, base, entry, fileSections(res.Sections))
			if werr != nil {
				fmt.Fprintln(os.Stderr, "elf:", werr)
				os.Exit(1)
			}

			blob = elfBlob
		} else {
			for _, sec := range res.Sections {
				blob = append(blob, sec.Data...)
			}
		}

		if err := os.WriteFile(*outPath, blob, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "write:", err)
			os.Exit(1)
		}
	}

	if *symPath != "" && len(errs) == 0 {
		if err := writeSymbols(*symPath, res); err != nil {
			fmt.Fprintln(os.Stderr, "write symbols:", err)
			os.Exit(1)
		}
	}

	if *hexDump || *outPath == "" {
		for _, sec := range res.Sections {
			fmt.Printf(".section %s @ %#x (%d bytes)\n", sec.Name, sec.Addr, len(sec.Data))
			for i := 0; i < len(sec.Data); i += 16 {
				end := min(i+16, len(sec.Data))
				fmt.Printf("%08x: %s\n", sec.Addr+uint64(i), hexLine(sec.Data[i:end]))
			}
		}
	}

	if len(errs) > 0 {
		os.Exit(1)
	}
}

// runDisasm prints the disassembly of a raw binary (the reverse direction of
// assembly-asm): disasm assembles the lines from the ObjDump text of the
// instructions; arm64 — byte-column style, riscv — hex word.
func runDisasm(data []byte, arch string, base uint64) error {
	switch arch {
	case "arm64", "aarch64":
		instrs, err := arm64.Parse(base)(bytes.Buffer(data))
		if err != nil {
			return err
		}

		return disasm.Write(
			os.Stdout,
			base,
			data,
			instrs,
			disasm.NewOptions(text.CodeBytes),
		)
	case "riscv64", "riscv":
		instrs, err := riscv.Parse(base)(bytes.Buffer(data))
		if err != nil {
			return err
		}

		return disasm.Write(
			os.Stdout,
			base,
			data,
			instrs,
			disasm.NewOptions(text.CodeWord),
		)
	case "loong64", "loongarch64":
		instrs, err := loong64.Parse(base)(bytes.Buffer(data))
		if err != nil {
			return err
		}

		return disasm.Write(
			os.Stdout,
			base,
			data,
			instrs,
			disasm.NewOptions(text.CodeWord),
		)
	}

	return nil
}

// elfMachine — e_machine and e_flags from the CLI's string architecture
// name. LoongArch carries the ABI bits (base + double-float): the kernel
// rejects an ELF without a valid float-ABI setting.
func elfMachine(arch string) (uint16, uint32, error) {
	switch arch {
	case "arm64", "aarch64":
		return file.EM_AARCH64, 0, nil
	case "riscv64", "riscv":
		return file.EM_RISCV, 0, nil
	case "loong64", "loongarch64":
		return file.EM_LOONGARCH, 0x43, nil
	}

	return 0, 0, fmt.Errorf("no ELF machine for %q", arch)
}

// fileSections converts the sections of the assembly result into sections
// of the file package (the WriteELF emitter needs data and size). A NOBITS
// (.bss) in the middle of the layout is materialized with zeros — file
// continuity of addresses; only a section after the last PROGBITS remains
// as a memsz tail.
func fileSections(secs []asm.Section) []file.Section {
	lastProgbits := -1
	for i, s := range secs {
		if !s.Nobits {
			lastProgbits = i
		}
	}

	out := make([]file.Section, len(secs))
	for i, s := range secs {
		data := s.Data
		if s.Nobits && i < lastProgbits {
			data = make([]byte, s.Size)
		}

		out[i] = *file.NewSection(s.Name, "", s.Addr, 0, uint64(s.Size), data)
	}

	return out
}

func parseAddr(s string) (uint64, error) {
	s = strings.TrimPrefix(strings.TrimPrefix(strings.ToLower(s), "0x"), "0x")
	if s == "" {
		return 0, nil
	}

	return strconv.ParseUint(s, 16, 64)
}

func readSource(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}

	return os.ReadFile(path)
}

func sourceName(path string) string {
	if path == "" || path == "-" {
		return "<stdin>"
	}

	return path
}

func writeSymbols(path string, res *asm.Result) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()
	for name, addr := range res.Symbols {
		if _, err := fmt.Fprintf(f, "%#x %s\n", addr, name); err != nil {
			return err
		}
	}

	return nil
}

func hexLine(b []byte) string {
	var sb strings.Builder
	for i, x := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}

		fmt.Fprintf(&sb, "%02x", x)
	}

	return sb.String()
}
