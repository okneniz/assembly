// Package macho is a self-contained parser of the Mach-O format (64/32-bit,
// LE/BE, including FAT/Universal archives), built on the okneniz/parsec
// parser combinators over a single positional buffer, without using
// debug/macho. The package knows nothing about ELF or the rest of assembly:
// it owns all of its types.
//
// Constants are taken from <mach-o/loader.h>, <mach-o/nlist.h>,
// <mach-o/stab.h>, <mach-o/arm64/reloc.h>, <mach-o/fat.h>,
// <mach/machine.h> of the macOS SDK.
package macho

import "strconv"

// --- Magic numbers ---

const (
	MH_MAGIC    uint32 = 0xfeedface // 32-bit Mach-O (LE)
	MH_CIGAM    uint32 = 0xcefaedfe // MH_MAGIC in reversed byte order (BE file)
	MH_MAGIC_64 uint32 = 0xfeedfacf // 64-bit Mach-O (LE)
	MH_CIGAM_64 uint32 = 0xcffaedfe // MH_MAGIC_64 in reversed byte order (BE file)

	FAT_MAGIC    uint32 = 0xcafebabe // FAT/Universal (BE header)
	FAT_CIGAM    uint32 = 0xbebafeca
	FAT_MAGIC_64 uint32 = 0xcafebabf // FAT with 64-bit offset/size
	FAT_CIGAM_64 uint32 = 0xbfbafeca
)

// --- File type (filetype) ---

// FileType is the filetype field of the header.
type FileType uint32

const (
	MH_OBJECT      FileType = 1  // object file (.o)
	MH_EXECUTE     FileType = 2  // executable file
	MH_FVMLIB      FileType = 3  // fixed VM shared library
	MH_CORE        FileType = 4  // core dump
	MH_PRELOAD     FileType = 5  // preloaded executable
	MH_DYLIB       FileType = 6  // dynamic library
	MH_DYLINKER    FileType = 7  // dynamic linker
	MH_BUNDLE      FileType = 8  // loadable bundle
	MH_DYLIB_STUB  FileType = 9  // dylib stub
	MH_DSYMFILE    FileType = 10 // companion file with debug symbols
	MH_KEXT_BUNDLE FileType = 11 // kernel extension
	MH_FILESET     FileType = 12 // collection of Mach-O files (dyld shared cache)
	MH_GPU_EXECUTE FileType = 13
	MH_GPU_DYLIB   FileType = 14
)

func (t FileType) String() string {
	switch t {
	case MH_OBJECT:
		return "MH_OBJECT"
	case MH_EXECUTE:
		return "MH_EXECUTE"
	case MH_FVMLIB:
		return "MH_FVMLIB"
	case MH_CORE:
		return "MH_CORE"
	case MH_PRELOAD:
		return "MH_PRELOAD"
	case MH_DYLIB:
		return "MH_DYLIB"
	case MH_DYLINKER:
		return "MH_DYLINKER"
	case MH_BUNDLE:
		return "MH_BUNDLE"
	case MH_DYLIB_STUB:
		return "MH_DYLIB_STUB"
	case MH_DSYMFILE:
		return "MH_DSYMFILE"
	case MH_KEXT_BUNDLE:
		return "MH_KEXT_BUNDLE"
	case MH_FILESET:
		return "MH_FILESET"
	case MH_GPU_EXECUTE:
		return "MH_GPU_EXECUTE"
	case MH_GPU_DYLIB:
		return "MH_GPU_DYLIB"
	}

	return "MH_<" + strconv.FormatUint(uint64(t), 10) + ">"
}

// --- Header flags (mach_header.flags) ---

// HeaderFlag is a bitmask of header flags.
type HeaderFlag uint32

const (
	MH_NOUNDEFS                      HeaderFlag = 0x1
	MH_INCRLINK                      HeaderFlag = 0x2
	MH_DYLDLINK                      HeaderFlag = 0x4
	MH_BINDATLOAD                    HeaderFlag = 0x8
	MH_PREBOUND                      HeaderFlag = 0x10
	MH_SPLIT_SEGS                    HeaderFlag = 0x20
	MH_LAZY_INIT                     HeaderFlag = 0x40
	MH_TWOLEVEL                      HeaderFlag = 0x80
	MH_FORCE_FLAT                    HeaderFlag = 0x100
	MH_NOMULTIDEFS                   HeaderFlag = 0x200
	MH_NOFIXPREBINDING               HeaderFlag = 0x400
	MH_PREBINDABLE                   HeaderFlag = 0x800
	MH_ALLMODSBOUND                  HeaderFlag = 0x1000
	MH_SUBSECTIONS_VIA_SYMBOLS       HeaderFlag = 0x2000
	MH_CANONICAL                     HeaderFlag = 0x4000
	MH_WEAK_DEFINES                  HeaderFlag = 0x8000
	MH_BINDS_TO_WEAK                 HeaderFlag = 0x10000
	MH_ALLOW_STACK_EXECUTION         HeaderFlag = 0x20000
	MH_ROOT_SAFE                     HeaderFlag = 0x40000
	MH_SETUID_SAFE                   HeaderFlag = 0x80000
	MH_NO_REEXPORTED_DYLIBS          HeaderFlag = 0x100000
	MH_PIE                           HeaderFlag = 0x200000
	MH_DEAD_STRIPPABLE_DYLIB         HeaderFlag = 0x400000
	MH_HAS_TLV_DESCRIPTORS           HeaderFlag = 0x800000
	MH_NO_HEAP_EXECUTION             HeaderFlag = 0x1000000
	MH_APP_EXTENSION_SAFE            HeaderFlag = 0x02000000
	MH_NLIST_OUTOFSYNC_WITH_DYLDINFO HeaderFlag = 0x04000000
	MH_SIM_SUPPORT                   HeaderFlag = 0x08000000
	MH_DYLIB_IN_CACHE                HeaderFlag = 0x80000000
)

