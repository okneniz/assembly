package arm64

// The producer's path through the output: a program built with the arch
// Builder and test-local deferred records (an adr, an la pair - unit
// ships no arch macros by design, a producer constructs its own),
// byte-identical to the prog chain world (the golden .s pipeline) and to
// the prog-built Mach-O image.

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	prog "github.com/okneniz/assembly/prog/arm64"
	"github.com/okneniz/assembly/unit"
)

const (
	msg          = "hello world\n"
	trapMach     = 0x80
	sysClassUnix = 0x2000000
	sysWrite     = 4
	sysExit      = 1
)

func TestHelloGolden(t *testing.T) {
	// the flat resolve of a text-only program: base 0, like the prog
	// chain's flat Assemble
	f := helloUnit(arch.Builder{}).Resolve(func(text, data, dataMem int) (uint64, uint64) {
		require.Equal(t, 52, text) // 10 instructions + 12 data bytes
		return 0, 0
	})
	require.Empty(t, f.Errs)

	code, codeErr := f.EncodeText()
	require.NoError(t, codeErr)

	// the same golden as the prog chain's: the .s pipeline output
	golden, err := os.ReadFile("../../prog/arm64/testdata/hello-macos.bin")
	require.NoError(t, err)
	require.Equal(t, golden, code)

	require.Equal(t, uint64(0), f.Syms["start"])
	require.Equal(t, uint64(40), f.Syms["msg"])
}

func TestMachoParityWithProg(t *testing.T) {
	// the same program in the two worlds: byte-identical images - the
	// unit output is a drop-in for the chain's Mach-O end
	progBin, progErrs := streamsProg().Build()
	require.Empty(t, progErrs)

	progImg, err := progBin.MachO("start")
	require.NoError(t, err)

	unitImg, err := Macho(streamsUnit(arch.Builder{}), "start")
	require.NoError(t, err)

	require.Equal(t, progImg.Bytes(), unitImg.Bytes())
}

// mustX pre-mints a register (0..30 is valid - the panic is unreachable).
func mustX(n int) arch.Reg {
	r, err := arch.X(n)
	if err != nil {
		panic(err) // unreachable: 0..30 is a valid register number
	}

	return r
}

// newAdrSym - the test's label-directed adr as the canonical record: the
// arch formula is an injected function, the resolve driver is the unit's
// (this is the dcc shape - no custom Resolve types).
func newAdrSym(b arch.Builder, rd arch.Reg, label string) unit.Sym {
	return unit.NewBranch("adr", label, 4, func(t, pc uint64) (unit.Resolved, error) {
		return b.Adr(rd, int64(t)-int64(pc))
	})
}

// newLaSym - the test's label-directed la pair (adrp+add) the same way.
func newLaSym(b arch.Builder, rd arch.Reg, label string) unit.Sym {
	return unit.NewPair("la", label, 8, func(t, pc uint64) ([]unit.Resolved, error) {
		target, at := int64(t), int64(pc)
		adrp, err := b.Adrp(rd, (target&^0xFFF-at&^0xFFF)>>12)
		if err != nil {
			return nil, err
		}

		lo, err := b.Imm12(target & 0xFFF)
		if err != nil {
			return nil, err
		}

		add, err := b.AddImm(rd, rd, lo, arch.NoSh12)
		if err != nil {
			return nil, err
		}

		return []unit.Resolved{adrp, add}, nil
	})
}

