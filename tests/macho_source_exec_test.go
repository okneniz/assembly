package tests

// The .s pipeline end to end with data: a source with .text, .data, and
// .bss - including a symbolic adrp over the page gap between the
// segments - assembled and wrapped by one library call (alias.MachO, the
// writer's own placement inside) and executed natively.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/arm64/alias"
)

// staticsSrc - val (.data, 4) incremented in place and re-read; zeroed
// (.bss) must read zero, accept a write, and read it back; exit(7).
// val sits at __data+0 (the page-aligned segment start), zeroed at +8
// (the bss reserve is 8-aligned past the 4 data bytes).
const staticsSrc = `
.global _start

.text
_start:
    adrp    x0, val                 // symbolic adrp over the page gap
    add     x0, x0, #0              // val at __data+0
    ldr     w1, [x0]
    add     w1, w1, #1
    str     w1, [x0]
    ldr     w0, [x0]                // 4+1 = 5

    adrp    x1, zeroed
    add     x1, x1, #8              // zeroed at __bss+0, 8 past val
    ldr     w2, [x1]
    add     w2, w2, #2
    str     w2, [x1]
    ldr     w2, [x1]                // 0+2 = 2
    add     w0, w0, w2              // 7

    movz    x16, #0x200, lsl #16
    movk    x16, #0x1               // exit
    svc     #0x80

.data
val:
    .word 4

.bss
zeroed:
    .zero 4
`

func TestMachOStaticSourceExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	img, errs := alias.MachO(staticsSrc, "_start")
	require.Empty(t, errs)

	path := filepath.Join(t.TempDir(), "statics")
	require.NoError(t, os.WriteFile(path, img.Bytes(), 0o755))

	out := exec.CommandContext(context.Background(), path)
	runErr := out.Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	require.Equal(t, 7, exitErr.ExitCode()) // 4+1 from data, 0+2 from bss
}

// TestMachOStaticLayout - the sections and symbols of the same source,
// placed by the policy: .text in __TEXT, .data and .bss in __DATA past
// the page-rounded text, the data symbols page-aligned.
func TestMachOStaticLayout(t *testing.T) {
	res, errs := alias.AssembleLayout(staticsSrc, alias.MachoLayout)
	require.Empty(t, errs)

	require.Len(t, res.Sections, 3)
	text, data, bss := res.Sections[0], res.Sections[1], res.Sections[2]

	require.Equal(t, ".text", text.Name)
	require.Zero(t, text.Addr % 4)
	require.Equal(t, ".data", data.Name)
	require.Zero(t, data.Addr % 16384) // the page-aligned __DATA start
	require.Equal(t, uint64(4), uint64(data.Size))
	require.Equal(t, ".bss", bss.Name)
	require.True(t, bss.Nobits)
	require.Equal(t, data.Addr+8, bss.Addr) // 8-aligned past the 4 data bytes

	require.Equal(t, data.Addr, res.Symbols["val"])
	require.Equal(t, bss.Addr, res.Symbols["zeroed"])
	require.Equal(t, text.Addr, res.Symbols["_start"])
}
