package arm64

// Logical immediate encoder: mask → (N, immr, imms). The inverse of
// decodeBitMasks (which in this decoder always uses esize = width: the mask
// is a rotation of a contiguous run of ones; the encoder agrees with it by
// construction — exactly what TestBitMasksRoundTrip checks by exhaustive
// enumeration).

import (
	"fmt"
)

// encodeBitMasks finds (N, immr, imms) for a 32/64-bit wide mask value.
// ok=false for non-encodable masks (0, all ones, non-periodic sets).
// Consistent with decodeBitMasks: the mask — a replication (esize is a
// power of two) of a rotated contiguous run of ones.
func encodeBitMasks(is64 bool, value uint64) (n, immr, imms uint32, ok bool) {
	width := uint32(64)
	if !is64 {
		width = 32
		value &= 0xFFFFFFFF
	}

	full := ^uint64(0) >> (64 - width)
	if value == 0 || value == full {
		return 0, 0, 0, false
	}

	// the smallest esize (a power of two) by which the mask replicates
	for esize := uint32(2); esize <= width; esize <<= 1 {
		if width%esize != 0 {
			continue
		}

		emask := ^uint64(0) >> (64 - esize)
		welem := value & emask
		// replication check — by doubling, as in decodeBitMasks
		replicated := welem
		for e := esize; e < width; e <<= 1 {
			replicated |= replicated << e
		}

		if replicated&full != value {
			continue
		}

		if welem == 0 || welem == emask {
			continue // degenerate (would give 0/all ones — cut off above, but just in case)
		}

		ones := uint32(0)
		for i := range esize {
			if welem>>i&1 == 1 {
				ones++
			}
		}

		if ones >= esize {
			continue
		}

		run := (uint64(1) << ones) - 1
		// welem = ror(run, r) within esize → search for r
		for r := range esize {
			rot := welem<<r | welem>>(esize-r)
			rot &= emask
			if rot != run {
				continue
			}

			nbits := log2pow2(esize)
			immr = r % esize
			if is64 && nbits == 6 {
				return 1, immr, ones - 1, true
			}

			if nbits >= 6 {
				return 0, 0, 0, false
			}

			// n=0: ~imms&0x3f must have the top bit nbits →
			// imms bits above nbits — ones, bit nbits — zero, the low nbits bits — s
			imms = ((uint32(0x3f) >> (nbits + 1)) << (nbits + 1)) | (ones - 1)
			imms &^= 1 << nbits
			return 0, immr, imms, true
		}
	}

	return 0, 0, 0, false
}

func log2pow2(v uint32) uint32 {
	n := uint32(0)
	for v > 1 {
		v >>= 1
		n++
	}

	return n
}

// verifyBitMasks — a property test: all (N, immr, imms) × widths pass
// decode → encode → decode with the same result.
func verifyBitMasks() error {
	for _, is64 := range []bool{true, false} {
		for n := range uint32(2) {
			if !is64 && n == 1 {
				continue
			}

			for immr := range uint32(64) {
				for imms := range uint32(64) {
					mask := decodeBitMasks(n == 1, immr, imms, is64)
					if !is64 {
						mask &= 0xFFFFFFFF
					}

					full := ^uint64(0)
					if !is64 {
						full = 0xFFFFFFFF
					}

					if mask == 0 || mask == full {
						continue
					}

					en, er, es, ok := encodeBitMasks(is64, mask)
					if !ok {
						return fmt.Errorf("encodeBitMasks(%v, %#x) failed", is64, mask)
					}

					back := decodeBitMasks(en == 1, er, es, is64)
					if !is64 {
						back &= 0xFFFFFFFF
					}

					if back != mask {
						return fmt.Errorf("round-trip %v n=%d r=%d s=%d: %#x → %#x",
							is64, n, immr, imms, mask, back)
					}
				}
			}
		}
	}

	return nil
}
