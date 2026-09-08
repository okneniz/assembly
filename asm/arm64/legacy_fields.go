package arm64

// Schema candidates by mnemonic (base names + curated alias map mirroring
// the formatters' pseudo-logic) and format-family handlers: operands →
// schema field values. The verify loop of encodeARM (decoding with our own
// decoder) confirms the chosen candidate is correct.

import (
	"errors"
	"fmt"
	"strings"
)

import (
	"math"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// armAliasMap — source mnemonic → base (Meta.Name, Formatter) pairs that can
// produce it. The candidate order repeats the arm64Schemas order — as decoder
// priority.
var armAliasMap = map[string][][2]string{
	"mov": {
		{"orr", "logicalShifted"}, // register form: GAS/clang emit ORR (byte-parity)
		{"add", "addSubImm"},
		{"movz", "movWide"}, {"movn", "movWide"},
		{"orr", "logicalImm"}, // bitmask-xzr form — only when MOVZ/MOVN cannot
	},
	"movz": {{"movz", "movWide"}},
	"movn": {{"movn", "movWide"}},
	"movk": {{"movk", "movWide"}},
	"cmp":  {{"subs", "addSubImm"}, {"subs", "shiftedReg"}, {"subs", "addSubExt"}},
	"cmn":  {{"adds", "addSubImm"}, {"adds", "shiftedReg"}, {"adds", "addSubExt"}},
	"neg":  {{"sub", "shiftedReg"}},
	"negs": {{"subs", "shiftedReg"}},
	"tst":  {{"ands", "logicalImm"}, {"ands", "logicalShifted"}},
	"mvn":  {{"orn", "logicalShifted"}},
	"and":  {{"and", "simd3"}, {"and", "logicalShifted"}, {"and", "logicalImm"}},
	"ands": {{"ands", "logicalShifted"}, {"ands", "logicalImm"}},
	"bic":  {{"bic", "logicalShifted"}},
	"bics": {{"bics", "logicalShifted"}},
	"orr":  {{"orr", "simd3"}, {"orr", "logicalShifted"}, {"orr", "logicalImm"}},
	"orn":  {{"orn", "logicalShifted"}},
	"eor":  {{"eor", "simd3"}, {"eor", "logicalShifted"}, {"eor", "logicalImm"}},
	"eon":  {{"eon", "logicalShifted"}},
	// order: immediate form (UBFM/SBFM — as GAS emits it), then register
	// form (LSLV etc., schemas named "lsl"/"lsr"/...)
	"lsl":   {{"ubfm", "bfmAlias"}, {"lsl", "op3"}, {"sbfm", "bfmAlias"}},
	"lsr":   {{"ubfm", "bfmAlias"}, {"lsr", "op3"}},
	"asr":   {{"sbfm", "bfmAlias"}, {"asr", "op3"}},
	"ror":   {{"extr", "extrFmt"}, {"ror", "op3"}},
	"sxtb":  {{"sbfm", "bfmAlias"}},
	"sxth":  {{"sbfm", "bfmAlias"}},
	"sxtw":  {{"sbfm", "bfmAlias"}},
	"uxtb":  {{"ubfm", "bfmAlias"}},
	"uxth":  {{"ubfm", "bfmAlias"}},
	"ubfiz": {{"ubfm", "bfmAlias"}},
	"ubfx":  {{"ubfm", "bfmAlias"}},
	"sbfiz": {{"sbfm", "bfmAlias"}},
	"sbfx":  {{"sbfm", "bfmAlias"}},
	"ubfm":  {{"ubfm", "bfmAlias"}},
	"sbfm":  {{"sbfm", "bfmAlias"}},
	"bfm":   {{"bfm", "bfmAlias"}},
	"cset":  {{"csinc", "csel"}},
	"csetm": {{"csinv", "csel"}},
	"cinc":  {{"csinc", "csel"}},
	"cinv":  {{"csinv", "csel"}},
	"cneg":  {{"csneg", "csel"}},
	"csel":  {{"csel", "csel"}},
	"csinc": {{"csinc", "csel"}},
	"csinv": {{"csinv", "csel"}},
	"csneg": {{"csneg", "csel"}},
	"mul":   {{"madd", "madd"}},
	"mneg":  {{"msub", "madd"}},
	"madd":  {{"madd", "madd"}},
	"msub":  {{"msub", "madd"}},
	"nop":   {{"nop", "operandsNone"}},
	// structural load/store: the name is chosen from the ld1 schema by opcode/L
	"ld1":     {{"ld1", "ldStruct"}},
	"ld2":     {{"ld1", "ldStruct"}},
	"ld3":     {{"ld1", "ldStruct"}},
	"ld4":     {{"ld1", "ldStruct"}},
	"st1":     {{"ld1", "ldStruct"}},
	"st2":     {{"ld1", "ldStruct"}},
	"st3":     {{"ld1", "ldStruct"}},
	"st4":     {{"ld1", "ldStruct"}},
	"ld1r":    {{"ld1", "ldStruct"}},
	"ld2r":    {{"ld1", "ldStruct"}},
	"ld3r":    {{"ld1", "ldStruct"}},
	"st1r":    {{"ld1", "ldStruct"}},
	"st4r":    {{"ld1", "ldStruct"}},
	"ld4r":    {{"ld1", "ldStruct"}},
	"mov.d":   {{"fmov", "movElem"}},
	"mov.16b": {{"orr", "simd3"}},
	"mov.8b":  {{"orr", "simd3"}},
	"ldr": {{"ldr", "ldrLiteral"}, {"ldr", "lsImm"}, {"ldr", "lsUnscaled"},
		{"ldr", "lsRegOff"}, {"ldr", "lsSigned"}, {"ldr", "lsFpRegOff"},
		{"ldr", "lsPreIndex"}, {"ldr", "lsPostIndex"}},
	"ldur": {{"ldur", "lsUnscaled"}},
	"str": {{"str", "lsImm"}, {"str", "lsUnscaled"}, {"str", "lsRegOff"}, {"str", "lsFpRegOff"},
		{"str", "lsPreIndex"}, {"str", "lsPostIndex"}},
	"stur": {{"stur", "lsUnscaled"}},
	"ldp":  {{"ldp", "lsPair"}, {"ldp", "lsPairPP"}},
	"stp": {
		{"ldp", "lsPair"},
		{"ldp", "lsPairPP"},
	}, // stp = ldp schema with L=0 (the formatter picks the name)
	"ldpsw": {{"ldpsw", "lsPairPP"}},
	"adr":   {{"adr", "adrFromFields"}},
	"adrp":  {{"adrp", "adrpFromFields"}},
	"ret":   {{"ret", "retFormatter"}},
	"brk":   {{"brk", "brkFormatter"}},
	"hlt":   {{"hlt", "brkFormatter"}}, // bare-metal semihosting exit
	"hvc":   {{"hvc", "brkFormatter"}}, // bare-metal PSCI call (qemu virt)
	"b":     {{"b", "bAbs"}},
	"bl":    {{"bl", "bAbs"}},
}

// candidatesFor — candidate schemas for a mnemonic.
func candidatesFor(mnem string) []*Schema {
	schemas := getSchemas()
	base, _, _ := strings.Cut(mnem, ".") // b.eq → b

	type key struct {
		name, formatter string
	}
	var wanted []key
	if pairs, ok := armAliasMap[mnem]; ok {
		for _, p := range pairs {
			wanted = append(wanted, key{
				p[0],
				p[1],
			})
		}
	}

	if pairs, ok := armAliasMap[base]; ok && len(wanted) == 0 {
		for _, p := range pairs {
			wanted = append(wanted, key{
				p[0],
				p[1],
			})
		}
	}

	if suffix := mnem[strings.Index(mnem, ".")+1:]; strings.Contains(mnem, ".") &&
		isCondSuffix(suffix) {
		wanted = append(wanted, key{
			"b." + suffix,
			"bAbsImm19",
		})
	}

	var out []*Schema
	// candidate order = alias priority order (not the schemas' file order):
	// the first one checked is the preferred form (e.g. mov → ORR, as
	// GAS/clang emit it)
	for _, w := range wanted {
		for i := range schemas {
			s := &schemas[i]
			if s.Meta.Name == w.name && (w.formatter == "" || s.Formatter == w.formatter) {
				out = append(out, s)
			}
		}
	}

	if len(out) > 0 {
		return out
	}

	// fallback: exact name or base name (without .Arr()/.cond suffix: add.4s → add)
	for i := range schemas {
		if schemas[i].Meta.Name == mnem || schemas[i].Meta.Name == base {
			out = append(out, &schemas[i])
		}
	}

	return out
}

// isCondSuffix — a b.<cond> suffix.
func isCondSuffix(s string) bool {
	_, ok := invCondNames[s]
	return ok
}

// armFieldsFor — operands → field values by format family.
func armFieldsFor(s *Schema, in resolvedInstr, ctx ctx) (map[string]any, error) {
	switch s.Formatter {
	case "op2":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn}, nil

	case "op3":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm}, nil

	case "op4":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		ra, err := opRegOf(in, 3)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "Ra": ra}, nil

	case "adcFmt", "op2plain":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm}, nil

	case "ccmpFmt":
		rn, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		imm, err := opImmOf(in, 2)
		if err != nil {
			return nil, err
		}

		if len(in.ops) < 4 {
			return nil, errors.New("condition expected")
		}

		cond := in.ops[3]
		if cond.Kind() != arch.ArmOpImm || cond.Sym() == "" {
			return nil, errors.New("condition expected")
		}

		return map[string]any{"Rn": rn, "Rm": rm, "imm": imm, "cond": cond.Sym()}, nil

	case "addSubImm":
		// forms: add Rd, Rn, #imm[, lsl #12] | cmp Rn, #imm | mov Rd, Rn
		switch len(in.ops) {
		case 2: // mov Rd, Rn | cmp Rn, #imm | cmn
			if in.mnem == "mov" {
				rd, err := opRegOf(in, 0)
				if err != nil {
					return nil, err
				}

				rn, err := opRegOf(in, 1)
				if err != nil {
					return nil, err
				}

				// the ADD #0 form is the sp path (and a GPR
				// fallback); a vector/FP register here is a
				// misspelled SIMD form - it must not feed addSub
				// silently as the other register operand
				if rd != "sp" && rn != "sp" &&
					(!isGPR(rd) || !isGPR(rn)) {
					return nil, errors.New("mov: integer registers expected")
				}

				rdN, err := addSubRegNum(rd)
				if err != nil {
					return nil, err
				}

				rnN, err := addSubRegNum(rn)
				if err != nil {
					return nil, err
				}

				return map[string]any{"Rd": rdN, "Rn": rnN, "imm12": 0, "shift": "lsl #0"}, nil
			}

			// cmp/cmn Rn, #imm
			rn, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			imm, err := opImmOf(in, 1)
			if err != nil {
				return nil, err
			}

			rnN, err := addSubRegNum(rn)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": 31, "Rn": rnN, "imm12": imm, "shift": "lsl #0"}, nil
		case 3, 4:
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rn, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			imm, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			shift := "lsl #0"
			if len(in.ops) == 4 && in.ops[3].Kind() == arch.ArmOpShift {
				shift = fmt.Sprintf("%s #%d", in.ops[3].ShiftName(), shiftAmt(in.ops[3]))
				if shift != "lsl #0" && shift != "lsl #12" {
					return nil, fmt.Errorf("bad immediate shift %q", shift)
				}
			}

			rdN, err := addSubRegNum(rd)
			if err != nil {
				return nil, err
			}

			rnN, err := addSubRegNum(rn)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rdN, "Rn": rnN, "imm12": imm, "shift": shift}, nil
		}

		return nil, errors.New("bad operand count")

	case "shiftedReg":
		// add Rd, Rn, Rm[, shift #n] | cmp Rn, Rm | neg Rd, Rm
		shiftName, amt, hasShift := opShiftOf(in, 0)
		sh := "lsl"
		imm6 := int64(0)
		if hasShift {
			sh = shiftName
			imm6 = amt
		}

		switch len(in.ops) - boolToInt(hasShift) {
		case 2:
			switch in.mnem {
			case "cmp", "cmn":
				rn, err := opRegOf(in, 0)
				if err != nil {
					return nil, err
				}

				rm, err := opRegOf(in, 1)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"Rd":    "xzr",
					"Rn":    rn,
					"Rm":    rm,
					"shift": sh,
					"imm6":  imm6,
				}, nil
			case "neg", "negs":
				rd, err := opRegOf(in, 0)
				if err != nil {
					return nil, err
				}

				rm, err := opRegOf(in, 1)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"Rd":    rd,
					"Rn":    "xzr",
					"Rm":    rm,
					"shift": sh,
					"imm6":  imm6,
				}, nil
			}
		case 3:
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rn, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 2)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "shift": sh, "imm6": imm6}, nil
		}

		return nil, errors.New("bad operand count")

	case "addSubExt":
		// add Rd, Rn, Rm(, ext #imm3) | cmp Rn, Rm, ext
		opt := ""
		amt := int64(0)
		for _, op := range in.ops {
			if op.Kind() == arch.ArmOpExtend {
				opt = op.ShiftName()
				if op.HasAmt() {
					amt = op.Num()
				}
			}
		}

		// 64-bit Rm without extension — uxtx (lsl); 32-bit — uxtw
		rmWidth := "x"
		if len(in.ops) >= 2 && in.ops[1].Kind() == arch.ArmOpReg && regIsW(in.ops[1].Reg()) {
			rmWidth = "w"
		}

		if len(in.ops) >= 3 && in.ops[2].Kind() == arch.ArmOpReg {
			if regIsW(in.ops[2].Reg()) {
				rmWidth = "w"
			} else {
				rmWidth = "x"
			}
		}

		if opt == "" || opt == "lsl" {
			if rmWidth == "x" {
				opt = "uxtx"
			} else {
				opt = "uxtw"
			}
		}

		fields := map[string]any{"option": opt, "imm3": amt}
		switch len(in.ops) {
		case 2: // cmp/cmn
			rn, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			fields["Rd"] = 31 // S form: Rd = 31 (xzr)
			rnN, err := armRegNum(rn)
			if err != nil {
				return nil, err
			}

			rmN, err := armRegNum(rm)
			if err != nil {
				return nil, err
			}

			fields["Rn"], fields["Rm"] = rnN, rmN
			return fields, nil
		case 3, 4:
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rn, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 2)
			if err != nil {
				return nil, err
			}

			rdN, err := armRegNum(rd)
			if err != nil {
				return nil, err
			}

			rnN, err := armRegNum(rn)
			if err != nil {
				return nil, err
			}

			rmN, err := armRegNum(rm)
			if err != nil {
				return nil, err
			}

			fields["Rd"], fields["Rn"], fields["Rm"] = rdN, rnN, rmN
			return fields, nil
		}

		return nil, errors.New("bad operand count")

	case "logicalShifted":
		// and Rd, Rn, Rm[, shift] | tst Rn, Rm | mov Rd, Rm | mvn Rd, Rm
		shiftName, amt, hasShift := opShiftOf(in, 0)
		sh := "lsl"
		imm6 := int64(0)
		if hasShift {
			sh = shiftName
			imm6 = amt
		}

		n := len(in.ops) - boolToInt(hasShift)
		switch in.mnem {
		case "tst":
			if n != 2 {
				return nil, errors.New("tst expects 2 operands")
			}

			rn, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":    zrFor(rn),
				"Rn":    rn,
				"Rm":    rm,
				"shift": sh,
				"imm6":  imm6,
			}, nil
		case "mov":
			if n != 2 {
				return nil, errors.New("mov expects 2 operands")
			}

			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":    rd,
				"Rn":    zrFor(rd),
				"Rm":    rm,
				"shift": sh,
				"imm6":  imm6,
			}, nil
		case "mvn":
			if n != 2 {
				return nil, errors.New("mvn expects 2 operands")
			}

			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":    rd,
				"Rn":    zrFor(rd),
				"Rm":    rm,
				"shift": sh,
				"imm6":  imm6,
				"N":     1,
			}, nil
		}

		if n != 3 {
			return nil, errors.New("expects 3 operands")
		}

		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "shift": sh, "imm6": imm6}, nil

	case "logicalImm":
		// and Rd, Rn, #imm | tst Rn, #imm | mov Rd, #imm
		switch len(in.ops) {
		case 2:
			if in.mnem == "tst" {
				rn, err := opRegOf(in, 0)
				if err != nil {
					return nil, err
				}

				imm, err := opImmOf(in, 1)
				if err != nil {
					return nil, err
				}

				n, immr, imms, ok := encodeBitMasks(regIs64(rn), uint64(imm))
				if !ok {
					return nil, fmt.Errorf("immediate %#x is not encodable as logical", imm)
				}

				return map[string]any{
					"Rd":   zrFor(rn),
					"Rn":   rn,
					"N":    n,
					"immr": immr,
					"imms": imms,
				}, nil
			}

			// mov Rd, #imm
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			imm, err := opImmOf(in, 1)
			if err != nil {
				return nil, err
			}

			n, immr, imms, ok := encodeBitMasks(regIs64(rd), uint64(imm))
			if !ok {
				return nil, fmt.Errorf("immediate %#x is not encodable as logical", imm)
			}

			return map[string]any{
				"Rd":   rd,
				"Rn":   zrFor(rd),
				"N":    n,
				"immr": immr,
				"imms": imms,
			}, nil
		case 3:
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rn, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			imm, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			n, immr, imms, ok := encodeBitMasks(regIs64(rd), uint64(imm))
			if !ok {
				return nil, fmt.Errorf("immediate %#x is not encodable as logical", imm)
			}

			return map[string]any{"Rd": rd, "Rn": rn, "N": n, "immr": immr, "imms": imms}, nil
		}

		return nil, errors.New("bad operand count")

	case "movWide":
		// movz/movn/movk Rd, #imm16[, lsl #hw*16] | mov Rd, #imm (alias)
		if len(in.ops) != 2 && (len(in.ops) != 3 || in.ops[2].Kind() != arch.ArmOpShift) {
			return nil, errors.New("bad operand count")
		}

		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		imm, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		hw := int64(0)
		if len(in.ops) == 3 && in.ops[2].Kind() == arch.ArmOpShift {
			sh := shiftAmt(in.ops[2])
			if sh%16 != 0 || sh > 48 {
				return nil, fmt.Errorf("bad movk shift %d", sh)
			}

			hw = sh / 16
		}

		if in.mnem == "mov" {
			v := uint64(imm)
			if s.Meta.Name == "movn" {
				if imm >= 0 {
					// positive movn pattern: v == ~(imm16 << 16hw) —
					// all ones except one 16-bit window of zeros
					// (~v — a window of ones); objdump prints these as mov
					mask := ^uint64(0)
					if rd[0] != 'x' {
						mask = 0xFFFFFFFF
					}

					nv := ^v & mask
					if nv == 0 {
						hw, imm = 0, 0 // all ones: movn #0
					} else {
						found := false
						for h := range int64(4) {
							lane := nv >> uint(16*h)
							if lane <= 0xFFFF && nv == lane<<uint(16*h) && lane != 0 {
								hw, imm, found = h, int64(lane), true
								break
							}
						}

						if !found {
							return nil, fmt.Errorf("immediate %#x not movn-encodable", v)
						}
					}
				} else {
					// mov Rd, #negative → MOVN: val = ~(imm16 << 16hw)
					x := uint64(-imm)
					switch {
					case x <= 0x10000 && x > 0:
						hw, imm = 0, int64(x-1)
					case x <= 0x100000000 && x%0x10000 == 0:
						hw, imm = 1, int64(x/0x10000-1)
					case x <= 0x1000000000000 && x%0x100000000 == 0:
						hw, imm = 2, int64(x/0x100000000-1)
					case x%0x1000000000000 == 0:
						hw, imm = 3, int64(x/0x1000000000000-1)
					default:
						return nil, fmt.Errorf("immediate %#x not movn-encodable", imm)
					}
				}
			} else {
				// mov Rd, #imm (movz): pick hw based on the value
				switch {
				case v <= 0xFFFF:
					hw = 0
				case v <= 0xFFFF0000 && v&0xFFFF == 0:
					hw = 1
					imm = int64(v >> 16)
				case v <= 0xFFFF00000000 && v&0xFFFFFFFF == 0:
					hw = 2
					imm = int64(v >> 32)
				case v&0xFFFFFFFFFFFF == 0:
					hw = 3
					imm = int64(v >> 48)
				default:
					return nil, fmt.Errorf("immediate %#x not movz-encodable", v)
				}
			}
		}

		if imm < 0 || imm > 0xFFFF {
			return nil, fmt.Errorf("imm16 out of range: %d", imm)
		}

		return map[string]any{"Rd": rd, "imm16": imm, "hw": hw}, nil

	case "bAbs":
		target, err := opImmOf(in, 0)
		if err != nil {
			return nil, err
		}

		bits, err := brBits(target, int64(ctx.Addr), 26)
		if err != nil {
			return nil, err
		}

		return map[string]any{"imm26": bits}, nil

	case "bAbsImm19":
		target, err := opImmOf(in, 0)
		if err != nil {
			return nil, err
		}

		bits, err := brBits(target, int64(ctx.Addr), 19)
		if err != nil {
			return nil, err
		}

		return map[string]any{"imm19": bits}, nil

	case "cbz":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		target, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		bits, err := brBits(target, int64(ctx.Addr), 19)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rt": rt, "imm19": bits}, nil

	case "tbBit":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		bit, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		target, err := opImmOf(in, 2)
		if err != nil {
			return nil, err
		}

		bits, err := brBits(target, int64(ctx.Addr), 14)
		if err != nil {
			return nil, err
		}

		if bit < 0 || bit > 63 {
			return nil, fmt.Errorf("bit %d out of range", bit)
		}

		num, err := armRegNum(rt)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"immRt": num, "b5": uint32(bit) >> 5, "b40": uint32(bit) & 0x1f, "imm14": bits,
		}, nil

	case "adrFromFields", "adr":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		off, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		// Our formatter prints the signed offset (#<off>) itself — exactly what
		// we encode. resolveOps normalized the symbolic target into an offset.
		rel := off

		if rel < -(1<<20) || rel >= (1<<20) {
			return nil, fmt.Errorf("adr offset %d out of range", rel)
		}

		u := uint32(rel) & 0x1FFFFF
		return map[string]any{"Rd": rd, "immlo": u & 3, "immhi": u >> 2}, nil

	case "adrpFromFields", "adrp":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		// Our formatter prints the page offset (decimal) — we encode it as is;
		// the absolute address (after ';') is cut off by the comment.
		off, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		u := uint32(int32(off)) & 0x1FFFFF
		return map[string]any{"Rd": rd, "immlo": u & 3, "immhi": u >> 2}, nil

	case "lsImm", "lsUnscaled", "lsSigned":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 1)
		if err != nil {
			return nil, err
		}

		fields := map[string]any{"Rt": rt, "Rn": m.Base()}
		if m.HasOff() {
			off := m.Off()
			switch s.Formatter {
			case "lsImm", "lsSigned":
				scale := int64(1) << (uint64(s.Value>>30) & 3)
				if off%scale != 0 {
					return nil, fmt.Errorf("offset %d not multiple of %d", off, scale)
				}

				fields["imm12"] = off / scale
			default: // unscaled
				if off < -256 || off > 255 {
					return nil, fmt.Errorf("offset %d out of imm9 range", off)
				}

				fields["imm9"] = off
			}
		} else {
			switch s.Formatter {
			case "lsImm", "lsSigned":
				fields["imm12"] = 0
			default:
				fields["imm9"] = 0
			}
		}

		return fields, nil

	case "lsRegOff", "lsFpRegOff":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 1)
		if err != nil {
			return nil, err
		}

		if m.OffReg() == "" {
			return nil, errors.New("register offset expected")
		}

		fields := map[string]any{"Rt": rt, "Rn": m.Base(), "Rm": m.OffReg()}
		if m.Opt() != "" {
			fields["option"] = m.Opt()
			sBit := int64(0)
			if m.HasOpt() {
				sBit = 1
			}

			fields["S"] = sBit
		} else {
			fields["option"] = "lsl"
			fields["S"] = 0
		}

		return fields, nil

	case "lsPreIndex", "lsPostIndex":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 1)
		if err != nil {
			return nil, err
		}

		var off int64
		switch {
		case m.Pre() && s.Formatter == "lsPreIndex" && m.HasOff():
			off = m.Off()
		case m.HasPost() && s.Formatter == "lsPostIndex":
			off = m.Post()
		case m.HasOff() && s.Formatter == "lsPreIndex":
			off = m.Off()
		default:
			return nil, errors.New("indexing form mismatch")
		}

		if off < -256 || off > 255 {
			return nil, fmt.Errorf("offset %d out of imm9 range", off)
		}

		return map[string]any{"Rt": rt, "Rn": m.Base(), "imm9": off}, nil

	case "lsPair", "lsPairPP":
		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rt2, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 2)
		if err != nil {
			return nil, err
		}

		fields := map[string]any{"Rt": rt, "Rt2": rt2, "Rn": m.Base()}
		var off int64
		switch {
		case m.HasPost():
			off = m.Post()
			fields["idx"] = 1
		case m.Pre():
			off = m.Off()
			fields["idx"] = 3
		default:
			if m.HasOff() {
				off = m.Off()
			}

			fields["idx"] = 2
		}

		scale := int64(8)
		if regIsW(rt) {
			scale = 4
		}

		fields["imm7"] = off / scale // raw bits: the lsPairImm7_* context transform
		// L: 1 for loads (ldp/ldpsw), 0 for stores; opc: 1 for ldpsw
		if in.mnem == "ldp" || in.mnem == "ldpsw" {
			fields["L"] = 1
		} else {
			fields["L"] = 0
		}

		if in.mnem == "ldpsw" {
			fields["opc"] = 1
		} else {
			fields["opc"] = 0
		}

		return fields, nil

	case "lsAtomic", "lsExclusive":
		name := s.Meta.Name
		if name == "ldxr" || name == "ldaxr" || name == "ldxp" || name == "ldaxp" ||
			s.Formatter == "lsAtomic" {
			rt, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			m, err := opMemOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rt": rt, "Rn": m.Base()}, nil
		}

		// stxr Rs, Rt, [Rn]
		rs, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rt, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 2)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rs": rs, "Rt": rt, "Rn": m.Base()}, nil

	case "stxrbFmt":
		rs, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rt, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 2)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rs": rs, "Rt": rt, "Rn": m.Base()}, nil

	case "ldrLiteral":
		// ldr Rt, target (no brackets)
		if len(in.ops) != 2 || in.ops[1].Kind() != arch.ArmOpImm {
			return nil, errors.New("literal form expects Rt, target")
		}

		rt, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		target, err := opImmOf(in, 1)
		if err != nil {
			return nil, err
		}

		bits, err := brBits(target, int64(ctx.Addr), 19)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rt": rt, "imm19": bits}, nil

	case "bfmAlias", "bfm":
		// ubfm/sbfm/bfm Rd, Rn, #immr, #imms | lsl/lsr/asr Rd, Rn, #sh
		// | sxtb/sxth/sxtw Rd, Rn | ubfiz/ubfx/sbfiz/sbfx Rd, Rn, #lsb, #width
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		regsize := int64(64)
		if regIsW(rd) {
			regsize = 32
		}

		switch in.mnem {
		case "lsl":
			sh, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":   rd,
				"Rn":   rn,
				"immr": (regsize - sh) % regsize,
				"imms": regsize - 1,
			}, nil
		case "lsr":
			sh, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rd, "Rn": rn, "immr": sh, "imms": regsize - 1}, nil
		case "asr":
			sh, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rd, "Rn": rn, "immr": sh, "imms": regsize - 1}, nil
		case "sxtb", "sxth", "sxtw", "uxtb", "uxth":
			widths := map[string]int64{"sxtb": 8, "sxth": 16, "sxtw": 32, "uxtb": 8, "uxth": 16}
			rnNum, nerr := armRegNum(rn)
			if nerr != nil {
				return nil, nerr
			}

			// sbfmAlias prints the W name from the immRn field; Rn — the x name of the same number
			return map[string]any{
				"Rd": rd, "Rn": "x" + strconvFormat(rnNum), "immRn": rnNum,
				"immr": 0, "imms": widths[in.mnem] - 1,
			}, nil
		case "ubfiz", "sbfiz":
			lsb, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			width, err := opImmOf(in, 3)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":   rd,
				"Rn":   rn,
				"immr": (regsize - lsb) % regsize,
				"imms": width - 1,
			}, nil
		case "ubfx", "sbfx":
			lsb, err := opImmOf(in, 2)
			if err != nil {
				return nil, err
			}

			width, err := opImmOf(in, 3)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rd, "Rn": rn, "immr": lsb, "imms": lsb + width - 1}, nil
		}

		if len(in.ops) != 4 {
			return nil, errors.New("bad operand count")
		}

		immr, err := opImmOf(in, 2)
		if err != nil {
			return nil, err
		}

		imms, err := opImmOf(in, 3)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "immr": immr, "imms": imms}, nil

	case "csel":
		// csel Rd, Rn, Rm, cond | cset Rd, cond | cinc Rd, Rm, cond ...
		// cset/csetm/cinc/cinv/cneg carry an inverted condition (that is how
		// the formatter prints them from csinc/csinv/csneg)
		inv := invertCond
		switch in.mnem {
		case "cset", "csetm":
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			if len(in.ops) < 2 {
				return nil, errors.New("condition expected")
			}

			return map[string]any{
				"Rd":   rd,
				"Rn":   "xzr",
				"Rm":   "xzr",
				"cond": inv(condText(in.ops[1])),
			}, nil
		}

		switch len(in.ops) {
		case 3: // cinc/cinv/cneg Rd, Rm, cond
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"Rd":   rd,
				"Rn":   rm,
				"Rm":   rm,
				"cond": inv(condText(in.ops[2])),
			}, nil
		case 4:
			rd, err := opRegOf(in, 0)
			if err != nil {
				return nil, err
			}

			rn, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			rm, err := opRegOf(in, 2)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "cond": condText(in.ops[3])}, nil
		}

		return nil, errors.New("bad operand count")

	case "madd":
		// madd Rd, Rn, Rm, Ra | mul Rd, Rn, Rm | mneg Rd, Rn, Rm
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		ra := "xzr"
		if regIsW(rd) {
			ra = "wzr"
		}

		if len(in.ops) >= 4 {
			ra, err = opRegOf(in, 3)
			if err != nil {
				return nil, err
			}
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "Ra": ra}, nil

	case "extrFmt", "extr":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm := rn
		immIdx := 2
		if len(in.ops) >= 4 {
			rm, err = opRegOf(in, 2)
			if err != nil {
				return nil, err
			}

			immIdx = 3
		}

		imms, err := opImmOf(in, immIdx)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "imms": imms}, nil

	case "retFormatter":
		if len(in.ops) == 0 {
			return map[string]any{"Rn": "x30"}, nil
		}

		rn, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rn": rn}, nil

	case "brReg":
		rn, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rn": rn}, nil

	case "operandsNone":
		return map[string]any{}, nil

	case "brkFormatter":
		imm := int64(0)
		if len(in.ops) > 0 {
			v, ierr := opImmOf(in, 0)
			if ierr != nil {
				return nil, ierr
			}

			imm = v
		}

		return map[string]any{"imm16": imm}, nil

	case "udfFormatter":
		imm, err := opImmOf(in, 0)
		if err != nil {
			return nil, err
		}

		return map[string]any{"imm16": imm}, nil

	case "svcFmt":
		imm, err := opImmOf(in, 0)
		if err != nil {
			return nil, err
		}

		return map[string]any{"imm16": imm}, nil

	case "mrsFmt":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		if len(in.ops) < 2 {
			return nil, errors.New("sysreg expected")
		}

		sr := in.ops[1]
		if sr.Kind() != arch.ArmOpImm || sr.Sym() == "" {
			return nil, errors.New("sysreg name expected")
		}

		return map[string]any{"Rd": rd, "sysreg": sr.Sym()}, nil

	case "msrFmt":
		if len(in.ops) < 2 {
			return nil, errors.New("msr expects sysreg, Rt")
		}

		sr := in.ops[0]
		if sr.Kind() != arch.ArmOpImm || sr.Sym() == "" {
			return nil, errors.New("sysreg name expected")
		}

		rt, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rt": rt, "sysreg": sr.Sym()}, nil

	case "fmovImmD":
		if len(in.ops) != 2 {
			return nil, errors.New("fmov expects 2 operands")
		}

		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		if in.ops[1].Kind() != arch.ArmOpFloat {
			return nil, errors.New("float immediate expected")
		}

		imm8, ok := encodeVFPImm64(in.ops[1].Float())
		if !ok {
			return nil, fmt.Errorf("immediate %v not VFP-encodable", in.ops[1].Float())
		}

		return map[string]any{"Rd": rd, "imm8": imm8}, nil

	case "fmovImmS":
		if len(in.ops) != 2 {
			return nil, errors.New("fmov expects 2 operands")
		}

		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		if in.ops[1].Kind() != arch.ArmOpFloat {
			return nil, errors.New("float immediate expected")
		}

		imm8, ok := encodeVFPImm32(in.ops[1].Float())
		if !ok {
			return nil, fmt.Errorf("immediate %v not VFP-encodable", in.ops[1].Float())
		}

		return map[string]any{"Rd": rd, "imm8": imm8}, nil

	case "simdShiftImm":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		amt, err := opImmOf(in, 2)
		if err != nil {
			return nil, err
		}

		q, size, err := arrangementOf(arrOf(in))
		if err != nil {
			return nil, err
		}

		// Inversion of simdShiftAmount/decodeShiftImm: immh = (1<<size) | (raw>>3),
		// immb = raw&7, where raw is the 6-bit immh:immb value without the size
		// marker; ushr/sshr/sri print esize - raw.
		esize := int64(8) << size
		raw := amt
		switch s.Meta.Name {
		case "ushr", "sshr", "sri":
			raw = esize - amt
		}

		if raw < 0 || raw >= esize {
			return nil, fmt.Errorf("shift %d out of range for .%s", amt, arrOf(in))
		}

		immh := uint32(1)<<uint(size) | uint32(raw>>3)
		immb := uint32(raw) & 7
		return map[string]any{"Rd": rd, "Rn": rn, "Q": q, "immh": immh, "immb": immb}, nil

	case "fmovGen":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn}, nil

	case "fcmp0", "fcmp2":
		rn, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		if len(in.ops) >= 2 && in.ops[1].Kind() == arch.ArmOpReg {
			rm, err := opRegOf(in, 1)
			if err != nil {
				return nil, err
			}

			return map[string]any{"Rn": rn, "Rm": rm}, nil
		}

		return map[string]any{"Rn": rn}, nil

	case "simd3":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		rm, err := opRegOf(in, 2)
		if err != nil {
			return nil, err
		}

		q, size, err := arrangementOf(arrOf(in))
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Rm": rm, "Q": q, "size": size}, nil

	case "simd2":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		q, size, err := arrangementOf(arrOf(in))
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Q": q, "size": size}, nil

	case "v1arr":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn}, nil

	case "dupFmt":
		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		q, size, err := arrangementOf(arrOf(in))
		if err != nil {
			return nil, err
		}

		return map[string]any{"Rd": rd, "Rn": rn, "Q": q, "size": size}, nil

	case "movElem":
		if len(in.ops) != 2 {
			return nil, errors.New("mov.d expects 2 operands")
		}

		rd, err := opRegOf(in, 0)
		if err != nil {
			return nil, err
		}

		rn, err := opRegOf(in, 1)
		if err != nil {
			return nil, err
		}

		idx := int64(0)
		if in.ops[1].LaneIdx() {
			idx = in.ops[1].Num()
		}

		// .d element: imm5 = 1<<(3+idx) (mirror of movElem: idx = imm5>>4)
		return map[string]any{"Rd": rd, "Rn": rn, "imm5": 1 << (3 + idx)}, nil

	case "ldStruct", "ld1Fmt", "ld1rFmt":
		rt, err := listFirstReg(in, 0)
		if err != nil {
			return nil, err
		}

		m, err := opMemOf(in, 1)
		if err != nil {
			return nil, err
		}

		arr := arrOf(in)
		q, size, err := arrangementOf(arr)
		if err != nil {
			return nil, err
		}

		count := 1
		if len(in.ops) > 0 && in.ops[0].Kind() == arch.ArmOpList {
			count = len(in.ops[0].List())
		}

		l := 0
		if strings.HasPrefix(in.mnem, "ld") {
			l = 1
		}

		// opcode by register count (mirror of ldStructDecode)
		opc := map[int]uint32{4: 0x0, 3: 0x4, 2: 0xa, 1: 0x7}[count]
		switch {
		case strings.HasSuffix(in.mnem, "4r"):
			opc = 0xe
		case strings.HasSuffix(in.mnem, "1r"):
			opc = 0xc
		}

		fields := map[string]any{
			"Rt": rt, "Rn": m.Base(), "opcode": opc, "size": size, "Q": q, "L": l,
		}
		if m.HasPost() {
			// post-index immediate: Rm=11111 (the value is NOT encoded —
			// it is implicit: regBytes×count / element size)
			fields["s"] = uint32(0x1f)
		}

		return fields, nil
	}

	return nil, fmt.Errorf("formatter %q not supported by assembler", s.Formatter)
}

