package file

// The ad-hoc code signature of the Mach-O writer.

import "crypto/sha256"

// machoSignature - the ad-hoc CodeDirectory superblob over the 4K pages of
// everything before it (big-endian fields, the cs_blobs.h layout):
// CD v0x20400, flags adhoc|linker-signed, SHA-256, page size 2^12,
// execSeg covering __TEXT - the exact scheme of codesign -s - and of the
// Go linker, which is what lets the unsigned-toolchain image run.
func machoSignature(image []byte, codeLimit, execSegLimit uint32) []byte {
	const (
		magicSuper = 0xfade0cc0
		magicCD    = 0xfade0c02
		ident      = "assembly\x00"
		cdHdrLen   = 88 // magic+len+version..execSegFlags
	)

	n := (codeLimit + machoHashPage - 1) / machoHashPage
	cdLen := cdHdrLen + len(ident) + 32*int(n)

	sig := make([]byte, 0, 12+8+cdLen)
	put32 := func(v uint32) {
		sig = append(sig, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}

	put32(magicSuper)
	put32(uint32(12 + 8 + cdLen))
	put32(1)  // one blob
	put32(0)  // slot 0 type: CodeDirectory
	put32(20) // slot 0 offset
	put32(magicCD)
	put32(uint32(cdLen))

	head := make([]byte, 0, 72)
	hput32 := func(v uint32) {
		head = append(head, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}

	hput32(0x20400) // version
	hput32(
		0x20002,
	) // flags: CS_ADHOC | CS_LINKER_SIGNED
	hput32(uint32(cdHdrLen + len(ident))) // hashOffset
	hput32(cdHdrLen)                      // identOffset
	hput32(0)                             // nSpecialSlots
	hput32(n)                             // nCodeSlots
	hput32(codeLimit)                     // codeLimit
	head = append(
		head,
		32,
		2,
		0,
		12,
	) // hashSize, hashType SHA-256, platform, pageSize 2^12
	hput32(0)                                                   // spare2
	hput32(0)                                                   // scatterOffset
	hput32(0)                                                   // teamOffset
	hput32(0)                                                   // spare3
	for _, v := range []uint64{0, 0, uint64(execSegLimit), 1} { // codeLimit64, execSegBase, execSegLimit, MAIN_BINARY
		head = append(head, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}

	sig = append(sig, head...)
	sig = append(sig, ident...)

	for off := uint32(0); off < codeLimit; off += machoHashPage {
		end := min(off+machoHashPage, codeLimit)
		h := sha256.Sum256(image[off:end])
		sig = append(sig, h[:]...)
	}

	return sig
}
