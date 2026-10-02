package arm64

import (
	"fmt"
	"io"
	"strings"

	"github.com/okneniz/assembly/disasm"
)

// The IC/DC/TLBI/AT system operations: sys instructions whose first
// operand is an operation keyword (ic iallu, dc zva, x0, tlbi alle1,
// at s1e2r, x0). The word formula is 0xd5080000 | op1<<16 | CRn<<12 |
// CRm<<8 | op2<<5 | Rt - the class (ic/dc/at vs tlbi) is CRn 7 vs 8,
// the operation is the op1/CRm/op2 triple. Every spelling below pins
// its full word (clang-oracled); the Rt-bearing forms keep the base
// with Rt=0 and join the register at assembly, the operandless ones
// encode Rt=31.

// sysOpRec - one operation spelling: its mnemonic class, its word (the
// Rt forms carry Rt=0) and whether an Xt operand joins the word.
type sysOpRec struct {
	mnem  string
	enc   uint32
	hasRt bool
}

// sysOps - the operation spellings, lower-case keys (the lookup is
// case-insensitive, as the sysreg one).
var sysOps = map[string]sysOpRec{
	"iallu":        {mnem: "ic", enc: 0xD508751F},
	"ialluis":      {mnem: "ic", enc: 0xD508711F},
	"ivau":         {mnem: "ic", enc: 0xD50B7520, hasRt: true},
	"zva":          {mnem: "dc", enc: 0xD50B7420, hasRt: true},
	"ivac":         {mnem: "dc", enc: 0xD5087620, hasRt: true},
	"isw":          {mnem: "dc", enc: 0xD5087640, hasRt: true},
	"cvac":         {mnem: "dc", enc: 0xD50B7A20, hasRt: true},
	"cvau":         {mnem: "dc", enc: 0xD50B7B20, hasRt: true},
	"civac":        {mnem: "dc", enc: 0xD50B7E20, hasRt: true},
	"csw":          {mnem: "dc", enc: 0xD5087A40, hasRt: true},
	"cisw":         {mnem: "dc", enc: 0xD5087E40, hasRt: true},
	"cvadp":        {mnem: "dc", enc: 0xD50B7D20, hasRt: true},
	"alle1":        {mnem: "tlbi", enc: 0xD50C879F},
	"alle1is":      {mnem: "tlbi", enc: 0xD50C839F},
	"alle1os":      {mnem: "tlbi", enc: 0xD50C819F},
	"alle2":        {mnem: "tlbi", enc: 0xD50C871F},
	"alle2is":      {mnem: "tlbi", enc: 0xD50C831F},
	"alle2os":      {mnem: "tlbi", enc: 0xD50C811F},
	"alle3":        {mnem: "tlbi", enc: 0xD50E871F},
	"alle3is":      {mnem: "tlbi", enc: 0xD50E831F},
	"alle3os":      {mnem: "tlbi", enc: 0xD50E811F},
	"vmalle1":      {mnem: "tlbi", enc: 0xD508871F},
	"vmalle1is":    {mnem: "tlbi", enc: 0xD508831F},
	"vmalle1os":    {mnem: "tlbi", enc: 0xD508811F},
	"vmalls12e1":   {mnem: "tlbi", enc: 0xD50C87DF},
	"vmalls12e1is": {mnem: "tlbi", enc: 0xD50C83DF},
	"vmalls12e1os": {mnem: "tlbi", enc: 0xD50C81DF},
	"vae1":         {mnem: "tlbi", enc: 0xD5088720, hasRt: true},
	"vae1is":       {mnem: "tlbi", enc: 0xD5088320, hasRt: true},
	"vae1os":       {mnem: "tlbi", enc: 0xD5088120, hasRt: true},
	"rvae1":        {mnem: "tlbi", enc: 0xD5088620, hasRt: true},
	"rvae1is":      {mnem: "tlbi", enc: 0xD5088220, hasRt: true},
	"rvae1os":      {mnem: "tlbi", enc: 0xD5088520, hasRt: true},
	"aside1":       {mnem: "tlbi", enc: 0xD5088740, hasRt: true},
	"aside1is":     {mnem: "tlbi", enc: 0xD5088340, hasRt: true},
	"aside1os":     {mnem: "tlbi", enc: 0xD5088140, hasRt: true},
	"vae2":         {mnem: "tlbi", enc: 0xD50C8720, hasRt: true},
	"vae2is":       {mnem: "tlbi", enc: 0xD50C8320, hasRt: true},
	"vae2os":       {mnem: "tlbi", enc: 0xD50C8120, hasRt: true},
	"rvae2":        {mnem: "tlbi", enc: 0xD50C8620, hasRt: true},
	"rvae2is":      {mnem: "tlbi", enc: 0xD50C8220, hasRt: true},
	"rvae2os":      {mnem: "tlbi", enc: 0xD50C8520, hasRt: true},
	"vae3":         {mnem: "tlbi", enc: 0xD50E8720, hasRt: true},
	"vae3is":       {mnem: "tlbi", enc: 0xD50E8320, hasRt: true},
	"vae3os":       {mnem: "tlbi", enc: 0xD50E8120, hasRt: true},
	"rvae3":        {mnem: "tlbi", enc: 0xD50E8620, hasRt: true},
	"rvae3is":      {mnem: "tlbi", enc: 0xD50E8220, hasRt: true},
	"rvae3os":      {mnem: "tlbi", enc: 0xD50E8520, hasRt: true},
	"ipas2e1":      {mnem: "tlbi", enc: 0xD50C8420, hasRt: true},
	"ipas2e1is":    {mnem: "tlbi", enc: 0xD50C8020, hasRt: true},
	"ipas2e1os":    {mnem: "tlbi", enc: 0xD50C8400, hasRt: true},
	"ripas2e1":     {mnem: "tlbi", enc: 0xD50C8440, hasRt: true},
	"ripas2e1is":   {mnem: "tlbi", enc: 0xD50C8040, hasRt: true},
	"ripas2e1os":   {mnem: "tlbi", enc: 0xD50C8460, hasRt: true},
	"s1e1r":        {mnem: "at", enc: 0xD5087800, hasRt: true},
	"s1e1w":        {mnem: "at", enc: 0xD5087820, hasRt: true},
	"s1e0r":        {mnem: "at", enc: 0xD5087840, hasRt: true},
	"s1e0w":        {mnem: "at", enc: 0xD5087860, hasRt: true},
	"s1e1rp":       {mnem: "at", enc: 0xD5087900, hasRt: true},
	"s1e1wp":       {mnem: "at", enc: 0xD5087920, hasRt: true},
	"s12e1r":       {mnem: "at", enc: 0xD50C7880, hasRt: true},
	"s12e1w":       {mnem: "at", enc: 0xD50C78A0, hasRt: true},
	"s12e0r":       {mnem: "at", enc: 0xD50C78C0, hasRt: true},
	"s1e2r":        {mnem: "at", enc: 0xD50C7800, hasRt: true},
	"s1e2w":        {mnem: "at", enc: 0xD50C7820, hasRt: true},
	"s1e3r":        {mnem: "at", enc: 0xD50E7800, hasRt: true},
	"s1e3w":        {mnem: "at", enc: 0xD50E7820, hasRt: true},
}

