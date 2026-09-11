// bare-metal RISC-V hello written directly in Go - no .s source: one
// string through the ns16550 UART of the QEMU virt machine, then the
// sifive_test poweroff write (the machine stops itself - no idle loop,
// unlike the LoongArch twin). The emitted code is the uncompressed
// (NoRVC) encoding of the assembly pipeline over hello-riscv.s.
//
// Build:  go run . -o hello-go-riscv.elf
// Run:    make -C .. vm-riscv
// Debug:  go run . -debug   (serves a DAP session on stdio - the
//
//	editor adapter relays into it; the line map points at this
//	file, breakpoints land on the chain calls below)
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/okneniz/assembly/debug/dap"
	"github.com/okneniz/assembly/debug/riscv"
	"github.com/okneniz/assembly/debug/session"
	"github.com/okneniz/assembly/file"
	core "github.com/okneniz/assembly/prog"
	prog "github.com/okneniz/assembly/prog/riscv"
)

// The machine dictionary: the ns16550 UART data register of the virt
// machine, the sifive_test poweroff device, and the reset base.
const (
	uartData   = 0x10000000
	finisher   = 0x100000
	finishPass = 0x5555
	base       = 0x80000000
)

const msg = "hello world\n"

// callerPos - the position resolver of the chain: invoked inside a
// chain method, two frames up is the calling code (this file).
func callerPos() core.Pos {
	_, file, line, _ := runtime.Caller(2)
	return core.NewPos(file, line)
}

func main() {
	out := flag.String("o", "hello-go-riscv.elf", "the ELF image to write")
	serve := flag.Bool(
		"debug",
		false,
		"serve a DAP debug session on stdio instead of writing the image",
	)
	flag.Parse()

	// a0 = the UART data register; a1 = the message cursor; a2 = the
	// end; the loop stores every byte, then the finisher write powers
	// the machine off.
	p := prog.New().
		WithPos(callerPos).
		Label("start").
		Lui(prog.A0, uartData>>12).
		La(prog.A1, "msg").
		La(prog.A2, "end").
		Label("loop").
		Bgeu(prog.A1, prog.A2, "done").
		Lb(prog.A5, prog.A1, 0).
		Sb(prog.A5, prog.A0, 0).
		Addi(prog.A1, prog.A1, 1).
		J("loop").
		Label("done").
		Lui(prog.A5, finisher>>12).
		Lui(prog.A6, finishPass>>12).
		Addi(prog.A6, prog.A6, finishPass&0xfff).
		Sw(prog.A6, prog.A5, 0).
		// the safety idle loop: reached only when the finisher write
		// does not power the machine off.
		Label("hang").
		J("hang").
		Label("msg").
		Ascii(msg).
		Label("end").
		Entry("start")

	bin, errs := p.Build()
	if len(errs) > 0 {
		fail(errs)
	}

	res := bin.Assemble(base)
	if len(res.Errs) > 0 {
		fail(res.Errs)
	}

	entry := res.Syms["start"]
	blob, err := file.WriteELF(
		file.EM_RISCV,
		0,
		base,
		entry,
		[]file.Section{*file.NewSection(".text", "", base, 0, uint64(len(res.Code)), res.Code)},
	)
	if err != nil {
		fail([]error{err})
	}

	if *serve {
		debugServe(blob, res)
		return
	}

	if err := os.WriteFile(*out, blob, 0o755); err != nil {
		fail([]error{err})
	}

	fmt.Printf("wrote %s (%d bytes, entry %#x)\n", *out, len(res.Code), entry)
}

// debugServe is the editor session: the assembled image, the symbols,
// and the line map (the Go lines of the chain calls above) boot under
// qemu at the dap launch request and speak DAP on stdio.
func debugServe(blob []byte, res *core.Result) {
	lines := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		lines = append(lines, session.NewLine(e.Pos.File, e.Pos.Line, e.Addr, e.Size))
	}

	launcher := dap.NewImageLauncher(riscv.NewTarget(), blob, res.Syms, lines)
	srv, err := dap.NewServer(os.Stdin, os.Stdout, launcher)
	if err != nil {
		fail([]error{err})
	}

	if err := srv.Serve(); err != nil {
		fail([]error{err})
	}
}

func fail(errs []error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	os.Exit(1)
}
