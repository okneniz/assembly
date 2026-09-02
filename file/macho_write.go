package file

// The Mach-O writer: a minimal arm64 MH_EXECUTE that macOS runs natively -
// the Darwin counterpart of WriteELF. Beyond the segments (PAGEZERO, __TEXT
// with one __text section, __LINKEDIT) the image carries the full set of
// load commands ld emits: strict validation in AMFI/codesign (enforced for
// arm64 main executables since macOS 13) rejects anything leaner - every
// command of the reference set is mandatory, down to an empty chained-fixups
// header and a one-function FUNCTION_STARTS table.
//
// The same validation requires an ad-hoc code signature embedded in the
// file (LC_CODE_SIGNATURE + a CodeDirectory over the 4K pages of everything
// before it), so the writer signs the image itself - the byte-identical
// scheme codesign -s - produces (CD v0x20400, flags adhoc|linker-signed,
// SHA-256). __TEXT is padded to whole 16K kernel pages: every mapped page
// must be fully file-backed, or the kernel kills the process at exec.

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

// Mach-O constants the writer needs (the macho package stays a pure parser).
const (
	machoMagic64  = 0xfeedfacf
	machoCPUArm64 = 0x0100000c // CPU_TYPE_ARM64 = 12 | CPU_ARCH_ABI64
	machoExecute  = 2          // MH_EXECUTE

	machoSegment64    = 0x19       // LC_SEGMENT_64
	machoChainedFix   = 0x80000034 // LC_DYLD_CHAINED_FIXUPS
	machoExportsTrie  = 0x80000033 // LC_DYLD_EXPORTS_TRIE
	machoSymtab       = 0x2        // LC_SYMTAB
	machoDyldInfoOnly = 0xb        // LC_DYLD_INFO_ONLY (all-zero: nothing to do)
	machoDylinker     = 0xe        // LC_LOAD_DYLINKER
	machoUUID         = 0x1b       // LC_UUID
	machoBuildVersion = 0x32       // LC_BUILD_VERSION
	machoSourceVer    = 0x2a       // LC_SOURCE_VERSION
	machoMain         = 0x80000028 // LC_MAIN | LC_REQ_DYLD: dyld calls the entry
	machoLoadDylib    = 0xc        // LC_LOAD_DYLIB
	machoFuncStarts   = 0x26       // LC_FUNCTION_STARTS
	machoDataInCode   = 0x29       // LC_DATA_IN_CODE
	machoCodeSig      = 0x1d       // LC_CODE_SIGNATURE
)

// Page granularity: Apple Silicon maps and validates __TEXT at the 16K
// kernel page; the CodeDirectory hashes at the 4K granularity codesign uses.
const (
	machoVMAddr    = 0x100000000 // the arm64 convention for __TEXT
	machoKPage     = 16384
	machoHashPage  = 4096
	machoFlags     = 0x200085   // MH_PIE | MH_TWOLEVEL | MH_DYLDLINK | MH_NOUNDEFS
	machoMinOS     = 0x000F0500 // 15.5: what the reference ld declares
	machoLdTool    = 3          // TOOL_LD
	machoLdVersion = 0x00058F04
)

// MachoTextBase - the vmaddr the writer places __TEXT at: the arm64
// convention. Programs with absolute operands must be assembled at this
// base; pc-relative code runs at any base.
const MachoTextBase = machoVMAddr

// MachoCodeOff - the file offset the code starts at: right after the header
// and the fixed command set (the commands themselves are size-constant, so
// the offset never moves). Assemble position-independent code or account
// for it in absolute operands.
const MachoCodeOff = 696

