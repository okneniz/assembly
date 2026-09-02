package riscv

// System instructions (no operands): ecall/ebreak/mret/sret/wfi and fence
// (an optional fm operand). The system word = match; fence encodes fm.

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// systemInstr - an operand-less system instruction: word = match.
type systemInstr struct {
	base

	name  string
	group string
}

// cSystem - compressed forms (c.ebreak): base - halfword, length 2.
func cSystem(h uint32, addr uint64, name, group string) systemInstr {
	return systemInstr{
		base:  newHalfBase(h, addr),
		name:  name,
		group: group,
	}
}

func (i systemInstr) ObjDump(_ disasm.ViewCtx) string {
	return i.name
}

func (i systemInstr) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings[i.name][0])
}

func (i systemInstr) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(i.name, i.ObjDump(disasm.DefaultViewCtx()), i.group, nil)
}

func decodeSystem(name, group string) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return systemInstr{
			base:  newBase(addr, w),
			name:  name,
			group: group,
		}
	}
}

func newSystem(name string) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if len(ops) != 0 {
			return nil, fmt.Errorf("%s expects no operands", name)
		}

		return systemInstr{name: name}, nil
	}
}

// Fence - fence pred,succ (an optional numeric fm; the bare word 0xf
// is printed as "fence").
type Fence struct {
	base

	fm imm
}

// Fence - fence fm; fm is a 4-bit fence modifier 0..15
// (fm 0 is printed as the bare "fence").
func (Builder) Fence(fm uint8) Instr {
	return Fence{
		fm: immNum(int64(fm)),
	}
}

func decodeFence(w uint32, addr uint64) Instr {
	return Fence{
		base: newBase(addr, w),
		fm:   immNum(int64(w >> 20 & 0xf)),
	}
}

func (i Fence) ObjDump(_ disasm.ViewCtx) string {
	if i.fm.val == 0 {
		return "fence"
	}

	return "fence " + i.fm.text()
}

func (i Fence) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	fm := i.fm.val

	if fm < 0 || fm > 0xf {
		return 0, fmt.Errorf("fence fm %#x out of range", fm)
	}

	return writeWord(w, riscvEncodings["fence"][0]|uint32(fm)<<20)
}

func (i Fence) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("fence", i.ObjDump(disasm.DefaultViewCtx()), "RV32I", nil)
}

func newFence(ops []Op) (Instr, error) {
	f := Fence{}
	if len(ops) > 1 {
		return nil, errors.New("fence: want at most one operand")
	}

	if len(ops) == 1 {
		e, err := wantExpr(ops[0])
		if err != nil {
			return nil, fmt.Errorf("fence: %w", err)
		}

		f.fm = e
	}

	return f, nil
}