func (f HeaderFlag) String() string {
	if f == 0 {
		return "0"
	}

	names := []struct {
		bit  HeaderFlag
		name string
	}{
		{
			MH_NOUNDEFS,
			"MH_NOUNDEFS",
		},
		{
			MH_INCRLINK,
			"MH_INCRLINK",
		},
		{
			MH_DYLDLINK,
			"MH_DYLDLINK",
		},
		{
			MH_BINDATLOAD,
			"MH_BINDATLOAD",
		},
		{
			MH_PREBOUND,
			"MH_PREBOUND",
		},
		{
			MH_SPLIT_SEGS,
			"MH_SPLIT_SEGS",
		},
		{
			MH_LAZY_INIT,
			"MH_LAZY_INIT",
		},
		{
			MH_TWOLEVEL,
			"MH_TWOLEVEL",
		},
		{
			MH_FORCE_FLAT,
			"MH_FORCE_FLAT",
		},
		{
			MH_NOMULTIDEFS,
			"MH_NOMULTIDEFS",
		},
		{
			MH_NOFIXPREBINDING,
			"MH_NOFIXPREBINDING",
		},
		{
			MH_PREBINDABLE,
			"MH_PREBINDABLE",
		},
		{
			MH_ALLMODSBOUND,
			"MH_ALLMODSBOUND",
		},
		{
			MH_SUBSECTIONS_VIA_SYMBOLS,
			"MH_SUBSECTIONS_VIA_SYMBOLS",
		},
		{
			MH_CANONICAL,
			"MH_CANONICAL",
		},
		{
			MH_WEAK_DEFINES,
			"MH_WEAK_DEFINES",
		},
		{
			MH_BINDS_TO_WEAK,
			"MH_BINDS_TO_WEAK",
		},
		{
			MH_ALLOW_STACK_EXECUTION,
			"MH_ALLOW_STACK_EXECUTION",
		},
		{
			MH_ROOT_SAFE,
			"MH_ROOT_SAFE",
		},
		{
			MH_SETUID_SAFE,
			"MH_SETUID_SAFE",
		},
		{
			MH_NO_REEXPORTED_DYLIBS,
			"MH_NO_REEXPORTED_DYLIBS",
		},
		{
			MH_PIE,
			"MH_PIE",
		},
		{
			MH_DEAD_STRIPPABLE_DYLIB,
			"MH_DEAD_STRIPPABLE_DYLIB",
		},
		{
			MH_HAS_TLV_DESCRIPTORS,
			"MH_HAS_TLV_DESCRIPTORS",
		},
		{
			MH_NO_HEAP_EXECUTION,
			"MH_NO_HEAP_EXECUTION",
		},
		{
			MH_APP_EXTENSION_SAFE,
			"MH_APP_EXTENSION_SAFE",
		},
		{
			MH_NLIST_OUTOFSYNC_WITH_DYLDINFO,
			"MH_NLIST_OUTOFSYNC_WITH_DYLDINFO",
		},
		{
			MH_SIM_SUPPORT,
			"MH_SIM_SUPPORT",
		},
		{
			MH_DYLIB_IN_CACHE,
			"MH_DYLIB_IN_CACHE",
		},
	}
	s := ""
	for _, n := range names {
		if f&n.bit != 0 {
			if s != "" {
				s += "|"
			}

			s += n.name
		}
	}

	if rest := f &^ 0x8fffffff; rest != 0 {
		if s != "" {
			s += "|"
		}

		s += "0x" + strconv.FormatUint(uint64(rest), 16)
	}

	return s
}

// --- CPU ---

// CPU types: int32 with the signed ABI64/ABI64_32 flags in the high bits.
const (
	CPU_ARCH_ABI64    int32 = 0x01000000
	CPU_ARCH_ABI64_32 int32 = 0x02000000

	CPU_TYPE_X86       int32 = 7
	CPU_TYPE_I386      int32 = 7
	CPU_TYPE_X86_64    int32 = 7 | CPU_ARCH_ABI64
	CPU_TYPE_MC98000   int32 = 10
	CPU_TYPE_HPPA      int32 = 11
	CPU_TYPE_ARM       int32 = 12
	CPU_TYPE_ARM64     int32 = 12 | CPU_ARCH_ABI64
	CPU_TYPE_ARM64_32  int32 = 12 | CPU_ARCH_ABI64_32
	CPU_TYPE_MC88000   int32 = 13
	CPU_TYPE_SPARC     int32 = 14
	CPU_TYPE_I860      int32 = 15
	CPU_TYPE_POWERPC   int32 = 18
	CPU_TYPE_POWERPC64 int32 = 18 | CPU_ARCH_ABI64
)

// CPU_SUBTYPE: common subtypes.
const (
	CPU_SUBTYPE_X86_64_ALL   int32 = 3
	CPU_SUBTYPE_X86_64_H     int32 = 8
	CPU_SUBTYPE_ARM_ALL      int32 = 0
	CPU_SUBTYPE_ARM64_ALL    int32 = 0
	CPU_SUBTYPE_ARM64_V8     int32 = 1
	CPU_SUBTYPE_ARM64E       int32 = 2
	CPU_SUBTYPE_ARM64_32_ALL int32 = 0
	CPU_SUBTYPE_POWERPC_ALL  int32 = 0

	CPU_SUBTYPE_MASK                uint32 = 0xff000000 // capabilities bits
	CPU_SUBTYPE_LIB64               uint32 = 0x80000000
	CPU_SUBTYPE_ARM64_PTR_AUTH_MASK uint32 = 0x0f000000
)

