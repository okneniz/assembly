package tests

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

// instrText returns the normalized text of the first instruction of a binary
// (for decode-equivalence comparison).
func instrText(b []byte, addr uint64) string {
	insts, err := arm64.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return objdump.Normalize(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// loongInstrText - the same for LoongArch.
func loongInstrText(b []byte, addr uint64) string {
	insts, err := loong64.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return objdump.Normalize(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// riscvInstrText - the same for RISC-V.
func riscvInstrText(b []byte, addr uint64) string {
	insts, err := riscv.Parse(addr)(parsecbytes.Buffer(b))
	if err != nil {
		return ""
	}

	if len(insts) == 0 {
		return ""
	}

	return objdump.Normalize(insts[0].ObjDump(disasm.DefaultViewCtx()))
}

// fileSections converts sections of an assembly result into file-package
// sections.
func fileSections(secs []asm.Section) []file.Section {
	lastProgbits := -1
	for i, s := range secs {
		if !s.Nobits {
			lastProgbits = i
		}
	}

	out := make([]file.Section, len(secs))
	for i, s := range secs {
		data := s.Data
		if s.Nobits && i < lastProgbits {
			data = make([]byte, s.Size)
		}

		out[i] = *file.NewSection(s.Name, "", s.Addr, 0, uint64(s.Size), data)
	}

	return out
}

// TestHelloAsmExample - the demo patient examples/hello-asm: the source
// assembles without errors, and its binary passes a full byte-exact
// round-trip (decode → ObjDump text → assemble → same bytes), including the
// .word fallback for data.
func TestHelloAsmExample(t *testing.T) {
	cases := []struct {
		src      string
		base     uint64
		assemble func(string, uint64) (*asm.Result, []asm.AsmError)
	}{
		{"examples/hello-asm/hello-macos.s", 0, alias.Assemble},
		{"examples/hello-asm/hello-linux.s", 0, alias.Assemble},
		{"examples/hello-asm/hello-arm-vm.s", 0x40100000, alias.Assemble},
		{"examples/hello-asm/hello-riscv.s", 0x80000000, pseudo.Assemble},
		{"examples/hello-asm/hello-loongarch.s", 0x1c000000, lpseudo.Assemble},
	}
	for _, c := range cases {
		t.Run(filepath.Base(c.src), func(t *testing.T) {
			data, err := os.ReadFile(c.src)
			if err != nil {
				t.Skipf("example source not available: %v", err)
			}

			res, errs := c.assemble(string(data), c.base)
			require.Empty(t, errs)
			bin := res.Sections[0].Data
			// minimum: code + "hello world\n"
			require.GreaterOrEqual(t, len(bin), 40)
			require.Equal(t, c.base, res.Symbols["start"], "symbols: %+v", res.Symbols)
			require.NotZero(t, res.Symbols["msg"], "symbols: %+v", res.Symbols)

			riscvCase := strings.Contains(c.src, "riscv")
			loongCase := strings.Contains(c.src, "loong")
			type slot struct {
				addr uint64
				line string
				size int // length of the source instruction
			}
			var instrs []slot
			switch {
			case loongCase:
				insts, err := loong64.Parse(c.base)(parsecbytes.Buffer(bin))
				require.NoError(t, err)
				for _, in := range insts {
					line := objdump.StripComments(
						objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())),
					)
					if line != "" && line != "<unknown>" {
						instrs = append(instrs, slot{
							in.Addr(),
							line,
							in.Len(),
						})
					}
				}
			case riscvCase:
				insts, err := riscv.Parse(c.base)(parsecbytes.Buffer(bin))
				require.NoError(t, err)
				for _, in := range insts {
					line := objdump.StripComments(
						objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())),
					)
					if line != "" && line != "<unknown>" {
						instrs = append(instrs, slot{
							in.Addr(),
							line,
							in.Len(),
						})
					}
				}
			default:
				insts, err := arm64.Parse(c.base)(parsecbytes.Buffer(bin))
				require.NoError(t, err)
				for _, in := range insts {
					if _, ok := in.(arm64.Generic); ok {
						continue // decode-only (armISA tail): the text does not assemble
					}

					line := objdump.StripComments(
						objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())),
					)
					if line != "" {
						instrs = append(instrs, slot{
							in.Addr(),
							line,
							4,
						})
					}
				}
			}

			for _, st := range instrs {
				off := st.addr - c.base
				if off+uint64(st.size) > uint64(len(bin)) {
					continue
				}

				want := bin[off : off+uint64(st.size)] // the whole source instruction
				res2, errs2 := c.assemble(st.line, st.addr)
				require.Empty(t, errs2, "addr %#x: %q", st.addr, st.line)
				got := res2.Sections[0].Data
				if bytes.Equal(got, want) {
					continue // byte-exact
				}

				// lengths may differ: a symbolic target of the source is not
				// compressed, while a numeric one from the round-trip is;
				// these are decode-equivalent forms of the same text
				var equiv func(b []byte, addr uint64) string
				switch {
				case riscvCase:
					equiv = riscvInstrText
				case loongCase:
					equiv = loongInstrText
				default:
					equiv = instrText
				}

				if equiv(got, st.addr) == equiv(want, st.addr) {
					continue
				}

				require.Equal(t, want, got, "addr %#x: %q", st.addr, st.line)
			}
		})
	}
}

