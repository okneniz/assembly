package arm64

// Encoding with self-verify: negotiation of constructors (injected
// aliases -> the arch's armCtors) and legacy paths (format-family
// handlers via the decode table); each candidate is encoded and decoded
// back - the decoder must reproduce the source text (looseNormalize is
// insensitive to the number format).

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ctx is the internal orchestration context of encodeARM: the
// evaluation environment (address + resolver, like the core's asm.Ctx)
// plus the injected alias constructors (extra; nil for the pure
// backend).
type ctx struct {
	Addr    uint64
	Resolve func(name string) (uint64, bool)
	extra   map[string]arch.ArmCtor
}

// encodeARM picks a candidate scheme, encodes, and verifies with the
// decoder.
func encodeARM(in armAsmInstr, ctx ctx) (uint32, error) {
	// expression slots -> numbers BEFORE dispatch: constructors and
	// legacy handlers consume only evaluated operands (instruction
	// structures contain no holes); the placeholder pass yields rel=0.
	ops, rerr := resolveOps(in.mnem, in.ops, ctx)
	if rerr != nil {
		return 0, fmt.Errorf("%s: %w", in.mnem, rerr)
	}

	res := resolvedInstr{mnem: in.mnem, ops: ops}
	rendered := renderInstr(res)
	loose := looseNormalize(rendered)

	// objdump style: an arrangement suffix on the mnemonic (orr.16b
	// v6, v3, v4). SIMD ctors expect it on the first operand - move it
	// onto A COPY of ops (the backing array is shared between the
	// layout/encoding passes: a mutation would poison the repeated
	// render+verify). The mnemonic keeps the suffix (armCtors keys like
	// "orr.16b"). Conditions: the first operand is a bare v-register
	// without its own suffix, and the suffix is a valid arrangement
	// (b.eq / fadd.d / ld1.16b {...} are not affected).
	if dot := strings.IndexByte(res.mnem, '.'); dot >= 0 && len(res.ops) > 0 {
		if op := res.ops[0]; op.IsReg() && op.Arr() == "" {
			if _, _, err := arch.ArrQSize(res.mnem[dot+1:]); err == nil {
				vops := append([]arch.VOp(nil), res.ops...)
				vops[0] = arch.VOpReg(op.Reg(), res.mnem[dot+1:], op.LaneIdx(), op.Num())
				res.ops = vops
			}
		}
	}

	// Per-instruction structures: the constructor builds the structure,
	// it encodes itself and passes self-verify (verifyWord; the
	// exceptions are the SkipVerify marker, see there). First the
	// injected constructors (aliases), then the own ones
	// (arch.BuildInstr), then the legacy path via the decode table
	// below.
	if ctor, ok := ctx.extra[res.mnem]; ok {
		st, cerr := ctor(res.ops)
		if cerr == nil {
			if w, ok2 := verifyWord(st, ctx.Addr, loose); ok2 {
				return w, nil
			}
		}
	}

	if st, cerr := arch.BuildInstr(res.mnem, res.ops); cerr == nil {
		if w, ok := verifyWord(st, ctx.Addr, loose); ok {
			return w, nil
		}
	}

	var lastErr error
	encoded := false
	var firstWord uint32
	for _, cand := range arch.LegacyCandidates(res.mnem) {
		w, err := arch.BuildLegacy(cand, res.mnem, res.ops, ctx.Addr)
		if err != nil {
			if lastErr == nil {
				lastErr = err
			}

			continue
		}

		if !encoded {
			encoded, firstWord = true, w
		}

		// verify: our decoder reproduces the source text
		decText := instrTextOf(arch.DecodeWord(w, ctx.Addr))
		if decText != "" && looseNormalize(decText) == loose {
			return w, nil
		}
	}

	if encoded {
		return firstWord, nil // the only successful encoding (the text may have been an alias variant)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("unknown mnemonic %q", res.mnem)
	}

	return 0, fmt.Errorf("%q: %w", rendered, lastErr)
}

// resolvedInstr is an instruction with evaluated operands (arch.VOp):
// the input of the self-verify render.
type resolvedInstr struct {
	mnem string
	ops  []arch.VOp
}