func cpuName(t int32) string {
	switch t {
	case CPU_TYPE_X86:
		return "CPU_TYPE_X86"
	case CPU_TYPE_X86_64:
		return "CPU_TYPE_X86_64"
	case CPU_TYPE_ARM:
		return "CPU_TYPE_ARM"
	case CPU_TYPE_ARM64:
		return "CPU_TYPE_ARM64"
	case CPU_TYPE_ARM64_32:
		return "CPU_TYPE_ARM64_32"
	case CPU_TYPE_POWERPC:
		return "CPU_TYPE_POWERPC"
	case CPU_TYPE_POWERPC64:
		return "CPU_TYPE_POWERPC64"
	case CPU_TYPE_SPARC:
		return "CPU_TYPE_SPARC"
	}

	return "CPU_TYPE_<" + strconv.FormatInt(int64(t), 10) + ">"
}

// --- Load commands ---

// Cmd is the cmd field of a load command. The LC_REQ_DYLD bit (0x80000000)
// is part of the number: dyld needs to know the commands with this bit.
type Cmd uint32

const (
	LC_REQ_DYLD uint32 = 0x80000000

	LC_SEGMENT                  Cmd = 0x1
	LC_SYMTAB                   Cmd = 0x2
	LC_SYMSEG                   Cmd = 0x3 // obsolete
	LC_THREAD                   Cmd = 0x4
	LC_UNIXTHREAD               Cmd = 0x5
	LC_LOADFVMLIB               Cmd = 0x6
	LC_IDFVMLIB                 Cmd = 0x7
	LC_IDENT                    Cmd = 0x8 // obsolete
	LC_FVMFILE                  Cmd = 0x9
	LC_PREPAGE                  Cmd = 0xa
	LC_DYSYMTAB                 Cmd = 0xb
	LC_LOAD_DYLIB               Cmd = 0xc
	LC_ID_DYLIB                 Cmd = 0xd
	LC_LOAD_DYLINKER            Cmd = 0xe
	LC_ID_DYLINKER              Cmd = 0xf
	LC_PREBOUND_DYLIB           Cmd = 0x10
	LC_ROUTINES                 Cmd = 0x11
	LC_SUB_FRAMEWORK            Cmd = 0x12
	LC_SUB_UMBRELLA             Cmd = 0x13
	LC_SUB_CLIENT               Cmd = 0x14
	LC_SUB_LIBRARY              Cmd = 0x15
	LC_TWOLEVEL_HINTS           Cmd = 0x16
	LC_PREBIND_CKSUM            Cmd = 0x17
	LC_LOAD_WEAK_DYLIB          Cmd = 0x18 | Cmd(LC_REQ_DYLD)
	LC_SEGMENT_64               Cmd = 0x19
	LC_ROUTINES_64              Cmd = 0x1a
	LC_UUID                     Cmd = 0x1b
	LC_RPATH                    Cmd = 0x1c | Cmd(LC_REQ_DYLD)
	LC_CODE_SIGNATURE           Cmd = 0x1d
	LC_SEGMENT_SPLIT_INFO       Cmd = 0x1e
	LC_REEXPORT_DYLIB           Cmd = 0x1f | Cmd(LC_REQ_DYLD)
	LC_LAZY_LOAD_DYLIB          Cmd = 0x20
	LC_ENCRYPTION_INFO          Cmd = 0x21
	LC_DYLD_INFO                Cmd = 0x22
	LC_DYLD_INFO_ONLY           Cmd = 0x22 | Cmd(LC_REQ_DYLD)
	LC_LOAD_UPWARD_DYLIB        Cmd = 0x23 | Cmd(LC_REQ_DYLD)
	LC_VERSION_MIN_MACOSX       Cmd = 0x24
	LC_VERSION_MIN_IPHONEOS     Cmd = 0x25
	LC_FUNCTION_STARTS          Cmd = 0x26
	LC_DYLD_ENVIRONMENT         Cmd = 0x27
	LC_MAIN                     Cmd = 0x28 | Cmd(LC_REQ_DYLD)
	LC_DATA_IN_CODE             Cmd = 0x29
	LC_SOURCE_VERSION           Cmd = 0x2a
	LC_DYLIB_CODE_SIGN_DRS      Cmd = 0x2b
	LC_ENCRYPTION_INFO_64       Cmd = 0x2c
	LC_LINKER_OPTION            Cmd = 0x2d
	LC_LINKER_OPTIMIZATION_HINT Cmd = 0x2e
	LC_VERSION_MIN_TVOS         Cmd = 0x2f
	LC_VERSION_MIN_WATCHOS      Cmd = 0x30
	LC_NOTE                     Cmd = 0x31
	LC_BUILD_VERSION            Cmd = 0x32
	LC_DYLD_EXPORTS_TRIE        Cmd = 0x33 | Cmd(LC_REQ_DYLD)
	LC_DYLD_CHAINED_FIXUPS      Cmd = 0x34 | Cmd(LC_REQ_DYLD)
	LC_FILESET_ENTRY            Cmd = 0x35 | Cmd(LC_REQ_DYLD)
	LC_ATOM_INFO                Cmd = 0x36
	LC_FUNCTION_VARIANTS        Cmd = 0x37
	LC_FUNCTION_VARIANT_FIXUPS  Cmd = 0x38
	LC_TARGET_TRIPLE            Cmd = 0x39
)

