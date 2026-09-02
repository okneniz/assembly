package tests

// The -format macho pipeline end to end: hello-macos.s assembled in
// process, wrapped by file.WriteMachO, and executed natively - the whole
// toolchain is the one assembly binary, no cc/ld/codesign anywhere.

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/file"
)

func TestMachoHelloExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	src, err := os.ReadFile("examples/hello-asm/hello-macos.s")
	require.NoError(t, err)

	res, errs := alias.Assemble(string(src), 0)
	require.Empty(t, errs)

	var raw []byte
	for _, sec := range res.Sections {
		raw = append(raw, sec.Data...)
	}

	entry := uint64(0) // the offset of the entry inside the code
	for _, name := range []string{"start", "_start"} {
		if a, ok := res.Symbols[name]; ok {
			entry = a
			break
		}
	}

	bin, err := file.WriteMachO(raw, entry)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "hello")
	require.NoError(t, os.WriteFile(path, bin, 0o755))

	out, err := exec.Command(path).Output()
	require.NoError(t, err, "the macho executable must run as-is")
	require.Equal(t, "hello world\n", string(out))
}
