package tests

// The no-contract pipeline end to end: a chain program with a data static
// and a bss one, assembled and wrapped by one library call (Binary.MachO
// - the writer's own placement policy inside), and executed natively. No
// address is computed anywhere outside the file package - the program
// just uses labels.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	prog "github.com/okneniz/assembly/prog/arm64"
)

func TestMachOProgramExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	p := prog.New().Entry("start")

	// counter (data, 7) incremented in place and re-read; scratch (bss)
	// must read zero, accept a write, and read it back; exit(9)
	p.Label("start")
	p.La(prog.X0, "counter")
	p.Ldr(prog.X1, prog.X0, 0)
	p.AddImm(prog.X1, prog.X1, 1, arch.NoSh12)
	p.Str(prog.X1, prog.X0, 0)
	p.Ldr(prog.X0, prog.X0, 0)
	p.La(prog.X2, "scratch")
	p.Ldr(prog.X3, prog.X2, 0)
	p.AddImm(prog.X3, prog.X3, 1, arch.NoSh12)
	p.Str(prog.X3, prog.X2, 0)
	p.Ldr(prog.X3, prog.X2, 0)
	p.AddShift(prog.X0, prog.X0, prog.X3, 0, arch.LSL)
	p.Movz(prog.X16, 0x200, arch.Hw1)
	p.Movk(prog.X16, 1, arch.Hw0)
	p.Svc(0x80)

	p.Data()
	p.Label("counter").Quad(7)
	p.Label("scratch").Bss(16)

	bin, errs := p.Build()
	require.Empty(t, errs)

	img, err := bin.MachO("start")
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "statics")
	require.NoError(t, os.WriteFile(path, img.Bytes(), 0o755))

	out := exec.CommandContext(context.Background(), path)
	runErr := out.Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	require.Equal(t, 9, exitErr.ExitCode()) // 7+1 from data, 0+1 from bss
}