// The __LINKEDIT tables: 176 constant bytes in the layout
//
//	chained fixups [0:56)   an empty header (no rebase/bind targets)
//	exports trie   [56:104) __mh_execute_header and main
//	function strts [104:112) the single entry function (the entry point)
//	symbol table   [112:144) the same two symbols as nlist_64
//	strings        [144:176)
//
// The template is extracted verbatim from a minimal ld-produced executable
// (clang -c + ld -e _main, macOS 15) with the entry-dependent spots left as
// the ULEB/value of the reference entry (696 = MachoCodeOff, so a program
// entering at its first byte needs no patching at all); WriteMachO rewrites
// those spots for other entry offsets. Reproducing the trie/nlist bytes by
// hand is not worth it: the encoding is validated byte-for-byte by AMFI.
const (
	machoLinkeditSize  = 176
	machoTrieEntryAt   = 71  // 2-byte ULEB inside the trie
	machoFstartEntryAt = 104 // 2-byte ULEB of FUNCTION_STARTS
	machoSymtabEntryAt = 136 // u64 nlist value of main
	machoSymtabMHAt    = 120 // u64 nlist value of __mh_execute_header
)

var machoLinkedit = [machoLinkeditSize]byte{
	0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00,
	0x30, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x5f, 0x00, 0x12, 0x00, 0x00, 0x00,
	0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0xb8,
	0x05, 0x00, 0x00, 0x02, 0x5f, 0x6d, 0x68, 0x5f,
	0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x5f,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x00, 0x09,
	0x6d, 0x61, 0x69, 0x6e, 0x00, 0x0d, 0x00, 0x00,
	0xb8, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x02, 0x00, 0x00, 0x00, 0x0f, 0x01, 0x10, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x16, 0x00, 0x00, 0x00, 0x0f, 0x01, 0x00, 0x00,
	0xb8, 0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x20, 0x00, 0x5f, 0x5f, 0x6d, 0x68, 0x5f, 0x65,
	0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x5f, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x00, 0x5f, 0x6d,
	0x61, 0x69, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00,
}

