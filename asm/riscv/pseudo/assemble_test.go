package pseudo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"

	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

func assembleOne(t *testing.T, src string, addr uint64) []byte {
	t.Helper()
	res, errs := asm.Assemble(src, addr, NewASMBackend())
	require.Empty(t, errs, "assemble %q: %v", src, errs)
	require.NotEmpty(t, res.Sections, "assemble %q: no sections", src)
	return res.Sections[0].Data
}

func wordOf(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}

	return binary.LittleEndian.Uint32(b)
}

func instrTextSource(inst arch.Instr) string {
	return inst.ObjDump(disasm.DefaultViewCtx())
}

func TestAssemblePseudo(t *testing.T) {
	words := []struct {
		src  string
		word uint32
	}{
		{
			"li a0, 0x42",
			0x04200513,
		}, // -> addi a0, zero, 0x42 (not compressed: rs1=zero != rd)
		{
			"mv a0, a1",
			0,
		}, // -> c.mv (add form), 2 bytes
		{
			"not a0, a1",
			0xfff5c513,
		}, // → xori a0, a1, -1
		{
			"neg a0, a1",
			0x40b00533,
		}, // → sub a0, zero, a1
		{
			"snez a0, a1",
			0x00b03533,
		}, // → sltu a0, zero, a1
		{
			"seqz a0, a1",
			0x015b513,
		}, // → sltiu a0, a1, 1
		{
			"sext.w a0, a1",
			0x0005851b,
		}, // → addiw a0, a1, 0
	}
	for _, c := range words {
		got := assembleOne(t, c.src, 0)
		require.Equal(t, c.word, wordOf(got), "case %q = % x", c.src, got)
	}

	// Compressible pseudo-instructions: nop -> c.nop, li small -> c.li, ret -> c.jr.
	structural := []struct {
		src  string
		want string // the expected expanded text from our decoder
	}{
		{
			"nop",
			"nop",
		},
		{
			"li a0, 0x1",
			"li a0, 0x1",
		},
		{
			"ret",
			"ret",
		},
		{
			"csrr a0, fflags",
			"frflags a0",
		},
		{
			"csrw fflags, a0",
			"fsflags a0",
		},
		{
			"rdcycle a0",
			"rdcycle a0",
		},
	}
	for _, c := range structural {
		got := assembleOne(t, c.src, 0)
		insts, err := arch.Parse(0)(parsecbytes.Buffer(got))
		require.NoError(t, err)
		require.NotEmpty(t, insts, "case %q: nothing decoded", c.src)
		back := objdump.Normalize(instrTextSource(insts[0]))
		require.Equal(t, c.want, back, "case %q → % x → %q", c.src, got, back)
	}
}

func TestAssembleLiLarge(t *testing.T) {
	// li a0, 0x12345678 → lui a0, 0x12345 + addi a0, a0, 0x678
	got := assembleOne(t, "li a0, 0x12345678", 0)
	require.Len(t, got, 8, "li large")
	hi := binary.LittleEndian.Uint32(got[0:4])
	lo := binary.LittleEndian.Uint32(got[4:8])
	require.Equal(t, uint32(0x12345537), hi)
	require.Equal(t, uint32(0x67850513), lo)
}

func TestLabelsAndCalls(t *testing.T) {
	src := `
start:
  li a0, 0x123
  la a1, start
  call func
  j start
func:
  ret
`
	res, errs := asm.Assemble(src, 0x1000, NewASMBackend())
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data
	// li a0, 0x123: imm in [32..2047] -> addi without c.li -> 4; la =
	// auipc+addi = 8; call = auipc+jalr = 8; j start -> c.j = 2; ret ->
	// c.jr = 2. func @ +24.
	require.Equal(t, uint64(0x1000+24), res.Symbols["func"], "func")
	require.Len(t, d, 26, "total: % x", d)
	// la a1, start @0x1004: rel = -4 → hi=0, lo=-4: auipc a1, 0 + addi a1, a1, -4
	laHi := binary.LittleEndian.Uint32(d[4:8])
	laLo := binary.LittleEndian.Uint32(d[8:12])
	require.Equal(t, uint32(0x00000597), laHi, "la") // auipc a1,0
	require.Equal(t, uint32(0xffc58593), laLo, "la") // addi a1,a1,-4
}

