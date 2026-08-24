package riscv

import (
	"encoding/binary"
	"testing"

	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/riscv"
	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/tests/cmd/objdump"
)

func assembleOne(t *testing.T, src string, addr uint64) []byte {
	t.Helper()
	res, errs := asm.Assemble(src, addr, New())
	require.Empty(t, errs, "assemble %q: %v", src, errs)
	require.NotEmpty(t, res.Sections, "assemble %q: no sections", src)
	return res.Sections[0].Data
}

func TestAssembleWords(t *testing.T) {
	// Exact words (verified manually/against the specification).
	words := []struct {
		src  string
		word uint32
	}{
		{
			"addi a0, a1, 0x42",
			0x04258513,
		},
		{
			"addi a0,a1,1",
			0x00158513,
		}, // no spaces (GNU)
		{
			"addi x10, x11, -1",
			0xfff58513,
		}, // xN notation, negative imm
		{
			"lui a0, 0x12345",
			0x12345537,
		}, // U: imm20<<12
		{
			"add a0, a1, a2",
			0x00c58533,
		}, // R: rs1=a1, rs2=a2, rd=a0
		{
			"auipc a0, 0x1",
			0x00001517,
		},
		{
			"jalr a0, 0x0(a1)",
			0x00058567,
		},
		{
			"slli a0, a1, 4",
			0x00459513,
		},
		{
			"ecall",
			0x00000073,
		},
	}
	for _, c := range words {
		got := assembleOne(t, c.src, 0)
		require.Len(t, got, 4, "case %q", c.src)
		require.Equal(
			t,
			c.word,
			binary.LittleEndian.Uint32(got),
			"case %q = % x",
			c.src,
			got,
		)
	}

	// Structural: assemble -> decode -> decoder text (with its aliases).
	texts := []struct {
		src, want string
	}{
		{
			"lw a0, 0x8(sp)",
			"lw a0, 0x8(sp)",
		},
		{
			"sw a0, 0x8(sp)",
			"sw a0, 0x8(sp)",
		},
		{
			"sd a0, 0x10(sp)",
			"sd a0, 0x10(sp)",
		},
		{
			"flw fa0, 0x4(a0)",
			"flw fa0, 0x4(a0)",
		},
		{
			"fsd fa0, 0x8(a0)",
			"fsd fa0, 0x8(a0)",
		},
		{
			"amoadd.w a0, a1, (a2)",
			"amoadd.w a0, a1, (a2)",
		},
		{
			"csrrw a0, fflags, a1",
			"csrrw a0, fflags, a1",
		},
		{
			"csrrs a0, sstatus, zero",
			"csrr a0, sstatus",
		},
		{
			"csrrwi a0, frm, 3",
			"csrrwi a0, frm, 0x3",
		},
		{
			"fadd.s fa0, fa1, fa2",
			"fadd.s fa0, fa1, fa2",
		},
		{
			"fmadd.d fa0, fa1, fa2, fa3",
			"fmadd.d fa0, fa1, fa2, fa3",
		},
		{
			"mulw a0, a1, a2",
			"mulw a0, a1, a2",
		},
		{
			"fence",
			"fence",
		},
		{
			"mret",
			"mret",
		},
		{
			"wfi",
			"wfi",
		},
		{
			"ebreak",
			"ebreak",
		},
	}
	for _, c := range texts {
		got := assembleOne(t, c.src, 0)
		insts, err := arch.Parse(0)(bytes.Buffer(got))
		require.NoError(t, err)
		require.NotEmpty(t, insts, "case %q: nothing decoded from % x", c.src, got)
		for _, in := range insts {
			back := objdump.Normalize(instrTextSource(in))
			require.Equal(
				t,
				objdump.Normalize(c.want),
				back,
				"case %q → % x → %q",
				c.src,
				got,
				back,
			)
		}
	}
}

func wordOf(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}

	return binary.LittleEndian.Uint32(b)
}

func TestAssembleBranch(t *testing.T) {
	// beq a0, a1, 0x1008 @ 0x1000 → off=8: imm[4:1]=0100
	got := assembleOne(t, "beq a0, a1, 0x1008", 0x1000)
	want := uint32(0x63) | 0b000<<12 | 10<<15 | 11<<20 | 0b0100<<8
	require.Equal(t, want, wordOf(got))
}

func TestRVCCompression(t *testing.T) {
	cases := []struct {
		src    string
		addr   uint64
		twoLen bool // a 16-bit form expected
	}{
		{
			"addi a0, a0, 1",
			0,
			true,
		}, // c.addi
		{
			"addi s0, sp, 0x10",
			0,
			true,
		}, // c.addi4spn
		{
			"addi sp, sp, 0x10",
			0,
			true,
		}, // c.addi16sp
		{
			"lw a0, 0x8(sp)",
			0,
			true,
		}, // c.lwsp
		{
			"lw s0, 0x8(s1)",
			0,
			true,
		}, // c.lw
		{
			"sw a0, 0x8(sp)",
			0,
			true,
		}, // c.swsp
		{
			"add s0, s0, s1",
			0,
			true,
		}, // c.add
		{
			"sub s0, s0, s1",
			0,
			true,
		}, // c.sub
		{
			"andi a0, a0, 0xf",
			0,
			true,
		}, // c.andi (a0 = x10 ∈ x8-x15!)
		{
			"slli a0, a0, 3",
			0,
			true,
		}, // c.slli
		{
			"srli s0, s0, 3",
			0,
			true,
		}, // c.srli
		{
			"addi a0, a1, 1",
			0,
			false,
		}, // rd != rs1 -> not compressed
		{
			"lui a0, 0x1",
			0,
			true,
		}, // c.lui
		{
			"lui sp, 0x1",
			0,
			false,
		}, // rd=sp -> not compressed
		{
			"addiw a0, a0, 1",
			0,
			true,
		}, // c.addiw
		{
			"ld a0, 0x8(sp)",
			0,
			true,
		}, // c.ldsp
		{
			"sd a0, 0x8(sp)",
			0,
			true,
		}, // c.sdsp
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

// instrTextSource returns the instruction's mnemonic+operands - its
// own ObjDump text.
func instrTextSource(inst arch.Instr) string {
	return inst.ObjDump(disasm.DefaultViewCtx())
}