// helloUnit - the Go counterpart of prog/arm64's hello-macos chain: the
// same ten instructions and the same string, built on the unit output.
func helloUnit(b arch.Builder) *unit.Unit {
	x0, x1, x2, x16 := mustX(0), mustX(1), mustX(2), mustX(16)

	u := unit.New()
	emit := func(i arch.Instr, err error) { u.Instr(unit.Pos{}, i, err) }
	u.Label("start")
	emit(movz(b, x0, 1, arch.Hw0))                 // write(fd=stdout, ...)
	u.Sym(unit.Pos{}, newAdrSym(b, x1, "msg"))     // ... buf - the string address
	emit(movz(b, x2, int64(len(msg)), arch.Hw0))   // ... len
	emit(movz(b, x16, sysClassUnix>>16, arch.Hw1)) // x16 = 0x2000000 | ...
	emit(movk(b, x16, sysWrite, arch.Hw0))         // ... 0x4 = write
	emit(svc(b, trapMach))
	emit(movz(b, x0, 0, arch.Hw0)) // exit(return code = 0)
	emit(movz(b, x16, sysClassUnix>>16, arch.Hw1))
	emit(movk(b, x16, sysExit, arch.Hw0))
	emit(svc(b, trapMach))
	u.Label("msg")
	u.Ascii(unit.Pos{}, msg)
	u.Entry("start")
	return u
}

// movz/movk/svc - the Builder calls of the chain twins, one instruction
// each (the immediate conversions the twins do inline).
func movz(b arch.Builder, rd arch.Reg, imm int64, hw arch.Hw) (arch.Instr, error) {
	v, err := b.Imm16(imm)
	if err != nil {
		return nil, err
	}

	return b.Movz(rd, v, hw)
}

func movk(b arch.Builder, rd arch.Reg, imm int64, hw arch.Hw) (arch.Instr, error) {
	v, err := b.Imm16(imm)
	if err != nil {
		return nil, err
	}

	return b.Movk(rd, v, hw)
}

func svc(b arch.Builder, imm int64) (arch.Instr, error) {
	v, err := b.Imm16(imm)
	if err != nil {
		return nil, err
	}

	return b.Svc(v), nil
}

// streamsUnit - the Go counterpart of prog/arm64's streams program: a
// data static through an la pair, a bss tail, the layout mode.
func streamsUnit(b arch.Builder) *unit.Unit {
	x0, x1, x16 := mustX(0), mustX(1), mustX(16)

	u := unit.New().Entry("start")
	emit := func(i arch.Instr, err error) { u.Instr(unit.Pos{}, i, err) }
	u.Label("start")
	u.Sym(unit.Pos{}, newLaSym(b, x0, "counter"))
	emit(mustLdr(b, x1, x0))
	emit(mustAddImm(b, x1, x1, 1))
	emit(mustStr(b, x1, x0))
	emit(mustLdr(b, x0, x0))
	emit(movz(b, x16, 0x200, arch.Hw1))
	emit(movk(b, x16, 1, arch.Hw0))
	emit(svc(b, trapMach))

	u.Data()
	u.Label("counter").Quad(unit.Pos{}, 7)
	u.Label("buf").Bss(unit.Pos{}, 16)
	return u
}

// mustLdr/mustAddImm/mustStr - the Builder calls of the load/store twins.
func mustLdr(b arch.Builder, rt, rn arch.Reg) (arch.Instr, error) {
	return b.Ldr(rt, rn, arch.Off(0))
}

func mustAddImm(b arch.Builder, rd, rn arch.Reg, imm int64) (arch.Instr, error) {
	v, err := b.Imm12(imm)
	if err != nil {
		return nil, err
	}

	return b.AddImm(rd, rn, v, arch.NoSh12)
}

func mustStr(b arch.Builder, rt, rn arch.Reg) (arch.Instr, error) {
	return b.Str(rt, rn, arch.Off(0))
}

// streamsProg - prog/arm64's streams program, rebuilt here for the
// image parity (the chain world).
func streamsProg() *prog.Program {
	x0, x1, x16 := prog.X0, prog.X1, prog.X16

	p := prog.New(unit.New()).Entry("start")
	p.Label("start")
	p.La(x0, "counter")
	p.Ldr(x1, x0, 0)
	p.AddImm(x1, x1, 1, arch.NoSh12)
	p.Str(x1, x0, 0)
	p.Ldr(x0, x0, 0)
	p.Movz(x16, 0x200, arch.Hw1)
	p.Movk(x16, 1, arch.Hw0)
	p.Svc(0x80)

	p.Data()
	p.Label("counter").Quad(7)
	p.Label("buf").Bss(16)
	return p
}