// TestRiscvNumericLabels tests the GAS numeric local labels: beq/j via
// 1:/1b/1f; numeric labels never get into Symbols.
func TestRiscvNumericLabels(t *testing.T) {
	src := `
1:
  beq a0, a1, 1f
  j 1b
1:
  j 1b
`
	res, errs := asm.Assemble(src, 0x1000, NewASMBackend())
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data
	// branches to symbols (including numeric labels) are not compressed
	// into RVC - deterministic sizes between passes; beq 4 + two jal, 4
	// each
	require.Len(t, d, 12, "total: % x", d)
	// beq a0, a1, +8 @0x1000 → 1: @0x1008
	require.Equal(t, uint32(0x00B50463), binary.LittleEndian.Uint32(d[0:4]), "beq +8")
	// jal x0, -4 @0x1004 → 1: @0x1000
	require.Equal(t, uint32(0xFFDFF06F), binary.LittleEndian.Uint32(d[4:8]), "jal -4")
	// jal x0, 0 @0x1008 -> 1: @0x1008 (the label before the instruction)
	require.Equal(t, uint32(0x0000006F), binary.LittleEndian.Uint32(d[8:12]), "jal 0")
	require.NotContains(t, res.Symbols, "1", "numeric label must not be a symbol")
}

// TestOptionNorvc checks that .option norvc forbids RVC auto-compression
// (literal close targets stay uncompressed); rvc/pop restore it;
// out-of-model values (norelax/pic) are harmlessly ignored. Symbolic
// targets are always uncompressed - hence the literal targets here.
func TestOptionNorvc(t *testing.T) {
	src := `
  j 0x1000
.option norvc
  j 0x1000
.option rvc
.option push
.option norvc
  j 0x1000
.option pop
  j 0x1000
.option norelax
.option pic
  nop
`
	res, errs := asm.Assemble(src, 0x1000, NewASMBackend())
	require.Empty(t, errs, "errors: %v", errs)
	d := res.Sections[0].Data
	// c.j @0x1000 + jal @0x1002 + jal @0x1006 + c.j @0x100A + c.nop = 14
	require.Len(t, d, 14, "total: % x", d)
	require.Equal(t, uint16(0xA001), binary.LittleEndian.Uint16(d[0:2]), "c.j 0")
	// j 0x1000 @0x1002 → off -2, norvc → jal
	require.Equal(t, uint32(0xFFFFF06F), binary.LittleEndian.Uint32(d[2:6]), "norvc: jal -2")
	// j 0x1000 @0x1006 → off -6, push/norvc → jal
	require.Equal(t, uint32(0xFFBFF06F), binary.LittleEndian.Uint32(d[6:10]), "push/norvc: jal -6")
	// pop restored compression: c.j -10
	require.Equal(t, uint16(0xBFDD), binary.LittleEndian.Uint16(d[10:12]), "pop: c.j -10")
	// norelax/pic are outside the model, they do not touch compression
	require.Equal(t, uint16(0x0001), binary.LittleEndian.Uint16(d[12:14]), "nop → c.nop")
}

func TestRVCPseudoCompression(t *testing.T) {
	cases := []struct {
		src    string
		addr   uint64
		twoLen bool // a 16-bit form expected
	}{
		{
			"li a0, 1",
			0,
			true,
		}, // c.li
		{
			"mv a0, a1",
			0,
			true,
		}, // -> add form -> c.mv (as GNU as)
		{
			"beqz a0, 0x10",
			0,
			true,
		}, // c.beqz (a0 ∈ x8-x15)
		{
			"beqz t0, 0x10",
			0,
			false,
		}, // t0 = x5 ∉ x8-x15
		{
			"beqz s0, 0x10",
			0,
			true,
		}, // c.beqz
		{
			"bltz s0, 0x10",
			0,
			false,
		}, // -> blt s0, zero -> no c form (rs2=zero != rs1') -> 32 bits
		{
			"j 0x10",
			0,
			true,
		}, // c.j
		{
			"j 0x1000",
			0,
			false,
		}, // within +/-1MB, outside +/-2KB -> 32-bit jal
	}
	for _, c := range cases {
		got := assembleOne(t, c.src, c.addr)
		require.Equal(
			t,
			c.twoLen,
			len(got) == 2,
			"case %q @ %#x = %d bytes (% x), compressed want %v",
			c.src,
			c.addr,
			len(got),
			got,
			c.twoLen,
		)
	}
}

