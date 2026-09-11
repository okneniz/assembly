package riscv

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/prog"
)

// callerPos - a caller-based resolver for the tests: invoked inside a
// chain method, two frames up is the code calling the chain (this file).
func callerPos() prog.Pos {
	_, file, line, _ := runtime.Caller(2)
	return prog.NewPos(file, line)
}

// The Go counterpart of tests/examples/hello-asm/hello-riscv.s: one
// string through the ns16550 UART of the virt machine, then the
// sifive_test poweroff write.
func helloRiscv() *Program {
	return New().
		WithPos(callerPos).
		Label("start").
		// a0 = the UART data register (ns16550, byte-wide at 0x10000000).
		Lui(A0, 0x10000).
		// a1 = the message cursor; a2 = the end.
		La(A1, "msg").
		La(A2, "end").
		Label("loop").
		Bgeu(A1, A2, "done"). // while cursor < end
		Lb(A5, A1, 0).        // the string byte
		Sb(A5, A0, 0).        // → UART
		Addi(A1, A1, 1).
		J("loop").
		// done: a5 = the sifive_test device; 0x5555 (FINISHER_PASS)
		// powers the machine off.
		Label("done").
		Lui(A5, 0x100).
		Lui(A6, 0x5).
		Addi(A6, A6, 0x555).
		Sw(A6, A5, 0).
		Label("hang").
		J("hang").
		Label("msg").
		Ascii("hello world\n").
		Label("end").
		Entry("start")
}

func TestAssembleRiscvGolden(t *testing.T) {
	bin, buildErrs := helloRiscv().Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0x80000000)
	require.Empty(t, res.Errs)

	golden, err := os.ReadFile("testdata/hello-riscv.bin")
	if err != nil {
		t.Skip("golden binary not built: run the assembly CLI")
	}

	require.Equal(t, golden, res.Code)
	require.Equal(t, uint64(0x80000000), res.Syms["start"])

	// the injected caller resolver reports this test file (labels emit
	// nothing; each la pair is one entry)
	require.Len(t, res.Lines, 14)
	for _, e := range res.Lines {
		require.Contains(t, e.Pos.File, "prog/riscv/program_test.go")
		require.Positive(t, e.Pos.Line)
	}
}

func TestAssembleRiscvErrors(t *testing.T) {
	// undefined branch target
	bin, _ := New().J("nowhere").Build()
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")

	// undefined la target
	laBin, _ := New().La(T0, "gone").Build()
	laErrs := laBin.Assemble(0).Errs
	require.Len(t, laErrs, 1)
	require.Contains(t, laErrs[0].Error(), "gone")

	// deferred immediate error surfaces at Build
	_, buildErrs := New().Lui(T0, 1<<20).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "lui")
}

func TestLineMap(t *testing.T) {
	// the injected resolver: the la pair is ONE entry of 8 bytes (one
	// chain call), data lines carry their byte size
	p := New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 3) }).
		Label("start").
		Lui(T0, 0x10000).
		La(T1, "msg").
		Label("msg").
		Ascii("hi").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)
	res := bin.Assemble(0)
	require.Empty(t, res.Errs)
	require.Equal(t, []prog.LineEntry{
		prog.NewLineEntry(0, 4, prog.NewPos("synthetic.go", 3)),
		prog.NewLineEntry(4, 8, prog.NewPos("synthetic.go", 3)),
		prog.NewLineEntry(12, 2, prog.NewPos("synthetic.go", 3)),
	}, res.Lines)
}

func TestTwins(t *testing.T) {
	// jal/bnez/ecall: three fixed 4-byte lines, no compression (ecall
	// is the word 0x00000073)
	bin, buildErrs := New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 1) }).
		Jal(Ra, "f").
		Bnez(T0, "f").
		Ecall().
		Label("f").
		Entry("f").
		Build()
	require.Empty(t, buildErrs)

	assembled := bin.Assemble(0)
	require.Empty(t, assembled.Errs)
	require.Len(t, assembled.Code, 12)
	require.Equal(t, []byte{0x73, 0x00, 0x00, 0x00}, assembled.Code[8:])
}
