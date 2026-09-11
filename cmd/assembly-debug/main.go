// Command assembly-debug runs a program under the debugger: one static
// binary over a qemu executor, with the project's own decoders,
// symbols, and source lines. The engine is the debug/session library -
// this REPL is one consumer of it (the DAP adapter is the other); a
// debugging script is an ordinary Go program importing it.
//
// The input is either a source file (assembled in-process; the line map
// and the symbols come along) or a binary with an optional symbol
// sidecar (-sym, the assembly CLI's -sym output). The REPL commands
// are listed by its help (b, c, s, regs, x, dis, list, syms, help, q).
//
// Usage:
//
//	assembly-debug -arch arm64 hello.s
//	assembly-debug -arch arm64 -bin hello.elf -sym hello.sym
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/okneniz/assembly/debug/load"
	"github.com/okneniz/assembly/debug/qemu"
	"github.com/okneniz/assembly/debug/rsp"
	"github.com/okneniz/assembly/debug/session"
)

func main() {
	archFlag := flag.String("arch", "arm64", "architecture: arm64|riscv64|loong64")
	baseFlag := flag.String("base", "", "base address (hex or dec; the arch default otherwise)")
	binPath := flag.String("bin", "", "debug a binary instead of assembling a source")
	symPath := flag.String("sym", "", "symbol sidecar of -bin (the assembly CLI's -sym output)")
	icount := flag.Bool("icount", false, "deterministic virtual clock (-icount shift=auto)")
	execArg := flag.String(
		"e",
		"",
		"batch mode: semicolon-separated commands (\"b start; c; regs\") instead of the REPL",
	)
	flag.Usage = func() {
		fmt.Fprintln(
			os.Stderr,
			"usage: assembly-debug -arch ARCH [SOURCE] [-bin FILE] [-sym FILE] [-base ADDR] [-e CMDS]",
		)
		flag.PrintDefaults()
	}
	flag.Parse()
	if err := run(
		*archFlag,
		*baseFlag,
		*binPath,
		*symPath,
		*icount,
		*execArg,
		flag.Arg(0),
	); err != nil {
		fmt.Fprintln(os.Stderr, "assembly-debug:", err)
		os.Exit(1)
	}
}

func run(
	archName, baseFlag, binPath, symPath string,
	icount bool,
	execArg, srcPath string,
) (err error) {
	setup, err := load.SetupFor(archName)
	if err != nil {
		return err
	}

	res, err := load.Build(setup, load.Input{
		SrcPath: srcPath,
		BinPath: binPath,
		SymPath: symPath,
		Base:    baseFlag,
	})
	if err != nil {
		return err
	}

	tgt := setup.Tgt
	opts := qemu.NewOptions()
	opts.Serial = os.Stdout // the program's console: the debugged code's own output
	opts.Icount = icount
	machine, err := qemu.Start(context.Background(), tgt, res.Img, opts)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := machine.Close(); err == nil {
			err = cerr
		}
	}()

	s, err := session.New(machine.Conn(), tgt, res.Syms, res.Lines)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "assembly-debug: %s %s\n", tgt.QemuBinary(), res.SrcName)
	if err := report(s, res.SrcLines); err != nil {
		return err
	}

	if execArg != "" {
		return batch(s, execArg, res.SrcLines)
	}

	return repl(s, res.SrcLines)
}

// batch runs semicolon-separated commands and leaves (the gate form:
// no prompt, one line of output per command).
func batch(s *session.Session, cmds string, srcLines []string) error {
	for cmd := range strings.SplitSeq(cmds, ";") {
		if _, err := dispatch(s, strings.TrimSpace(cmd), srcLines); err != nil {
			return err
		}
	}

	return nil
}

// repl reads commands until quit or EOF: b, c, s, regs, x, dis, list,
// syms, help, q. The program's own output goes to stdout (the executor's
// console); the debugging UI stays on stderr.
func repl(s *session.Session, srcLines []string) error {
	fmt.Fprintln(os.Stderr, "type help for the commands")

	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "(asmdb) ")
		if !sc.Scan() {
			fmt.Fprintln(os.Stderr)
			return nil
		}

		quit, err := dispatch(s, sc.Text(), srcLines)
		if quit {
			return nil
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "assembly-debug:", err)
		}
	}
}

// dispatch runs one command line (the shared vocabulary of the REPL
// and the -e batch mode); quit reports a requested end.
func dispatch(s *session.Session, line string, srcLines []string) (bool, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, nil
	}

	switch fields[0] {
	case "b":
		if len(fields) < 2 {
			return false, errors.New("usage: b <symbol|0xaddr>")
		}

		addr, err := s.BreakAt(fields[1])
		if err != nil {
			return false, err
		}

		fmt.Fprintf(os.Stderr, "breakpoint at %#x\n", addr)
		return false, nil
	case "c":
		return false, continueAndReport(s, srcLines)
	case "s":
		stop, err := s.Step()
		if err != nil {
			return false, err
		}

		return false, stopContext(s, stop, srcLines)
	case "regs":
		return false, dumpRegs(s)
	case "x":
		return false, dumpMem(s, fields[1:])
	case "dis":
		return false, dumpDisasm(s, fields[1:])
	case "list":
		return false, dumpSource(s, srcLines)
	case "syms":
		dumpSymbols(s)
		return false, nil
	case "help", "h", "?":
		usage()
		return false, nil
	case "q":
		return true, nil
	}

	return false, fmt.Errorf("unknown command %q (help lists the commands)", fields[0])
}

