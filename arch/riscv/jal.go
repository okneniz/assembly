package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Jal - jal rd, off; pseudo: j (rd=zero), "jal off" (rd=ra, rd
// omitted). Compression: c.j (rd=zero).
type Jal struct {
	base

	rd  string
	off imm // pc-relative byte offset
}

// Jal - jal rd, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Jal(rd Reg, off int64) Instr {
	return Jal{
		rd:  rd.name(),
		off: immNum(off),
	}
}

func decodeJal(w uint32) Instr {
	return Jal{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		off:  immNum(jImm(w)),
	}
}

// cJal - compressed forms (c.j): base - halfword, length 2.
func cJal(h uint32, rd string, off int64) Jal {
	return Jal{
		base: newHalfBase(h),
		rd:   rd,
		off:  immNum(off),
	}
}

func (i Jal) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	switch i.rd {
	case "zero":
		return "j " + target.text()
	case "ra":
		return "jal " + target.text() // rd omitted
	}

	return fmt.Sprintf("jal %s, %s", i.rd, target.text())
}

func (i Jal) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encJ(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["jal"][0] | regBits(i.rd)<<7 | bits
	if i.rd == "zero" && !o.NoRVC {
		if half, ok := cjal(i.off.val); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

// newJal - constructor from parsing: jal off | jal rd, off.
func newJal(ops []Op) (Instr, error) {
	switch len(ops) {
	case 1: // jal off → jal ra, off
		e, err := wantExpr(ops[0])
		if err != nil {
			return nil, fmt.Errorf("jal: %w", err)
		}

		return Jal{
			rd:  "ra",
			off: e,
		}, nil
	case 2:
		rd, err := wantReg(ops[0], false)
		if err != nil {
			return nil, fmt.Errorf("jal: %w", err)
		}

		e, err := wantExpr(ops[1])
		if err != nil {
			return nil, fmt.Errorf("jal: %w", err)
		}

		return Jal{
			rd:  rd,
			off: e,
		}, nil
	}

	return nil, errors.New("jal expects 1 or 2 operands")
}
