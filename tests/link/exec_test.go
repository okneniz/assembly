package linktests

// The execution law of the generated specimens, one subtest per
// architecture: the linked program runs NATIVELY on its OS and exits
// with the generator's bookkeeping value - assembly semantics and Go
// semantics computed from one seed must agree. arm64 runs as a Mach-O
// on macOS and as an ELF on Linux (the docker gate); riscv64 and
// loong64 are Linux-only ELFs (they wake up when a matching container
// joins the gates; today they skip everywhere the suite runs).

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
	linkloong "github.com/okneniz/assembly/link/loong64"
	linkriscv "github.com/okneniz/assembly/link/riscv"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
)

func TestLinkExecMatchesGenerator(t *testing.T) {
	for _, ac := range arches() {
		t.Run(ac.name, func(t *testing.T) {
			if !runsNatively(ac.name) {
				t.Skipf("no native %s on %s/%s", ac.name, runtime.GOOS, runtime.GOARCH)
			}

			arb := linkarb.NewArchProgramArb(seedRnd(t), ac.arch)
			ohsnap.Check(t, 100, arb, func(a linkarb.ArchProgram) bool {
				a.Prog.ExitLinux = runtime.GOOS == "linux" // macOS keeps its own idiom

				blob, err := imageOf(ac.name, a)
				if err != nil {
					t.Logf("link: %v", err)
					return false
				}

				path := filepath.Join(t.TempDir(), "linked")
				if err := os.WriteFile(path, blob, 0o755); err != nil {
					t.Fatalf("write: %v", err)
				}

				got := runExit(t, path)
				if got != a.Prog.Exit() {
					t.Logf("exit %d, want %d", got, a.Prog.Exit())
					return false
				}

				return true
			})
		})
	}
}

// runsNatively - the gate matrix: an arch executes on its own OS/ARCH
// pair only (arm64 both ways - Mach-O on macOS, ELF on Linux).
func runsNatively(arch string) bool {
	switch arch {
	case "arm64":
		return runtime.GOARCH == "arm64" &&
			(runtime.GOOS == "darwin" || runtime.GOOS == "linux")
	case "riscv64":
		return runtime.GOOS == "linux" && runtime.GOARCH == "riscv64"
	case "loong64":
		return runtime.GOOS == "linux" && runtime.GOARCH == "loong64"
	}

	return false
}

// imageOf links the specimen through the native writer of the target.
func imageOf(arch string, a linkarb.ArchProgram) ([]byte, error) {
	srcs := sources(a)
	switch arch {
	case "arm64":
		if runtime.GOOS == "darwin" {
			img, err := arm64.Macho(srcs, "")
			if err != nil {
				return nil, err
			}

			return img.Bytes(), nil
		}

		return arm64.ELF(srcs, "", 0x10000)
	case "riscv64":
		return linkriscv.ELF(srcs, "", 0x10000)
	default:
		return linkloong.ELF(srcs, "", 0x10000)
	}
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
