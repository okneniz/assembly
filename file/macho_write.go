package file

// The Mach-O writer: a minimal arm64 MH_EXECUTE that macOS runs natively -
// the Darwin counterpart of WriteELF. The universal entry is NewMachOImage
// (arbitrary sections, symbols, a symbol entry point); WriteMachO is its
// legacy text-only face, kept byte-identical by the golden test. Beyond
// the segments (PAGEZERO, __TEXT with its sections, __DATA and custom
// rw- segments when present, __LINKEDIT with generated tables) the image
// carries the full set of load commands ld emits: strict validation in
// AMFI/codesign (enforced for arm64 main executables since macOS 13)
// rejects anything leaner - every command of the reference set is
// mandatory, down to an empty chained-fixups header and a one-function
// FUNCTION_STARTS table.
//
// The same validation requires an ad-hoc code signature embedded in the
// file (LC_CODE_SIGNATURE + a CodeDirectory over the 4K pages of everything
// before it), so the writer signs the image itself - the byte-identical
// scheme codesign -s - produces (CD v0x20400, flags adhoc|linker-signed,
// SHA-256). Every mapped page is fully file-backed, or the kernel kills
// the process at exec: segments are padded to whole 16K pages, and a NOBITS
// reserve only ever trails file-backed data.

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

// MachoCodeOff - the file offset the code starts at in a WriteMachO image:
// right after the header and the fixed command set (the commands of the
// legacy shape are size-constant, so the offset never moves). Images with
// data sections have more commands and a different offset - their geometry
// lives inside MachOImage, and label references resolve through the
// exported placement policies, never through this constant.
const MachoCodeOff = 696

// WriteMachO wraps an arm64 machine code blob into a MH_EXECUTE that macOS
// (arm64 hosts) executes as-is: no linker, no codesign step. text is the raw
// program (the concatenation the assembler produces), entry the offset of
// the first instruction inside text (0 for a program entering at its start).
// The program sees argc/argv/envp in the C main registers but must not
// return - exit via a syscall (there is no libc to come back to).
func WriteMachO(text []byte, entry uint64) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("macho: no text")
	}

	if entry >= uint64(len(text)) {
		return nil, fmt.Errorf(
			"macho: entry offset %#x is outside the text (%d bytes)",
			entry,
			len(text),
		)
	}

	img, err := NewMachOImage(
		[]MachOSection{{Segment: "__TEXT", Name: "__text", Data: text, Align: 4}},
		[]MachOSym{{Name: "_main", Section: "__text", Off: entry, Global: true}},
		"_main",
	)
	if err != nil {
		return nil, err
	}

	return img.Bytes(), nil
}

