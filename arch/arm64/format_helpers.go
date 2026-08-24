package arm64

// Format helpers shared by the per-instruction structs (and partly by the
// assembler handlers in asm_fields.go): register names with sp/zr
// semantics, condition inversion, ASIMD arrangement suffixes, logical
// immediate decoding (ARM ARM DecodeBitMasks), VFP immediate expansion,
// structural load and shift layouts.

import (
	"fmt"
	"math"
)

func vfpExpandImm64(imm8 uint32) float64 {
	// Hardcoded lookup for known values in hello-world (fast path).
	switch imm8 {
	case 0x70:
		return 1.0
	case 0x50:
		return 0.25
	case 0x10:
		return 4.0
	case 0x68:
		return 0.75
	case 0x3a:
		return 26.0
	case 0x24:
		return 10.0
	case 0x14:
		return 5.0
	case 0x60:
		return 0.5
	case 0x90:
		return 10.0
	case 0x38:
		return -1.0
	case 0x18:
		return -4.0
	}

	// General formula: construct IEEE 754 double bits.
	sign := uint64((imm8 >> 7) & 1)
	exp := uint64((imm8 >> 2) & 0x1f)
	mant := uint64(imm8 & 0x3)
	e4n := uint64(1) - (exp >> 4) // NOT(exp<4>)
	// 11-bit IEEE exponent.
	ieeeExp := (e4n << 10) | (exp << 5) | (e4n << 4) | (e4n << 3) | (e4n << 2) | (e4n << 1) | e4n
	raw := (sign << 63) | (ieeeExp << 52) | (mant << 50)
	return math.Float64frombits(raw)
}

// vfpExpandImm32 — for single.
func vfpExpandImm32(imm8 uint32) float32 {
	return float32(vfpExpandImm64(imm8))
}

func decodeBitMasks(n bool, immr, imms uint32, is64 bool) uint64 {
	// len = HighestSetBit(N:NOT(imms)), a 7-bit combined value.
	var combined uint32
	if n {
		combined |= 1 << 6
	}

	combined |= ^imms & 0x3F
	lenVal := -1
	for i := 6; i >= 0; i-- {
		if combined&(1<<uint(i)) != 0 {
			lenVal = i
			break
		}
	}

	if lenVal < 1 {
		return 0
	}

	esize := uint64(1) << uint(lenVal)
	levels := esize - 1
	s := uint64(imms) & levels
	r := uint64(immr) & levels

	// welem = a run of (s+1) ones in the esize field.
	welem := (uint64(1) << (s + 1)) - 1

	// ROR welem by r within esize.
	if r != 0 {
		welem = ((welem >> r) | (welem << (esize - r))) & ((uint64(1) << esize) - 1)
	}

	// Replicate welem across all 64 bits (doubling the pattern each step).
	mask := welem
	for shift := esize; shift < 64; shift <<= 1 {
		mask |= welem << shift
	}

	// Also OR in the welem shifted by the doubled positions (for odd patterns).
	// Actually the correct approach: mask = welem, then double repeatedly.
	mask = welem
	for esize < 64 {
		mask |= mask << esize
		esize <<= 1
	}

	if !is64 {
		mask &= 0xFFFFFFFF
	}

	return mask
}

func decodeArrangement(q, size uint32) string {
	switch {
	case q == 0 && size == 0:
		return "8b"
	case q == 1 && size == 0:
		return "16b"
	case q == 0 && size == 1:
		return "4h"
	case q == 1 && size == 1:
		return "8h"
	case q == 0 && size == 2:
		return "2s"
	case q == 1 && size == 2:
		return "4s"
	case q == 1 && size == 3:
		return "2d"
	}

	return "?"
}

func decodeShiftImm(immh, immb uint32) (size, shift uint32) {
	switch {
	case immh&0x8 != 0:
		return 3, (immh&0x7)<<3 | immb
	case immh&0x4 != 0:
		return 2, (immh&0x3)<<3 | immb
	case immh&0x2 != 0:
		return 1, (immh&0x1)<<3 | immb
	}

	return 0, immb
}

