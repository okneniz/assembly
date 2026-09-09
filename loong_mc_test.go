package assembly_test

// Byte-parity differential against the LLVM toolchain (the encoding
// oracle), loong64. Two levels:
//   - reloc-free sources: our .text vs the .text of the llvm-mc object,
//     byte for byte (llvm-mc resolves branch/alu fixups at assembly
//     time);
//   - la.pcrel/la.abs sources: llvm-mc leaves the pcala/abs fixups as
//     relocations, so the oracle is the linked image - ld.lld resolves
//     them, llvm-objcopy dumps the flat binary, and it must equal our
//     raw output (sections concatenated).
//
// llvm-mc's plain la is the GOT-indirect form (pcalau12i+ld.d through a
// linker-built GOT slot): it cannot resolve in this standalone
// assembler, so here la IS la.pcrel (explicit spelling accepted), and
// la.got is rejected with a dedicated error. Skipped without the tools;
// the docker toolchain container (make -C tests docker) has them all.

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

	"github.com/okneniz/assembly/asm/loong64/pseudo"
)

// toolPath - a tool from fixed candidates then PATH ("" when absent),
// the llvmMcPath way.
func toolPath(candidates []string) string {
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}

	return ""
}

func ldLldPath() string {
	return toolPath([]string{
		"/opt/homebrew/opt/lld/bin/ld.lld",
		"/opt/homebrew/bin/ld.lld",
		"/usr/local/opt/lld/bin/ld.lld",
		"/usr/local/bin/ld.lld",
		"ld.lld",
	})
}

func llvmObjcopyPath() string {
	return toolPath([]string{
		"/opt/homebrew/opt/llvm/bin/llvm-objcopy",
		"/opt/homebrew/bin/llvm-objcopy",
		"/usr/local/opt/llvm/bin/llvm-objcopy",
		"/usr/local/bin/llvm-objcopy",
		"llvm-objcopy",
	})
}

// loRunTool - a timed tool run; fails the test on an error.
func loRunTool(t *testing.T, name string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v: %s", name, args, err, out)
	}
}

// loLinkedImage - the flat binary of src assembled by llvm-mc and linked
// by ld.lld: .text at 0, .data at the SAME address our core gives it
// (-Tdata; lld's default puts .data into its own page-aligned RW segment,
// which is not the flat standalone layout).
func loLinkedImage(t *testing.T, src string, dataAddr uint64) []byte {
	t.Helper()
	dir := t.TempDir()
	obj, elf, bin := dir+"/mc.o", dir+"/mc.elf", dir+"/mc.bin"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mcCmd := exec.CommandContext(ctx, llvmMcPath(),
		"-assemble", "--triple=loongarch64", "-filetype=obj", "-o", obj)
	mcCmd.Stdin = strings.NewReader(src)
	if out, err := mcCmd.CombinedOutput(); err != nil {
		t.Fatalf("llvm-mc %q: %v: %s", src, err, out)
	}

	loRunTool(t, ldLldPath(),
		"-Ttext=0", fmt.Sprintf("-Tdata=%#x", dataAddr), "-o", elf, obj)
	loRunTool(t, llvmObjcopyPath(), "-O", "binary", elf, bin)

	img, err := os.ReadFile(bin)
	if err != nil {
		t.Fatalf("read %s: %v", bin, err)
	}

	return img
}

// loRequireImageEqual - the parity assertion with a failure message that
// stays readable for whole images (lengths + the first difference, not a
// full hex dump).
func loRequireImageEqual(t *testing.T, src string, ours, want []byte) {
	t.Helper()
	if bytes.Equal(ours, want) {
		return
	}

	n := min(len(ours), len(want))
	off := n
	for i := range n {
		if ours[i] != want[i] {
			off = i
			break
		}
	}

	t.Fatalf("our image ≠ lld-linked image for %q: ours %d bytes, lld %d bytes"+
		" (first difference at %#x: ours % x, lld % x)",
		src, len(ours), len(want), off,
		ours[min(off, len(ours)):min(off+8, len(ours))],
		want[min(off, len(want)):min(off+8, len(want))])
}

// loOursText - our .text of the source assembled at base 0.
func loOursText(t *testing.T, src string) []byte {
	t.Helper()
	res, errs := pseudo.Assemble(src, 0)
	require.Empty(t, errs, "assemble %q", src)
	require.NotEmpty(t, res.Sections, "no sections")
	return res.Sections[0].Data
}