var cmdNames = map[Cmd]string{
	LC_SEGMENT:                  "LC_SEGMENT",
	LC_SYMTAB:                   "LC_SYMTAB",
	LC_SYMSEG:                   "LC_SYMSEG",
	LC_THREAD:                   "LC_THREAD",
	LC_UNIXTHREAD:               "LC_UNIXTHREAD",
	LC_LOADFVMLIB:               "LC_LOADFVMLIB",
	LC_IDFVMLIB:                 "LC_IDFVMLIB",
	LC_IDENT:                    "LC_IDENT",
	LC_FVMFILE:                  "LC_FVMFILE",
	LC_PREPAGE:                  "LC_PREPAGE",
	LC_DYSYMTAB:                 "LC_DYSYMTAB",
	LC_LOAD_DYLIB:               "LC_LOAD_DYLIB",
	LC_ID_DYLIB:                 "LC_ID_DYLIB",
	LC_LOAD_DYLINKER:            "LC_LOAD_DYLINKER",
	LC_ID_DYLINKER:              "LC_ID_DYLINKER",
	LC_PREBOUND_DYLIB:           "LC_PREBOUND_DYLIB",
	LC_ROUTINES:                 "LC_ROUTINES",
	LC_SUB_FRAMEWORK:            "LC_SUB_FRAMEWORK",
	LC_SUB_UMBRELLA:             "LC_SUB_UMBRELLA",
	LC_SUB_CLIENT:               "LC_SUB_CLIENT",
	LC_SUB_LIBRARY:              "LC_SUB_LIBRARY",
	LC_TWOLEVEL_HINTS:           "LC_TWOLEVEL_HINTS",
	LC_PREBIND_CKSUM:            "LC_PREBIND_CKSUM",
	LC_LOAD_WEAK_DYLIB:          "LC_LOAD_WEAK_DYLIB",
	LC_SEGMENT_64:               "LC_SEGMENT_64",
	LC_ROUTINES_64:              "LC_ROUTINES_64",
	LC_UUID:                     "LC_UUID",
	LC_RPATH:                    "LC_RPATH",
	LC_CODE_SIGNATURE:           "LC_CODE_SIGNATURE",
	LC_SEGMENT_SPLIT_INFO:       "LC_SEGMENT_SPLIT_INFO",
	LC_REEXPORT_DYLIB:           "LC_REEXPORT_DYLIB",
	LC_LAZY_LOAD_DYLIB:          "LC_LAZY_LOAD_DYLIB",
	LC_ENCRYPTION_INFO:          "LC_ENCRYPTION_INFO",
	LC_DYLD_INFO:                "LC_DYLD_INFO",
	LC_DYLD_INFO_ONLY:           "LC_DYLD_INFO_ONLY",
	LC_LOAD_UPWARD_DYLIB:        "LC_LOAD_UPWARD_DYLIB",
	LC_VERSION_MIN_MACOSX:       "LC_VERSION_MIN_MACOSX",
	LC_VERSION_MIN_IPHONEOS:     "LC_VERSION_MIN_IPHONEOS",
	LC_FUNCTION_STARTS:          "LC_FUNCTION_STARTS",
	LC_DYLD_ENVIRONMENT:         "LC_DYLD_ENVIRONMENT",
	LC_MAIN:                     "LC_MAIN",
	LC_DATA_IN_CODE:             "LC_DATA_IN_CODE",
	LC_SOURCE_VERSION:           "LC_SOURCE_VERSION",
	LC_DYLIB_CODE_SIGN_DRS:      "LC_DYLIB_CODE_SIGN_DRS",
	LC_ENCRYPTION_INFO_64:       "LC_ENCRYPTION_INFO_64",
	LC_LINKER_OPTION:            "LC_LINKER_OPTION",
	LC_LINKER_OPTIMIZATION_HINT: "LC_LINKER_OPTIMIZATION_HINT",
	LC_VERSION_MIN_TVOS:         "LC_VERSION_MIN_TVOS",
	LC_VERSION_MIN_WATCHOS:      "LC_VERSION_MIN_WATCHOS",
	LC_NOTE:                     "LC_NOTE",
	LC_BUILD_VERSION:            "LC_BUILD_VERSION",
	LC_DYLD_EXPORTS_TRIE:        "LC_DYLD_EXPORTS_TRIE",
	LC_DYLD_CHAINED_FIXUPS:      "LC_DYLD_CHAINED_FIXUPS",
	LC_FILESET_ENTRY:            "LC_FILESET_ENTRY",
	LC_ATOM_INFO:                "LC_ATOM_INFO",
	LC_FUNCTION_VARIANTS:        "LC_FUNCTION_VARIANTS",
	LC_FUNCTION_VARIANT_FIXUPS:  "LC_FUNCTION_VARIANT_FIXUPS",
	LC_TARGET_TRIPLE:            "LC_TARGET_TRIPLE",
}

func (c Cmd) String() string {
	if s, ok := cmdNames[c]; ok {
		return s
	}

	return "LC_<" + strconv.FormatUint(uint64(c), 16) + ">"
}

// --- Segments ---

// SegFlag is the flags of segment_command(_64).
type SegFlag uint32

const (
	SG_HIGHVM              SegFlag = 0x1 // file is not mapped to the start of the segment
	SG_FVMLIB              SegFlag = 0x2
	SG_NORELOC             SegFlag = 0x4
	SG_PROTECTED_VERSION_1 SegFlag = 0x8
	SG_READ_ONLY           SegFlag = 0x10 // read-only after fixups
)

