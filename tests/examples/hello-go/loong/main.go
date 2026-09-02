// bare-metal LoongArch64 hello written directly in Go - no .s source:
// one string through the 16550A UART of the QEMU virt machine, then an
// idle loop (the gate stops the machine once the line is out). The
// emitted code is byte-identical to the assembly pipeline over
// hello-loongarch.s.
//
// Build:  go run . -o hello-go-loong.bin
// Run:    make -C .. run-loong
package main

import (
	"flag"
	"fmt"
	"os"

	prog "github.com/okneniz/assembly/prog/loong64"
)

// The machine dictionary: the UART data register of the virt machine
// (16550A, byte-wide) and the flash base LoongArch resets into.
const (
	uartData = 0x1fe001e0
	base     = 0x1c000000
)

const msg = "hello world\n"

func main() {
	out := flag.String("o", "hello-go-loong.bin", "the raw image to write")
	flag.Parse()

	// $t0 = the UART data register; $t1 = the message cursor; the loop
	// stores every byte until the terminating zero.
	p := prog.New().
		Label("start").
		Lu12iW(prog.T0, uartData>>12).
		Ori(prog.T0, prog.T0, uartData&0xfff).
		La(prog.T1, "msg").
		Label("loop").
		LdBu(prog.T2, prog.T1, 0).
		Beq(prog.T2, prog.Zero, "idle").
		StB(prog.T2, prog.T0, 0).
		AddiW(prog.T1, prog.T1, 1).
		B("loop").
		Label("idle").
		B("idle").
		Label("msg").
		Ascii(msg).
		Bytes(0).
		Entry("start")

	bin, errs := p.Build()
	if len(errs) > 0 {
		fail(errs)
	}

	code, syms, errs := bin.Assemble(base)
	if len(errs) > 0 {
		fail(errs)
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
