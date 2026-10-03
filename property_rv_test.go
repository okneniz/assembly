package assembly_test

// Property tests (oh-snap), riscv: round trip of single instructions and
// compositions (variable length 2/4, canonicalization of pseudo-forms),
// decoder robustness, differential against the real objdump. Generators are
// arb/riscv on top of arch/riscv constructors. The common part of the suite
// is property_test.go.

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	mrnd "math/rand/v2"
	"strconv"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec"
	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	rv "github.com/okneniz/assembly/arb/riscv"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/tests/cmd/objdump"
	"github.com/okneniz/assembly/text"
)

// rvText - normalized ObjDump text of a riscv instruction.
func rvText(in riscv.Instr) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.DefaultViewCtx())))
}

// rvTextAt - rvText in the context of an explicit address.
func rvTextAt(in riscv.Instr, addr uint64) string {
	return objdump.StripComments(objdump.Normalize(in.ObjDump(disasm.ViewCtxAt(addr))))
}

// rvAssemblesTo - bytes from assembling the text (false on assembly error).
func rvAssemblesTo(t *testing.T, src string) ([]byte, bool) {
	t.Helper()
	res, errs := pseudo.Assemble(src, propAddr)
	if len(errs) != 0 {
		t.Logf("%q: assemble: %v", src, errs)
		return nil, false
	}

	return res.Sections[0].Data, true
}

// rvBytesOf - instruction bytes (2 or 4 - RVC compression).
func rvBytesOf(t *testing.T, in riscv.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	if _, err := in.Encode(&buf, riscv.EncOpts{}); err != nil {
		t.Logf("%s: Encode: %v", in.ObjDump(disasm.DefaultViewCtx()), err)
		return nil, false
	}

	return buf.Bytes(), true
}

// rvEnc - the encoding context of the "bytes" property: a fixed address
// (PC-relative forms), unrestricted modes, no symbols.
type rvEnc struct {
	addr uint64
}

// rvEncodeAll - encodes a list sequentially starting at propAddr (the address
// of each is propAddr + bytes written; compression allowed, no symbols).
func rvEncodeAll(t *testing.T, ins []riscv.Instr) ([]byte, bool) {
	t.Helper()
	var buf bytes.Buffer
	addr := propAddr
	for _, in := range ins {
		n, err := in.Encode(&buf, riscv.EncOpts{})
		if err != nil {
			t.Logf("encode: %v", err)
			return nil, false
		}

		addr += int(n)
	}

	return buf.Bytes(), true
}

// rvBytesRoundTrip - the "bytes" property: the law enc∘dec∘enc == enc
// (RoundTrip) in the encoding context propAddr without symbols - bytes are
// stable after the round trip.
func rvBytesRoundTrip(t *testing.T, in riscv.Instr) bool {
	t.Helper()
	return RoundTrip[rvEnc, riscv.Instr, []byte](
		rvEnc{addr: propAddr},
		func(_ rvEnc, x riscv.Instr) ([]byte, bool) {
			return rvBytesOf(t, x)
		},
		func(ctx rvEnc, b []byte) (riscv.Instr, bool) {
			back, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(b))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(back) != 1 {
				t.Logf("% x: decode = %d instr", b, len(back))
				return nil, false
			}

			return back[0], true
		},
		bytes.Equal,
	)(in)
}

// rvTextRoundTrip - the "text" property: the RoundTrip law in the
// DefaultViewCtx context - a fixed point of the "text → assembly → decode"
// cycle (pseudo-forms collapse into a single canon: mv and its base form
// print identically). The input is canonicalized through bytes and decode:
// the fixed point is sought for the decoder's text, while structures from
// parsing may be pseudo-forms whose text is not fixed.
func rvTextRoundTrip(t *testing.T, in riscv.Instr) bool {
	t.Helper()
	b, ok := rvBytesOf(t, in)
	if !ok {
		return false
	}

	d1, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(b))
	if err != nil {
		t.Logf("decode: %v", err)
		return false
	}

	if len(d1) != 1 {
		t.Logf("% x: decode = %d instr", b, len(d1))
		return false
	}

	// the text renders and assembles at the SAME base (propAddr): the
	// pc-relative texts carry absolute targets, the two sides of the
	// law must not shift them apart
	return RoundTrip[disasm.ViewCtx, riscv.Instr, string](
		disasm.ViewCtxAt(propAddr),
		func(_ disasm.ViewCtx, y riscv.Instr) (string, bool) {
			return objdump.StripComments(objdump.Normalize(
				y.ObjDump(disasm.ViewCtxAt(propAddr)))), true
		},
		func(_ disasm.ViewCtx, src string) (riscv.Instr, bool) {
			data, ok := rvAssemblesTo(t, src)
			if !ok {
				return nil, false
			}

			d2, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
			if err != nil {
				t.Logf("decode: %v", err)
				return nil, false
			}

			if len(d2) != 1 {
				t.Logf("%q: decode = %d instr", src, len(d2))
				return nil, false
			}

			return d2[0], true
		},
		func(a, b string) bool { return a == b },
	)(d1[0])
}

