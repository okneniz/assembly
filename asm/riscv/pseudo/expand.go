package pseudo

// Expansion of single-word pseudo-instructions into base forms (the
// inverse of the arch formatters' pseudo logic): nop<->addi
// zero,zero,0; li<->addi rd,zero,imm; mv<->addi rd,rs,0;
// not/neg/sext.w/seqz/snez/zext.b; beqz..blez; j/jal/ret/jr; csr
// pseudo-instructions (csrr/csrw/fs*/rd*). Multi-word ones
// (la/call/tail) have their own structures (la.go/call.go/tail.go).
// The base forms are encoded as real syntax layer forms
// (ResolveForm): RVC compression by the general rules.

import (
	"bytes"
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
	riscv "github.com/okneniz/assembly/asm/riscv"
)

// expandable is the dictionary of single-word pseudo-instructions
// (expansion into base forms).
var expandable = []string{
	"nop", "li", "mv", "not", "neg", "negw", "sext.w", "seqz", "snez",
	"sltz", "sgtz", "zext.b",
	"beqz", "bnez", "bgtz", "bltz", "bgez", "blez",
	"j", "ret", "jr",
	"csrr", "csrw", "csrwi",
	"frcsr", "frrm", "frflags", "rdcycle", "rdtime", "rdinstret",
	"fscsr", "fsrm", "fsflags", "fscsri", "fsrmi", "fsflagsi",
}

// form is a base form after expansion (mnemonic + operands).
type form struct {
	mnem string
	ops  []riscv.Op
}

func newForm(mnem string, ops []riscv.Op) form {
	return form{
		mnem: mnem,
		ops:  ops,
	}
}

// expandPseudo expands pseudo-forms into base ones.
func expandPseudo(p pInstr) ([]form, error) {
	ops := p.ops
	switch p.mnem {
	case "nop":
		return rebase("addi", riscv.OpReg("zero"), riscv.OpReg("zero"), riscv.OpNum(0)), nil

	case "mv":
		if len(ops) != 2 {
			return nil, errors.New("mv expects 2 operands")
		}

		return rebase("addi", ops[0], ops[1], riscv.OpNum(0)), nil
	case "not":
		if len(ops) != 2 {
			return nil, errors.New("not expects 2 operands")
		}

		return rebase("xori", ops[0], ops[1], riscv.OpNum(-1)), nil
	case "seqz":
		if len(ops) != 2 {
			return nil, errors.New("seqz expects 2 operands")
		}

		return rebase("sltiu", ops[0], ops[1], riscv.OpNum(1)), nil
	case "zext.b":
		if len(ops) != 2 {
			return nil, errors.New("zext.b expects 2 operands")
		}

		return rebase("andi", ops[0], ops[1], riscv.OpNum(0xff)), nil
	case "sext.w":
		if len(ops) != 2 {
			return nil, errors.New("sext.w expects 2 operands")
		}

		return rebase("addiw", ops[0], ops[1], riscv.OpNum(0)), nil
	case "neg":
		if len(ops) != 2 {
			return nil, errors.New("neg expects 2 operands")
		}

		return rebase("sub", ops[0], riscv.OpReg("zero"), ops[1]), nil
	case "negw":
		if len(ops) != 2 {
			return nil, errors.New("negw expects 2 operands")
		}

		return rebase("subw", ops[0], riscv.OpReg("zero"), ops[1]), nil
	case "snez":
		if len(ops) != 2 {
			return nil, errors.New("snez expects 2 operands")
		}

		return rebase("sltu", ops[0], riscv.OpReg("zero"), ops[1]), nil
	case "sltz":
		if len(ops) != 2 {
			return nil, errors.New("sltz expects 2 operands")
		}

		return rebase("slt", ops[0], ops[1], riscv.OpReg("zero")), nil
	case "sgtz":
		if len(ops) != 2 {
			return nil, errors.New("sgtz expects 2 operands")
		}

		return rebase("slt", ops[0], riscv.OpReg("zero"), ops[1]), nil

	case "beqz", "bnez":
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s expects 2 operands", p.mnem)
		}

		return branchPseudo("beq", "bne", p.mnem == "bnez", ops), nil
	case "bgtz":
		if len(ops) != 2 {
			return nil, errors.New("bgtz expects 2 operands")
		}

		return rebase("blt", riscv.OpReg("zero"), ops[0], ops[1]), nil
	case "bltz":
		if len(ops) != 2 {
			return nil, errors.New("bltz expects 2 operands")
		}

		return rebase("blt", ops[0], riscv.OpReg("zero"), ops[1]), nil
	case "blez":
		if len(ops) != 2 {
			return nil, errors.New("blez expects 2 operands")
		}

		return rebase("bge", riscv.OpReg("zero"), ops[0], ops[1]), nil
	case "bgez":
		if len(ops) != 2 {
			return nil, errors.New("bgez expects 2 operands")
		}

		return rebase("bge", ops[0], riscv.OpReg("zero"), ops[1]), nil

	case "j":
		if len(ops) != 1 {
			return nil, errors.New("j expects 1 operand")
		}

		return []form{newForm("jal", []riscv.Op{riscv.OpReg("zero"), ops[0]})}, nil

	case "ret":
		return []form{
			newForm("jalr", []riscv.Op{riscv.OpReg("zero"), riscv.OpMem("ra", 0)}),
		}, nil

	case "jr":
		if len(ops) != 1 {
			return nil, errors.New("jr expects 1 operand")
		}

		if ops[0].IsMem() { // jr imm(rs)
			return []form{newForm("jalr", []riscv.Op{riscv.OpReg("zero"), ops[0]})}, nil
		}

		return []form{
			newForm("jalr", []riscv.Op{riscv.OpReg("zero"), riscv.OpMem(ops[0].Reg(), 0)}),
		}, nil

	case "li":
		return expandLi(p)

	case "csrr":
		if len(ops) != 2 {
			return nil, errors.New("csrr expects 2 operands")
		}

		return csrRw("csrrs", ops[0], ops[1], riscv.OpReg("zero")), nil
	case "csrw":
		if len(ops) != 2 {
			return nil, errors.New("csrw expects 2 operands")
		}

		return csrRw("csrrw", riscv.OpReg("zero"), ops[0], ops[1]), nil
	case "csrwi":
		if len(ops) != 2 {
			return nil, errors.New("csrwi expects 2 operands")
		}

		return csrRw("csrrwi", riscv.OpReg("zero"), ops[0], ops[1]), nil

	case "frcsr", "frrm", "frflags", "rdcycle", "rdtime", "rdinstret":
		csrs := map[string]string{
			"frcsr": "fcsr", "frrm": "frm", "frflags": "fflags",
			"rdcycle": "cycle", "rdtime": "time", "rdinstret": "instret",
		}
		if len(ops) != 1 {
			return nil, fmt.Errorf("%s expects 1 operand", p.mnem)
		}

		return csrRw(
			"csrrs",
			ops[0],
			riscv.OpExpr(expr.Sym(csrs[p.mnem])),
			riscv.OpReg("zero"),
		), nil
	case "fscsr", "fsrm", "fsflags":
		csrs := map[string]string{"fscsr": "fcsr", "fsrm": "frm", "fsflags": "fflags"}
		if len(ops) != 1 {
			return nil, fmt.Errorf("%s expects 1 operand", p.mnem)
		}

		return csrRw(
			"csrrw",
			riscv.OpReg("zero"),
			riscv.OpExpr(expr.Sym(csrs[p.mnem])),
			ops[0],
		), nil
	case "fscsri", "fsrmi", "fsflagsi":
		csrs := map[string]string{"fscsri": "fcsr", "fsrmi": "frm", "fsflagsi": "fflags"}
		if len(ops) != 1 {
			return nil, fmt.Errorf("%s expects 1 operand", p.mnem)
		}

		return csrRw(
			"csrrwi",
			riscv.OpReg("zero"),
			riscv.OpExpr(expr.Sym(csrs[p.mnem])),
			ops[0],
		), nil
	}

	return []form{newForm(p.mnem, ops)}, nil
}