// --- Sections ---

// SectionType is the low byte of the section flags.
type SectionType uint8

const (
	S_REGULAR                             SectionType = 0x0
	S_ZEROFILL                            SectionType = 0x1
	S_CSTRING_LITERALS                    SectionType = 0x2
	S_4BYTE_LITERALS                      SectionType = 0x3
	S_8BYTE_LITERALS                      SectionType = 0x4
	S_LITERAL_POINTERS                    SectionType = 0x5
	S_NON_LAZY_SYMBOL_POINTERS            SectionType = 0x6
	S_LAZY_SYMBOL_POINTERS                SectionType = 0x7
	S_SYMBOL_STUBS                        SectionType = 0x8
	S_MOD_INIT_FUNC_POINTERS              SectionType = 0x9
	S_MOD_TERM_FUNC_POINTERS              SectionType = 0xa
	S_COALESCED                           SectionType = 0xb
	S_GB_ZEROFILL                         SectionType = 0xc
	S_INTERPOSING                         SectionType = 0xd
	S_16BYTE_LITERALS                     SectionType = 0xe
	S_DTRACE_DOF                          SectionType = 0xf
	S_LAZY_DYLIB_SYMBOL_POINTERS          SectionType = 0x10
	S_THREAD_LOCAL_REGULAR                SectionType = 0x11
	S_THREAD_LOCAL_ZEROFILL               SectionType = 0x12
	S_THREAD_LOCAL_VARIABLES              SectionType = 0x13
	S_THREAD_LOCAL_VARIABLE_POINTERS      SectionType = 0x14
	S_THREAD_LOCAL_INIT_FUNCTION_POINTERS SectionType = 0x15
	S_INIT_FUNC_OFFSETS                   SectionType = 0x16
)

func (t SectionType) String() string {
	switch t {
	case S_REGULAR:
		return "S_REGULAR"
	case S_ZEROFILL:
		return "S_ZEROFILL"
	case S_CSTRING_LITERALS:
		return "S_CSTRING_LITERALS"
	case S_4BYTE_LITERALS:
		return "S_4BYTE_LITERALS"
	case S_8BYTE_LITERALS:
		return "S_8BYTE_LITERALS"
	case S_LITERAL_POINTERS:
		return "S_LITERAL_POINTERS"
	case S_NON_LAZY_SYMBOL_POINTERS:
		return "S_NON_LAZY_SYMBOL_POINTERS"
	case S_LAZY_SYMBOL_POINTERS:
		return "S_LAZY_SYMBOL_POINTERS"
	case S_SYMBOL_STUBS:
		return "S_SYMBOL_STUBS"
	case S_MOD_INIT_FUNC_POINTERS:
		return "S_MOD_INIT_FUNC_POINTERS"
	case S_MOD_TERM_FUNC_POINTERS:
		return "S_MOD_TERM_FUNC_POINTERS"
	case S_COALESCED:
		return "S_COALESCED"
	case S_GB_ZEROFILL:
		return "S_GB_ZEROFILL"
	case S_INTERPOSING:
		return "S_INTERPOSING"
	case S_16BYTE_LITERALS:
		return "S_16BYTE_LITERALS"
	case S_DTRACE_DOF:
		return "S_DTRACE_DOF"
	case S_LAZY_DYLIB_SYMBOL_POINTERS:
		return "S_LAZY_DYLIB_SYMBOL_POINTERS"
	case S_THREAD_LOCAL_REGULAR:
		return "S_THREAD_LOCAL_REGULAR"
	case S_THREAD_LOCAL_ZEROFILL:
		return "S_THREAD_LOCAL_ZEROFILL"
	case S_THREAD_LOCAL_VARIABLES:
		return "S_THREAD_LOCAL_VARIABLES"
	case S_THREAD_LOCAL_VARIABLE_POINTERS:
		return "S_THREAD_LOCAL_VARIABLE_POINTERS"
	case S_THREAD_LOCAL_INIT_FUNCTION_POINTERS:
		return "S_THREAD_LOCAL_INIT_FUNCTION_POINTERS"
	case S_INIT_FUNC_OFFSETS:
		return "S_INIT_FUNC_OFFSETS"
	}

	return "S_<" + strconv.Itoa(int(t)) + ">"
}

// AttrFlag is the high bits of the section flags (S_ATTR_*).
type AttrFlag uint32

const (
	S_ATTR_PURE_INSTRUCTIONS   AttrFlag = 0x80000000
	S_ATTR_NO_TOC              AttrFlag = 0x40000000
	S_ATTR_STRIP_STATIC_SYMS   AttrFlag = 0x20000000
	S_ATTR_NO_DEAD_STRIP       AttrFlag = 0x10000000
	S_ATTR_LIVE_SUPPORT        AttrFlag = 0x08000000
	S_ATTR_SELF_MODIFYING_CODE AttrFlag = 0x04000000
	S_ATTR_DEBUG               AttrFlag = 0x02000000
	S_ATTR_SOME_INSTRUCTIONS   AttrFlag = 0x00000400
	S_ATTR_EXT_RELOC           AttrFlag = 0x00000200
	S_ATTR_LOC_RELOC           AttrFlag = 0x00000100
)

// --- nlist ---

// Bits of the nlist.n_type field.
const (
	N_STAB uint8 = 0xe0 // if set, a debug stab symbol (see below)
	N_PEXT uint8 = 0x10 // private external
	N_TYPE uint8 = 0x0e // type mask
	N_EXT  uint8 = 0x01 // external (global)
)

