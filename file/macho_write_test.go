package file

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/macho"
)

// The emitted image must parse back with the package's own parser: the
// header fields, the three segments, the __text section, the LC_MAIN entry,
// and the code verbatim at MachoCodeOff.
func TestMachoWriterRoundTrip(t *testing.T) {
	code := "\x1f\x20\x03\xd5" // nop

	bin, err := WriteMachO([]byte(code), 0)
	require.NoError(t, err)
	require.Greater(t, len(bin), MachoCodeOff+len(code))

	path := filepath.Join(t.TempDir(), "hello")
	require.NoError(t, os.WriteFile(path, bin, 0o755))

	f, err := macho.Open(path)
	require.NoError(t, err)

	hdr := f.Header()
	require.Equal(t, uint32(machoMagic64), hdr.Magic)
	require.Equal(t, int32(machoCPUArm64), hdr.CpuType)
	require.Equal(t, macho.MH_EXECUTE, hdr.FileType)
	require.Equal(t, uint32(16), hdr.Ncmds)

	segs := f.Segments()
	require.Len(t, segs, 3)
	require.Equal(t, "__PAGEZERO", segs[0].SegName)
	require.Equal(t, "__TEXT", segs[1].SegName)
	require.Equal(t, uint64(machoVMAddr), segs[1].Vmaddr)
	require.Equal(t, uint64(0), segs[1].Fileoff)
	require.Equal(t, uint64(16384), segs[1].Vmsize)
	require.Equal(t, "__LINKEDIT", segs[2].SegName)

	entry, ok := f.Entry()
	require.True(t, ok, "LC_MAIN entry")
	require.Equal(t, uint64(machoVMAddr+MachoCodeOff), entry)

	// the code sits verbatim after the load commands
	require.Equal(t, code, string(bin[MachoCodeOff:MachoCodeOff+4]))
}

// TestMachoWriterExec - the whole point of the writer: on an arm64 Mac the
// emitted binary must run as-is (AMFI strict validation passes, dyld calls
// the entry). The program exits with code 42 via a direct syscall - no
// libc, no return.
func TestMachoWriterExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	// Darwin arm64 exit(42): movz x16, #0x200, lsl 16; movk x16, #1
	// (0x2000001 = SYSCALL_CLASS_UNIX|exit); mov x0, #42; svc #0x80.
	code := []byte{
		0x10, 0x40, 0xa0, 0xd2,
		0x30, 0x00, 0x80, 0xf2,
		0x40, 0x05, 0x80, 0xd2,
		0x01, 0x10, 0x00, 0xd4,
	}

	bin, err := WriteMachO(code, 0)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "prog")
	require.NoError(t, os.WriteFile(path, bin, 0o755))

	cmd := exec.Command(path)
	runErr := cmd.Run()

	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	require.Equal(t, 42, exitErr.ExitCode())
}

func TestMachoWriterErrors(t *testing.T) {
	// no text
	_, err := WriteMachO(nil, 0)
	require.ErrorContains(t, err, "no text")

	// entry outside the text
	_, err = WriteMachO([]byte{0}, 4)
	require.ErrorContains(t, err, "outside")
}
