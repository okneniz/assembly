package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Jal - jal rd, target; pseudo: j (rd=zero), "jal target" (rd=ra, rd
// omitted). Compression: c.j (rd=zero).
type Jal struct {
	base

	rd     string
	target imm // absolute target address
}

func decodeJal(w uint32, addr uint64) Instr {
	return Jal{
		base:   newBase(addr, w),
		rd:     rvRegNames[w>>7&0x1f],
		target: immNum(int64(addr) + jImm(w)),
	}
}

// cJal - compressed forms (c.j): base - halfword, length 2.
func cJal(h uint32, addr uint64, rd string, target int64) Jal {
	return Jal{
		base:   newHalfBase(h, addr),
		rd:     rd,
		target: immNum(target),
	}
}

func (i Jal) ObjDump(_ disasm.ViewCtx) string {
	switch i.rd {
	case "zero":
		return "j " + i.target.text()
	case "ra":
		return "jal " + i.target.text() // rd omitted
	}

	return fmt.Sprintf("jal %s, %s", i.rd, i.target.text())
}

func (i Jal) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encJ(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["jal"][0] | regBits(i.rd)<<7 | bits
	if i.rd == "zero" && !o.NoRVC {
		if half, ok := cjal(target - int64(pc)); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func (i Jal) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"jal",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd},
	)
}

// newJal - constructor from parsing: jal target | jal rd, target.
func newJal(ops []Op) (Instr, error) {
	switch len(ops) {
	case 1: // jal target → jal ra, target
		e, err := wantExpr(ops[0])
		if err != nil {
			return nil, fmt.Errorf("jal: %w", err)
		}

		return Jal{
			rd:     "ra",
			target: e,
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
			rd:     rd,
			target: e,
		}, nil
	}

	return nil, errors.New("jal expects 1 or 2 operands")
}