// Values of the n_type & N_TYPE field.
const (
	N_UNDF uint8 = 0x0 // undefined (external)
	N_ABS  uint8 = 0x2 // absolute (segment 0)
	N_SECT uint8 = 0xe // defined in section n_sect
	N_PBUD uint8 = 0xc // prebound undefined
	N_INDR uint8 = 0xa // indirect
)

// Stab is the stab symbol types (N_STAB set).
type Stab uint8

const (
	N_GSYM    Stab = 0x20 // global symbol
	N_FNAME   Stab = 0x22 // procedure name (f77)
	N_FUN     Stab = 0x24 // procedure
	N_STSYM   Stab = 0x26 // static symbol
	N_LCSYM   Stab = 0x28 // .l
	N_BNSYM   Stab = 0x2e // begin nsect sym
	N_PC      Stab = 0x30 // global pascal symbol
	N_AST     Stab = 0x32
	N_OPT     Stab = 0x3c
	N_RSYM    Stab = 0x40 // register symbol
	N_SLINE   Stab = 0x44 // src line
	N_ENSYM   Stab = 0x4e // end nsect sym
	N_SSYM    Stab = 0x60 // struct elt
	N_SO      Stab = 0x64 // source file
	N_OSO     Stab = 0x66 // object file
	N_LSYM    Stab = 0x80 // local sym
	N_BINCL   Stab = 0x82
	N_SOL     Stab = 0x84
	N_PARAMS  Stab = 0x86
	N_VERSION Stab = 0x88
	N_OLEVEL  Stab = 0x8a
	N_PSYM    Stab = 0xa0 // parameter
	N_EINCL   Stab = 0xa2
	N_ENTRY   Stab = 0xa4
	N_LBRAC   Stab = 0xc0
	N_EXCL    Stab = 0xc2
	N_RBRAC   Stab = 0xe0
	N_BCOMM   Stab = 0xe2
	N_ECOMM   Stab = 0xe4
	N_ECOML   Stab = 0xe8
	N_LENG    Stab = 0xfe
)

var stabNames = map[Stab]string{
	N_GSYM: "N_GSYM", N_FNAME: "N_FNAME", N_FUN: "N_FUN", N_STSYM: "N_STSYM",
	N_LCSYM: "N_LCSYM", N_BNSYM: "N_BNSYM", N_PC: "N_PC", N_AST: "N_AST",
	N_OPT: "N_OPT", N_RSYM: "N_RSYM", N_SLINE: "N_SLINE", N_ENSYM: "N_ENSYM",
	N_SSYM: "N_SSYM", N_SO: "N_SO", N_OSO: "N_OSO", N_LSYM: "N_LSYM",
	N_BINCL: "N_BINCL", N_SOL: "N_SOL", N_PARAMS: "N_PARAMS", N_VERSION: "N_VERSION",
	N_OLEVEL: "N_OLEVEL", N_PSYM: "N_PSYM", N_EINCL: "N_EINCL", N_ENTRY: "N_ENTRY",
	N_LBRAC: "N_LBRAC", N_EXCL: "N_EXCL", N_RBRAC: "N_RBRAC", N_BCOMM: "N_BCOMM",
	N_ECOMM: "N_ECOMM", N_ECOML: "N_ECOML", N_LENG: "N_LENG",
}

func (s Stab) String() string {
	if n, ok := stabNames[s]; ok {
		return n
	}

	return "N_<" + strconv.FormatUint(uint64(s), 16) + ">"
}

// Bits of the nlist.n_desc field.
const (
	N_ARM_THUMB_DEF        uint16 = 0x0008
	REFERENCED_DYNAMICALLY uint16 = 0x0010
	N_DESC_DISCARDED       uint16 = 0x0020
	N_WEAK_REF             uint16 = 0x0040
	N_WEAK_DEF             uint16 = 0x0080
	N_SYMBOL_RESOLVER      uint16 = 0x0100
	N_ALT_ENTRY            uint16 = 0x0200
	N_COLD_FUNC            uint16 = 0x0400

	GET_LIBRARY_ORDINAL_MASK uint16 = 0xff00 // library ordinal in the high byte
	SELF_LIBRARY_ORDINAL     uint16 = 0x0000
	DYNAMIC_LOOKUP_ORDINAL   uint16 = 0xfe00
	EXECUTABLE_ORDINAL       uint16 = 0xff00
)

// LibraryOrdinal extracts the library ordinal from n_desc.
func LibraryOrdinal(desc uint16) int {
	return int(desc&GET_LIBRARY_ORDINAL_MASK) >> 8
}

// --- Relocations ---

// RelocType is a relocation type (the architecture-dependent part of r_info).
type RelocType uint32

// ARM64 (enum reloc_type_arm64).
const (
	ARM64_RELOC_UNSIGNED              RelocType = 0 // pointers (64/32-bit)
	ARM64_RELOC_SUBTRACTOR            RelocType = 1 // followed by UNSIGNED
	ARM64_RELOC_BRANCH26              RelocType = 2 // B/BL
	ARM64_RELOC_PAGE21                RelocType = 3 // ADRP
	ARM64_RELOC_PAGEOFF12             RelocType = 4 // offset within page
	ARM64_RELOC_GOT_LOAD_PAGE21       RelocType = 5 // ADRP of GOT
	ARM64_RELOC_GOT_LOAD_PAGEOFF12    RelocType = 6 // LDR from GOT
	ARM64_RELOC_POINTER_TO_GOT        RelocType = 7
	ARM64_RELOC_TLVP_LOAD_PAGE21      RelocType = 8
	ARM64_RELOC_TLVP_LOAD_PAGEOFF12   RelocType = 9
	ARM64_RELOC_ADDEND                RelocType = 10 // followed by PAGE21/PAGEOFF12
	ARM64_RELOC_AUTHENTICATED_POINTER RelocType = 11 // arm64e
)

