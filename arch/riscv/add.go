package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Add — add rd, rs1, rs2.
type Add struct {
	base

	rd, rs1, rs2 string
}

// NewAdd — add rd, rs1, rs2.
func NewAdd(rd, rs1, rs2 Reg) Instr {
	return Add{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAdd(w uint32, addr uint64) Instr {
	return Add{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

// cAdd - compressed forms (c.add): base - halfword, length 2.
func cAdd(h uint32, addr uint64, rd, rs1, rs2 string) Add {
	return Add{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs1:  rs1,
		rs2:  rs2,
	}
}

func (i Add) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Add) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	word := riscvEncodings["add"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	rs2, ok := asmRegNum[i.rs2]
	rd := r5(i.rd)
	if ok && !rs2.fp && !o.NoRVC {
		switch {
		case i.rs1 == "zero" && rd != 0 && i.rs2 != "zero":
			return writeHalf(w, 0x8002|rd<<7|uint16(rs2.num)<<2) // c.mv
		case i.rd == i.rs1 && rd != 0 && i.rs2 != "zero":
			return writeHalf(w, 0x9002|rd<<7|uint16(rs2.num)<<2) // c.add
		}
	}

	return writeWord(w, word)
}

func (i Add) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"add",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAdd(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "add")
	if err != nil {
		return nil, err
	}

	return Add{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