// loOursImage - our whole raw image (the CLI raw format): the section
// data concatenated; the core places the sections without gaps, and the
// compared sources keep the lld layout identical (.text sizes leave .data
// aligned already).
func loOursImage(t *testing.T, src string) []byte {
	t.Helper()
	res, errs := pseudo.Assemble(src, 0)
	require.Empty(t, errs, "assemble %q", src)

	var out []byte
	next := uint64(0)
	for _, s := range res.Sections {
		require.Equal(t, next, s.Addr, "section %s: gap in the raw layout", s.Name)
		out = append(out, s.Data...)
		next += uint64(len(s.Data))
	}

	return out
}

// TestLoongVsLlvmMcText - reloc-free sources: direct .text byte parity
// with the llvm-mc object (branches and alu resolve at assembly).
func TestLoongVsLlvmMcText(t *testing.T) {
	mc := llvmMcPath()
	if mc == "" {
		t.Skip("no llvm-mc found on this host")
	}

	for _, tc := range []struct {
		name string
		src  string
	}{
		{
			"alu",
			"\n\taddi.w $t1, $t1, 1\n\txori $t2, $t1, 0xff\n\tslli.w $t3, $t2, 4\n",
		},
		{
			"uart-loop",
			"\n1:\n\tld.bu $t2, $t1, 0\n\tbeq $t2, $zero, 2f\n\tst.b $t2, $t0, 0\n\taddi.w $t1, $t1, 1\n\tb 1b\n2:\n\tb 2b\n",
		},
		{
			"lu12i-ori",
			"\n\tlu12i.w $t0, 0x1fe00\n\tori $t0, $t0, 0x1e0\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := mcAssemble(t, mc, tc.src)
			require.Equal(t, want, loOursText(t, tc.src),
				"our .text ≠ llvm-mc .text for %q\n  ours % x\n  llvm % x",
				tc.src, loOursText(t, tc.src), want)
		})
	}
}

// TestLoongVsLlvmLink - the linked-image parity: llvm-mc leaves
// la.pcrel/la.abs as relocations, ld.lld resolves them (with .data
// pinned to our flat-layout address), and the image must equal our raw
// output byte for byte.
func TestLoongVsLlvmLink(t *testing.T) {
	if llvmMcPath() == "" {
		t.Skip("no llvm-mc found on this host")
	}

	if ldLldPath() == "" || llvmObjcopyPath() == "" {
		t.Skip("no ld.lld/llvm-objcopy found on this host")
	}

	parity := func(t *testing.T, src string) {
		t.Helper()
		ours := loOursImage(t, src)
		res, errs := pseudo.Assemble(src, 0)
		require.Empty(t, errs, "re-assemble %q", src)

		var dataAddr uint64
		for _, s := range res.Sections {
			if s.Name == ".data" {
				dataAddr = s.Addr
			}
		}

		loRequireImageEqual(t, src, ours, loLinkedImage(t, src, dataAddr))
	}

	for _, tc := range []struct {
		name string
		src  string
	}{
		{
			"la-pcrel",
			"\n1:\n\tld.bu $t2, $t1, 0\n\tbne $t2, $zero, 1b\n\tla.pcrel $t1, msg\n\tb 1b\n\n\t.data\nmsg:\n\t.ascii \"hi\"\n\t.byte 0\n",
		},
		{
			"la-abs",
			"\n\tla.abs $t1, msg\n\tb 1f\n1:\n\tb 1b\n\n\t.data\nmsg:\n\t.ascii \"lo\"\n\t.byte 0\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parity(t, tc.src)
		})
	}

	t.Run("hello-loongarch", func(t *testing.T) {
		src, err := os.ReadFile("tests/examples/hello-asm/hello-loongarch.s")
		if err != nil {
			t.Skipf("example not available: %v", err)
		}

		// label-relative code + constants: base 0 matches -Ttext=0
		parity(t, string(src))
	})
}

// TestLaGotRejected - la.got is the one la form this assembler must not
// silently mangle: the GOT slot does not exist without a linker.
func TestLaGotRejected(t *testing.T) {
	_, errs := pseudo.Assemble("\n\tla.got $t1, msg\n\n\t.data\nmsg:\n\t.byte 0\n", 0)
	require.NotEmpty(t, errs, "la.got must be rejected")
	require.Contains(t, fmt.Sprint(errs), "la.got", "the error must name la.got")
}