// Generic (x86 etc.).
const (
	GENERIC_RELOC_VANILLA        RelocType = 0
	GENERIC_RELOC_PAIR           RelocType = 1
	GENERIC_RELOC_SECTDIFF       RelocType = 2
	GENERIC_RELOC_PB_LAZY_PTR    RelocType = 3
	GENERIC_RELOC_LOCAL_SECTDIFF RelocType = 4
	GENERIC_RELOC_TLV            RelocType = 5
)

// R_SCATTERED is the scattered relocation bit in r0 (obsolete mechanism).
const R_SCATTERED uint32 = 0x80000000

func arm64RelocName(r RelocType) string {
	switch r {
	case ARM64_RELOC_UNSIGNED:
		return "ARM64_RELOC_UNSIGNED"
	case ARM64_RELOC_SUBTRACTOR:
		return "ARM64_RELOC_SUBTRACTOR"
	case ARM64_RELOC_BRANCH26:
		return "ARM64_RELOC_BRANCH26"
	case ARM64_RELOC_PAGE21:
		return "ARM64_RELOC_PAGE21"
	case ARM64_RELOC_PAGEOFF12:
		return "ARM64_RELOC_PAGEOFF12"
	case ARM64_RELOC_GOT_LOAD_PAGE21:
		return "ARM64_RELOC_GOT_LOAD_PAGE21"
	case ARM64_RELOC_GOT_LOAD_PAGEOFF12:
		return "ARM64_RELOC_GOT_LOAD_PAGEOFF12"
	case ARM64_RELOC_POINTER_TO_GOT:
		return "ARM64_RELOC_POINTER_TO_GOT"
	case ARM64_RELOC_TLVP_LOAD_PAGE21:
		return "ARM64_RELOC_TLVP_LOAD_PAGE21"
	case ARM64_RELOC_TLVP_LOAD_PAGEOFF12:
		return "ARM64_RELOC_TLVP_LOAD_PAGEOFF12"
	case ARM64_RELOC_ADDEND:
		return "ARM64_RELOC_ADDEND"
	case ARM64_RELOC_AUTHENTICATED_POINTER:
		return "ARM64_RELOC_AUTHENTICATED_POINTER"
	}

	return "ARM64_RELOC_<" + strconv.FormatUint(uint64(r), 10) + ">"
}

// --- dyld: rebase/bind opcodes and export flags ---

// Rebase opcodes (the LC_DYLD_INFO.rebase_off stream).
const (
	REBASE_TYPE_POINTER         uint8 = 0x1
	REBASE_TYPE_TEXT_ABSOLUTE32 uint8 = 0x2
	REBASE_TYPE_TEXT_PCREL32    uint8 = 0x3

	REBASE_OPCODE_MASK                               uint8 = 0xF0
	REBASE_OPCODE_DONE                               uint8 = 0x00
	REBASE_OPCODE_SET_TYPE_IMM                       uint8 = 0x10
	REBASE_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB        uint8 = 0x20
	REBASE_OPCODE_ADD_ADDR_ULEB                      uint8 = 0x30
	REBASE_OPCODE_ADD_ADDR_IMM_SCALED                uint8 = 0x40
	REBASE_OPCODE_DO_REBASE_IMM_TIMES                uint8 = 0x50
	REBASE_OPCODE_DO_REBASE_ULEB_TIMES               uint8 = 0x60
	REBASE_OPCODE_DO_REBASE_ADD_ADDR_ULEB            uint8 = 0x70
	REBASE_OPCODE_DO_REBASE_ULEB_TIMES_SKIPPING_ULEB uint8 = 0x80
)

// Bind opcodes (the bind/weak_bind/lazy_bind streams).
const (
	BIND_TYPE_POINTER         uint8 = 0x1
	BIND_TYPE_TEXT_ABSOLUTE32 uint8 = 0x2
	BIND_TYPE_TEXT_PCREL32    uint8 = 0x3

	BIND_OPCODE_MASK                                     uint8 = 0xF0
	BIND_OPCODE_DONE                                     uint8 = 0x00
	BIND_OPCODE_SET_DYLIB_ORDINAL_IMM                    uint8 = 0x10
	BIND_OPCODE_SET_DYLIB_ORDINAL_ULEB                   uint8 = 0x20
	BIND_OPCODE_SET_DYLIB_SPECIAL_IMM                    uint8 = 0x30
	BIND_OPCODE_SET_SYMBOL_TRAILING_FLAGS_IMM            uint8 = 0x40
	BIND_OPCODE_SET_TYPE_IMM                             uint8 = 0x50
	BIND_OPCODE_SET_ADDEND_SLEB                          uint8 = 0x60
	BIND_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB              uint8 = 0x70
	BIND_OPCODE_ADD_ADDR_ULEB                            uint8 = 0x80
	BIND_OPCODE_DO_BIND                                  uint8 = 0x90
	BIND_OPCODE_DO_BIND_ADD_ADDR_ULEB                    uint8 = 0xA0
	BIND_OPCODE_DO_BIND_ADD_ADDR_IMM_SCALED              uint8 = 0xB0
	BIND_OPCODE_DO_BIND_ULEB_TIMES_SKIPPING_ULEB         uint8 = 0xC0
	BIND_SUBOPCODE_THREADED                              uint8 = 0xD0 // subopcode: chained fixups in the bind stream
	BIND_SUBOPCODE_THREADED_APPLY_SET_BIND_ORDINAL_TABLE uint8 = 0x00 // with BIND_SUBOPCODE_THREADED
	BIND_SUBOPCODE_THREADED_APPLY_APPLY                  uint8 = 0x01
)