func rebase(mnem string, ops ...riscv.Op) []form {
	return []form{newForm(mnem, ops)}
}

func branchPseudo(eq, ne string, useNe bool, ops []riscv.Op) []form {
	name := eq
	if useNe {
		name = ne
	}

	return []form{newForm(name, []riscv.Op{ops[0], riscv.OpReg("zero"), ops[1]})}
}

func csrRw(mnem string, rd, csr, src riscv.Op) []form {
	return []form{newForm(mnem, []riscv.Op{rd, csr, src})}
}

// expandLi expands li rd, imm: the shortest sequence.
// imm must be numeric (for symbolic targets use la).
func expandLi(p pInstr) ([]form, error) {
	if len(p.ops) != 2 || p.ops[1].Expr() == nil {
		return nil, errors.New("li expects rd, imm")
	}

	v, err := p.ops[1].Expr().Eval(nil)
	if err != nil {
		return nil, fmt.Errorf("li immediate must be a constant (use la for symbols): %w", err)
	}

	rd := p.ops[0]
	if arch.Fits12(v) {
		return rebase("addi", rd, riscv.OpReg("zero"), riscv.OpNum(v)), nil
	}

	hi, lo := arch.PcrelHiLo(v)
	if lo == 0 {
		return rebase("lui", rd, riscv.OpNum(hi&0xfffff)), nil
	}

	return []form{
		newForm("lui", []riscv.Op{rd, riscv.OpNum(hi & 0xfffff)}),
		newForm("addi", []riscv.Op{rd, rd, riscv.OpNum(lo)}),
	}, nil
}

// preferAddFormMv prefers the add form (c.mv) for mv rd, rs when
// possible, as GNU as does; otherwise the addi form (as the decoder
// prints). The decision is deterministic in both passes.
func (p pInstr) preferAddFormMv(forms []form, ctx asm.Ctx) []form {
	if p.mnem != "mv" || len(p.ops) != 2 || !p.ops[0].IsReg() || !p.ops[1].IsReg() ||
		p.ops[0].Reg() == "zero" || p.ops[1].Reg() == "zero" {
		return forms
	}

	addOps := []riscv.Op{p.ops[0], riscv.OpReg("zero"), p.ops[1]}
	// the add form is preferred when it compresses into c.mv (as GNU as does)
	if res, err := p.be.ResolveForm("add", addOps, ctx); err == nil {
		var cbuf bytes.Buffer
		if _, werr := res.Encode(&cbuf); werr == nil && cbuf.Len() == 2 {
			return []form{newForm("add", addOps)}
		}
	}

	return forms
}
