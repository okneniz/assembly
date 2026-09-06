package assembly_test

// Byte-parity differential against llvm-mc (the encoding oracle), arm64:
// the three hello examples must assemble into exactly the bytes llvm-mc
// chooses. The compared bytes are the whole code section of the llvm-mc
// object (.text of the ELF triples, __text of the Mach-O one -
// file.Detect + CodeSection handles both; llvm-objcopy's Mach-O
// --only-section is not usable). The examples are position-independent
// with in-section adr targets, so base 0 matches the object's code
// section at VMA 0. Skipped without llvm-mc; the docker toolchain
// container (make -C tests docker) runs it.

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/file"
)

// a64McText - the code section bytes of the llvm-mc object for src.
func a64McText(t *testing.T, triple, src string) []byte {
	t.Helper()
	path := t.TempDir() + "/mc.o"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, llvmMcPath(),
		"-assemble", "-triple="+triple, "-filetype=obj", "-o", path)
	cmd.Stdin = strings.NewReader(src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("llvm-mc %q (%s): %v: %s", src, triple, err, out)
	}

	f, err := file.Detect(path)
	if err != nil {
		t.Fatalf("detect %s: %v", path, err)
	}

	sec, err := f.CodeSection()
	if err != nil {
		t.Fatalf("code section of %s: %v", path, err)
	}

	return sec.Data
}

// TestArm64VsLlvmMc - the hello examples vs llvm-mc, whole code sections,
// byte for byte.
func TestArm64VsLlvmMc(t *testing.T) {
	if llvmMcPath() == "" {
		t.Skip("no llvm-mc found on this host")
	}

	for _, tc := range []struct {
		name   string
		path   string
		triple string
	}{
		{"hello-macos", "tests/examples/hello-asm/hello-macos.s", "arm64-apple-macos"},
		{"hello-linux", "tests/examples/hello-asm/hello-linux.s", "aarch64-linux-gnu"},
		{"hello-arm-vm", "tests/examples/hello-asm/hello-arm-vm.s", "aarch64"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src, err := os.ReadFile(tc.path)
			if err != nil {
				t.Skipf("example not available: %v", err)
			}

			res, errs := alias.Assemble(string(src), 0)
			require.Empty(t, errs, "assemble %s", tc.path)
			require.NotEmpty(t, res.Sections, "no sections")

			want := a64McText(t, tc.triple, string(src))
			require.Equal(t, want, res.Sections[0].Data,
				"our .text ≠ llvm-mc code section for %s\n  ours % x\n  llvm % x",
				tc.path, res.Sections[0].Data, want)
		})
	}
}