// TestHelloVM - running the bare-metal demo in qemu virtual machines
// (arm64, riscv64 and loong64 virt). The gate is the "hello world" line
// on the serial console: machines with a bare-metal poweroff exit by
// themselves; the idling one (loong64) is stopped as soon as the line is
// observed - the deadline only bounds a hang.
func TestHelloVM(t *testing.T) {
	src, err := os.ReadFile("examples/hello-asm/hello-riscv.s")
	if err != nil {
		t.Skipf("example source not available: %v", err)
	}

	res, errs := pseudo.Assemble(string(src), 0x80000000)
	require.Empty(t, errs)
	riscvELF := mustWriteELF(t,
		file.EM_RISCV,
		0x80000000,
		res.Symbols["start"],
		fileSections(res.Sections),
	)
	elfPath := filepath.Join(t.TempDir(), "hello-riscv.elf")
	err = os.WriteFile(elfPath, riscvELF, 0o755)
	require.NoError(t, err)

	srcL, err := os.ReadFile("examples/hello-asm/hello-loongarch.s")
	if err != nil {
		t.Skipf("example source not available: %v", err)
	}

	resL, errsL := lpseudo.Assemble(string(srcL), 0x1c000000)
	require.Empty(t, errsL, "assemble loong")
	loongBin := filepath.Join(t.TempDir(), "hello-loong.bin")
	loongImage := make([]byte, 0)
	for _, sec := range resL.Sections {
		loongImage = append(loongImage, sec.Data...)
	}

	err = os.WriteFile(loongBin, loongImage, 0o755)
	require.NoError(t, err)

	src2, err := os.ReadFile("examples/hello-asm/hello-arm-vm.s")
	require.NoError(t, err)
	res2, errs2 := alias.Assemble(string(src2), 0x40100000)
	require.Empty(t, errs2, "assemble arm")
	armELF := mustWriteELF(t,
		file.EM_AARCH64,
		0x40100000,
		res2.Symbols["start"],
		fileSections(res2.Sections),
	)
	armPath := filepath.Join(t.TempDir(), "hello-arm.elf")
	err = os.WriteFile(armPath, armELF, 0o755)
	require.NoError(t, err)

	cases := []struct {
		name string
		qemu string
		args []string
	}{
		{
			"riscv64-virt",
			"qemu-system-riscv64",
			[]string{"-machine", "virt", "-bios", "none", "-nographic", "-kernel", elfPath},
		},
		{
			"arm64-virt",
			"qemu-system-aarch64",
			[]string{"-machine", "virt", "-cpu", "cortex-a53", "-nographic",
				"-device", "loader,file=" + armPath + ",cpu-num=0"},
		},
		{
			// LoongArch resets into the flash (0x1c000000), not RAM; the
			// machine has no bare-metal poweroff, so the program idles after
			// printing - the test stops it once the line is out.
			"loong64-virt",
			"qemu-system-loongarch64",
			[]string{"-machine", "virt", "-nographic",
				"-device", "loader,file=" + loongBin + ",addr=0x1c000000,cpu-num=0"},
		},
	}

	const (
		helloLine = "hello world"

		// helloVMDelay is the hang bound, not the expected runtime: the
		// passing boots finish in a couple of seconds.
		helloVMDelay = 15 * time.Second
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := exec.LookPath(c.qemu); err != nil {
				t.Skipf("%s not on PATH (brew install qemu)", c.qemu)
			}

			ctx, cancel := context.WithTimeout(context.Background(), helloVMDelay)
			defer cancel()

			console, serialOut, err := os.Pipe()
			require.NoError(t, err)

			cmd := exec.CommandContext(ctx, c.qemu, c.args...)
			cmd.Stdout = serialOut
			cmd.Stderr = serialOut
			require.NoError(t, cmd.Start())

			// the child owns the only write end now
			require.NoError(t, serialOut.Close())

			serial := &bytes.Buffer{}
			saw := make(chan struct{})
			scanned := make(chan struct{})
			go func() {
				defer close(scanned)

				scanner := bufio.NewScanner(console)
				for scanner.Scan() {
					line := scanner.Bytes()
					serial.Write(line)
					serial.WriteByte('\n')

					if bytes.Contains(line, []byte(helloLine)) {
						close(saw)
					}
				}
			}()

			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()

			var exitErr error
			select {
			case <-saw:
				// the line is out - stop the idling machine
				cancel()
				<-done
			case exitErr = <-done:
			}

			<-scanned

			// a qemu timeout is not an error: by the deadline the output has
			// already been produced (or not - the Contains below decides)
			if exitErr != nil && ctx.Err() != context.DeadlineExceeded {
				require.NoError(t, exitErr, "%s: %v\n%s", c.qemu, exitErr, serial.String())
			}

			require.Contains(t, serial.String(), helloLine)
		})
	}
}

// mustWriteELF - WriteELF for tests: a file-build error fails the test.
func mustWriteELF(
	t *testing.T,
	machine uint16,
	base, entry uint64,
	sections []file.Section,
) []byte {
	t.Helper()
	blob, err := file.WriteELF(machine, 0, base, entry, sections)
	require.NoError(t, err, "WriteELF")
	return blob
}