// --- helpers ---

func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

func regIs64(name string) bool {
	return len(name) > 0 && name[0] == 'x' || name == "sp"
}

// zrFor — the zero-register name of the same width as ref.
func zrFor(ref string) string {
	if regIsW(ref) {
		return "wzr"
	}

	return "xzr"
}

func regIsW(name string) bool {
	return len(name) > 0 && name[0] == 'w'
}

// addSubRegNum — register number honoring add/sub semantics (for fields
// without a Transform: raw uint32): the S form (adds/subs → cmp) uses xzr, otherwise sp.
func addSubRegNum(name string) (uint32, error) {
	n, err := armRegNum(name)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// condText — condition text from an operand.
func condText(op vOp) string {
	return op.Sym()
}

// arrangementOf — ".8b"/".2d"/... → (Q, size).
func arrangementOf(arr string) (uint32, uint32, error) {
	switch arr {
	case "8b":
		return 0, 0, nil
	case "16b":
		return 1, 0, nil
	case "4h":
		return 0, 1, nil
	case "8h":
		return 1, 1, nil
	case "2s":
		return 0, 2, nil
	case "4s":
		return 1, 2, nil
	case "1d":
		return 0, 3, nil
	case "2d":
		return 1, 3, nil
	}

	return 0, 0, fmt.Errorf("unknown arrangement %q", arr)
}

// arrOf — arrangement from the first operand's register suffix or from the
// mnemonic suffix ("cmeq.2d v8, ..." — suffix on the mnemonic).
func arrOf(in resolvedInstr) string {
	if len(in.ops) > 0 && in.ops[0].Arr() != "" {
		return in.ops[0].Arr()
	}

	if i := strings.IndexByte(in.mnem, '.'); i >= 0 {
		return in.mnem[i+1:]
	}

	return ""
}

// listFirstReg — the first register of a list (or a single register).
func listFirstReg(in resolvedInstr, i int) (string, error) {
	if i >= len(in.ops) {
		return "", fmt.Errorf("operand %d expected", i+1)
	}

	op := in.ops[i]
	switch op.Kind() {
	case arch.ArmOpReg:
		return op.Reg(), nil
	case arch.ArmOpList:
		if len(op.List()) == 0 {
			return "", errors.New("empty register list")
		}

		return op.List()[0].Reg(), nil
	default:
		return "", errors.New("register or register list expected")
	}
}

// strconvFormat — decimal formatting without importing strconv on the hot path.
func strconvFormat(n uint32) string {
	if n == 0 {
		return "0"
	}

	var b [4]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return string(b[i:])
}

// encodeVFPImm — float → imm8 by brute force (VFPExpandImm is not invertible
// for all values; 256 candidates are cheap to check).
func encodeVFPImm64(v float64) (uint32, bool) {
	for i := range uint32(256) {
		if vfpExpandImm64(i) == v {
			return i, true
		}
	}

	return 0, false
}

func encodeVFPImm32(v float64) (uint32, bool) {
	for i := range uint32(256) {
		if float64(vfpExpandImm32(i)) == v {
			return i, true
		}
	}

	return 0, false
}

// --- operand access + word packing (moved from arch assemble.go) ---

// contextTransforms — transforms that require addr (branches, literal
// offsets): their inverses are not in the registry, format handlers put
// ready bits into the fields, and packFields passes them through the empty
// transform.
var contextTransforms = map[string]bool{
	"brOff26": true, "brOff19": true, "brOff14": true,
	"sext19": true, "sext9": true,
	"lsPairImm7_64": true, "lsPairImm7_32": true,
}

// packFields assembles the word: Schema.Value | fields (inverse transforms).
func packFields(s *Schema, fields map[string]any) (uint32, error) {
	w := s.Value
	for _, f := range s.Fields {
		v, ok := fields[f.Name]
		if !ok {
			// field not set by a handler — must be part of Value
			continue
		}

		var bits uint32
		var err error
		if contextTransforms[f.Transform] {
			bits, err = applyInverseTransform("", v) // value already = bits
			if err == nil && f.Width < 32 {
				bits &= (1 << f.Width) - 1 // signed raw bits (imm7=-2 → 0x7e)
			}
		} else {
			bits, err = applyInverseTransform(f.Transform, v)
		}

		if err != nil {
			return 0, fmt.Errorf("field %s: %w", f.Name, err)
		}

		if f.Width < 32 && bits>>f.Width != 0 {
			return 0, fmt.Errorf("field %s: %#x does not fit in %d bits", f.Name, bits, f.Width)
		}

		w |= bits << f.Offset
	}

	return w, nil
}

// opRegOf — the register operand at position i (or an error).
func opRegOf(in resolvedInstr, i int) (string, error) {
	if i >= len(in.ops) {
		return "", fmt.Errorf("operand %d: register expected", i+1)
	}

	op := in.ops[i]
	if op.Kind() != arch.ArmOpReg {
		return "", fmt.Errorf("operand %d: register expected", i+1)
	}

	return op.Reg(), nil
}

// opImmOf — the number operand at position i (the value is already computed).
func opImmOf(in resolvedInstr, i int) (int64, error) {
	if i >= len(in.ops) {
		return 0, fmt.Errorf("operand %d: immediate expected", i+1)
	}

	op := in.ops[i]
	if op.Kind() == arch.ArmOpImm {
		return op.Num(), nil
	}

	if op.Kind() == arch.ArmOpFloat {
		return int64(math.Float64bits(op.Float())), nil
	}

	return 0, fmt.Errorf("operand %d: immediate expected", i+1)
}

// opMemOf — the memory operand at position i.
func opMemOf(in resolvedInstr, i int) (vMem, error) {
	if i >= len(in.ops) || in.ops[i].Kind() != arch.ArmOpMem {
		return vMem{}, fmt.Errorf("operand %d: memory operand expected", i+1)
	}

	return in.ops[i].Mem(), nil
}

// opShiftOf — a trailing lsl/... modifier (after position from).
func opShiftOf(in resolvedInstr, from int) (name string, amt int64, ok bool) {
	for i := from; i < len(in.ops); i++ {
		if in.ops[i].Kind() == arch.ArmOpShift {
			return in.ops[i].ShiftName(), in.ops[i].Num(), true
		}
	}

	return "", 0, false
}

// legacyBuild — legacy encoding of a candidate from computed operands:
// format-family handlers + word assembly via inverse transforms (the
// old arch.BuildLegacy seam, now local).
func legacyBuild(s *Schema, mnem string, ops []vOp, pc uint64) (uint32, error) {
	fields, err := armFieldsFor(s, resolvedInstr{mnem: mnem, ops: ops}, ctx{Addr: pc})
	if err != nil {
		return 0, err
	}

	return packFields(s, fields)
}