// TestRoundTripExample is a byte-exact round-trip test of the test
// binary: decode -> Text -> assemble -> the same bytes. The decoder
// text contains pseudo-forms (nop/li/mv/ret/...) - the round-trip
// goes through the pseudo layer.
func TestRoundTripExample(t *testing.T) {
	ff, err := file.Detect("../../../tests/examples/hello-riscv/hello-arch.o")
	if err != nil {
		t.Skipf("example not available: %v", err)
	}

	ts, err := ff.CodeSection()
	if err != nil {
		t.Skipf("example not available: %v", err)
	}

	insts, err := arch.Parse(ts.Addr)(parsecbytes.Buffer(ts.Data))
	require.NoError(t, err)
	matched, rmLossy, mismatched := 0, 0, 0
	var failures []string
	for _, in := range insts {
		src := instrTextSource(in)
		if src == "" || src == "<unknown>" {
			continue
		}

		off := in.Addr() - ts.Addr
		want := ts.Data[off : off+uint64(in.Len())]
		res, errs := asm.Assemble(src, in.Addr(), NewASMBackend())
		if len(errs) != 0 {
			mismatched++
			failures = append(failures, fmt.Sprintf("addr %#x: %q: %v", in.Addr(), src, errs))
			continue
		}

		got := res.Sections[0].Data
		switch {
		case bytes.Equal(got, want):
			matched++
		case len(got) == 4 && len(want) == 4 &&
			binary.LittleEndian.Uint32(got)&^0x7000 == binary.LittleEndian.Uint32(want)&^0x7000:
			// the difference is only in funct3 - the FP rounding mode,
			// which the objdump-like text does not carry (a known lossy
			// class)
			rmLossy++
		default:
			mismatched++
			if len(failures) < 5 {
				failures = append(
					failures,
					fmt.Sprintf("addr %#x: %q\n  got  % x\n  want % x", in.Addr(), src, got, want),
				)
			}
		}
	}

	require.Empty(t, failures, "hard round-trip mismatches (%d)", mismatched)
	require.NotZero(t, matched, "nothing round-tripped")
	t.Logf("round-trip: %d/%d byte-exact + %d rm-only (FP rounding, lossy text), %d hard",
		matched, matched+rmLossy+mismatched, rmLossy, mismatched)
	require.Zero(t, mismatched, "hard mismatches")
}

