package arm64

// Inverse transforms for assembly: register name → number (honoring the
// x31 semantics: xzr or sp — decided by the field's Transform or the format
// handler), condition → index, sysreg name → 15-bit key, shift/extend names.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// armRegNum parses a register name: prefix + number (x0..x30, w3, v31, d7,
// s2, b0, h5, q1). Named ones (sp/xzr/wzr/wsp) → 31.
func armRegNum(name string) (uint32, error) {
	switch name {
	case "sp", "xzr", "wzr", "wsp":
		return 31, nil
	}

	if len(name) < 2 {
		return 0, fmt.Errorf("bad register %q", name)
	}

	n, err := strconv.Atoi(name[1:])
	if err != nil || n < 0 || n > 31 {
		return 0, fmt.Errorf("bad register %q", name)
	}

	return uint32(n), nil
}

// invRegX: names for the regX Transform (x31 = xzr).
func invRegX(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("register expected")
	}

	if s == "xzr" || s == "x31" {
		return 31, nil
	}

	if len(s) < 2 || s[0] != 'x' {
		return 0, fmt.Errorf("x-register expected, got %q", s)
	}

	return armRegNum(s)
}

// invRegW: names for the regW Transform (w31 = wzr).
func invRegW(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("register expected")
	}

	if s == "wzr" || s == "w31" {
		return 31, nil
	}

	if len(s) < 2 || s[0] != 'w' {
		return 0, fmt.Errorf("w-register expected, got %q", s)
	}

	return armRegNum(s)
}

// invRegXSP: x31 = sp.
func invRegXSP(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("register expected")
	}

	if s == "sp" {
		return 31, nil
	}

	if s == "xzr" {
		return 0, errors.New("sp-capable register expected, got xzr")
	}

	return invRegX(v)
}

// invRegWSP: w31 = wsp.
func invRegWSP(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("register expected")
	}

	if s == "wsp" {
		return 31, nil
	}

	return invRegW(v)
}

// invRegV: SIMD v registers.
func invRegV(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("register expected")
	}

	if len(s) < 2 || s[0] != 'v' {
		return 0, fmt.Errorf("v-register expected, got %q", s)
	}

	return armRegNum(s)
}

// invFpReg: FP registers with a given prefix (b/h/s/d/q).
func invFpReg(prefix byte) func(any) (uint32, error) {
	return func(v any) (uint32, error) {
		s, ok := v.(string)
		if !ok {
			return 0, errors.New("register expected")
		}

		if len(s) < 2 || s[0] != prefix {
			return 0, fmt.Errorf("%c-register expected, got %q", prefix, s)
		}

		return armRegNum(s)
	}
}

// invCond: condition name → index (inverse of condNames; the cs/cc
// synonyms from GNU syntax are also accepted).
var invCondNames = func() map[string]uint32 {
	m := map[string]uint32{}
	for i, n := range condNames {
		m[n] = uint32(i)
	}

	m["cs"] = 2 // hs
	m["cc"] = 3 // lo
	return m
}()

func invCond(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("condition expected")
	}

	if i, ok := invCondNames[s]; ok {
		return i, nil
	}

	return 0, fmt.Errorf("unknown condition %q", s)
}

// invSysReg: system register name → 15-bit key (inverse of sysregNames;
// the objdump form S<op0>_<op1>_C<n>_C<m>_<op2> is also accepted).
var invSysRegNames = func() map[string]uint32 {
	m := map[string]uint32{}
	for k, name := range sysregNames {
		if _, exists := m[name]; !exists {
			m[name] = k
		}
	}

	return m
}()

func invSysReg(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("sysreg expected")
	}

	if k, ok := invSysRegNames[s]; ok {
		return k, nil
	}

	if strings.HasPrefix(s, "S") && strings.Count(s, "_") == 4 {
		parts := strings.Split(s[1:], "_")
		op0, err0 := strconv.Atoi(parts[0])
		op1, err1 := strconv.Atoi(parts[1])
		crn, err2 := strconv.Atoi(strings.TrimPrefix(parts[2], "C"))
		crm, err3 := strconv.Atoi(strings.TrimPrefix(parts[3], "C"))
		op2, err4 := strconv.Atoi(parts[4])
		if err0 == nil && err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
			op0 >= 2 && op0 <= 3 && crn >= 0 && crn <= 15 && crm >= 0 && crm <= 15 {
			return uint32(
				op0&1,
			)<<14 | uint32(
				op1,
			)<<11 | uint32(
				crn,
			)<<7 | uint32(
				crm,
			)<<3 | uint32(
				op2,
			), nil
		}
	}

	return 0, fmt.Errorf("unknown system register %q", s)
}

var invShiftNames = map[string]uint32{"lsl": 0, "lsr": 1, "asr": 2, "ror": 3}

func invShiftName(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("shift name expected")
	}

	if i, ok := invShiftNames[s]; ok {
		return i, nil
	}

	return 0, fmt.Errorf("unknown shift %q", s)
}

// invImmShiftLSL: "lsl #0" → 0, "lsl #12" → 1.
func invImmShiftLSL(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("shift expected")
	}

	switch s {
	case "lsl #0", "":
		return 0, nil
	case "lsl #12":
		return 1, nil
	}

	return 0, fmt.Errorf("bad immediate shift %q", s)
}

// invExtOpt: extend names.
var invExtNames = map[string]uint32{
	"uxtb": 0, "uxth": 1, "uxtw": 2, "uxtx": 3,
	"sxtb": 4, "sxth": 5, "sxtw": 6, "sxtx": 7,
}

func invExtOpt(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("extend expected")
	}

	if i, ok := invExtNames[s]; ok {
		return i, nil
	}

	return 0, fmt.Errorf("unknown extend %q", s)
}

// invLsOptNames: load/store options (3-bit codes).
var invLsOptNames = map[string]uint32{
	"uxtw": 0b010, "lsl": 0b011, "sxtw": 0b110, "sxtx": 0b111,
}

func invLsOpt(v any) (uint32, error) {
	s, ok := v.(string)
	if !ok {
		return 0, errors.New("option expected")
	}

	if i, ok := invLsOptNames[s]; ok {
		return i, nil
	}

	return 0, fmt.Errorf("unknown option %q", s)
}

// invInt: a numeric field value (int64/uint32/int).
func invInt(v any) (uint32, error) {
	switch n := v.(type) {
	case int:
		return uint32(n), nil
	case int64:
		return uint32(n), nil
	case uint32:
		return n, nil
	case uint64:
		return uint32(n), nil
	}

	return 0, errors.New("numeric field expected")
}
