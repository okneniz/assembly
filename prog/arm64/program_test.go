package arm64

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// callerPos - a caller-based resolver for the tests: invoked inside a
// chain method, two frames up is the code calling the chain (this file).
func callerPos() prog.Pos {
	_, file, line, _ := runtime.Caller(2)
	return prog.NewPos(file, line)
}

// hello-macos written on the chain: the Go counterpart of
// tests/examples/hello-asm/hello-macos.s, byte-identical when assembled.
func helloMacOS() *Program {
	return New().
		WithPos(callerPos).
		Label("start").
		Mov(X0, fdStdout).                     // write(fd=stdout, ...)
		Adr(X1, "msg").                        // ... buf - the string address
		Mov(X2, int64(len(msg))).              // ... len
		Movz(X16, sysClassUnix>>16, arch.Hw1). // x16 = 0x2000000 | ...
		Movk(X16, sysWrite, arch.Hw0).         // ... 0x4 = write
		Svc(trapMach).
		Mov(X0, 0). // exit(return code = 0)
		Movz(X16, sysClassUnix>>16, arch.Hw1).
		Movk(X16, sysExit, arch.Hw0).
		Svc(trapMach).
		Label("msg").
		Ascii(msg).
		Entry("start")
}

const (
	msg          = "hello world\n"
	trapMach     = 0x80
	fdStdout     = 1
	sysClassUnix = 0x2000000
	sysWrite     = 4
	sysExit      = 1
)

func TestAssembleHelloGolden(t *testing.T) {
	bin, buildErrs := helloMacOS().Build()
	require.Empty(t, buildErrs)
	res := bin.Assemble(0)
	require.Empty(t, res.Errs)
	require.Len(t, res.Code, 52) // 10 instructions + 12 data bytes

	golden, err := os.ReadFile("testdata/hello-macos.bin")
	require.NoError(t, err)

	// the Go-written program is byte-identical to the .s pipeline output
	require.Equal(t, golden, res.Code)
	require.Equal(t, uint64(0), res.Syms["start"])
	require.Equal(t, uint64(40), res.Syms["msg"])

	// the injected caller resolver reports this test file (10
	// instructions + 1 data line, labels emit nothing)
	require.Len(t, res.Lines, 11)
	for _, e := range res.Lines {
		require.Contains(t, e.Pos.File, "prog/arm64/program_test.go")
		require.Positive(t, e.Pos.Line)
	}
}

func TestAssembleLabels(t *testing.T) {
	// backward branch loop: three instructions between the label and the b.
	p := New().
		Label("loop").
		Mov(X0, 1).
		Mov(X1, 2).
		Mov(X2, 3).
		B("loop").
		Entry("loop")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)
	res := bin.Assemble(0x1000)
	require.Empty(t, res.Errs)
	require.Len(t, res.Code, 16)
	require.Equal(t, uint64(0x1000), res.Syms["loop"])

	// no resolver injected: lines carry the zero position
	require.Empty(t, res.Lines[0].Pos.File)
	require.Zero(t, res.Lines[0].Pos.Line)

	// b -12: imm26 of -3 words.
	require.Equal(t, uint32(0x17fffffd), le32(res.Code[12:16]))
}

func TestAssembleErrors(t *testing.T) {
	// undefined branch target
	bin, _ := New().B("nowhere").Build()
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")

	// undefined entry
	entryBin, _ := New().Label("a").Mov(X0, 0).Entry("gone").Build()
	entryErrs := entryBin.Assemble(0).Errs
	require.Len(t, entryErrs, 1)
	require.Contains(t, entryErrs[0].Error(), "gone")

	// deferred immediate error surfaces at Build
	_, buildErrs := New().Mov(X0, 1<<20).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "mov")
}

func TestLineMap(t *testing.T) {
	// the injected resolver: every line reports the synthetic position
	// (what a macro pointing at its own call site would do)
	p := New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 7) }).
		Label("start").
		Mov(X0, 1).
		Mov(X1, 2).
		Ascii("ab").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)
	res := bin.Assemble(0x100)
	require.Empty(t, res.Errs)
	require.Equal(t, []prog.LineEntry{
		prog.NewLineEntry(0x100, 4, prog.NewPos("synthetic.go", 7)),
		prog.NewLineEntry(0x104, 4, prog.NewPos("synthetic.go", 7)),
		prog.NewLineEntry(0x108, 2, prog.NewPos("synthetic.go", 7)),
	}, res.Lines)
}

func TestWithPosSwap(t *testing.T) {
	// WithPos applies from its call on: a later call swaps the resolver
	// for the lines after it (the caller-based one reports this test
	// file; the line is pinned loosely - it moves with edits)
	bin, _ := New().
		WithPos(func() prog.Pos { return prog.NewPos("macro.go", 1) }).
		Mov(X0, 1).
		WithPos(callerPos).
		Mov(X1, 2).
		Build()
	res := bin.Assemble(0)
	require.Empty(t, res.Errs)
	require.Len(t, res.Lines, 2)
	require.Equal(t, prog.NewPos("macro.go", 1), res.Lines[0].Pos)
	require.Contains(t, res.Lines[1].Pos.File, "prog/arm64/program_test.go")
	require.Positive(t, res.Lines[1].Pos.Line)
}

func le32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