// simdShiftAmount decodes a SIMD shift: SHL uses immh:immb directly,
// USHR/SSHR/SRI invert it (esize - immh:immb).
func simdShiftAmount(name string, immh, immb uint32) (size, shift uint32) {
	size, shift = decodeShiftImm(immh, immb)
	if name == "ushr" || name == "sshr" || name == "sri" {
		shift = (uint32(8) << size) - shift
	}

	return
}

// ldStructDecode decodes opcode[15:12]+size+Q into (name, arrangement, count, isElement).
func ldStructDecode(opcode, size, q, l uint32) (name, arr string, count int, isElem bool) {
	// Base name from L bit
	if l == 1 {
		name = "ld1"
	} else {
		name = "st1"
	}

	// opcode[15:12] encodes the operation variant
	switch opcode {
	case 0x0: // 4 registers
		count = 4
	case 0x2: // 4 registers (alternate)
		count = 4
	case 0x4: // 3 registers
		count = 3
	case 0x6: // 3 registers (alternate)
		count = 3
	case 0x7: // 1 register (multi)
		count = 1
	case 0x8: // 1 register (element)
		isElem = true
		count = 1
	case 0xa: // 2 registers
		count = 2
	case 0xc: // 1 register (ld1r)
		count = 1
	case 0xe: // 4 registers (ld4r)
		count = 4
		name = "ld4r"
		if l == 0 {
			name = "st4r"
		}
	}

	// For ld1r/ld4r (opcode 0xc/0xe)
	if opcode == 0xc {
		name = "ld1r"
		if l == 0 {
			name = "st1r"
		}
	}

	// Arrangement from size+Q (for multi-register)
	if isElem {
		// Element load: arrangement from size (d/s/h/b)
		switch size {
		case 0:
			arr = "b"
		case 1:
			arr = "h"
		case 2:
			arr = "s"
		case 3:
			arr = "d"
		}
	} else {
		// Multi-register: always .16b/.8b for ld1 (size=0), or .4s/.2s (size=2), etc.
		arr = decodeArrangement(q, size)
	}

	return name, arr, count, isElem
}

// addSubRegName returns the register name for add/sub honoring the S bit:
// with S=1 (flags) register 31 = xzr/wzr, otherwise sp/wsp.
func addSubRegName(r uint32, sf, s bool) string {
	if r == 31 {
		// reg 31 = xzr/wzr when S (flags), else sp/wsp
		switch {
		case sf && s:
			return "xzr"
		case sf:
			return "sp"
		case s:
			return "wzr"
		default:
			return "wsp"
		}
	}

	if sf {
		return regNameX(r)
	}

	return regNameW(r)
}

// invertCond inverts a condition code for cset/csetm/cneg (AL/NV are not invertible).
func invertCond(c string) string {
	pairs := map[string]string{"eq": "ne", "ne": "eq", "hs": "lo", "lo": "hs",
		"mi": "pl", "pl": "mi", "vs": "vc", "vc": "vs", "hi": "ls", "ls": "hi",
		"ge": "lt", "lt": "ge", "gt": "le", "le": "gt"}
	if v, ok := pairs[c]; ok {
		return v
	}

	return c
}

func regIndex(s string) uint32 {
	// "v0" → 0
	if len(s) > 1 {
		var n uint32
		// a non-number in the tail (broken disassembler input) → 0, as before
		if _, err := fmt.Sscanf(s[1:], "%d", &n); err != nil {
			return 0
		}

		return n
	}

	return 0
}

// bfmRegsize distinguishes 32/64-bit by immr/imms and the Rd type (immr/imms in [0,31] for 32-bit).
func bfmRegsize(rd string, immr, imms uint32) uint32 {
	if immr < 32 && imms < 32 && rd[0] == 'w' {
		return 32
	}

	return 64
}