// rvInstrParam - parameters of any riscv family.
type rvInstrParam interface {
	Instr() riscv.Instr
}

// rvFamilyEntry - one riscv family in the TestPropertyRiscvSingleInstrRoundTrip table.
type rvFamilyEntry struct {
	name string
	run  func(t *testing.T, rnd *mrnd.Rand)
}

// newRvFamily - a riscv family entry: closes over the generic instantiation
// of the properties (families have different parameter types, Arbitrary is
// invariant - see newPropFamily in property_a64_test.go; each property is a
// separate named subtest, the check runs over the parameters for the sake of
// shrinking).
func newRvFamily[P rvInstrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
) rvFamilyEntry {
	return rvFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return rvBytesRoundTrip(t, p.Instr())
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					return rvTextRoundTrip(t, p.Instr())
				})
			})
		},
	}
}

// TestPropertyRiscvSingleInstrRoundTrip - round trip of each riscv family.
func TestPropertyRiscvSingleInstrRoundTrip(t *testing.T) {
	families := []rvFamilyEntry{
		newRvFamily("Add", rv.Add),
		newRvFamily("Sub", rv.Sub),
		newRvFamily("Addi", rv.Addi),
		newRvFamily("Addiw", rv.Addiw),
		newRvFamily("Lui", rv.Lui),
		newRvFamily("Lw", rv.Lw),
		newRvFamily("Ld", rv.Ld),
		newRvFamily("Sw", rv.Sw),
		newRvFamily("Sd", rv.Sd),
		newRvFamily("And", rv.And),
		newRvFamily("Or", rv.Or),
		newRvFamily("Xor", rv.Xor),
		newRvFamily("Sll", rv.Sll),
		newRvFamily("Srl", rv.Srl),
		newRvFamily("Sra", rv.Sra),
		newRvFamily("Slt", rv.Slt),
		newRvFamily("Sltu", rv.Sltu),
		newRvFamily("Addw", rv.Addw),
		newRvFamily("Subw", rv.Subw),
		newRvFamily("Sllw", rv.Sllw),
		newRvFamily("Srlw", rv.Srlw),
		newRvFamily("Sraw", rv.Sraw),
		newRvFamily("Mul", rv.Mul),
		newRvFamily("Mulh", rv.Mulh),
		newRvFamily("Mulhsu", rv.Mulhsu),
		newRvFamily("Mulhu", rv.Mulhu),
		newRvFamily("Mulw", rv.Mulw),
		newRvFamily("Div", rv.Div),
		newRvFamily("Divu", rv.Divu),
		newRvFamily("Divw", rv.Divw),
		newRvFamily("Divuw", rv.Divuw),
		newRvFamily("Rem", rv.Rem),
		newRvFamily("Remu", rv.Remu),
		newRvFamily("Remw", rv.Remw),
		newRvFamily("Remuw", rv.Remuw),
		newRvFamily("AmoaddW", rv.AmoaddW),
		newRvFamily("AmoaddD", rv.AmoaddD),
		newRvFamily("AmoandW", rv.AmoandW),
		newRvFamily("AmoandD", rv.AmoandD),
		newRvFamily("AmomaxW", rv.AmomaxW),
		newRvFamily("AmomaxD", rv.AmomaxD),
		newRvFamily("AmomaxuW", rv.AmomaxuW),
		newRvFamily("AmomaxuD", rv.AmomaxuD),
		newRvFamily("AmominW", rv.AmominW),
		newRvFamily("AmominD", rv.AmominD),
		newRvFamily("AmominuW", rv.AmominuW),
		newRvFamily("AmominuD", rv.AmominuD),
		newRvFamily("AmoorW", rv.AmoorW),
		newRvFamily("AmoorD", rv.AmoorD),
		newRvFamily("AmoswapW", rv.AmoswapW),
		newRvFamily("AmoswapD", rv.AmoswapD),
		newRvFamily("AmoxorW", rv.AmoxorW),
		newRvFamily("AmoxorD", rv.AmoxorD),
		newRvFamily("Andi", rv.Andi),
		newRvFamily("Ori", rv.Ori),
		newRvFamily("Xori", rv.Xori),
		newRvFamily("Slti", rv.Slti),
		newRvFamily("Sltiu", rv.Sltiu),
		newRvFamily("Slli", rv.Slli),
		newRvFamily("Srli", rv.Srli),
		newRvFamily("Srai", rv.Srai),
		newRvFamily("Slliw", rv.Slliw),
		newRvFamily("Srliw", rv.Srliw),
		newRvFamily("Sraiw", rv.Sraiw),
		newRvFamily("Lb", rv.Lb),
		newRvFamily("Lbu", rv.Lbu),
		newRvFamily("Lh", rv.Lh),
		newRvFamily("Lhu", rv.Lhu),
		newRvFamily("Lwu", rv.Lwu),
		newRvFamily("Flw", rv.Flw),
		newRvFamily("Fld", rv.Fld),
		newRvFamily("Sb", rv.Sb),
		newRvFamily("Sh", rv.Sh),
		newRvFamily("Fsw", rv.Fsw),
		newRvFamily("Fsd", rv.Fsd),
		newRvFamily("Beq", rv.Beq),
		newRvFamily("Bne", rv.Bne),
		newRvFamily("Blt", rv.Blt),
		newRvFamily("Bge", rv.Bge),
		newRvFamily("Bltu", rv.Bltu),
		newRvFamily("Bgeu", rv.Bgeu),
		newRvFamily("Jal", rv.Jal),
		newRvFamily("Jalr", rv.Jalr),
		newRvFamily("JalrReg", rv.JalrReg),
		newRvFamily("Mv", rv.Mv),
		newRvFamily("Auipc", rv.Auipc),
		newRvFamily("Csrrw", rv.Csrrw),
		newRvFamily("Csrrs", rv.Csrrs),
		newRvFamily("Csrrc", rv.Csrrc),
		newRvFamily("Csrrwi", rv.Csrrwi),
		newRvFamily("Csrrsi", rv.Csrrsi),
		newRvFamily("Csrrci", rv.Csrrci),
		newRvFamily("Fence", rv.Fence),
		newRvFamily("FaddS", rv.FaddS),
		newRvFamily("FaddD", rv.FaddD),
		newRvFamily("FsubS", rv.FsubS),
		newRvFamily("FsubD", rv.FsubD),
		newRvFamily("FmulS", rv.FmulS),
		newRvFamily("FmulD", rv.FmulD),
		newRvFamily("FdivS", rv.FdivS),
		newRvFamily("FdivD", rv.FdivD),
		newRvFamily("FmaddS", rv.FmaddS),
		newRvFamily("FmaddD", rv.FmaddD),
		newRvFamily("FmsubS", rv.FmsubS),
		newRvFamily("FmsubD", rv.FmsubD),
		newRvFamily("FnmaddS", rv.FnmaddS),
		newRvFamily("FnmaddD", rv.FnmaddD),
		newRvFamily("FnmsubS", rv.FnmsubS),
		newRvFamily("FnmsubD", rv.FnmsubD),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// newRvAliasFamily - like newRvFamily, but the parameters are pinned
// (the alias condition of a pseudo-form); combinations invalid after pinning
// are outside the property. See newAliasFamily in property_a64_test.go.
func newRvAliasFamily[P rvInstrParam](
	name string,
	mk func(rnd *mrnd.Rand) ohsnap.Arbitrary[P],
	pin func(P) P,
) rvFamilyEntry {
	return rvFamilyEntry{
		name: name,
		run: func(t *testing.T, rnd *mrnd.Rand) {
			t.Helper()
			pinned := func(p P) (riscv.Instr, bool) {
				q := pin(p)
				if q.Instr() == nil {
					return nil, false
				}

				return q.Instr(), true
			}

			t.Run("bytes", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return rvBytesRoundTrip(t, in)
					}

					return true
				})
			})

			t.Run("text", func(t *testing.T) {
				ohsnap.Check(t, 300, mk(rnd), func(p P) bool {
					if in, ok := pinned(p); ok {
						return rvTextRoundTrip(t, in)
					}

					return true
				})
			})
		},
	}
}

