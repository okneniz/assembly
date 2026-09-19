package arm64

import "fmt"

// Shared bits of the by-element family (one abstraction per
// instruction lives in its own file; this is the field layout they
// encode against - the family's Vm/index bit scatter is common).

// byElemBits - the family's word: q, U (bit 29), size, Vm, opc, the
// MSB-aligned lane index in the field {b11,b21,b20,b19}, Rn, Rd.
func byElemBits(q, u, size, rm, opc, idx, rn, rd uint32) uint32 {
	f := idx << size // MSB-aligned in the 4-bit field {b11,b21,b20,b19}
	return q<<30 | u<<29 | 0x0F<<24 | size<<22 | rm<<16 | opc<<12 |
		(f>>3&1)<<11 | (f>>2&1)<<21 | (f>>1&1)<<20 | (f&1)<<19 | rn<<5 | rd
}

// byElemIndex - the lane index read from the word: MSB-aligned bits of
// the field {b11,b21,b20,b19} (lanes: .b→4 .h→3 .s→2 .d→1).
func byElemIndex(w, size uint32) uint32 {
	lanes := 4 - size
	field := w>>11&1<<3 | w>>21&1<<2 | w>>20&1<<1 | w>>19&1
	return field >> (4 - lanes)
}

// byElemVm - the Vm register number: bits [19:16] for a .h source
// (4-bit), [20:16] otherwise.
func byElemVm(w, size uint32) uint32 {
	if size == 1 {
		return w >> 16 & 0xf
	}

	return w >> 16 & 0x1f
}

// requireByElemLane - the lane index fits the lane count, and a .h
// source stays inside the 4-bit Vm field.
func requireByElemLane(instr string, size, idx uint32, rm VReg) error {
	if idx >= 1<<(4-size) {
		return fmt.Errorf(
			"arm64.New%s: lane index %d out of range (0..%d)",
			instr, idx, 1<<(4-size)-1,
		)
	}

	if size == 1 && rm.Num() > 15 {
		return fmt.Errorf(
			"arm64.New%s: register %s is out of the 4-bit Vm field of .h lanes",
			instr, rm.name(),
		)
	}

	return nil
}