// emitMachO - the byte emitter: header, load commands, section data, the
// linkedit tables, the signature. Everything is validated and placed, so
// the emitter only serializes.
func emitMachO(m *MachOImage) []byte {
	p := &m.place
	le := &m.le
	out := make([]byte, le.sigOff+le.sigLen)
	b := binary.LittleEndian

	// mach_header_64.
	b.PutUint32(out[0:], machoMagic64)
	b.PutUint32(out[4:], machoCPUArm64)
	b.PutUint32(out[8:], 0) // cpusubtype ARM64_ALL
	b.PutUint32(out[12:], machoExecute)
	b.PutUint32(out[16:], p.ncmds)
	b.PutUint32(out[20:], p.cmdSize)
	b.PutUint32(out[24:], machoFlags)
	b.PutUint32(out[28:], 0)

	pos := 32

	// LC_SEGMENT_64 __PAGEZERO: the guard trap range at zero.
	b.PutUint32(out[pos:], machoSegment64)
	b.PutUint32(out[pos+4:], 72)
	copy(out[pos+8:], "__PAGEZERO\x00\x00\x00\x00\x00\x00")
	b.PutUint64(out[pos+32:], 0x100000000) // vmsize: the lower 4GB
	pos += 72

	// LC_SEGMENT_64 per placement (__TEXT first, then the rw- ones), each
	// with its sections.
	for si := range p.segs {
		seg := &p.segs[si]

		b.PutUint32(out[pos:], machoSegment64)
		b.PutUint32(out[pos+4:], uint32(72+80*len(seg.sects)))
		copy(out[pos+8:], seg.name) // the 16-byte name field is zero-padded
		b.PutUint64(out[pos+24:], seg.vmaddr)
		b.PutUint64(out[pos+32:], seg.vmsize)
		b.PutUint64(out[pos+40:], seg.fileoff)
		b.PutUint64(out[pos+48:], seg.filesize)
		b.PutUint32(out[pos+56:], seg.maxprot)
		b.PutUint32(out[pos+60:], seg.initprot)
		b.PutUint32(out[pos+64:], uint32(len(seg.sects)))

		sec := out[pos+72:]
		for _, sp := range seg.sects {
			s := &m.sections[sp.idx]

			copy(sec, s.Name)
			copy(sec[16:], seg.name)
			b.PutUint64(sec[32:], sp.addr)
			b.PutUint64(sec[40:], sp.size)
			b.PutUint32(sec[48:], sp.offset)
			b.PutUint32(sec[52:], sp.align)
			b.PutUint32(sec[64:], machoSectFlags(seg.name, &sp))
			sec = sec[80:]
		}

		pos += 72 + 80*len(seg.sects)
	}

	// LC_SEGMENT_64 __LINKEDIT: the tables + the signature. Its vmaddr is
	// the memory continuation (a pure-bss segment before it has file
	// bytes but still owns pages), its file offset the file continuation.
	b.PutUint32(out[pos:], machoSegment64)
	b.PutUint32(out[pos+4:], 72)
	copy(out[pos+8:], "__LINKEDIT\x00\x00\x00\x00\x00\x00")
	b.PutUint64(out[pos+24:], machoVMAddr+p.vmTotal)
	b.PutUint64(out[pos+32:], roundMachoKPage(le.leSize))
	b.PutUint64(out[pos+40:], le.leOff)
	b.PutUint64(out[pos+48:], le.leSize)
	b.PutUint32(out[pos+56:], 1)
	b.PutUint32(out[pos+60:], 1)
	pos += 72

	// the fixed command set, in the reference file's order
	lc16 := func(cmd uint32, a, b2 uint32) {
		b.PutUint32(out[pos:], cmd)
		b.PutUint32(out[pos+4:], 16)
		b.PutUint32(out[pos+8:], a)
		b.PutUint32(out[pos+12:], b2)
		pos += 16
	}

	trieLen := le.fstartAt - le.trieAt
	fstartLen := le.symAt - le.fstartAt

	lc16(machoChainedFix, uint32(le.leOff), 56)                  // empty fixups header
	lc16(machoExportsTrie, uint32(le.leOff)+le.trieAt, trieLen)  // __mh_execute_header + globals

	// LC_SYMTAB: __mh_execute_header plus the input symbols.
	b.PutUint32(out[pos:], machoSymtab)
	b.PutUint32(out[pos+4:], 24)
	b.PutUint32(out[pos+8:], uint32(le.leOff)+le.symAt)
	b.PutUint32(out[pos+12:], le.nsyms)
	b.PutUint32(out[pos+16:], uint32(le.leOff)+le.strAt)
	b.PutUint32(out[pos+20:], le.strSize)
	pos += 24

	// LC_DYLD_INFO_ONLY: every offset zero - nothing to rebase or bind.
	b.PutUint32(out[pos:], machoDyldInfoOnly)
	b.PutUint32(out[pos+4:], 80)
	pos += 80

	// LC_LOAD_DYLINKER: macOS has no static executables; dyld calls the
	// LC_MAIN entry. The binary loads libSystem (below) but never calls
	// into it - the program talks to the kernel directly.
	b.PutUint32(out[pos:], machoDylinker)
	b.PutUint32(out[pos+4:], 32)
	b.PutUint32(out[pos+8:], 12)
	copy(out[pos+12:], "/usr/lib/dyld\x00")
	pos += 32

	// LC_UUID: derived from the section data - stable for the same program.
	b.PutUint32(out[pos:], machoUUID)
	b.PutUint32(out[pos+4:], 24)
	id := sha256.New()
	for i := range m.sections {
		id.Write(m.sections[i].Data)
	}
	copy(out[pos+8:], id.Sum(nil)[:16])
	pos += 24

	// LC_BUILD_VERSION: platform macos, the reference minos, ld as the tool.
	b.PutUint32(out[pos:], machoBuildVersion)
	b.PutUint32(out[pos+4:], 32)
	b.PutUint32(out[pos+8:], 1) // PLATFORM_MACOS
	b.PutUint32(out[pos+12:], machoMinOS)
	b.PutUint32(out[pos+16:], 0) // sdk: n/a
	b.PutUint32(out[pos+20:], 1) // ntools
	b.PutUint32(out[pos+24:], machoLdTool)
	b.PutUint32(out[pos+28:], machoLdVersion)
	pos += 32

	// LC_SOURCE_VERSION: none.
	b.PutUint32(out[pos:], machoSourceVer)
	b.PutUint32(out[pos+4:], 16)
	pos += 16

	// LC_MAIN: dyld calls vmaddr+entryoff with main-style arguments.
	b.PutUint32(out[pos:], machoMain)
	b.PutUint32(out[pos+4:], 24)
	b.PutUint64(out[pos+8:], m.entryOff)
	b.PutUint64(out[pos+16:], 0) // stacksize: the default
	pos += 24

	// LC_LOAD_DYLIB: strict validation requires at least one; the reference
	// version constants are carried over verbatim.
	b.PutUint32(out[pos:], machoLoadDylib)
	b.PutUint32(out[pos+4:], 56)
	b.PutUint32(out[pos+8:], 24)          // name offset
	b.PutUint32(out[pos+12:], 2)          // timestamp
	b.PutUint32(out[pos+16:], 0x05470000) // current_version
	b.PutUint32(out[pos+20:], 0x00010000) // compatibility_version
	copy(out[pos+24:], "/usr/lib/libSystem.B.dylib\x00")
	pos += 56

	lc16(machoFuncStarts, uint32(le.leOff)+le.fstartAt, fstartLen) // the __TEXT symbols
	lc16(machoDataInCode, uint32(le.leOff)+le.symAt, 0)            // empty
	lc16(machoCodeSig, uint32(le.sigOff), uint32(le.sigLen))       // the signature closing the file

	// the section data at the placed offsets
	for i := range m.sections {
		if sp := &p.sect[i]; !sp.nobits {
			copy(out[sp.offset:], m.sections[i].Data)
		}
	}

	// the linkedit tables, then the signature over everything before it
	copy(out[le.leOff:], le.data)
	copy(out[le.sigOff:], machoSignature(out[:le.sigOff], uint32(le.sigOff), le.execSeg))

	return out
}

// machoSectFlags - section attribute flags: executable for __TEXT data
// sections, S_ZEROFILL for the nobits reserves.
func machoSectFlags(seg string, sp *machoSectPlace) uint32 {
	if sp.nobits {
		return 1 // S_ZEROFILL
	}

	if seg == "__TEXT" {
		return 0x80000400 // PURE_INSTRUCTIONS | SOME_INSTRUCTIONS
	}

	return 0
}