// TestPropertyRiscvAliasRoundTrip - the "alias" property (pseudo-forms):
// an instruction written as a pseudo-form survives the cycle
// struct → alias → struct without loss; each entry is its own subtest.
// nop is pinned entirely (the only instruction of the form) - the check
// degenerates into a constant one, but keeps the common framework in place.
func TestPropertyRiscvAliasRoundTrip(t *testing.T) {
	families := []rvFamilyEntry{
		newRvAliasFamily("mv", rv.Addi, func(p rv.AddiParams) rv.AddiParams {
			p.Imm = riscv.Imm12{}
			return p
		}),
		newRvAliasFamily("nop", rv.Addi, func(p rv.AddiParams) rv.AddiParams {
			p.Rd = riscv.Zero
			p.Rs1 = riscv.Zero
			p.Imm = riscv.Imm12{}
			return p
		}),
		newRvAliasFamily("not", rv.Xori, func(p rv.RiParams) rv.RiParams {
			if v, err := riscv.New().Imm12(-1); err == nil {
				p.Imm = v
			}

			return p
		}),
		newRvAliasFamily("seqz", rv.Sltiu, func(p rv.RiParams) rv.RiParams {
			if v, err := riscv.New().Imm12(1); err == nil {
				p.Imm = v
			}

			return p
		}),
		newRvAliasFamily("zext.b", rv.Andi, func(p rv.RiParams) rv.RiParams {
			if v, err := riscv.New().Imm12(0xff); err == nil {
				p.Imm = v
			}

			return p
		}),
		newRvAliasFamily("sext.w", rv.Addiw, func(p rv.RiParams) rv.RiParams {
			p.Imm = riscv.Imm12{}
			return p
		}),
		newRvAliasFamily("neg", rv.Sub, func(p rv.SubParams) rv.SubParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("negw", rv.Subw, func(p rv.RrrParams) rv.RrrParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("snez", rv.Sltu, func(p rv.RrrParams) rv.RrrParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("sltz", rv.Slt, func(p rv.RrrParams) rv.RrrParams {
			p.Rs2 = riscv.Zero
			return p
		}),
		newRvAliasFamily("sgtz", rv.Slt, func(p rv.RrrParams) rv.RrrParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("beqz", rv.Beq, func(p rv.BranchParams) rv.BranchParams {
			p.Rs2 = riscv.Zero
			return p
		}),
		newRvAliasFamily("bnez", rv.Bne, func(p rv.BranchParams) rv.BranchParams {
			p.Rs2 = riscv.Zero
			return p
		}),
		newRvAliasFamily("bgtz", rv.Blt, func(p rv.BranchParams) rv.BranchParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("bltz", rv.Blt, func(p rv.BranchParams) rv.BranchParams {
			p.Rs2 = riscv.Zero
			return p
		}),
		newRvAliasFamily("bgez", rv.Bge, func(p rv.BranchParams) rv.BranchParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("blez", rv.Bge, func(p rv.BranchParams) rv.BranchParams {
			p.Rs2 = riscv.Zero
			return p
		}),
		newRvAliasFamily("j", rv.Jal, func(p rv.JalParams) rv.JalParams {
			p.Rd = riscv.Zero
			return p
		}),
		newRvAliasFamily("ret", rv.Jalr, func(p rv.JalrParams) rv.JalrParams {
			p.Rd = riscv.Zero
			p.Rs1 = riscv.Ra
			p.Off = riscv.Off{}
			return p
		}),
		newRvAliasFamily("jr", rv.Jalr, func(p rv.JalrParams) rv.JalrParams {
			p.Rd = riscv.Zero
			return p
		}),
		newRvAliasFamily("csrr", rv.Csrrs, func(p rv.CsrParams) rv.CsrParams {
			p.Rs1 = riscv.Zero
			return p
		}),
		newRvAliasFamily("csrw", rv.Csrrw, func(p rv.CsrParams) rv.CsrParams {
			p.Rd = riscv.Zero
			return p
		}),
		newRvAliasFamily("csrwi", rv.Csrrwi, func(p rv.CsrParams) rv.CsrParams {
			p.Rd = riscv.Zero
			return p
		}),
		newRvAliasFamily("frcsr", rv.Csrrs, csrRead(riscvCSRfcsr)),
		newRvAliasFamily("frrm", rv.Csrrs, csrRead(riscvCSRfrm)),
		newRvAliasFamily("frflags", rv.Csrrs, csrRead(riscvCSRfflags)),
		newRvAliasFamily("rdcycle", rv.Csrrs, csrRead(riscvCSRcycle)),
		newRvAliasFamily("rdtime", rv.Csrrs, csrRead(riscvCSRtime)),
		newRvAliasFamily("rdinstret", rv.Csrrs, csrRead(riscvCSRinstret)),
		newRvAliasFamily("fscsr", rv.Csrrw, csrWrite(riscvCSRfcsr)),
		newRvAliasFamily("fsrm", rv.Csrrw, csrWrite(riscvCSRfrm)),
		newRvAliasFamily("fsflags", rv.Csrrw, csrWrite(riscvCSRfflags)),
		newRvAliasFamily("fscsri", rv.Csrrwi, csrWrite(riscvCSRfcsr)),
		newRvAliasFamily("fsrmi", rv.Csrrwi, csrWrite(riscvCSRfrm)),
		newRvAliasFamily("fsflagsi", rv.Csrrwi, csrWrite(riscvCSRfflags)),
	}
	for _, f := range families {
		t.Run(f.name, func(t *testing.T) {
			f.run(t, seedRnd(t))
		})
	}
}

// The CSR addresses of the pseudo read/write families (the named pool
// of the csr core).
const (
	riscvCSRfflags  uint16 = 0x001
	riscvCSRfrm     uint16 = 0x002
	riscvCSRfcsr    uint16 = 0x003
	riscvCSRcycle   uint16 = 0xc00
	riscvCSRtime    uint16 = 0xc01
	riscvCSRinstret uint16 = 0xc02
)

// TestPropertyRiscvSymbolPseudoRoundTrip - the resolve-level pseudo
// forms: li (the decoded ladder must compute the value - the signed
// 32-bit domain of the expansion) and la/call/tail (the decoded pair
// must land exactly on the target, at the property base).
func TestPropertyRiscvSymbolPseudoRoundTrip(t *testing.T) {
	t.Run("Li", func(t *testing.T) {
		ohsnap.Check(t, 300, rv.Li(seedRnd(t)), func(p rv.LiParams) bool {
			res, errs := pseudo.Assemble(p.String(), propAddr)
			if len(errs) != 0 {
				t.Logf("%q: assemble: %v", p, errs[0])
				return false
			}

			got, ok := rvLiComputes(res.Sections[0].Data)
			if !ok {
				t.Logf("%q: unexpected ladder shape", p)
				return false
			}

			return got == p.Val
		})
	})

	t.Run("La", func(t *testing.T) {
		ohsnap.Check(t, 300, rv.La(seedRnd(t)), func(p rv.LaParams) bool {
			src := fmt.Sprintf("la %s, %#x", p.Rd, uint64(int64(propAddr)+p.Off))
			return rvPcrelLands(t, src, p.Off)
		})
	})

	t.Run("Call", func(t *testing.T) {
		ohsnap.Check(t, 300, rv.Call(seedRnd(t)), func(p rv.CallParams) bool {
			src := fmt.Sprintf("call %#x", uint64(int64(propAddr)+p.Off))
			return rvPcrelLands(t, src, p.Off)
		})
	})

	t.Run("Tail", func(t *testing.T) {
		ohsnap.Check(t, 300, rv.Tail(seedRnd(t)), func(p rv.TailParams) bool {
			src := fmt.Sprintf("tail %#x", uint64(int64(propAddr)+p.Off))
			return rvPcrelLands(t, src, p.Off)
		})
	})
}

// rvSymText - the assembled bytes of a pseudo text at the property
// base, with the absolute target substituted for the offset.
func rvSymText(t *testing.T, src string) ([]byte, bool) {
	t.Helper()
	res, errs := pseudo.Assemble(src, propAddr)
	if len(errs) != 0 {
		t.Logf("%q: assemble: %v", src, errs[0])
		return nil, false
	}

	return res.Sections[0].Data, true
}

// rvDecodeAt - the decoded instructions of bytes at the property base.
func rvDecodeAt(t *testing.T, data []byte) []riscv.Instr {
	t.Helper()
	ins, err := riscv.MakeDecoder()(
		parsec.Stateless{},
		parsecbytes.Buffer(data),
	)
	if err != nil {
		t.Logf("decode: %v", err)
		return nil
	}

	return ins
}

// rvPcrelLands - the text's decoded sequence lands exactly on
// base+off: tail is the printed absolute target, la and call are the
// auipc+addi/jalr pair arithmetic.
func rvPcrelLands(t *testing.T, src string, off int64) bool {
	t.Helper()
	target := int64(propAddr) + off
	data, ok := rvSymText(t, src)
	if !ok {
		return false
	}

	ins := rvDecodeAt(t, data)
	if len(ins) == 0 {
		return false
	}

	texts := make([]string, 0, len(ins))
	for _, in := range ins {
		texts = append(texts, rvTextAt(in, propAddr))
	}

	if strings.HasPrefix(texts[0], "j ") { // tail: the absolute target
		if len(ins) != 1 {
			t.Logf("%q: %d instructions: %v", src, len(ins), texts)
			return false
		}

		v, ok := rvFieldImm(texts[0], "j")
		return ok && v == target
	}

	// la/call: auipc rd, hi + addi/jalr rd, lo(rd)
	if len(ins) != 2 {
		t.Logf("%q: %d instructions: %v", src, len(ins), texts)
		return false
	}

	hi, ok := rvFieldImm(texts[0], "auipc")
	if !ok {
		t.Logf("%q: %q is not an auipc pair head", src, texts[0])
		return false
	}

	if hi >= 1<<19 {
		hi -= 1 << 20 // the 20-bit field reads signed
	}

	lo, ok := rvPairLo(texts[1])
	if !ok {
		t.Logf("%q: %q is not a pair tail", src, texts[1])
		return false
	}

	return int64(propAddr)+hi<<12+lo == target
}

// rvPairLo - the low half of a pcrel pair from its canonical text: the
// addi shape prints as addi/mv/nop/li depending on the operands, the
// jalr shape keeps the mem form (bare "jalr rd" is the lo=0 shape).
func rvPairLo(text string) (int64, bool) {
	switch {
	case text == "nop":
		return 0, true
	case strings.HasPrefix(text, "mv "):
		return 0, true // mv rd, rs is the addi rd, rs, 0 shape
	case strings.HasPrefix(text, "jalr ") && !strings.Contains(text, ","):
		return 0, true // bare "jalr rd" = jalr rd, 0(rd)
	}

	for _, mnem := range []string{"li", "addi", "jalr"} {
		if v, ok := rvFieldImm(text, mnem); ok {
			return v, true
		}
	}

	return 0, false
}

// rvFieldImm - the immediate of a canonical one-operand text
// ("auipc a0, 0x12344", "addi a0, a0, 0x678", "jalr ra, 0x234(ra)"):
// a mem operand keeps only the part before the paren.
func rvFieldImm(text, mnem string) (int64, bool) {
	if !strings.HasPrefix(text, mnem+" ") {
		return 0, false
	}

	f := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ','
	})

	last := f[len(f)-1]
	last, _, _ = strings.Cut(last, "(")

	v, err := strconv.ParseInt(last, 0, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}

// rvLiComputes - the value the decoded li ladder computes (the arch
// operand fields are not exported - the law reads the canonical
// texts: lui sext32(imm20<<12), li/addi the imm, addiw the
// sign-extending sum, slli the shift).
func rvLiComputes(data []byte) (int64, bool) {
	ins, err := riscv.MakeDecoder()(
		parsec.Stateless{},
		parsecbytes.Buffer(data),
	)
	if err != nil {
		return 0, false
	}

	sext32 := func(x int64) int64 { return int64(int32(x)) }

	val := int64(0)
	for _, in := range ins {
		text := rvTextAt(in, propAddr)
		f := strings.FieldsFunc(text, func(r rune) bool {
			return r == ' ' || r == ','
		})

		v, err := strconv.ParseInt(f[len(f)-1], 0, 64)
		if err != nil {
			return 0, false
		}

		switch {
		case strings.HasPrefix(text, "lui "):
			val = sext32(v << 12)
		case strings.HasPrefix(text, "li "):
			val = v
		case strings.HasPrefix(text, "addiw "):
			val = sext32(val + v)
		case strings.HasPrefix(text, "addi "):
			if len(f) >= 3 && f[2] == "zero" {
				val = v // the li form: addi rd, zero, imm
			} else {
				val += v // a ladder step: addi rd, rd, imm
			}
		case strings.HasPrefix(text, "slli "):
			val <<= uint(v)
		default:
			return 0, false
		}
	}

	return val, true
}

// csrRead - the pin of the read pseudo-forms (csrrs with rs1 = x0).
func csrRead(csr uint16) func(rv.CsrParams) rv.CsrParams {
	return func(p rv.CsrParams) rv.CsrParams {
		p.Rs1 = riscv.Zero
		p.Csr = csr
		return p
	}
}

// csrWrite - the pin of the write pseudo-forms (csrrw/csrrwi with
// rd = x0).
func csrWrite(csr uint16) func(rv.CsrParams) rv.CsrParams {
	return func(p rv.CsrParams) rv.CsrParams {
		p.Rd = riscv.Zero
		p.Csr = csr
		return p
	}
}

// rvInstrOf - a sampler function over a riscv family generator (the generator
// is created once; see instrOf in property_a64_test.go).
func rvInstrOf[P rvInstrParam](a ohsnap.Arbitrary[P]) func() riscv.Instr {
	return func() riscv.Instr {
		return ohsnap.First(a.Generate()).Instr()
	}
}

// rvPropFamilies - riscv family generators for composition and the differential.
func rvPropFamilies(rnd *mrnd.Rand) []func() riscv.Instr {
	return []func() riscv.Instr{
		rvInstrOf(rv.Add(rnd)),
		rvInstrOf(rv.Sub(rnd)),
		rvInstrOf(rv.Addi(rnd)),
		rvInstrOf(rv.Lui(rnd)),
		rvInstrOf(rv.Lw(rnd)),
		rvInstrOf(rv.Ld(rnd)),
		rvInstrOf(rv.Sw(rnd)),
		rvInstrOf(rv.Sd(rnd)),
		rvInstrOf(rv.And(rnd)),
		rvInstrOf(rv.Or(rnd)),
		rvInstrOf(rv.Xor(rnd)),
		rvInstrOf(rv.Mul(rnd)),
		rvInstrOf(rv.Div(rnd)),
		rvInstrOf(rv.Rem(rnd)),
		rvInstrOf(rv.Andi(rnd)),
		rvInstrOf(rv.Slli(rnd)),
		rvInstrOf(rv.Lb(rnd)),
		rvInstrOf(rv.Lbu(rnd)),
		rvInstrOf(rv.Sb(rnd)),
		rvInstrOf(rv.Sh(rnd)),
		rvInstrOf(rv.Beq(rnd)),
		rvInstrOf(rv.Bne(rnd)),
		rvInstrOf(rv.Jal(rnd)),
		rvInstrOf(rv.Jalr(rnd)),
		rvInstrOf(rv.Mv(rnd)),
		rvInstrOf(rv.Auipc(rnd)),
		rvInstrOf(rv.Csrrw(rnd)),
		rvInstrOf(rv.Fence(rnd)),
		rvInstrOf(rv.FaddS(rnd)),
		rvInstrOf(rv.FmulD(rnd)),
		rvInstrOf(rv.Fld(rnd)),
		rvInstrOf(rv.Fsw(rnd)),
		rvInstrOf(rv.AmoaddW(rnd)),
	}
}

// TestPropertyRiscvBytesRoundTripList - the "bytes" property for a list of
// instructions: the list (variable length 2/4) is encoded, decoded line by
// line, and encoded again into the same bytes.
func TestPropertyRiscvBytesRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, rvPropFamilies(rnd))
	ohsnap.Check(t, 100, seq, func(ins []riscv.Instr) bool {
		raw, ok := rvEncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(buf.Bytes()))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		raw2, ok := rvEncodeAll(t, back)
		if !ok {
			return false
		}

		buf2 := *bytes.NewBuffer(raw2)

		if !bytes.Equal(buf2.Bytes(), buf.Bytes()) {
			t.Logf("re-encode: % x ≠ % x", buf2.Bytes(), buf.Bytes())
			return false
		}

		return true
	})
}

// TestPropertyRiscvTextRoundTripList - the "text" property for a list of
// instructions: the joined canonical text of the list assembles and decodes
// line by line into the same texts (canonical - see rvTextRoundTrip).
func TestPropertyRiscvTextRoundTripList(t *testing.T) {
	rnd := seedRnd(t)
	seq := arb.Seq(rnd, rvPropFamilies(rnd))
	ohsnap.Check(t, 100, seq, func(ins []riscv.Instr) bool {
		raw, ok := rvEncodeAll(t, ins)
		if !ok {
			return false
		}

		buf := *bytes.NewBuffer(raw)

		back, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(buf.Bytes()))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back) != len(ins) {
			t.Logf("decoded %d of %d", len(back), len(ins))
			return false
		}

		// the text renders at the instruction's own address in the list
		// (pc-relative targets print absolute, the lengths vary with
		// RVC) - the joined text assembles at propAddr, the same base
		// the render used; the addresses step by the stream rule over
		// the encoded bytes
		texts := make([]string, len(back))
		addr := uint64(propAddr)
		off := 0
		for i := range back {
			texts[i] = rvTextAt(back[i], addr)
			n := riscv.InstrLen(raw[off:])
			addr += uint64(n)
			off += n
		}

		data, ok := rvAssemblesTo(t, strings.Join(texts, "\n"))
		if !ok {
			return false
		}

		back2, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
		if err != nil {
			t.Logf("decode: %v", err)
			return false
		}

		if len(back2) != len(ins) {
			t.Logf("decoded %d of %d", len(back2), len(ins))
			return false
		}

		addr = uint64(propAddr)
		off = 0
		for i := range back2 {
			if rvTextAt(back2[i], addr) != texts[i] {
				t.Logf("[%d] text %q ≠ %q", i, rvTextAt(back2[i], addr), texts[i])
				return false
			}

			n := riscv.InstrLen(data[off:])
			addr += uint64(n)
			off += n
		}

		return true
	})
}

