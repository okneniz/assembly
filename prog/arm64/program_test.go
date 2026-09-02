package arm64

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// hello-macos written on the chain: the Go counterpart of
// tests/examples/hello-asm/hello-macos.s, byte-identical when assembled.
func helloMacOS() *Program {
	return New().
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
	code, syms, errs := bin.Assemble(0)
	require.Empty(t, errs)
	require.Len(t, code, 52) // 10 instructions + 12 data bytes

	golden, err := os.ReadFile("testdata/hello-macos.bin")
	require.NoError(t, err)

	// the Go-written program is byte-identical to the .s pipeline output
	require.Equal(t, golden, code)
	require.Equal(t, uint64(0), syms["start"])
	require.Equal(t, uint64(40), syms["msg"])
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
	code, syms, errs := bin.Assemble(0x1000)
	require.Empty(t, errs)
	require.Len(t, code, 16)
	require.Equal(t, uint64(0x1000), syms["loop"])

	// b -12: imm26 of -3 words.
	require.Equal(t, uint32(0x17fffffd), le32(code[12:16]))
}

func TestAssembleErrors(t *testing.T) {
	// undefined branch target
	bin, _ := New().B("nowhere").Build()
	_, _, errs := bin.Assemble(0)
	require.Len(t, errs, 1)
	require.Contains(t, errs[0].Error(), "nowhere")

	// undefined entry
	entryBin, _ := New().Label("a").Mov(X0, 0).Entry("gone").Build()
	_, _, entryErrs := entryBin.Assemble(0)
	require.Len(t, entryErrs, 1)
	require.Contains(t, entryErrs[0].Error(), "gone")

	// deferred immediate error surfaces at Build
	_, buildErrs := New().Mov(X0, 1<<20).Build()
	require.Len(t, buildErrs, 1)
	require.Contains(t, buildErrs[0].Error(), "mov")
}

func le32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