// usage is the REPL's help: every command with what it takes and what
// it does, plus a first-session example.
func usage() {
	fmt.Fprint(os.Stderr, `commands:
  b <label|0xaddr>   set a breakpoint at a source label or an address
  c                  continue: run to the next breakpoint or the program end
  s                  step one instruction
  regs               show the core registers
  x <addr> [len]     read memory: a hexdump of len bytes at addr (default 16)
  dis [n]            disassemble n instructions at pc (default 4)
  list               the source window around pc
  syms               the symbol table (labels and their addresses)
  help               this help
  q                  quit

a session starts frozen at the entry; the pc, the source line, and the
instruction print at every stop. example:

  (asmdb) b done          break at the label done
  (asmdb) c               run there
  pc 0x40100024  hello.s:21  movz x0, #0x8400, lsl #16
`)
}

// continueAndReport continues and prints the stop context.
func continueAndReport(s *session.Session, srcLines []string) error {
	stop, err := s.Continue()
	if err != nil {
		// a powered-off machine closes the stub: the run is over, not
		// broken
		return fmt.Errorf("continue: %w (machine off?)", err)
	}

	return stopContext(s, stop, srcLines)
}

// stopContext prints the stop: a finished program (it powered itself
// off - PSCI SYSTEM_OFF and friends) reports its exit, and there is no
// pc to ask from a dead machine.
func stopContext(s *session.Session, stop rsp.StopReply, srcLines []string) error {
	if stop.Exited() {
		fmt.Fprintf(os.Stderr, "the program finished: %s\n", stop)
		return nil
	}

	return report(s, srcLines)
}

// report prints the stop context: pc, the source line, the instruction.
func report(s *session.Session, srcLines []string) error {
	pc, err := s.PC()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "pc %#x", pc)
	if line, ok := s.LineAt(pc); ok {
		fmt.Fprintf(os.Stderr, "  %s:%d", line.File, line.Line)
		if i := line.Line - 1; i >= 0 && i < len(srcLines) && line.File != "" {
			fmt.Fprintf(os.Stderr, "  %s", strings.TrimSpace(srcLines[i]))
		}
	}

	fmt.Fprintln(os.Stderr)

	dis, err := s.Disasm(pc, 1)
	if err == nil && len(dis) > 0 {
		fmt.Fprintln(os.Stderr, dis[0])
	}

	return nil
}

func dumpRegs(s *session.Session) error {
	regs, err := s.Regs()
	if err != nil {
		return err
	}

	for _, r := range regs {
		fmt.Fprintf(os.Stderr, "%-4s %#018x\n", r.Name, r.Value)
	}

	return nil
}

func dumpMem(s *session.Session, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: x <addr> [len]")
	}

	addr, err := strconv.ParseUint(args[0], 0, 64)
	if err != nil {
		return fmt.Errorf("bad address %q", args[0])
	}

	n := 16
	if len(args) > 1 {
		if n, err = strconv.Atoi(args[1]); err != nil || n <= 0 || n > 4096 {
			return fmt.Errorf("bad length %q", args[1])
		}
	}

	mem, err := s.Read(addr, n)
	if err != nil {
		return err
	}

	for i := 0; i < len(mem); i += 16 {
		end := min(i+16, len(mem))
		fmt.Fprintf(os.Stderr, "%08x: % x\n", addr+uint64(i), mem[i:end])
	}

	return nil
}

func dumpDisasm(s *session.Session, args []string) error {
	n := 4
	if len(args) > 0 {
		var err error
		if n, err = strconv.Atoi(args[0]); err != nil || n <= 0 || n > 64 {
			return fmt.Errorf("bad count %q", args[0])
		}
	}

	pc, err := s.PC()
	if err != nil {
		return err
	}

	dis, err := s.Disasm(pc, n)
	if err != nil {
		return err
	}

	for _, l := range dis {
		fmt.Fprintln(os.Stderr, l)
	}

	return nil
}

func dumpSource(s *session.Session, srcLines []string) error {
	pc, err := s.PC()
	if err != nil {
		return err
	}

	line, ok := s.LineAt(pc)
	if !ok {
		return errors.New("no line map: run from a source (or use dis)")
	}

	lo := max(line.Line-3, 0)
	hi := min(line.Line+2, len(srcLines))
	for i := lo; i < hi; i++ {
		mark := "  "
		if i == line.Line-1 {
			mark = "->"
		}

		fmt.Fprintf(os.Stderr, "%s %s:%d  %s\n", mark, line.File, i+1, srcLines[i])
	}

	return nil
}

func dumpSymbols(s *session.Session) {
	for _, name := range s.Symbols() {
		addr, _ := s.Symbol(name)
		fmt.Fprintf(os.Stderr, "%#x %s\n", addr, name)
	}
}