// TestPropertyRiscvDecodeRobustness - arbitrary bytes do not crash the decoder.
func TestPropertyRiscvDecodeRobustness(t *testing.T) {
	rnd := seedRnd(t)
	ohsnap.CheckWith(t, 100000, arb.Word(rnd), func(w uint32) bool {
		ok := true
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic at %#08x: %v", w, r)
					ok = false
				}
			}()
			data := binary.LittleEndian.AppendUint32(nil, w)
			ins, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(data))
			if err != nil {
				t.Errorf("parse %#08x: %v", w, err)
				ok = false
				return
			}

			// the stream rule must walk the decoded instructions inside
			// the 4 input bytes, never past the buffer (render errors are
			// acceptable - garbage words - only the panic check matters)
			off := 0
			for _, in := range ins {
				_ = in.ObjDump(disasm.DefaultViewCtx())
				off += riscv.InstrLen(data[off:])
				if off > len(data) {
					t.Logf("word %#08x: stream overruns at %d", w, off)
					ok = false
					break
				}
			}

			// the tail may be a truncated 32-bit form - arbitrary bytes are
			// not required to parse entirely; the requirement is not to crash
		}()
		return ok
	}, ohsnap.CheckOptions{ProgressEvery: checkProgressEvery, LogShrinkSteps: true})
}

