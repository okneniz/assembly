package linktests

// The clang differential: the same generated sources through TWO
// independent toolchains - our in-memory link (assemble each source,
// one resolve) and clang (assemble each .s to an .o, ld links with an
// explicit entry) - must produce programs that exit with the same code,
// the generator's value. Catches the blind spots of the internal laws:
// both sides of those run through OUR assembler.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/link"
	"github.com/okneniz/assembly/link/arm64"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
)

func TestLinkClangDifferential(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("no clang on the host")
	}

	ohsnap.Check(t, 40, linkarb.NewProgramArb(seedRnd(t)), func(p linkarb.Program) bool {
		srcs := sources(p)

		ours := filepath.Join(t.TempDir(), "ours")
		img, err := arm64.Macho(srcs, "")
		if err != nil {
			t.Logf("our link: %v", err)
			return false
		}

		if err := os.WriteFile(ours, img.Bytes(), 0o755); err != nil {
			t.Fatalf("write: %v", err)
		}

		ref, err := clangBuild(t, srcs)
		if err != nil {
			t.Logf("clang: %v", err)
			return false
		}

		got, want := runExit(t, ours), runExit(t, ref)
		if got != p.Exit() || want != p.Exit() {
			t.Logf("ours %d, clang %d, want %d", got, want, p.Exit())
			return false
		}

		return true
	})
}

// clangBuild assembles every source separately (the source on stdin,
// `-x assembler -`) and links them once: the default link (libSystem
// and the crt glue ride along - harmless, the entry is ours) with
// -Wl,-e,_start naming it. -nostdlib is not an option: the newer ld
// refuses a dynamic image without libSystem, and -static is SIGKILLed
// by AMFI.
func clangBuild(t *testing.T, srcs []link.Source) (string, error) {
	t.Helper()

	dir := t.TempDir()
	objs := make([]string, len(srcs))
	for i := range srcs {
		objs[i] = filepath.Join(dir, fmt.Sprintf("t%d.o", i))

		cmd := exec.CommandContext(
			context.Background(),
			"clang", "-c", "-arch", "arm64", "-o", objs[i], "-x", "assembler", "-",
		)
		cmd.Stdin = strings.NewReader(srcs[i].Src)

		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("clang -c %s: %w: %s", srcs[i].File, err, out)
		}
	}

	out := filepath.Join(dir, "ref")
	args := append([]string{"-Wl,-e,_start", "-o", out}, objs...)
	if out, err := (exec.CommandContext(
		context.Background(), "clang", args...,
	)).CombinedOutput(); err != nil {
		return "", fmt.Errorf("clang link: %w: %s", err, out)
	}

	return out, nil
}