// Flags of exports trie nodes (LC_DYLD_INFO.export_info / LC_DYLD_EXPORTS_TRIE).
const (
	EXPORT_SYMBOL_FLAGS_KIND_MASK         uint64 = 0x03
	EXPORT_SYMBOL_FLAGS_KIND_REGULAR      uint64 = 0x00
	EXPORT_SYMBOL_FLAGS_KIND_THREAD_LOCAL uint64 = 0x01
	EXPORT_SYMBOL_FLAGS_KIND_ABSOLUTE     uint64 = 0x02
	EXPORT_SYMBOL_FLAGS_WEAK_DEFINITION   uint64 = 0x04
	EXPORT_SYMBOL_FLAGS_REEXPORT          uint64 = 0x08
	EXPORT_SYMBOL_FLAGS_STUB_AND_RESOLVER uint64 = 0x10
	EXPORT_SYMBOL_FLAGS_STATIC_RESOLVER   uint64 = 0x20
)

// --- Data-in-code (DICE) ---

// DiceKind is a data-in-code entry type.
type DiceKind uint16

const (
	DICE_KIND_DATA             DiceKind = 0x0001
	DICE_KIND_JUMP_TABLE8      DiceKind = 0x0002
	DICE_KIND_JUMP_TABLE16     DiceKind = 0x0003
	DICE_KIND_JUMP_TABLE32     DiceKind = 0x0004
	DICE_KIND_ABS_JUMP_TABLE32 DiceKind = 0x0005
)

func (k DiceKind) String() string {
	switch k {
	case DICE_KIND_DATA:
		return "DICE_KIND_DATA"
	case DICE_KIND_JUMP_TABLE8:
		return "DICE_KIND_JUMP_TABLE8"
	case DICE_KIND_JUMP_TABLE16:
		return "DICE_KIND_JUMP_TABLE16"
	case DICE_KIND_JUMP_TABLE32:
		return "DICE_KIND_JUMP_TABLE32"
	case DICE_KIND_ABS_JUMP_TABLE32:
		return "DICE_KIND_ABS_JUMP_TABLE32"
	}

	return "DICE_KIND_<" + strconv.FormatUint(uint64(k), 16) + ">"
}

// --- LC_BUILD_VERSION: platforms and tools ---

// Platform is the target build platform.
type Platform uint32

const (
	PLATFORM_UNKNOWN           Platform = 0
	PLATFORM_MACOS             Platform = 1
	PLATFORM_IOS               Platform = 2
	PLATFORM_TVOS              Platform = 3
	PLATFORM_WATCHOS           Platform = 4
	PLATFORM_BRIDGEOS          Platform = 5
	PLATFORM_MACCATALYST       Platform = 6
	PLATFORM_IOSSIMULATOR      Platform = 7
	PLATFORM_TVOSSIMULATOR     Platform = 8
	PLATFORM_WATCHOSSIMULATOR  Platform = 9
	PLATFORM_DRIVERKIT         Platform = 10
	PLATFORM_VISIONOS          Platform = 11
	PLATFORM_VISIONOSSIMULATOR Platform = 12
)

func (p Platform) String() string {
	switch p {
	case PLATFORM_MACOS:
		return "PLATFORM_MACOS"
	case PLATFORM_IOS:
		return "PLATFORM_IOS"
	case PLATFORM_TVOS:
		return "PLATFORM_TVOS"
	case PLATFORM_WATCHOS:
		return "PLATFORM_WATCHOS"
	case PLATFORM_BRIDGEOS:
		return "PLATFORM_BRIDGEOS"
	case PLATFORM_MACCATALYST:
		return "PLATFORM_MACCATALYST"
	case PLATFORM_IOSSIMULATOR:
		return "PLATFORM_IOSSIMULATOR"
	case PLATFORM_TVOSSIMULATOR:
		return "PLATFORM_TVOSSIMULATOR"
	case PLATFORM_WATCHOSSIMULATOR:
		return "PLATFORM_WATCHOSSIMULATOR"
	case PLATFORM_DRIVERKIT:
		return "PLATFORM_DRIVERKIT"
	case PLATFORM_VISIONOS:
		return "PLATFORM_VISIONOS"
	case PLATFORM_VISIONOSSIMULATOR:
		return "PLATFORM_VISIONOSSIMULATOR"
	default:
		return "PLATFORM_<" + strconv.FormatUint(uint64(p), 10) + ">"
	}
}

// Tool is a build tool (LC_BUILD_VERSION).
type Tool uint32

const (
	TOOL_CLANG Tool = 1
	TOOL_SWIFT Tool = 2
	TOOL_LD    Tool = 3
	TOOL_LLD   Tool = 4
)

func (t Tool) String() string {
	switch t {
	case TOOL_CLANG:
		return "TOOL_CLANG"
	case TOOL_SWIFT:
		return "TOOL_SWIFT"
	case TOOL_LD:
		return "TOOL_LD"
	case TOOL_LLD:
		return "TOOL_LLD"
	}

	return "TOOL_<" + strconv.FormatUint(uint64(t), 10) + ">"
}