// TestPropertyRiscvVsObjdump - the differential: words produced by the riscv
// constructors are disassembled identically by us and by objdump.
func TestPropertyRiscvVsObjdump(t *testing.T) {
	const perFamily = 48
	const threshold = 90.0

	rnd := seedRnd(t)
	gens := rvPropFamilies(rnd)
	ins := make([]riscv.Instr, 0, perFamily*len(gens))
	for _, gen := range gens {
		for range perFamily {
			ins = append(ins, gen())
		}
	}

	raw, ok := rvEncodeAll(t, ins)
	require.True(t, ok)
	code := raw

	path := writeRiscvELF(t, code)
	out, err := objdump.Run(context.Background(), objdump.Args("ELF", 0, path))
	if err != nil {
		t.Skipf("no objdump for ELF/riscv64: %v", err)
	}

	objLines := objdump.ParseByAddr(string(out))
	require.NotEmpty(t, objLines)

	style := text.StyleFor("ELF")
	opts := disasm.NewOptions(style)
	matched, mismatched, notInOurs := 0, 0, 0
	var samples []string
	for addr, objLine := range objLines {
		off := int(addr) - propAddr
		if off < 0 || off >= len(code) {
			notInOurs++
			continue
		}

		ins, err := riscv.MakeDecoder()(parsec.Stateless{}, parsecbytes.Buffer(code[off:]))
		if err != nil {
			notInOurs++
			continue
		}

		if len(ins) == 0 {
			notInOurs++
			continue
		}

		ourLine := objdump.StripComments(objdump.Normalize(
			disasm.Line(addr, code[off:off+riscv.InstrLen(code[off:])], ins[0], opts)))
		if ourLine == objdump.StripComments(objLine) {
			matched++
		} else {
			mismatched++
			if len(samples) < 10 {
				samples = append(samples, fmt.Sprintf(
					"0x%x\n    ours:    %s\n    objdump: %s",
					addr,
					ourLine,
					objdump.StripComments(objLine),
				))
			}
		}
	}

	total := matched + mismatched
	pct := 0.0
	if total > 0 {
		pct = float64(matched) * 100 / float64(total)
	}

	t.Logf("vs objdump: %d matched (%.2f%%), %d mismatched, %d foreign addresses",
		matched, pct, mismatched, notInOurs)
	for _, sm := range samples {
		t.Log(sm)
	}

	require.GreaterOrEqual(t, pct, threshold)
}

// writeRiscvELF - minimal ELF64 LE riscv (e_machine = EM_RISCV,
// e_flags = RVC, as in tests/examples/hello-riscv): .text @ 0x1000.
func writeRiscvELF(t *testing.T, code []byte) string {
	t.Helper()
	path, err := writeObjELF(t, code, 0xF3, 0x5)
	require.NoError(t, err)
	return path
}