// sysOpLookup - the operation spelling, case-insensitive.
func sysOpLookup(sym string) (sysOpRec, bool) {
	r, ok := sysOps[strings.ToLower(sym)]
	return r, ok
}

// sysOpOfWord - the operation a sys-class word carries (decode side):
// an Rt-bearing base matches with any Rt, an operandless word matches
// exactly.
func sysOpOfWord(w uint32) (sysOpRec, uint32, bool) {
	for _, r := range sysOps {
		if r.hasRt {
			if w&^uint32(0x1F) == r.enc {
				return r, w & 0x1F, true
			}

			continue
		}

		if w == r.enc {
			return r, 31, true
		}
	}

	return sysOpRec{}, 0, false
}

// SysOp - one IC/DC/TLBI operation: the keyword alone (ic iallu,
// tlbi alle1) or with the Xt operand (dc zva, x0, tlbi vae1, x3).
type SysOp struct {
	mnem  string
	op    string
	rt    uint32
	hasRt bool
	enc   uint32
}

// newSysOp - the SysOp constructor: the spelling must belong to the
// mnemonic's class and the register must come exactly when the
// operation takes one (an X-register: 31 reads as xzr).
func newSysOp(mnem, op string, rt Reg) (SysOp, error) {
	rec, ok := sysOpLookup(op)
	if !ok || rec.mnem != mnem {
		return SysOp{}, fmt.Errorf("%s: unknown operation %q", mnem, op)
	}

	enc := rec.enc
	if rec.hasRt {
		if err := requireClass(
			rt,
			"SysOp",
			"rt",
			"a 64-bit register expected (w-regs and sp do not fit)",
			classX,
			classXZR,
		); err != nil {
			return SysOp{}, err
		}

		enc |= rt.bits()
	}

	return SysOp{
		mnem:  mnem,
		op:    strings.ToLower(op),
		rt:    rt.bits(),
		hasRt: rec.hasRt,
		enc:   enc,
	}, nil
}

