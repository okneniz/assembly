package arm64

import (
	"fmt"
	"io"
)

// memKind — the kind of load/store immediate addressing.
type memKind int

const (
	memImm      memKind = iota // [rn, #imm12<<scale]
	memUnscaled                // [rn, #±imm9] (ldur/stur — separate mnemonics)
	memPost                    // [rn], #±imm9
	memPre                     // [rn, #±imm9]!
	memLiteral                 // ldr rt, =addr (imm19*4 from pc)
	memRegOff                  // [rn, rm{, ext #scale}]
)

// lsBase — the common record of load/store fields (immediate + register offset).
type lsBase struct {
	rt       string
	rn       string
	kind     memKind
	off      int64  // byte offset (imm12 already multiplied by the scale)
	tgt      uint64 // absolute address (literal)
	enc      uint32 // base word of the matched entry
	rm       string // index register (memRegOff)
	option   string // extension: lsl/uxtw/sxtw/sxtx ("" — none)
	shiftAmt uint32 // extension scale (0 — not printed)
}

func newLsBase(
	rt string,
	rn string,
	kind memKind,
	off int64,
	tgt uint64,
	enc uint32,
	rm string,
	option string,
	shiftAmt uint32,
) lsBase {
	return lsBase{
		rt:       rt,
		rn:       rn,
		kind:     kind,
		off:      off,
		tgt:      tgt,
		enc:      enc,
		rm:       rm,
		option:   option,
		shiftAmt: shiftAmt,
	}
}

// lsText — operands by addressing kind (objdump style).
func (m lsBase) lsText() string {
	switch m.kind {
	case memImm:
		if m.off == 0 {
			return fmt.Sprintf("[%s]", m.rn)
		}

		return fmt.Sprintf("[%s, #0x%x]", m.rn, m.off)
	case memUnscaled:
		if m.off == 0 {
			return fmt.Sprintf("[%s]", m.rn)
		}

		return fmt.Sprintf("[%s, #%#x]", m.rn, m.off)
	case memPost:
		return fmt.Sprintf("[%s], #%#x", m.rn, m.off)
	case memPre:
		return fmt.Sprintf("[%s, #%#x]!", m.rn, m.off)
	case memRegOff:
		if m.option == "" {
			return fmt.Sprintf("[%s, %s]", m.rn, m.rm)
		}

		if m.shiftAmt == 0 {
			return fmt.Sprintf("[%s, %s, %s]", m.rn, m.rm, m.option)
		}

		return fmt.Sprintf("[%s, %s, %s #%d]", m.rn, m.rm, m.option, m.shiftAmt)
	case memLiteral:
		// literal: llvm prints the absolute for near targets (symbol/Mach-O
		// location) and #offset for far ones — the heuristic is not
		// reproducible without section info; we print the absolute
		return fmt.Sprintf("0x%x", m.tgt)
	default:
		// all memKind values are listed above; defensive branch
		return fmt.Sprintf("[%s]", m.rn)
	}
}

// lsWrite — the encoding word by addressing kind (computed form: the
// literal target is already a number in tgt; pc — the instruction address).
func (m lsBase) lsWrite(w io.Writer, pc uint64, mnem string) (int64, error) {
	rt, err := armRegNum(m.rt)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	// the literal form has no base register (rn is empty)
	rn := uint32(0)
	if m.kind != memLiteral {
		rnNum, rnErr := armRegNum(m.rn)
		if rnErr != nil {
			return 0, fmt.Errorf("%s: %w", mnem, rnErr)
		}

		rn = rnNum
	}

	switch m.kind {
	case memImm:
		scale := m.enc >> 30 & 3
		if m.off < 0 || m.off&(int64(1)<<scale-1) != 0 || m.off>>scale > 0xfff {
			return 0, fmt.Errorf("%s: offset %#x out of imm12 range", mnem, m.off)
		}

		return writeWord(w, m.enc|rt|rn<<5|uint32(m.off>>scale)<<10)
	case memUnscaled, memPost, memPre:
		if m.off < -256 || m.off > 255 {
			return 0, fmt.Errorf("%s: offset %#x out of imm9 range", mnem, m.off)
		}

		return writeWord(w, m.enc|rt|rn<<5|uint32(m.off&0x1ff)<<12)
	case memRegOff:
		rm, err := armRegNum(m.rm)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", mnem, err)
		}

		opt := uint32(0b011) // plain [rn, rm] — this is lsl with S=0
		if m.option != "" {
			opt, err = lsOptNum(m.option)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", mnem, err)
			}
		}

		sh := uint32(0)
		if m.shiftAmt != 0 {
			sh = 1
		}

		return writeWord(w, m.enc|rt|rn<<5|sh<<12|opt<<13|rm<<16)
	case memLiteral: // literal
		rel := int64(m.tgt) - int64(pc)
		if rel%4 != 0 || rel < -(1<<20) || rel >= 1<<20 {
			return 0, fmt.Errorf("%s: literal out of range", mnem)
		}

		return writeWord(w, m.enc|rt|uint32(rel>>2&0x7ffff)<<5)
	}

	// all memKind values are listed above; defensive return for the compiler
	return 0, fmt.Errorf("%s: unsupported addressing kind", mnem)
}

// pairBase — the common fields of ldp/stp/ldpsw.
type pairBase struct {
	rt, rt2, rn string
	kind        memKind // memImm | memPost | memPre
	off         int64   // already multiplied by the scale
	scale       uint32  // 2 (32-bit pairs) | 3 (64-bit)
	enc         uint32
}