// TestRoundTripSynthetic is a synthetic corpus: random operands across
// all supported mnemonics; source -> bytes -> decode -> Text ->
// assemble -> the same bytes (assembler/decoder/formatter agree; the
// decoder text contains pseudo-forms - assembly goes through the
// pseudo layer).
func TestRoundTripSynthetic(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	regs := []string{"zero", "ra", "sp", "tp", "t0", "t1", "t2", "s0", "s1",
		"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "s2", "t3", "x28", "x31"}
	fregs := []string{"ft0", "fs0", "fa0", "fa1", "fs5", "ft8", "f31", "f2"}
	rreg := func() string {
		return regs[rnd.Intn(len(regs))]
	}
	rfreg := func() string {
		return fregs[rnd.Intn(len(fregs))]
	}
	rimm := func(lo, hi int64) int64 {
		return lo + rnd.Int63n(hi-lo+1)
	}

	type gen struct {
		format string
		make   func(addr uint64) string
	}
	gens := []gen{
		{
			"%s %s, %s, %d",
			func(uint64) string {
				return fmt.Sprintf("add %s, %s, %s", rreg(), rreg(), rreg())
			},
		},
		{
			"addi",
			func(addr uint64) string {
				return fmt.Sprintf("addi %s, %s, %d", rreg(), rreg(), rimm(-2048, 2047))
			},
		},
		{
			"ori",
			func(addr uint64) string {
				return fmt.Sprintf("ori %s, %s, %d", rreg(), rreg(), rimm(-2048, 2047))
			},
		},
		{
			"slli",
			func(addr uint64) string {
				return fmt.Sprintf("slli %s, %s, %d", rreg(), rreg(), rimm(1, 63))
			},
		},
		{
			"srliw",
			func(addr uint64) string {
				return fmt.Sprintf("srliw %s, %s, %d", rreg(), rreg(), rimm(1, 31))
			},
		},
		{
			"lui",
			func(addr uint64) string {
				return fmt.Sprintf("lui %s, %#x", rreg(), rimm(0, 0xfffff))
			},
		},
		{
			"lw",
			func(addr uint64) string {
				return fmt.Sprintf("lw %s, %d(%s)", rreg(), rimm(-2048, 2047), rreg())
			},
		},
		{
			"fsd",
			func(addr uint64) string {
				return fmt.Sprintf("fsd %s, %d(%s)", rfreg(), rimm(-2048, 2047), rreg())
			},
		},
		{
			"beq",
			func(addr uint64) string {
				target := int64(addr) + rimm(-512, 508)&^1
				return fmt.Sprintf("beq %s, %s, %#x", rreg(), rreg(), target)
			},
		},
		{
			"bge",
			func(addr uint64) string {
				target := int64(addr) + rimm(-512, 508)&^1
				return fmt.Sprintf("bge %s, %s, %#x", rreg(), rreg(), target)
			},
		},
		{
			"jal",
			func(addr uint64) string {
				target := int64(addr) + rimm(-4096, 4094)&^1
				return fmt.Sprintf("jal %s, %#x", rreg(), target)
			},
		},
		{
			"fadd.s",
			func(addr uint64) string {
				return fmt.Sprintf("fadd.s %s, %s, %s", rfreg(), rfreg(), rfreg())
			},
		},
		{
			"fmadd.d",
			func(addr uint64) string {
				return fmt.Sprintf("fmadd.d %s, %s, %s, %s", rfreg(), rfreg(), rfreg(), rfreg())
			},
		},
		{
			"amoadd.w",
			func(addr uint64) string {
				return fmt.Sprintf("amoadd.w %s, %s, (%s)", rreg(), rreg(), rreg())
			},
		},
		{
			"mulw",
			func(addr uint64) string {
				return fmt.Sprintf("mulw %s, %s, %s", rreg(), rreg(), rreg())
			},
		},
		{
			"csrrw",
			func(addr uint64) string {
				return fmt.Sprintf("csrrw %s, %#x, %s", rreg(), rimm(0, 0xfff), rreg())
			},
		},
		{
			"csrrs",
			func(addr uint64) string {
				return fmt.Sprintf("csrrs %s, sstatus, %s", rreg(), rreg())
			},
		},
	}

	matched, mismatched := 0, 0
	var failures []string
	for range 400 {
		g := gens[rnd.Intn(len(gens))]
		addr := uint64(0x10000 + rnd.Intn(64)*2)
		src := g.make(addr)
		res, errs := asm.Assemble(src, addr, NewASMBackend())
		if len(errs) != 0 {
			mismatched++
			failures = append(failures, fmt.Sprintf("assemble %q: %v", src, errs))
			continue
		}

		want := res.Sections[0].Data
		insts, err := arch.Parse(addr)(parsecbytes.Buffer(want))
		require.NoError(t, err)
		for _, in := range insts {
			src2 := instrTextSource(in)
			if src2 == "" || src2 == "<unknown>" {
				mismatched++
				failures = append(failures, fmt.Sprintf("%q → % x: decode failed", src, want))
				continue
			}

			res2, errs2 := asm.Assemble(src2, in.Addr(), NewASMBackend())
			if len(errs2) != 0 {
				mismatched++
				failures = append(failures, fmt.Sprintf("re-assemble %q: %v", src2, errs2))
				continue
			}

			got := res2.Sections[0].Data
			want2 := want[in.Addr()-addr : in.Addr()-addr+uint64(in.Len())]
			if bytes.Equal(got, want2) {
				matched++
			} else {
				mismatched++
				if len(failures) < 5 {
					failures = append(
						failures,
						fmt.Sprintf(
							"round-trip %q (→ %q)\n  got  % x\n  want % x",
							src,
							src2,
							got,
							want2,
						),
					)
				}
			}
		}
	}

	require.NotZero(t, matched, "nothing generated")
	t.Logf("synthetic round-trip: %d/%d byte-exact", matched, matched+mismatched)
	require.Empty(t, failures, "mismatches (%d)", mismatched)
}