// SysOpOf - the system operation by its spelling; rt is "" for the
// operandless forms, the register name for the Rt-bearing ones.
func SysOpOf(mnem, op, rt string) (SysOp, error) {
	rec, ok := sysOpLookup(op)
	if !ok || rec.mnem != mnem {
		return SysOp{}, fmt.Errorf("%s: unknown operation %q", mnem, op)
	}

	if !rec.hasRt && rt != "" {
		return SysOp{}, fmt.Errorf("%s %s: takes no register", mnem, op)
	}

	var r Reg
	if rec.hasRt {
		if rt == "" {
			return SysOp{}, fmt.Errorf("%s %s: want register", mnem, op)
		}

		var err error
		if r, err = RegOf(rt); err != nil {
			return SysOp{}, err
		}
	}

	return newSysOp(mnem, op, r)
}

func (i SysOp) ObjDump(_ disasm.ViewCtx) string {
	if !i.hasRt {
		return i.mnem + " " + i.op
	}

	return fmt.Sprintf("%s %s, %s", i.mnem, i.op, regNameX(i.rt))
}

func (i SysOp) Encode(w io.Writer) (int64, error) {
	return writeWord(w, i.enc)
}

// sysClassFields - the honest field layout of a sys-class word (the
// Generic print of a word the operation table does not know).
var sysClassFields = []Field{
	NewField("op1", 16, 3),
	NewField("CRm", 8, 4),
	NewField("op2", 5, 3),
	NewField("Rt", 0, 5),
}

// decodeSysClass - the sys-class decode (IC/DC at CRn 7, TLBI at CRn 8):
// a word the operation table knows prints its spelling, the rest keeps
// the honest raw fields (the class entry's own print, no regression).
func decodeSysClass(mnem string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		if rec, rt, ok := sysOpOfWord(w); ok {
			return SysOp{
				mnem:  rec.mnem,
				op:    sysOpName(rec),
				rt:    rt,
				hasRt: rec.hasRt,
				enc:   w,
			}, nil
		}

		return Generic{
			name:   mnem,
			fields: sysClassFields,
			word:   w,
		}, nil
	}
}

// sysOpName - the spelling of a table record (the map's reverse).
func sysOpName(rec sysOpRec) string {
	for name, r := range sysOps {
		if r == rec {
			return name
		}
	}

	return "?"
}
