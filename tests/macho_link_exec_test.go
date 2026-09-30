package tests

// The linker end to end: two sources linked in memory through the
// library (link/arm64.Macho) - the main file calls into the util file,
// the data of both merge into one stream, the bss of both aggregates
// into the single tail - and the image executes natively.
//
// The addressing is adrp+offset (the statics-test idiom): the offsets
// are the merged layout's - msg at __data+0 (main.S deposits first),
// val at +4 (util.S's data past it), the bss tail pad+scratch at +8/+12
// (main.S's reserve first). No symbolic literal pools here: a pool slot
// carries an absolute pointer, and the Mach-O writer emits no
// chained-fixup rebases yet - under PIE the unslid pointer faults. That
// gap (any absolute pointer in the image, .quad label included) is the
// file package's, surfaced by this round.

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

// linkMainSrc - two bumps of the util counter (val 2 -> 4), a read of
// this file's own static (msg, 3), the util file's bss scratch written
// and read back (0 + 9), exit(4 + 3 + 9).
const linkMainSrc = `
.global _start

.text
_start:
    bl bump
    bl bump
    adrp x1, val
    add x1, x1, #4           // val at __data+4, past this file's msg
    ldr w0, [x1]             // 4
    adrp x2, msg
    add x2, x2, #0           // msg at __data+0
    ldr w3, [x2]             // 3
    add w0, w0, w3

    adrp x4, scratch
    add x4, x4, #12          // scratch at the bss tail, past pad
    ldr w5, [x4]
    add w5, w5, #9
    str w5, [x4]
    ldr w5, [x4]             // 0 + 9
    add w0, w0, w5

    movz x16, #0x200, lsl #16
    movk x16, #0x1
    svc #0x80

.data
msg:
    .word 3

.bss
pad:
    .zero 4
`

// linkUtilSrc - the counter of the link (val) and the scratch the main
// file writes through.
const linkUtilSrc = `
.data
val:
    .word 2

.text
bump:
    adrp x1, val
    add x1, x1, #4
    ldr w0, [x1]
    add w0, w0, #1
    str w0, [x1]
    ret

.bss
scratch:
    .zero 4
`

func TestMachOLinkExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	img, err := arm64.Macho([]link.Source{
		{File: "main.S", Src: linkMainSrc},
		{File: "util.S", Src: linkUtilSrc},
	}, "")
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "linked")
	require.NoError(t, os.WriteFile(path, img.Bytes(), 0o755))

	runErr := exec.CommandContext(context.Background(), path).Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	require.Equal(t, 16, exitErr.ExitCode()) // 4 + 3 + 9
}