// WriteMachO wraps an arm64 machine code blob into a MH_EXECUTE that macOS
// (Apple Silicon and Intel Macs running arm64 binaries via Rosetta aside -
// arm64 hosts) executes as-is: no linker, no codesign step. text is the raw
// program (the concatenation the assembler produces), entry the offset of
// the first instruction inside text (0 for a program entering at its start).
// The program sees argc/argv/envp in the C main registers but must not
// return - exit via a syscall (there is no libc to come back to).
func WriteMachO(text []byte, entry uint64) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("macho: no text")
	}

	if entry >= uint64(len(text)) {
		return nil, fmt.Errorf("macho: entry offset %#x is outside the text (%d bytes)", entry, len(text))
	}

	// The template holds the reference entry 696 as a 2-byte ULEB; a longer
	// encoding would shift the trie (a patched value must stay < 8192).
	entryAbs := uint64(MachoCodeOff) + entry
	if entryAbs >= 8192 {
		return nil, fmt.Errorf("macho: entry offset %#x does not fit the linkedit template", entryAbs)
	}

	// __TEXT: header + commands + code, padded to whole 16K pages.
	textSize := MachoCodeOff + len(text)
	if textSize%machoKPage == 0 {
		textSize /= machoKPage
	} else {
		textSize = (textSize/machoKPage + 1) * machoKPage
	}

	leOff := textSize
	sigOff := leOff + machoLinkeditSize
	sigLen := 12 + 8 + 88 + 9 + 32*((sigOff+machoHashPage-1)/machoHashPage)
	leSize := machoLinkeditSize + int(sigLen)

	out := make([]byte, sigOff+int(sigLen))
	le := binary.LittleEndian

	// mach_header_64.
	le.PutUint32(out[0:], machoMagic64)
	le.PutUint32(out[4:], machoCPUArm64)
	le.PutUint32(out[8:], 0) // cpusubtype ARM64_ALL
	le.PutUint32(out[12:], machoExecute)
	le.PutUint32(out[16:], 16) // ncmds
	le.PutUint32(out[20:], 664)
	le.PutUint32(out[24:], machoFlags)
	le.PutUint32(out[28:], 0)

	pos := 32

	// LC_SEGMENT_64 __PAGEZERO: the guard trap range at zero.
	le.PutUint32(out[pos:], machoSegment64)
	le.PutUint32(out[pos+4:], 72)
	copy(out[pos+8:], "__PAGEZERO\x00\x00\x00\x00\x00\x00")
	le.PutUint64(out[pos+32:], 0x100000000) // vmsize: the lower 4GB
	pos += 72

	// LC_SEGMENT_64 __TEXT: the whole first pages, one __text section.
	le.PutUint32(out[pos:], machoSegment64)
	le.PutUint32(out[pos+4:], 152)
	copy(out[pos+8:], "__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	le.PutUint64(out[pos+24:], machoVMAddr)
	le.PutUint64(out[pos+32:], uint64(textSize))
	le.PutUint64(out[pos+40:], 0)                // fileoff: the image from 0
	le.PutUint64(out[pos+48:], uint64(textSize)) // filesize: every page file-backed
	le.PutUint32(out[pos+56:], 7)                // maxprot rwx
	le.PutUint32(out[pos+60:], 5)                // initprot r-x
	le.PutUint32(out[pos+64:], 1)                // nsects
	sec := out[pos+72:]
	copy(sec, "__text\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	copy(sec[16:], "__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	le.PutUint64(sec[32:], machoVMAddr+MachoCodeOff) // addr
	le.PutUint64(sec[40:], uint64(len(text)))        // size
	le.PutUint32(sec[48:], MachoCodeOff)             // offset
	le.PutUint32(sec[52:], 2)                        // align: 2^2
	le.PutUint32(sec[64:], 0x80000400)               // PURE_INSTRUCTIONS | SOME_INSTRUCTIONS
	pos += 152

	// LC_SEGMENT_64 __LINKEDIT: the tables + the signature.
	le.PutUint32(out[pos:], machoSegment64)
	le.PutUint32(out[pos+4:], 72)
	copy(out[pos+8:], "__LINKEDIT\x00\x00\x00\x00\x00\x00")
	le.PutUint64(out[pos+24:], machoVMAddr+uint64(textSize))
	le.PutUint64(out[pos+32:], uint64((leSize+machoKPage-1)&^(machoKPage-1)))
	le.PutUint64(out[pos+40:], uint64(leOff))
	le.PutUint64(out[pos+48:], uint64(leSize))
	le.PutUint32(out[pos+56:], 1)
	le.PutUint32(out[pos+60:], 1)
	pos += 72

	lc16 := func(cmd uint32, a, b uint32) {
		le.PutUint32(out[pos:], cmd)
		le.PutUint32(out[pos+4:], 16)
		le.PutUint32(out[pos+8:], a)
		le.PutUint32(out[pos+12:], b)
		pos += 16
	}

	lc16(machoChainedFix, uint32(leOff), 56)     // empty fixups header
	lc16(machoExportsTrie, uint32(leOff)+56, 48) // __mh_execute_header, main

	// LC_SYMTAB: two symbols (the linkedit order is fixups, trie, fstarts,
	// symtab, strings; the commands follow the reference file's order).
	le.PutUint32(out[pos:], machoSymtab)
	le.PutUint32(out[pos+4:], 24)
	le.PutUint32(out[pos+8:], uint32(leOff)+112)  // symoff
	le.PutUint32(out[pos+12:], 2)                 // nsyms
	le.PutUint32(out[pos+16:], uint32(leOff)+144) // stroff
	le.PutUint32(out[pos+20:], 32)                // strsize
	pos += 24

	// LC_DYLD_INFO_ONLY: every offset zero - nothing to rebase or bind.
	le.PutUint32(out[pos:], machoDyldInfoOnly)
	le.PutUint32(out[pos+4:], 80)
	pos += 80

	// LC_LOAD_DYLINKER: macOS has no static executables; dyld calls the
	// LC_MAIN entry. The binary loads libSystem (below) but never calls
	// into it - the program talks to the kernel directly.
	le.PutUint32(out[pos:], machoDylinker)
	le.PutUint32(out[pos+4:], 32)
	le.PutUint32(out[pos+8:], 12)
	copy(out[pos+12:], "/usr/lib/dyld\x00")
	pos += 32

	// LC_UUID: derived from the code - stable for the same program.
	le.PutUint32(out[pos:], machoUUID)
	le.PutUint32(out[pos+4:], 24)
	id := sha256.Sum256(text)
	copy(out[pos+8:], id[:16])
	pos += 24

	// LC_BUILD_VERSION: platform macos, the reference minos, ld as the tool.
	le.PutUint32(out[pos:], machoBuildVersion)
	le.PutUint32(out[pos+4:], 32)
	le.PutUint32(out[pos+8:], 1) // PLATFORM_MACOS
	le.PutUint32(out[pos+12:], machoMinOS)
	le.PutUint32(out[pos+16:], 0) // sdk: n/a
	le.PutUint32(out[pos+20:], 1) // ntools
	le.PutUint32(out[pos+24:], machoLdTool)
	le.PutUint32(out[pos+28:], machoLdVersion)
	pos += 32

	// LC_SOURCE_VERSION: none.
	le.PutUint32(out[pos:], machoSourceVer)
	le.PutUint32(out[pos+4:], 16)
	pos += 16

	// LC_MAIN: dyld calls vmaddr+entryoff with main-style arguments.
	le.PutUint32(out[pos:], machoMain)
	le.PutUint32(out[pos+4:], 24)
	le.PutUint64(out[pos+8:], entryAbs)
	le.PutUint64(out[pos+16:], 0) // stacksize: the default
	pos += 24

	// LC_LOAD_DYLIB: strict validation requires at least one; the reference
	// version constants are carried over verbatim.
	le.PutUint32(out[pos:], machoLoadDylib)
	le.PutUint32(out[pos+4:], 56)
	le.PutUint32(out[pos+8:], 24)          // name offset
	le.PutUint32(out[pos+12:], 2)          // timestamp
	le.PutUint32(out[pos+16:], 0x05470000) // current_version
	le.PutUint32(out[pos+20:], 0x00010000) // compatibility_version
	copy(out[pos+24:], "/usr/lib/libSystem.B.dylib\x00")
	pos += 56

	lc16(machoFuncStarts, uint32(leOff)+104, 8) // the entry function
	lc16(machoDataInCode, uint32(leOff)+112, 0) // empty

	// LC_CODE_SIGNATURE: the ad-hoc signature closing the file.
	lc16(machoCodeSig, uint32(sigOff), uint32(sigLen))

	copy(out[MachoCodeOff:], text)

	// The linkedit tables with the entry point patched in.
	copy(out[leOff:], machoLinkedit[:])
	uleb2 := func(v uint64) (byte, byte) {
		return byte(v) | 0x80, byte(v >> 7)
	}

	b0, b1 := uleb2(entryAbs)
	out[leOff+machoTrieEntryAt] = b0
	out[leOff+machoTrieEntryAt+1] = b1
	out[leOff+machoFstartEntryAt] = b0
	out[leOff+machoFstartEntryAt+1] = b1
	le.PutUint64(out[leOff+machoSymtabEntryAt:], machoVMAddr+entryAbs)

	copy(out[sigOff:], machoSignature(out[:sigOff], uint32(sigOff), uint32(textSize)))

	return out, nil
}

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

	hput32(0x20400)                                             // version
	hput32(0x20002)                                             // flags: CS_ADHOC | CS_LINKER_SIGNED
	hput32(uint32(cdHdrLen + len(ident)))                       // hashOffset
	hput32(cdHdrLen)                                            // identOffset
	hput32(0)                                                   // nSpecialSlots
	hput32(n)                                                   // nCodeSlots
	hput32(codeLimit)                                           // codeLimit
	head = append(head, 32, 2, 0, 12)                           // hashSize, hashType SHA-256, platform, pageSize 2^12
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
