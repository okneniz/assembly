// hello world written directly in Go - no .s source, no expression
// evaluator: the program is a chain of arch builder calls, everything is
// computed by ordinary Go. The emitted code is byte-identical to the
// assembly pipeline over hello-macos.s.
//
// Build & run:  go run . -o hello-go.bin && make -C .. run-macos
//
// The program talks to the kernel directly (svc), no libc: the Darwin
// arm64 convention keeps the syscall number in x16 and packs the class
// prefix into the high bits - the full number of write is
// sysClassUnix|sysWrite, assembled from two halves (movz + movk), because
// it does not fit into a single imm16.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	arch "github.com/okneniz/assembly/arch/arm64"
	prog "github.com/okneniz/assembly/prog/arm64"
)

// The program dictionary: syscall numbers, trap immediates, fds - the
// Go constants the chain below is spelled with.
const (
	sysClassUnix = 0x2000000 // the SYSCALL_CLASS_UNIX prefix (2 << 24)
	sysWrite     = 4
	sysExit      = 1
	trapMach     = 0x80 // the canonical Darwin trap immediate
	fdStdout     = 1
)

const msg = "hello world\n"

func main() {
	out := flag.String("o", "hello.bin", "the executable to write")
	flag.Parse()

	// The program: every chain method is a source line, labels resolve
	// at assembly time, macros would be ordinary Go functions returning
	// *prog.Program.
	p := prog.New().
		Label("start").
		Mov(prog.X0, fdStdout).                     // write(fd=stdout, ...)
		Adr(prog.X1, "msg").                        // ... buf = the string
		Mov(prog.X2, int64(len(msg))).              // ... len
		Movz(prog.X16, sysClassUnix>>16, arch.Hw1). // x16 = sysClassUnix | ...
		Movk(prog.X16, sysWrite, arch.Hw0).         // ... sysWrite
		Svc(trapMach).
		Mov(prog.X0, 0). // exit(return code = 0)
		Movz(prog.X16, sysClassUnix>>16, arch.Hw1).
		Movk(prog.X16, sysExit, arch.Hw0).
		Svc(trapMach).
		Label("msg").
		Ascii(msg).
		Entry("start")

	bin, errs := p.Build()
	if len(errs) > 0 {
		fail(errs)
	}

	// The raw program image at 0x1000 (the runner maps it there and
	// jumps to the entry).
	code, syms, errs := bin.Assemble(0x1000)
	if len(errs) > 0 {
		fail(errs)
	}

	if _, ok := syms["start"]; !ok {
		fail([]error{errors.New("no start symbol")})
	}

	if err := os.WriteFile(*out, code, 0o755); err != nil {
		fail([]error{err})
	}

	fmt.Printf("wrote %s (%d bytes, entry %#x)\n", *out, len(code), syms["start"])
}

func fail(errs []error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	os.Exit(1)
}