// verifyWord is encoding + self-verify (the SkipVerify marker is for
// those without a decoding scheme or with a keyword operand).
func verifyWord(st arch.Instr, addr uint64, loose string) (uint32, bool) {
	var sbuf bytes.Buffer
	if _, werr := st.Encode(&sbuf, addr); werr != nil || sbuf.Len() != 4 {
		return 0, false
	}

	w := binary.LittleEndian.Uint32(sbuf.Bytes())
	if _, skip := st.(interface{ SkipVerify() }); skip ||
		looseNormalize(instrTextOf(arch.DecodeWord(w, addr))) == loose {
		return w, true
	}

	return 0, false
}

// instrTextOf returns the mnemonic+operands of the decoded instruction -
// its own ObjDump text (compared with renderInstr via looseNormalize
// during self-verify).

func instrTextOf(inst arch.Instr) string {
	return inst.ObjDump(disasm.DefaultViewCtx())
}

// renderInstr renders the evaluated operands for self-verify: values
// are printed as numbers (the decoded text contains numbers).
func renderInstr(in resolvedInstr) string {
	var b strings.Builder
	b.WriteString(in.mnem)
	for _, op := range in.ops {
		b.WriteByte(' ')
		b.WriteString(renderVOp(op))
	}

	return b.String()
}

func renderVOp(op arch.VOp) string {
	switch {
	case op.IsReg():
		rs := op.Reg()
		if op.Arr() != "" {
			rs += "." + op.Arr()
		}

		if op.LaneIdx() { // lane suffix index: "v30[1]" (INS/elementwise)
			rs += "[" + strconv.FormatInt(op.Num(), 10) + "]"
		}

		return rs
	case op.IsImm():
		if op.Sym() != "" { // name operand: condition/sysreg/prfm hint
			return op.Sym()
		}

		return "#" + strconv.FormatInt(op.Num(), 10)
	case op.IsLit():
		return "=" + strconv.FormatInt(op.Num(), 10) // literal pool (does not pass verify)
	case op.IsFloat():
		return "#" + strconv.FormatFloat(op.Float(), 'f', 8, 64)
	case op.IsShift():
		return op.ShiftName() + " #" + strconv.FormatInt(op.Num(), 10)
	case op.IsExtend():
		if op.HasAmt() {
			return op.ShiftName() + " #" + strconv.FormatInt(op.Num(), 10)
		}

		return op.ShiftName()
	case op.IsMem():
		m := op.Mem()
		var b strings.Builder
		b.WriteString("[" + m.Base())
		switch {
		case m.HasOff():
			b.WriteString(", #" + strconv.FormatInt(m.Off(), 10))
		case m.OffReg() != "":
			b.WriteString(", " + m.OffReg())
			if m.Opt() != "" {
				b.WriteString(", " + m.Opt())
				if m.HasOpt() {
					b.WriteString(" #" + strconv.FormatInt(m.OptAmt(), 10))
				}
			}
		}

		b.WriteString("]")
		if m.Pre() {
			b.WriteString("!")
		}

		if m.HasPost() {
			b.WriteString(", #" + strconv.FormatInt(m.Post(), 10))
		}

		return b.String()
	case op.IsList():
		parts := make([]string, 0, len(op.List()))
		for _, r := range op.List() {
			s := r.Reg()
			if r.Arr() != "" {
				s += "." + r.Arr()
			}

			parts = append(parts, s)
		}

		return "{" + strings.Join(parts, ", ") + "}"
	}

	return "?"
}

// looseNormalize is a comparison insensitive to the number format: all
// numeric tokens are canonicalized to decimal ("#0x2a" == "#42" ==
// "#2a").
func looseNormalize(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		if c == ' ' || c == '\t' || c == ',' || c == '#' {
			i++
			continue
		}

		var j int
		if c == '0' && i+1 < len(s) && (s[i+1] == 'x' || s[i+1] == 'X') {
			j = i + 2
			for j < len(s) && isHexByte(s[j]) {
				j++
			}

			if v, err := strconv.ParseUint(s[i+2:j], 16, 64); err == nil && j > i+2 {
				fmt.Fprintf(&b, "%d", v)
				i = j
				continue
			}
		}

		if c >= '0' && c <= '9' {
			j = i
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}

			b.WriteString(s[i:j])
			i = j
			continue
		}

		b.WriteByte(lowerByte(c))
		i++
	}

	return b.String()
}

func lowerByte(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}

	return c
}

func isHexByte(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