func newPairBase(
	rt string,
	rt2 string,
	rn string,
	kind memKind,
	off int64,
	scale uint32,
	enc uint32,
) pairBase {
	return pairBase{
		rt:    rt,
		rt2:   rt2,
		rn:    rn,
		kind:  kind,
		off:   off,
		scale: scale,
		enc:   enc,
	}
}

func (p pairBase) pairText() string {
	switch p.kind {
	case memPost:
		return fmt.Sprintf("%s, %s, [%s], #%#x", p.rt, p.rt2, p.rn, p.off)
	case memPre:
		return fmt.Sprintf("%s, %s, [%s, #%#x]!", p.rt, p.rt2, p.rn, p.off)
	// pairs only ever have memImm; the other addressing kinds never reach
	// here, the explicit list is for switch completeness (they behave as a
	// signed offset)
	case memImm, memUnscaled, memLiteral, memRegOff:
	default:
		// non-empty memKind covered above; defensive branch for the compiler
	}

	if p.off == 0 {
		return fmt.Sprintf("%s, %s, [%s]", p.rt, p.rt2, p.rn)
	}

	return fmt.Sprintf("%s, %s, [%s, #%#x]", p.rt, p.rt2, p.rn, p.off)
}

func (p pairBase) pairWrite(w io.Writer, mnem string) (int64, error) {
	rt, err := armRegNum(p.rt)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	rt2, err := armRegNum(p.rt2)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	rn, err := armRegNum(p.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	if p.off&(int64(1)<<p.scale-1) != 0 || p.off>>p.scale < -64 || p.off>>p.scale > 63 {
		return 0, fmt.Errorf("%s: offset %#x out of imm7 range", mnem, p.off)
	}

	return writeWord(w, p.enc|rt|rn<<5|rt2<<10|uint32(p.off>>p.scale&0x7f)<<15)
}

// atomic — ldar/stlr/ldarb/stlrb/ldaxr/ldaxrb: rt, [rn].
type atomic struct {
	rt, rn string
}

func newAtomic(rt string, rn string) atomic {
	return atomic{
		rt: rt,
		rn: rn,
	}
}

func (a atomic) atText() string {
	return fmt.Sprintf("%s, [%s]", a.rt, a.rn)
}

func (a atomic) atWrite(w io.Writer, enc uint32, mnem string) (int64, error) {
	rt, rn, err := regNums2(a.rt, a.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	return writeWord(w, enc|rt|rn<<5)
}

// excl — stlxr/stlxrb/stxrb: rs, rt, [rn].
type excl struct {
	rs, rt, rn string
}

func newExcl(rs string, rt string, rn string) excl {
	return excl{
		rs: rs,
		rt: rt,
		rn: rn,
	}
}

func (e excl) exText() string {
	return fmt.Sprintf("%s, %s, [%s]", e.rs, e.rt, e.rn)
}

func (e excl) exWrite(w io.Writer, enc uint32, mnem string) (int64, error) {
	rs, err := armRegNum(e.rs)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	rt, rn, err := regNums2(e.rt, e.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	return writeWord(w, enc|rt|rn<<5|rs<<16)
}

// lsOptName — the load/store extension option name (like the lsOpt transform).
func lsOptName(v uint32) string {
	switch v {
	case 0b010:
		return "uxtw"
	case 0b011:
		return "lsl"
	case 0b110:
		return "sxtw"
	case 0b111:
		return "sxtx"
	}

	return "lsl"
}

// lsOptNum — the load/store extension option number (inverse of lsOptName).
func lsOptNum(name string) (uint32, error) {
	switch name {
	case "uxtw":
		return 0b010, nil
	case "lsl":
		return 0b011, nil
	case "sxtw":
		return 0b110, nil
	case "sxtx":
		return 0b111, nil
	}

	return 0, fmt.Errorf("unknown option %q", name)
}

// pairDecode — pair fields: idx (bits 24:23) selects the addressing kind, L — load.
// rtKind: "x"/"w"/"d" — the register-pair type.
func pairDecode(
	w uint32,
	scale uint32,
	rtKind string,
) (string, string, string, memKind, int64, bool) {
	rt := pairRegName(w&0x1f, rtKind)
	rt2 := pairRegName(w>>10&0x1f, rtKind)
	rn := regNameXSP(w >> 5 & 0x1f)
	kind := memImm
	switch w >> 23 & 3 {
	case 0b01:
		kind = memPost
	case 0b11:
		kind = memPre
	}

	return rt, rt2, rn, kind, signExtendN(w>>15&0x7f, 7) << scale, w>>22&1 == 1
}

// pairRegName — the pair register name by type ("x"/"w"/"d").
func pairRegName(n uint32, kind string) string {
	switch kind {
	case "d":
		return fpRegNameD(n)
	case "w":
		return regNameW(n)
	default:
		return regNameX(n)
	}
}

// lsSignedWrite — the signed-load word (imm12, scale from enc bits 31:30).
func lsSignedWrite(w io.Writer, enc uint32, rt, rn string, off int64, mnem string) (int64, error) {
	r, err := armRegNum(rt)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	n, err := armRegNum(rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	scale := enc >> 30 & 3
	if off < 0 || off&(int64(1)<<scale-1) != 0 || off>>scale > 0xfff {
		return 0, fmt.Errorf("%s: offset out of range", mnem)
	}

	return writeWord(w, enc|r|n<<5|uint32(off>>scale)<<10)
}
