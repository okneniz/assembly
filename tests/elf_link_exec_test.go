package tests

// The ELF end of the linker: the same two-source link through
// link/arm64.ELF, executed natively on Linux arm64 (the docker gate
// runs the suite in an arm64 container; the host skips).

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/link"
	"github.com/okneniz/assembly/link/arm64"
)

// elfLinkMainSrc - bump twice (val 2 -> 4 in the other file's data),
// read it back, exit. The addresses are flat: val at base+text+4 (past
// this file's msg).
const elfLinkMainSrc = `
.global _start

.text
_start:
    bl bump
    bl bump
    ldr x1, =val
    ldr w0, [x1]             // 4
    movz x16, #0x200, lsl #16
    movk x16, #0x1
    svc #0x80

.data
msg:
    .word 3
`

const elfLinkUtilSrc = `
.data
val:
    .word 2

.text
bump:
    ldr x1, =val
    ldr w0, [x1]
    add w0, w0, #1
    str w0, [x1]
    ret
`

func TestELFLinkExec(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "arm64" {
		t.Skip("native Linux arm64 only (the docker gate)")
	}

	blob, err := arm64.ELF([]link.Source{
		{File: "main.S", Src: elfLinkMainSrc},
		{File: "util.S", Src: elfLinkUtilSrc},
	}, "", 0x10000)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "linked.elf")
	require.NoError(t, os.WriteFile(path, blob, 0o755))

	runErr := exec.CommandContext(context.Background(), path).Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	require.Equal(t, 4, exitErr.ExitCode()) // val 2 bumped twice
}
