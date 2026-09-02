package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Jalr - jalr rd, off(rs1); pseudo: ret (zero,ra,0), jr (zero, imm=0 or
// offset), "jalr rs" (ra, imm=0). Compression: c.jr/c.jalr (imm=0).
type Jalr struct {
	base

	rd, rs1 string
	off     imm
}

// Jalr - jalr rd, off(rs1) (off = byte offset from rs1).
func (Builder) Jalr(rd, rs1 Reg, off Off) Instr {
	return Jalr{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeJalr(w uint32, addr uint64) Instr {
	return Jalr{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

// cJalr - compressed forms (c.jr/c.jalr): base - halfword, length 2.
func cJalr(h uint32, addr uint64, rd, rs1 string, off int64) Jalr {
	return Jalr{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		off:  immNum(off),
	}
}

func (i Jalr) ObjDump(_ disasm.ViewCtx) string {
	switch {
	case i.rd == "zero" && i.rs1 == "ra" && i.off.val == 0:
		return "ret"
	case i.rd == "zero" && i.off.val == 0:
		return "jr " + i.rs1
	case i.rd == "zero":
		return fmt.Sprintf("jr %s(%s)", i.off.text(), i.rs1)
	case i.rd == "ra" && i.off.val == 0:
		return "jalr " + i.rs1 // indirect call
	}

	return fmt.Sprintf("jalr %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Jalr) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["jalr"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if off == 0 && !o.NoRVC {
		if r := r5(i.rs1); r != 0 {
			switch i.rd {
			case "zero":
				return writeHalf(w, 0x8002|r<<7) // c.jr (rs2=0)
			case "ra":
				return writeHalf(w, 0x9002|r<<7) // c.jalr
			}
		}
	}

	return writeWord(w, word)
}

func (i Jalr) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"jalr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

// newJalr - constructor from parsing: jalr rd, off(rs) | jalr rs (indirect call,
// the 32-bit form jalr ra, 0(rs) - compression to c.jalr changes the text, so
// the form is fixed).
func newJalr(ops []Op) (Instr, error) {
	if len(ops) == 1 && ops[0].reg != "" {
		if _, err := rvRegNumOf(ops[0].reg, false); err != nil {
			return nil, fmt.Errorf("jalr: %w", err)
		}

		return JalrReg{rs1: ops[0].reg}, nil
	}

	if len(ops) != 2 {
		return nil, errors.New("jalr expects rd, off(rs) or rs")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("jalr: %w", err)
	}

	base, off, err := wantMem(ops[1], "jalr")
	if err != nil {
		return nil, err
	}

	return Jalr{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
