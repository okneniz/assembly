package assembly_test

// Byte-parity differential against llvm-mc (the encoding oracle), riscv:
// the assembler must choose exactly the bytes llvm-mc chooses for the
// hello example and for jump/branch snippets - including RVC compression
// of label targets (-mattr=+c, resolved by the layout relaxation). Skipped
// when no llvm-mc is installed; the docker toolchain container
// (make -C tests docker) runs it for real.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/file"
)

// rvMcAssemble - bytes of the .text llvm-mc assembles from src (riscv64,
// compressed). The object's .text starts at 0, so the snippets use only
// label-relative targets: the offsets are base-independent.
func rvMcAssemble(t *testing.T, src string) []byte {
	t.Helper()
	mc := llvmMcPath()
	path := t.TempDir() + "/mc.o"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, mc,
		"-assemble", "--triple=riscv64", "-mattr=+c", "-filetype=obj", "-o", path)
	cmd.Stdin = strings.NewReader(src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("llvm-mc %q: %v: %s", src, err, out)
		return nil
	}

	f, err := file.Detect(path)
	if err != nil {
		t.Fatalf("detect %s: %v", path, err)
		return nil
	}

	sec, err := f.CodeSection()
	if err != nil {
		t.Fatalf(".text of %s: %v", path, err)
		return nil
	}

	return sec.Data
}

// rvKnownLiDeviations - the llvm-mc li-ladder deviations, pinned per
// value so that any NEW deviation fails the differential while the known
// one stays visible (the larch convention, see knownLiDeviations):
//   - "shift-trick": 2048 - llvm-mc loads powers of two as
//     c.li rd, 1 + c.slli (4 bytes) when that beats the lui+addiw
//     ladder (8 bytes); ours always takes the ladder. Same loaded value.
var rvKnownLiDeviations = map[int64]string{
	2048: "shift-trick",
}

// TestRiscvVsLlvmMc - whole-snippet byte parity with llvm-mc. Both sides
// assemble the same source; every byte must match (sizes and encodings),
// no pinned deviations.
func TestRiscvVsLlvmMc(t *testing.T) {
	if llvmMcPath() == "" {
		t.Skip("no llvm-mc found on this host")
	}

	for _, tc := range []struct {
		name string
		src  string
	}{
		{
			"loop-jumps",
			"loop:\n  addi a1, a1, 1\n  j loop\nhang:\n  j hang\nmsg:\n  .ascii \"hi\"\n",
		},
		{
			"numeric-labels",
			"1:\n  beqz a0, 1f\n  j 1b\n1:\n  j 1b\n",
		},
		{
			"far-jump",
			"  j target\n  .space 4096\ntarget:\n  nop\n",
		},
		{
			"boundary-2046",
			"  j target\n  .space 2044\ntarget:\n  nop\n",
		},
		{
			"option-norvc",
			"  j target\n.option norvc\n  j target\n.option rvc\n  j target\ntarget:\n  nop\n",
		},
		{
			"branch-loop",
			"start:\n  beq a1, a2, done\n  lb a5, 0(a1)\n  sb a5, 0(a0)\n  addi a1, a1, 1\n  j start\ndone:\n  ret\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, errs := pseudo.Assemble(tc.src, 0)
			require.Empty(t, errs, "assemble %q", tc.src)
			require.NotEmpty(t, res.Sections, "no sections")

			want := rvMcAssemble(t, tc.src)
			require.Equal(t, want, res.Sections[0].Data,
				"our bytes ≠ llvm-mc bytes for %q\n  ours % x\n  llvm % x",
				tc.src, res.Sections[0].Data, want)
		})
	}

	// the li ladder: every rung must pick exactly the llvm-mc rung (the
	// addi rung, the lui rung, the lui+addiw pair, the negative-page
	// c.lui forms), except the pinned deviation below
	for _, v := range []int64{
		-2048, -1, 0, 1, 42, 2047, // addi rung
		4095, 0x1FFF, // first lui+addiw rung
		0x5555, 0x12345678, // the pair
		0x800, 0xFFF, 0x7FFF, 0x8000, // sign-edge los
		-0xBB8, -0x1000, -0x2000, -0x8000, // negative pages (c.lui)
		-0x7FFFFFFF,            // full negative pair
		0x10000000, 0x7FFFFFFF, // pure lui edges
	} {
		t.Run(fmt.Sprintf("li-%#x", v), func(t *testing.T) {
			src := fmt.Sprintf("li a0, %#x", v)
			res, errs := pseudo.Assemble(src, 0)
			require.Empty(t, errs, "assemble %q", src)
			want := rvMcAssemble(t, src)

			if class, known := rvKnownLiDeviations[v]; known &&
				!bytes.Equal(want, res.Sections[0].Data) {
				t.Logf("KNOWN llvm deviation (%s, reported): li of %#x"+
					" - ours % x, llvm % x",
					class, v, res.Sections[0].Data, want)
				return
			}

			require.Equal(t, want, res.Sections[0].Data,
				"li of %#x: our ladder ≠ llvm-mc ladder\n  ours % x\n  llvm % x",
				v, res.Sections[0].Data, want)
		})
	}

	t.Run("hello-riscv", func(t *testing.T) {
		src, err := os.ReadFile("tests/examples/hello-asm/hello-riscv.s")
		if err != nil {
			t.Skipf("example not available: %v", err)
		}

		// label-relative targets make the .text bytes base-independent;
		// base 0 matches the llvm-mc object's .text
		res, errs := pseudo.Assemble(string(src), 0)
		require.Empty(t, errs, "errors: %v", errs)

		want := rvMcAssemble(t, string(src))
		require.Equal(t, want, res.Sections[0].Data,
			"our bytes ≠ llvm-mc bytes\n  ours % x\n  llvm % x",
			res.Sections[0].Data, want)
	})
}
