package tests

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	rpseudo "github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/debug"
	darm64 "github.com/okneniz/assembly/debug/arm64"
	dloong "github.com/okneniz/assembly/debug/loong64"
	"github.com/okneniz/assembly/debug/qemu"
	driscv "github.com/okneniz/assembly/debug/riscv"
	"github.com/okneniz/assembly/debug/session"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/prog"
	aprog "github.com/okneniz/assembly/prog/arm64"
	rprog "github.com/okneniz/assembly/prog/riscv"
)

// callerPos - a caller-based position resolver for the prog chains:
// invoked inside a chain method, two frames up is the calling code
// (this file).
func callerPos() prog.Pos {
	_, file, line, _ := runtime.Caller(2)
	return prog.NewPos(file, line)
}

// asmDebugLines converts the asm line map into the session's normalized
// form (the one source file of the run).
func asmDebugLines(res *asm.Result, file string) []session.Line {
	out := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		out = append(out, session.NewLine(file, int(e.Line), e.Addr, e.Size))
	}

	return out
}

// readFile reads a test fixture (skip when absent).
func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("fixture not available: %v", err)
	}

	return string(data)
}

// startDebugged boots an ELF image under the arch's qemu executor and
// binds a session to it (the harness of the debugging gates).
func startDebugged(
	t *testing.T,
	tgt debug.Target,
	img []byte,
	syms map[string]uint64,
	lines []session.Line,
) *session.Session {
	t.Helper()

	if _, err := exec.LookPath(tgt.QemuBinary()); err != nil {
		t.Skipf("%s not on PATH (brew install qemu)", tgt.QemuBinary())
	}

	machine, err := qemu.Start(t.Context(), tgt, img, qemu.NewOptions())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, machine.Close()) })

	s, err := session.New(machine.Conn(), tgt, syms, lines)
	require.NoError(t, err)
	return s
}

// stepPC steps once and returns the pc after it.
func stepPC(t *testing.T, s *session.Session) uint64 {
	t.Helper()

	_, err := s.Step()
	require.NoError(t, err)

	pc, err := s.PC()
	require.NoError(t, err)
	return pc
}

// TestDebugSource - the .s path over all three arches: assemble the
// bare-metal hello, boot it frozen, and drive it: a breakpoint at the
// entry (the machines may sit in their reset trampolines before it),
// a step of the right length, a run to a label, memory matching the
// assembled bytes, and the line map pointing back into the source.
func TestDebugSource(t *testing.T) {
	cases := []struct {
		name     string
		src      string
		base     uint64
		tgt      debug.Target
		assemble func(string, uint64) (*asm.Result, []asm.AsmError)
		machine  uint16
		flags    uint32
		raw      bool // a raw blob image (the loong loader takes no ELF)
		// runTo resolves the address of the stop past the entry: a named
		// label where the source has one, the idle loop's own address
		// where only numeric locals live
		runTo func(res *asm.Result) uint64
	}{
		{
			name:     "arm64",
			src:      "examples/hello-asm/hello-arm-vm.s",
			base:     0x40100000,
			tgt:      darm64.NewTarget(),
			assemble: alias.Assemble,
			machine:  file.EM_AARCH64,
			runTo: func(res *asm.Result) uint64 {
				return res.Symbols["done"]
			},
		},
		{
			name:     "riscv64",
			src:      "examples/hello-asm/hello-riscv.s",
			base:     0x80000000,
			tgt:      driscv.NewTarget(),
			assemble: rpseudo.Assemble,
			machine:  file.EM_RISCV,
			runTo: func(res *asm.Result) uint64 {
				return res.Symbols["done"]
			},
		},
		{
			name:     "loong64",
			src:      "examples/hello-asm/hello-loongarch.s",
			base:     0x1c000000,
			tgt:      dloong.NewTarget(),
			assemble: lpseudo.Assemble,
			machine:  file.EM_LOONGARCH,
			flags:    0x43,
			raw:      true,
			// the idle loop "2: b 2b" is the last .text instruction
			// (numeric locals are not in Symbols)
			runTo: func(res *asm.Result) uint64 {
				text := res.Sections[0]
				return text.Addr + uint64(len(text.Data)) - 4
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := readFile(t, c.src)
			res, errs := c.assemble(src, c.base)
			require.Empty(t, errs, "assemble")

			img := []byte{}
			if c.raw {
				for _, sec := range res.Sections {
					img = append(img, sec.Data...)
				}
			} else {
				var err error
				img, err = file.WriteELF(
					c.machine,
					c.flags,
					c.base,
					res.Symbols["start"],
					fileSections(res.Sections),
				)
				require.NoError(t, err)
			}

			s := startDebugged(t, c.tgt, img, res.Symbols, asmDebugLines(res, c.src))

			// to the entry: the machine may sit in its reset trampoline
			start := res.Symbols["start"]
			addr, err := s.BreakAt("start")
			require.NoError(t, err)
			require.Equal(t, start, addr)

			_, err = s.Continue()
			require.NoError(t, err, "continue to start")

			pc, err := s.PC()
			require.NoError(t, err)
			require.Equal(t, start, pc)

			// the memory under the entry is the assembled code, and a
			// step moves by the target's instruction length
			code, err := s.Read(start, 8)
			require.NoError(t, err)
			require.Equal(t, res.Sections[0].Data[:8], code)
			require.Equal(t, start+uint64(c.tgt.InstrLen(code)), stepPC(t, s))

			// a run to the stop past the entry
			label := c.runTo(res)
			_, err = s.BreakAt(fmt.Sprintf("%#x", label))
			require.NoError(t, err)

			_, err = s.Continue()
			require.NoError(t, err, "continue to %#x", label)

			pc, err = s.PC()
			require.NoError(t, err)
			require.Equal(t, label, pc)

			// the line map points back into the .s source
			line, ok := s.LineAt(pc)
			require.True(t, ok)
			require.Equal(t, c.src, line.File)
			require.Positive(t, line.Line)

			// the register dump serves the core set with pc in it
			regs, err := s.Regs()
			require.NoError(t, err)
			pcSeen := false
			for _, r := range regs {
				if r.Name == "pc" {
					pcSeen = true
					require.Equal(t, label, r.Value)
				}
			}

			require.True(t, pcSeen, "pc in the register dump")
		})
	}
}

