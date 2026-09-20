package loong64

import (
	"github.com/okneniz/assembly/unit"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// callerPos - a caller-based resolver for the tests: invoked inside a
// chain method, two frames up is the code calling the chain (this file).
func callerPos() unit.Pos {
	_, file, line, _ := runtime.Caller(2)
	return unit.NewPos(file, line)
}

// The Go counterpart of tests/examples/hello-asm/hello-loongarch.s: one
// string through the 16550A UART, then an idle loop.
func helloLoong() *Program {
	return New(unit.New()).
		WithPos(callerPos).
		Label("start").
		// $t0 = the UART data register (16550A, byte-wide at 0x1fe001e0).
		Lu12iW(T0, 0x1fe00).
		Ori(T0, T0, 0x1e0).
		// $t1 = the message cursor; for (; *p; p++) *uart = *p;
		La(T1, "msg").
		Label("loop").
		LdBu(T2, T1, 0).
		Beq(T2, Zero, "idle").
		StB(T2, T0, 0).
		AddiW(T1, T1, 1).
		B("loop").
		// idle: the gate stops the machine once the line is out.
		Label("idle").
		B("idle").
		Label("msg").
		Ascii("hello world\n").
		Bytes(0).
		Entry("start")
}

func TestAssembleLoongGolden(t *testing.T) {
	bin, buildErrs := helloLoong().Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0x1c000000)
	require.Empty(t, res.Errs)

	golden, err := os.ReadFile("testdata/hello-loongarch.bin")
	if err != nil {
		t.Skip("golden binary not built: run the assembly CLI")
	}

	require.Equal(t, golden, res.Code)
	require.Equal(t, uint64(0x1c000000), res.Syms["start"])

	// the injected caller resolver reports this test file (labels emit
	// nothing; the la pair is one entry)
	require.Len(t, res.Lines, 11)
	for _, e := range res.Lines {
		require.Contains(t, e.Pos.File, "prog/loong64/program_test.go")
		require.Positive(t, e.Pos.Line)
	}
}

func TestAssembleLoongErrors(t *testing.T) {
	// undefined branch target
	bin, _ := New(unit.New()).B("nowhere").Build()
	errs := bin.Assemble(0).Errs
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")

	// undefined la target
	laBin, _ := New(unit.New()).La(T0, "gone").Build()
	laErrs := laBin.Assemble(0).Errs
	require.Len(t, laErrs, 1)
	require.Contains(t, laErrs[0].Error(), "gone")

	// deferred immediate error surfaces at Build
	_, buildErrs := New(unit.New()).Lu12iW(T0, 1<<20).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "lu12i.w")
}

func TestLineMap(t *testing.T) {
	// the injected resolver: the la pair is ONE entry of 8 bytes (one
	// chain call), data lines carry their byte size
	p := New(unit.New()).
		WithPos(func() unit.Pos { return unit.NewPos("synthetic.go", 3) }).
		Label("start").
		Lu12iW(T0, 0x1fe00).
		La(T1, "msg").
		Label("msg").
		Ascii("hi").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)
	res := bin.Assemble(0)
	require.Empty(t, res.Errs)
	require.Equal(t, []unit.LineEntry{
		unit.NewLineEntry(0, 4, unit.NewPos("synthetic.go", 3)),
		unit.NewLineEntry(4, 8, unit.NewPos("synthetic.go", 3)),
		unit.NewLineEntry(12, 2, unit.NewPos("synthetic.go", 3)),
	}, res.Lines)
}
