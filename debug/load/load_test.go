package load

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const helloSrc = `start:
    mov x0, #0
    b start
`

func TestSetupFor(t *testing.T) {
	cases := []struct {
		name   string
		arch   string
		base   uint64
		raw    bool
		qemu   string
		wantOK bool
	}{
		{name: "arm64", arch: "arm64", base: 0x40100000, qemu: "qemu-system-aarch64", wantOK: true},
		{
			name:   "aarch64 alias",
			arch:   "aarch64",
			base:   0x40100000,
			qemu:   "qemu-system-aarch64",
			wantOK: true,
		},
		{
			name:   "riscv64",
			arch:   "riscv64",
			base:   0x80000000,
			qemu:   "qemu-system-riscv64",
			wantOK: true,
		},
		{
			name:   "loong64 raw image",
			arch:   "loong64",
			base:   0x1c000000,
			raw:    true,
			qemu:   "qemu-system-loongarch64",
			wantOK: true,
		},
		{
			name:   "case insensitive",
			arch:   "ARM64",
			base:   0x40100000,
			qemu:   "qemu-system-aarch64",
			wantOK: true,
		},
		{name: "unknown", arch: "wasm"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			setup, err := SetupFor(c.arch)
			require.Equal(t, c.wantOK, err == nil)
			if !c.wantOK {
				return
			}

			require.Equal(t, c.base, setup.Base)
			require.Equal(t, c.raw, setup.Raw)
			require.Equal(t, c.qemu, setup.Tgt.QemuBinary())
		})
	}
}

func TestBuildSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.s")
	require.NoError(t, os.WriteFile(src, []byte(helloSrc), 0o644))

	setup, err := SetupFor("arm64")
	require.NoError(t, err)

	res, err := Build(setup, Input{SrcPath: src})
	require.NoError(t, err)

	require.Equal(t, src, res.SrcName)
	require.Equal(t, setup.Base, res.Syms["start"])
	require.NotEmpty(t, res.Img)
	require.Equal(t, "ELF", elfMagic(res.Img))

	require.NotEmpty(t, res.Lines)
	require.Equal(t, src, res.Lines[0].File)
	require.Equal(t, 2, res.Lines[0].Line) // the first instruction (line 1 is the label)
	require.Equal(t, setup.Base, res.Lines[0].Addr)

	require.Len(t, res.SrcLines, 4) // three lines + the empty tail after the last \n
	require.Contains(t, res.SrcLines[0], "start:")
}

func TestBuildBaseOverride(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.s")
	require.NoError(t, os.WriteFile(src, []byte(helloSrc), 0o644))

	setup, err := SetupFor("arm64")
	require.NoError(t, err)

	res, err := Build(setup, Input{SrcPath: src, Base: "0x1000"})
	require.NoError(t, err)
	require.Equal(t, uint64(0x1000), res.Syms["start"])
}

func TestBuildBinary(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "hello.elf")
	sym := filepath.Join(dir, "hello.sym")
	require.NoError(t, os.WriteFile(bin, []byte{0xde, 0xad, 0xbe, 0xef}, 0o644))
	require.NoError(t, os.WriteFile(sym, []byte("0x1000 start\n0x2008 loop\n"), 0o644))

	setup, err := SetupFor("arm64")
	require.NoError(t, err)

	res, err := Build(setup, Input{BinPath: bin, SymPath: sym})
	require.NoError(t, err)

	require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, res.Img)
	require.Equal(t, bin, res.SrcName)
	require.Equal(t, map[string]uint64{"start": 0x1000, "loop": 0x2008}, res.Syms)
	require.Empty(t, res.Lines)
}

func TestBuildErrors(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "broken.s")
	require.NoError(t, os.WriteFile(src, []byte("    frobnicate x0\n"), 0o644))

	setup, err := SetupFor("arm64")
	require.NoError(t, err)

	_, err = Build(setup, Input{SrcPath: src})
	require.Error(t, err)
	require.Contains(t, err.Error(), src+":1:")

	_, err = Build(setup, Input{})
	require.Error(t, err)

	_, err = Build(setup, Input{SrcPath: src, Base: "nothex"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad base")
}

// elfMagic reports "ELF" when the bytes open with the container's
// magic (the raw-image arches do not produce one).
func elfMagic(img []byte) string {
	if strings.HasPrefix(string(img), "\x7fELF") {
		return "ELF"
	}

	return ""
}