// TestDebugProgArm64 - the prog path: a Go-written program debugged with
// its line map pointing at the Go positions (an injected constant
// resolver for the first line, an injected caller resolver for the
// rest), the register state after a step matching the mov operand.
func TestDebugProgArm64(t *testing.T) {
	p := aprog.New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 42) }).
		Label("start").
		Mov(aprog.X0, 0x41).
		WithPos(callerPos).
		Label("loop").
		B("loop").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0x40100000)
	require.Empty(t, res.Errs)
	require.Len(t, res.Code, 8)

	lines := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		lines = append(lines, session.NewLine(e.Pos.File, e.Pos.Line, e.Addr, e.Size))
	}

	secs := []file.Section{
		*file.NewSection(".text", "", 0x40100000, 0, 8, res.Code),
	}
	img, err := file.WriteELF(
		file.EM_AARCH64,
		0,
		0x40100000,
		res.Syms["start"],
		secs,
	)
	require.NoError(t, err)

	s := startDebugged(t, darm64.NewTarget(), img, res.Syms, lines)

	// the arm64 loader starts cpu 0 at the entry directly
	pc, err := s.PC()
	require.NoError(t, err)
	require.Equal(t, res.Syms["start"], pc)

	// the injected resolver reported the synthetic position of the mov
	line, ok := s.LineAt(pc)
	require.True(t, ok)
	require.Equal(t, "synthetic.go", line.File)
	require.Equal(t, 42, line.Line)

	// after the step: pc at the loop label, x0 holding the operand
	require.Equal(t, res.Syms["loop"], stepPC(t, s))

	regs, err := s.Regs()
	require.NoError(t, err)
	x0 := uint64(0)
	for _, r := range regs {
		if r.Name == "x0" {
			x0 = r.Value
		}
	}

	require.Equal(t, uint64(0x41), x0)

	// the loop line comes from the injected caller resolver: this test file
	line, ok = s.LineAt(res.Syms["loop"])
	require.True(t, ok)
	require.Contains(t, line.File, "tests/debug_test.go")
	require.Positive(t, line.Line)
}

// TestDebugProgRiscv - the prog path on riscv: a Go-written program
// debugged with its line map pointing at the Go positions (an injected
// constant resolver for the first line, an injected caller resolver for
// the rest); the machine starts in its reset trampoline, so the first
// stop is a breakpoint at the entry. The register state after a step
// matches the addi operand.
func TestDebugProgRiscv(t *testing.T) {
	p := rprog.New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 42) }).
		Label("start").
		Addi(rprog.A0, rprog.Zero, 0x41).
		WithPos(callerPos).
		Label("loop").
		J("loop").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0x80000000)
	require.Empty(t, res.Errs)
	require.Len(t, res.Code, 8)

	lines := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		lines = append(lines, session.NewLine(e.Pos.File, e.Pos.Line, e.Addr, e.Size))
	}

	secs := []file.Section{
		*file.NewSection(".text", "", 0x80000000, 0, 8, res.Code),
	}
	img, err := file.WriteELF(
		file.EM_RISCV,
		0,
		0x80000000,
		res.Syms["start"],
		secs,
	)
	require.NoError(t, err)

	s := startDebugged(t, driscv.NewTarget(), img, res.Syms, lines)

	// to the entry: the machine starts in its reset trampoline
	_, err = s.BreakAt("start")
	require.NoError(t, err)

	_, err = s.Continue()
	require.NoError(t, err, "continue to start")

	pc, err := s.PC()
	require.NoError(t, err)
	require.Equal(t, res.Syms["start"], pc)

	// the injected resolver reported the synthetic position of the addi
	line, ok := s.LineAt(pc)
	require.True(t, ok)
	require.Equal(t, "synthetic.go", line.File)
	require.Equal(t, 42, line.Line)

	// after the step: pc at the loop label, a0 holding the operand
	require.Equal(t, res.Syms["loop"], stepPC(t, s))

	regs, err := s.Regs()
	require.NoError(t, err)
	a0 := uint64(0)
	for _, r := range regs {
		if r.Name == "x10" {
			a0 = r.Value
		}
	}

	require.Equal(t, uint64(0x41), a0)

	// the loop line comes from the injected caller resolver: this test file
	line, ok = s.LineAt(res.Syms["loop"])
	require.True(t, ok)
	require.Contains(t, line.File, "tests/debug_test.go")
	require.Positive(t, line.Line)
}
