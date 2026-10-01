package linktests

// The execution law of the generated specimens: the linked program runs
// natively and exits with the generator's bookkeeping value - assembly
// semantics and Go semantics computed from one seed must agree. The
// image is a Mach-O on macOS, an ELF on Linux (the docker gate); the
// exit idiom follows the OS, the value never changes.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/link/arm64"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
)

func TestLinkExecMatchesGenerator(t *testing.T) {
	if runtime.GOARCH != "arm64" ||
		(runtime.GOOS != "darwin" && runtime.GOOS != "linux") {
		t.Skip("native arm64 macOS/Linux only")
	}

	ohsnap.Check(t, 100, linkarb.NewProgramArb(seedRnd(t)), func(p linkarb.Program) bool {
		p.ExitLinux = runtime.GOOS == "linux"

		blob, err := imageOf(p)
		if err != nil {
			t.Logf("link: %v", err)
			return false
		}

		path := filepath.Join(t.TempDir(), "linked")
		if err := os.WriteFile(path, blob, 0o755); err != nil {
			t.Fatalf("write: %v", err)
		}

		got := runExit(t, path)
		if got != p.Exit() {
			t.Logf("exit %d, want %d", got, p.Exit())
			return false
		}

		return true
	})
}

// imageOf links the specimen through the native writer of the OS.
func imageOf(p linkarb.Program) ([]byte, error) {
	srcs := sources(p)
	if runtime.GOOS == "darwin" {
		img, err := arm64.Macho(srcs, "")
		if err != nil {
			return nil, err
		}

		return img.Bytes(), nil
	}

	return arm64.ELF(srcs, "", 0x10000)
}

// runExit runs the binary and returns its exit code (0 for a clean
// exit - the value may legitimately be zero).
func runExit(t *testing.T, path string) int {
	t.Helper()

	err := exec.CommandContext(context.Background(), path).Run()

	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return exitErr.ExitCode()
	}

	if err == nil {
		return 0
	}

	t.Fatalf("run %s: %v", path, err)
	return -1
}
