package loong64

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// The Go counterpart of tests/examples/hello-asm/hello-loongarch.s: one
// string through the 16550A UART, then an idle loop.
func helloLoong() *Program {
	return New().
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

	code, syms, errs := bin.Assemble(0x1c000000)
	require.Empty(t, errs)

	golden, err := os.ReadFile("testdata/hello-loongarch.bin")
	if err != nil {
		t.Skip("golden binary not built: run the assembly CLI")
	}

	require.Equal(t, golden, code)
	require.Equal(t, uint64(0x1c000000), syms["start"])
}

func TestAssembleLoongErrors(t *testing.T) {
	// undefined branch target
	bin, _ := New().B("nowhere").Build()
	_, _, errs := bin.Assemble(0)
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")

	// undefined la target
	laBin, _ := New().La(T0, "gone").Build()
	_, _, laErrs := laBin.Assemble(0)
	require.Len(t, laErrs, 1)
	require.Contains(t, laErrs[0].Error(), "gone")

	// deferred immediate error surfaces at Build
	_, buildErrs := New().Lu12iW(T0, 1<<20).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "lu12i.w")
}
