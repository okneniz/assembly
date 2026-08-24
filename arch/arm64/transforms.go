package arm64

import (
	"fmt"
)

// Register and arithmetic helpers of the ARM64 decode form: register names,
// sign extension, field literals. Consumed by the decode constructors of
// per-instruction structures.

// regNameX returns the name of a 64-bit register: x0..x30, x31 -> xzr.
func regNameX(v uint32) string {
	if v == 31 {
		return "xzr"
	}

	return regX(int(v))
}

// regNameW returns the name of a 32-bit register: w0..w30, w31 -> wzr.
func regNameW(v uint32) string {
	if v == 31 {
		return "wzr"
	}

	return regW(int(v))
}

// regNameXSP - like regNameX, but x31 -> sp (for instructions where 31 =
// SP: add/sub/cmp/mov with an immediate and load/store by the base
// address).
func regNameXSP(v uint32) string {
	if v == 31 {
		return "sp"
	}

	return regX(int(v))
}

func regX(n int) string {
	if n < 0 || n > 31 {
		return "x?"
	}

	return stringReg('x', n)
}

func regW(n int) string {
	if n < 0 || n > 31 {
		return "w?"
	}

	return stringReg('w', n)
}

// stringReg builds a register name from a prefix and a number without using
// fmt, to keep it out of the decoder's hot path.
func stringReg(prefix byte, n int) string {
	if n == 0 {
		return string([]byte{prefix, '0'})
	}

	var buf [3]byte
	buf[0] = prefix
	i := 1
	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i++
	}

	// Reverse the digits (at most two of them for 0..31).
	out := make([]byte, i)
	out[0] = buf[0]
	for j := 1; j < i; j++ {
		out[j] = buf[i-j]
	}

	return string(out)
}

// --- Sign extensions for branches ---

func signExtendN(v uint32, bits uint) int64 {
	signed := int64(v)
	if bits >= 64 {
		return signed
	}

	signBit := int64(1) << (bits - 1)
	if signed&signBit != 0 {
		signed |= ^int64((uint64(1) << bits) - 1)
	}

	return signed
}

// --- Conditions (for b.cond, csel, cset, etc.) ---

// condNames - condition names in objdump notation (uses hs/lo, not cs/cc,
// for conditions 2/3 = carry-set/clear).
var condNames = [16]string{
	"eq", "ne", "hs", "lo", "mi", "pl", "vs", "vc",
	"hi", "ls", "ge", "lt", "gt", "le", "al", "nv",
}

func condName(v uint32) string {
	if int(v) >= len(condNames) {
		return "??"
	}

	return condNames[v]
}

// sysRegName resolves a 15-bit system-register field value (instruction bits
// 19:5, as declared by the MSR/MRS schemas' "sysreg" field) to an architectural
// name from sysregNames (generated from the official ARM XML plus m1n1's Apple
// overlay). Unknown encodings fall back to the objdump-style
// S<op0>_<op1>_C<CRn>_C<CRm>_<op2> form. op0's high bit is pinned to 1 by the
// schema mask, so only its low bit is carried in v.
func sysRegName(v uint32) string {
	if name, ok := sysregNames[v]; ok {
		return name
	}

	op0 := 2 + (v>>14)&1
	op1 := (v >> 11) & 7
	crn := (v >> 7) & 0xf
	crm := (v >> 3) & 0xf
	op2 := v & 7
	return fmt.Sprintf("S%d_%d_C%d_C%d_%d", op0, op1, crn, crm, op2)
}

// --- Shifts for immediate arithmetic ---

// fpRegNameD - the name of a double FP register (d0..).
func fpRegNameD(r uint32) string {
	return fmt.Sprintf("d%d", r)
}

// fpRegNameS - the name of a single FP register (s0..).
func fpRegNameS(r uint32) string {
	return fmt.Sprintf("s%d", r)
}

var shiftNames = [4]string{"lsl", "lsr", "asr", "ror"}
