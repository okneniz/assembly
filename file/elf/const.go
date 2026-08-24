// Package elf is a self-contained ELF format parser (32/64-bit, LE/BE), built
// on the okneniz/parsec parser combinators over a single positional buffer,
// without using debug/elf. The package knows nothing about Mach-O or the rest
// of assembly: it owns all of its types.
//
// The constants come from the "Tool Interface Standard (TIS) Executable and
// Linking Format (ELF) Specification" plus GNU/Linux additions. The
// AArch64/RISC-V relocation numbers are verified against binutils (readelf)
// and LLVM.
package elf

import (
	"strconv"
	"strings"
)

// --- e_ident ---

// Indices of fields in the e_ident array.
const (
	EI_MAG0       = 0
	EI_MAG1       = 1
	EI_MAG2       = 2
	EI_MAG3       = 3
	EI_CLASS      = 4
	EI_DATA       = 5
	EI_VERSION    = 6
	EI_OSABI      = 7
	EI_ABIVERSION = 8
	EI_PAD        = 9
	EI_NIDENT     = 16
)

// ELF magic bytes.
const (
	Mag0 = 0x7f
	Mag1 = 'E'
	Mag2 = 'L'
	Mag3 = 'F'
)

// Class is the file class (32 or 64 bits), e_ident[EI_CLASS].
type Class uint8

const (
	CLASSNONE Class = 0
	CLASS32   Class = 1
	CLASS64   Class = 2
)

func (c Class) String() string {
	switch c {
	case CLASS32:
		return "ELFCLASS32"
	case CLASS64:
		return "ELFCLASS64"
	default:
		return "ELFCLASS<" + strconv.Itoa(int(c)) + ">"
	}
}

// Endianness of the file, e_ident[EI_DATA].
const (
	DATA2LSB = 1 // little-endian
	DATA2MSB = 2 // big-endian
)

// OSABI is the operating system/ABI, e_ident[EI_OSABI].
type OSABI uint8

const (
	OSABI_NONE       OSABI = 0   // SysV
	OSABI_HPUX       OSABI = 1   // Hewlett-Packard HP-UX
	OSABI_NETBSD     OSABI = 2   // NetBSD
	OSABI_GNU        OSABI = 3   // GNU/Linux
	OSABI_SOLARIS    OSABI = 6   // Sun Solaris
	OSABI_AIX        OSABI = 7   // AIX
	OSABI_IRIX       OSABI = 8   // IRIX
	OSABI_FREEBSD    OSABI = 9   // FreeBSD
	OSABI_TRU64      OSABI = 10  // Compaq TRU64 UNIX
	OSABI_MODESTO    OSABI = 11  // Novell Modesto
	OSABI_OPENBSD    OSABI = 12  // OpenBSD
	OSABI_ARM_AEABI  OSABI = 64  // ARM EABI
	OSABI_ARM        OSABI = 97  // ARM
	OSABI_STANDALONE OSABI = 255 // Standalone (embedded)
)

func (o OSABI) String() string {
	switch o {
	case OSABI_NONE:
		return "ELFOSABI_NONE"
	case OSABI_HPUX:
		return "ELFOSABI_HPUX"
	case OSABI_NETBSD:
		return "ELFOSABI_NETBSD"
	case OSABI_GNU:
		return "ELFOSABI_GNU"
	case OSABI_SOLARIS:
		return "ELFOSABI_SOLARIS"
	case OSABI_AIX:
		return "ELFOSABI_AIX"
	case OSABI_IRIX:
		return "ELFOSABI_IRIX"
	case OSABI_FREEBSD:
		return "ELFOSABI_FREEBSD"
	case OSABI_TRU64:
		return "ELFOSABI_TRU64"
	case OSABI_MODESTO:
		return "ELFOSABI_MODESTO"
	case OSABI_OPENBSD:
		return "ELFOSABI_OPENBSD"
	case OSABI_ARM_AEABI:
		return "ELFOSABI_ARM_AEABI"
	case OSABI_ARM:
		return "ELFOSABI_ARM"
	case OSABI_STANDALONE:
		return "ELFOSABI_STANDALONE"
	}

	return "ELFOSABI<" + strconv.Itoa(int(o)) + ">"
}

// --- e_type ---

// Type is the object type, e_type.
type Type uint16

const (
	ET_NONE   Type = 0
	ET_REL    Type = 1 // relocatable object (.o)
	ET_EXEC   Type = 2 // executable file
	ET_DYN    Type = 3 // shared object
	ET_CORE   Type = 4 // core dump
	ET_LOOS   Type = 0xfe00
	ET_HIOS   Type = 0xfeff
	ET_LOPROC Type = 0xff00
	ET_HIPROC Type = 0xffff
)

func (t Type) String() string {
	switch t {
	case ET_NONE:
		return "ET_NONE"
	case ET_REL:
		return "ET_REL"
	case ET_EXEC:
		return "ET_EXEC"
	case ET_DYN:
		return "ET_DYN"
	case ET_CORE:
		return "ET_CORE"
	default:
		if t >= ET_LOOS && t <= ET_HIOS {
			return "ET_LOOS+" + strconv.FormatUint(uint64(t-ET_LOOS), 16)
		}

		if t >= ET_LOPROC && t <= ET_HIPROC {
			return "ET_LOPROC+" + strconv.FormatUint(uint64(t-ET_LOPROC), 16)
		}

		return "ET<" + strconv.Itoa(int(t)) + ">"
	}
}

// --- e_machine (the commonly used part of the registry; the gaps are rare/historical ones) ---

// Machine is the machine architecture, e_machine.
type Machine uint16

const (
	EM_NONE          Machine = 0
	EM_M32           Machine = 1
	EM_SPARC         Machine = 2
	EM_386           Machine = 3
	EM_68K           Machine = 4
	EM_88K           Machine = 5
	EM_IAMCU         Machine = 6
	EM_860           Machine = 7
	EM_MIPS          Machine = 8
	EM_S370          Machine = 9
	EM_MIPS_RS3_LE   Machine = 10
	EM_PARISC        Machine = 15
	EM_VPP500        Machine = 17
	EM_SPARC32PLUS   Machine = 18
	EM_960           Machine = 19
	EM_PPC           Machine = 20
	EM_PPC64         Machine = 21
	EM_S390          Machine = 22
	EM_SPU           Machine = 23
	EM_V800          Machine = 36
	EM_FR20          Machine = 37
	EM_RH32          Machine = 38
	EM_RCE           Machine = 39
	EM_ARM           Machine = 40
	EM_SH            Machine = 42
	EM_SPARCV9       Machine = 43
	EM_TRICORE       Machine = 44
	EM_ARC           Machine = 45
	EM_H8_300        Machine = 46
	EM_H8_300H       Machine = 47
	EM_H8S           Machine = 48
	EM_H8_500        Machine = 49
	EM_IA_64         Machine = 50
	EM_MIPS_X        Machine = 51
	EM_COLDFIRE      Machine = 52
	EM_68HC12        Machine = 53
	EM_MMA           Machine = 54
	EM_PCP           Machine = 55
	EM_NCPU          Machine = 56
	EM_NDR1          Machine = 57
	EM_STARCORE      Machine = 58
	EM_ME16          Machine = 59
	EM_ST100         Machine = 60
	EM_TINYJ         Machine = 61
	EM_X86_64        Machine = 62
	EM_PDP11         Machine = 65
	EM_FX66          Machine = 66
	EM_ST9PLUS       Machine = 67
	EM_ST7           Machine = 68
	EM_68HC16        Machine = 69
	EM_68HC11        Machine = 70
	EM_68HC08        Machine = 71
	EM_68HC05        Machine = 72
	EM_SVX           Machine = 73
	EM_ST19          Machine = 74
	EM_VAX           Machine = 75
	EM_CRIS          Machine = 76
	EM_JAVELIN       Machine = 77
	EM_FIREPATH      Machine = 78
	EM_ZSP           Machine = 79
	EM_MMIX          Machine = 80
	EM_HUANY         Machine = 81
	EM_PRISM         Machine = 82
	EM_AVR           Machine = 83
	EM_FR30          Machine = 84
	EM_D10V          Machine = 85
	EM_D30V          Machine = 86
	EM_V850          Machine = 87
	EM_M32R          Machine = 88
	EM_MN10300       Machine = 89
	EM_MN10200       Machine = 90
	EM_PJ            Machine = 91
	EM_OPENRISC      Machine = 92
	EM_ARC_COMPACT   Machine = 93
	EM_XTENSA        Machine = 94
	EM_VIDEOCORE     Machine = 95
	EM_TMM_GPP       Machine = 96
	EM_NS32K         Machine = 97
	EM_TPC           Machine = 98
	EM_SNP1K         Machine = 99
	EM_ST200         Machine = 100
	EM_IP2K          Machine = 101
	EM_MAX           Machine = 102
	EM_CR            Machine = 103
	EM_F2MC16        Machine = 104
	EM_MSP430        Machine = 105
	EM_BLACKFIN      Machine = 106
	EM_SE_C33        Machine = 107
	EM_SEP           Machine = 108
	EM_ARCA          Machine = 109
	EM_UNICORE       Machine = 110
	EM_EXCESS        Machine = 111
	EM_DXP           Machine = 112
	EM_ALTERA_NIOS2  Machine = 113
	EM_CRX           Machine = 114
	EM_XGATE         Machine = 115
	EM_C166          Machine = 116
	EM_M16C          Machine = 117
	EM_DSPIC30F      Machine = 118
	EM_CE            Machine = 119
	EM_M32C          Machine = 120
	EM_TSK3000       Machine = 131
	EM_RS08          Machine = 132
	EM_SHARC         Machine = 133
	EM_ECOG2         Machine = 134
	EM_SCORE7        Machine = 135
	EM_DSP24         Machine = 136
	EM_VIDEOCORE3    Machine = 137
	EM_LATTICEMICO32 Machine = 138
	EM_SE_C17        Machine = 139
	EM_TI_C6000      Machine = 140
	EM_TI_C2000      Machine = 141
	EM_TI_C5500      Machine = 142
	EM_TI_ARP32      Machine = 143
	EM_TI_PRU        Machine = 144
	EM_MMDSP_PLUS    Machine = 160
	EM_CYPRESS_M8C   Machine = 161
	EM_R32C          Machine = 162
	EM_TRIMEDIA      Machine = 163
	EM_QDSP6         Machine = 164
	EM_8051          Machine = 165
	EM_STXP7X        Machine = 166
	EM_NDS32         Machine = 167
	EM_ECOG1X        Machine = 168
	EM_MAXQ30        Machine = 169
	EM_XIMO16        Machine = 170
	EM_MANIK         Machine = 171
	EM_CRAYNV2       Machine = 172
	EM_RX            Machine = 173
	EM_METAG         Machine = 174
	EM_MCST_ELBRUS   Machine = 175
	EM_ECOG16        Machine = 176
	EM_CR16          Machine = 177
	EM_ETPU          Machine = 178
	EM_SLE9X         Machine = 179
	EM_L10M          Machine = 180
	EM_K10M          Machine = 181
	EM_AARCH64       Machine = 183
	EM_AVR32         Machine = 185
	EM_STM8          Machine = 186
	EM_TILE64        Machine = 187
	EM_TILEPRO       Machine = 188
	EM_MICROBLAZE    Machine = 189
	EM_CUDA          Machine = 190
	EM_TILEGX        Machine = 191
	EM_CLOUDSHIELD   Machine = 192
	EM_COREA_1ST     Machine = 193
	EM_COREA_2ND     Machine = 194
	EM_ARC_COMPACT2  Machine = 195
	EM_OPEN8         Machine = 196
	EM_RL78          Machine = 197
	EM_VIDEOCORE5    Machine = 198
	EM_78KOR         Machine = 199
	EM_56800EX       Machine = 200
	EM_BA1           Machine = 201
	EM_BA2           Machine = 202
	EM_XCORE         Machine = 203
	EM_MCHP_PIC      Machine = 204
	EM_KM32          Machine = 210
	EM_KMX32         Machine = 211
	EM_EMX16         Machine = 212
	EM_EMX8          Machine = 213
	EM_KVARC         Machine = 214
	EM_CDP           Machine = 215
	EM_COGE          Machine = 216
	EM_COOL          Machine = 217
	EM_NORC          Machine = 218
	EM_CSR_KALIMBA   Machine = 219
	EM_Z80           Machine = 220
	EM_VISIUM        Machine = 221
	EM_FT32          Machine = 222
	EM_MOXIE         Machine = 223
	EM_AMDGPU        Machine = 224
	EM_RISCV         Machine = 243
	EM_BPF           Machine = 247
	EM_C_SKY         Machine = 252
)

var machNames = map[Machine]string{
	EM_NONE:          "EM_NONE",
	EM_M32:           "EM_M32",
	EM_SPARC:         "EM_SPARC",
	EM_386:           "EM_386",
	EM_68K:           "EM_68K",
	EM_88K:           "EM_88K",
	EM_IAMCU:         "EM_IAMCU",
	EM_860:           "EM_860",
	EM_MIPS:          "EM_MIPS",
	EM_S370:          "EM_S370",
	EM_MIPS_RS3_LE:   "EM_MIPS_RS3_LE",
	EM_PARISC:        "EM_PARISC",
	EM_VPP500:        "EM_VPP500",
	EM_SPARC32PLUS:   "EM_SPARC32PLUS",
	EM_960:           "EM_960",
	EM_PPC:           "EM_PPC",
	EM_PPC64:         "EM_PPC64",
	EM_S390:          "EM_S390",
	EM_SPU:           "EM_SPU",
	EM_V800:          "EM_V800",
	EM_FR20:          "EM_FR20",
	EM_RH32:          "EM_RH32",
	EM_RCE:           "EM_RCE",
	EM_ARM:           "EM_ARM",
	EM_SH:            "EM_SH",
	EM_SPARCV9:       "EM_SPARCV9",
	EM_TRICORE:       "EM_TRICORE",
	EM_ARC:           "EM_ARC",
	EM_H8_300:        "EM_H8_300",
	EM_H8_300H:       "EM_H8_300H",
	EM_H8S:           "EM_H8S",
	EM_H8_500:        "EM_H8_500",
	EM_IA_64:         "EM_IA_64",
	EM_MIPS_X:        "EM_MIPS_X",
	EM_COLDFIRE:      "EM_COLDFIRE",
	EM_68HC12:        "EM_68HC12",
	EM_MMA:           "EM_MMA",
	EM_PCP:           "EM_PCP",
	EM_NCPU:          "EM_NCPU",
	EM_NDR1:          "EM_NDR1",
	EM_STARCORE:      "EM_STARCORE",
	EM_ME16:          "EM_ME16",
	EM_ST100:         "EM_ST100",
	EM_TINYJ:         "EM_TINYJ",
	EM_X86_64:        "EM_X86_64",
	EM_PDP11:         "EM_PDP11",
	EM_FX66:          "EM_FX66",
	EM_ST9PLUS:       "EM_ST9PLUS",
	EM_ST7:           "EM_ST7",
	EM_68HC16:        "EM_68HC16",
	EM_68HC11:        "EM_68HC11",
	EM_68HC08:        "EM_68HC08",
	EM_68HC05:        "EM_68HC05",
	EM_SVX:           "EM_SVX",
	EM_ST19:          "EM_ST19",
	EM_VAX:           "EM_VAX",
	EM_CRIS:          "EM_CRIS",
	EM_JAVELIN:       "EM_JAVELIN",
	EM_FIREPATH:      "EM_FIREPATH",
	EM_ZSP:           "EM_ZSP",
	EM_MMIX:          "EM_MMIX",
	EM_HUANY:         "EM_HUANY",
	EM_PRISM:         "EM_PRISM",
	EM_AVR:           "EM_AVR",
	EM_FR30:          "EM_FR30",
	EM_D10V:          "EM_D10V",
	EM_D30V:          "EM_D30V",
	EM_V850:          "EM_V850",
	EM_M32R:          "EM_M32R",
	EM_MN10300:       "EM_MN10300",
	EM_MN10200:       "EM_MN10200",
	EM_PJ:            "EM_PJ",
	EM_OPENRISC:      "EM_OPENRISC",
	EM_ARC_COMPACT:   "EM_ARC_COMPACT",
	EM_XTENSA:        "EM_XTENSA",
	EM_VIDEOCORE:     "EM_VIDEOCORE",
	EM_TMM_GPP:       "EM_TMM_GPP",
	EM_NS32K:         "EM_NS32K",
	EM_TPC:           "EM_TPC",
	EM_SNP1K:         "EM_SNP1K",
	EM_ST200:         "EM_ST200",
	EM_IP2K:          "EM_IP2K",
	EM_MAX:           "EM_MAX",
	EM_CR:            "EM_CR",
	EM_F2MC16:        "EM_F2MC16",
	EM_MSP430:        "EM_MSP430",
	EM_BLACKFIN:      "EM_BLACKFIN",
	EM_SE_C33:        "EM_SE_C33",
	EM_SEP:           "EM_SEP",
	EM_ARCA:          "EM_ARCA",
	EM_UNICORE:       "EM_UNICORE",
	EM_EXCESS:        "EM_EXCESS",
	EM_DXP:           "EM_DXP",
	EM_ALTERA_NIOS2:  "EM_ALTERA_NIOS2",
	EM_CRX:           "EM_CRX",
	EM_XGATE:         "EM_XGATE",
	EM_C166:          "EM_C166",
	EM_M16C:          "EM_M16C",
	EM_DSPIC30F:      "EM_DSPIC30F",
	EM_CE:            "EM_CE",
	EM_M32C:          "EM_M32C",
	EM_TSK3000:       "EM_TSK3000",
	EM_RS08:          "EM_RS08",
	EM_SHARC:         "EM_SHARC",
	EM_ECOG2:         "EM_ECOG2",
	EM_SCORE7:        "EM_SCORE7",
	EM_DSP24:         "EM_DSP24",
	EM_VIDEOCORE3:    "EM_VIDEOCORE3",
	EM_LATTICEMICO32: "EM_LATTICEMICO32",
	EM_SE_C17:        "EM_SE_C17",
	EM_TI_C6000:      "EM_TI_C6000",
	EM_TI_C2000:      "EM_TI_C2000",
	EM_TI_C5500:      "EM_TI_C5500",
	EM_TI_ARP32:      "EM_TI_ARP32",
	EM_TI_PRU:        "EM_TI_PRU",
	EM_MMDSP_PLUS:    "EM_MMDSP_PLUS",
	EM_CYPRESS_M8C:   "EM_CYPRESS_M8C",
	EM_R32C:          "EM_R32C",
	EM_TRIMEDIA:      "EM_TRIMEDIA",
	EM_QDSP6:         "EM_QDSP6",
	EM_8051:          "EM_8051",
	EM_STXP7X:        "EM_STXP7X",
	EM_NDS32:         "EM_NDS32",
	EM_ECOG1X:        "EM_ECOG1X",
	EM_MAXQ30:        "EM_MAXQ30",
	EM_XIMO16:        "EM_XIMO16",
	EM_MANIK:         "EM_MANIK",
	EM_CRAYNV2:       "EM_CRAYNV2",
	EM_RX:            "EM_RX",
	EM_METAG:         "EM_METAG",
	EM_MCST_ELBRUS:   "EM_MCST_ELBRUS",
	EM_ECOG16:        "EM_ECOG16",
	EM_CR16:          "EM_CR16",
	EM_ETPU:          "EM_ETPU",
	EM_SLE9X:         "EM_SLE9X",
	EM_L10M:          "EM_L10M",
	EM_K10M:          "EM_K10M",
	EM_AARCH64:       "EM_AARCH64",
	EM_AVR32:         "EM_AVR32",
	EM_STM8:          "EM_STM8",
	EM_TILE64:        "EM_TILE64",
	EM_TILEPRO:       "EM_TILEPRO",
	EM_MICROBLAZE:    "EM_MICROBLAZE",
	EM_CUDA:          "EM_CUDA",
	EM_TILEGX:        "EM_TILEGX",
	EM_CLOUDSHIELD:   "EM_CLOUDSHIELD",
	EM_COREA_1ST:     "EM_COREA_1ST",
	EM_COREA_2ND:     "EM_COREA_2ND",
	EM_ARC_COMPACT2:  "EM_ARC_COMPACT2",
	EM_OPEN8:         "EM_OPEN8",
	EM_RL78:          "EM_RL78",
	EM_VIDEOCORE5:    "EM_VIDEOCORE5",
	EM_78KOR:         "EM_78KOR",
	EM_56800EX:       "EM_56800EX",
	EM_BA1:           "EM_BA1",
	EM_BA2:           "EM_BA2",
	EM_XCORE:         "EM_XCORE",
	EM_MCHP_PIC:      "EM_MCHP_PIC",
	EM_KM32:          "EM_KM32",
	EM_KMX32:         "EM_KMX32",
	EM_EMX16:         "EM_EMX16",
	EM_EMX8:          "EM_EMX8",
	EM_KVARC:         "EM_KVARC",
	EM_CDP:           "EM_CDP",
	EM_COGE:          "EM_COGE",
	EM_COOL:          "EM_COOL",
	EM_NORC:          "EM_NORC",
	EM_CSR_KALIMBA:   "EM_CSR_KALIMBA",
	EM_Z80:           "EM_Z80",
	EM_VISIUM:        "EM_VISIUM",
	EM_FT32:          "EM_FT32",
	EM_MOXIE:         "EM_MOXIE",
	EM_AMDGPU:        "EM_AMDGPU",
	EM_RISCV:         "EM_RISCV",
	EM_BPF:           "EM_BPF",
	EM_C_SKY:         "EM_C_SKY",
}

func (m Machine) String() string {
	if s, ok := machNames[m]; ok {
		return s
	}

	return "EM<" + strconv.Itoa(int(m)) + ">"
}

// --- Program header ---

// ProgType is the type of a program header table entry, p_type.
type ProgType uint32

const (
	PT_NULL             ProgType = 0
	PT_LOAD             ProgType = 1
	PT_DYNAMIC          ProgType = 2
	PT_INTERP           ProgType = 3
	PT_NOTE             ProgType = 4
	PT_SHLIB            ProgType = 5
	PT_PHDR             ProgType = 6
	PT_TLS              ProgType = 7
	PT_LOOS             ProgType = 0x60000000
	PT_GNU_EH_FRAME     ProgType = 0x6474e550
	PT_GNU_STACK        ProgType = 0x6474e551
	PT_GNU_RELRO        ProgType = 0x6474e552
	PT_GNU_PROPERTY     ProgType = 0x6474e553
	PT_GNU_SFRAME       ProgType = 0x6474e554
	PT_HIOS             ProgType = 0x6fffffff
	PT_LOPROC           ProgType = 0x70000000
	PT_ARM_EXIDX        ProgType = 0x70000001
	PT_RISCV_ATTRIBUTES ProgType = 0x70000003
	PT_HIPROC           ProgType = 0x7fffffff
)

func (p ProgType) String() string {
	switch p {
	case PT_NULL:
		return "PT_NULL"
	case PT_LOAD:
		return "PT_LOAD"
	case PT_DYNAMIC:
		return "PT_DYNAMIC"
	case PT_INTERP:
		return "PT_INTERP"
	case PT_NOTE:
		return "PT_NOTE"
	case PT_SHLIB:
		return "PT_SHLIB"
	case PT_PHDR:
		return "PT_PHDR"
	case PT_TLS:
		return "PT_TLS"
	case PT_GNU_EH_FRAME:
		return "PT_GNU_EH_FRAME"
	case PT_GNU_STACK:
		return "PT_GNU_STACK"
	case PT_GNU_RELRO:
		return "PT_GNU_RELRO"
	case PT_GNU_PROPERTY:
		return "PT_GNU_PROPERTY"
	case PT_GNU_SFRAME:
		return "PT_GNU_SFRAME"
	case PT_ARM_EXIDX:
		return "PT_ARM_EXIDX"
	case PT_RISCV_ATTRIBUTES:
		return "PT_RISCV_ATTRIBUTES"
	default:
		if p >= PT_LOOS && p <= PT_HIOS {
			return "PT_LOOS+" + strconv.FormatUint(uint64(p-PT_LOOS), 16)
		}

		if p >= PT_LOPROC && p <= PT_HIPROC {
			return "PT_LOPROC+" + strconv.FormatUint(uint64(p-PT_LOPROC), 16)
		}

		return "PT<" + strconv.FormatUint(uint64(p), 16) + ">"
	}
}

// ProgFlag is the segment flags, p_flags (a bit mask).
type ProgFlag uint32

const (
	PF_X        ProgFlag = 0x1 // execute
	PF_W        ProgFlag = 0x2 // write
	PF_R        ProgFlag = 0x4 // read
	PF_MASKOS   ProgFlag = 0x0ff00000
	PF_MASKPROC ProgFlag = 0xf0000000
)

func (f ProgFlag) String() string {
	s := ""
	if f&PF_R != 0 {
		s += "R"
	}

	if f&PF_W != 0 {
		s += "W"
	}

	if f&PF_X != 0 {
		s += "X"
	}

	if rest := f &^ (PF_R | PF_W | PF_X); rest != 0 {
		s += "|0x" + strconv.FormatUint(uint64(rest), 16)
	}

	if s == "" {
		s = "0"
	}

	return s
}

// --- Section header ---

// SectionType is the section type, sh_type.
type SectionType uint32

const (
	SHT_NULL          SectionType = 0
	SHT_PROGBITS      SectionType = 1
	SHT_SYMTAB        SectionType = 2
	SHT_STRTAB        SectionType = 3
	SHT_RELA          SectionType = 4
	SHT_HASH          SectionType = 5
	SHT_DYNAMIC       SectionType = 6
	SHT_NOTE          SectionType = 7
	SHT_NOBITS        SectionType = 8
	SHT_REL           SectionType = 9
	SHT_SHLIB         SectionType = 10
	SHT_DYNSYM        SectionType = 11
	SHT_INIT_ARRAY    SectionType = 14
	SHT_FINI_ARRAY    SectionType = 15
	SHT_PREINIT_ARRAY SectionType = 16
	SHT_GROUP         SectionType = 17
	SHT_SYMTAB_SHNDX  SectionType = 18
	SHT_RELR          SectionType = 19
	SHT_LOOS          SectionType = 0x60000000
	// GNU extensions (values from the LOOS..LOSUNW range).
	SHT_GNU_ATTRIBUTES SectionType = 0x6ffffff5
	SHT_GNU_HASH       SectionType = 0x6ffffff6
	SHT_GNU_LIBLIST    SectionType = 0x6ffffff7
	SHT_CHECKSUM       SectionType = 0x6ffffff8
	SHT_LOSUNW         SectionType = 0x6ffffffa
	SHT_SUNW_COMDAT    SectionType = 0x6ffffffb
	SHT_GNU_VERDEF     SectionType = 0x6ffffffd // verdef
	SHT_GNU_VERNEED    SectionType = 0x6ffffffe // verneed
	SHT_GNU_VERSYM     SectionType = 0x6fffffff // versym
	SHT_HIOS           SectionType = 0x6fffffff
	SHT_LOPROC         SectionType = 0x70000000
	// 0x70000003 is _attributes for both ARM and RISC-V.
	SHT_RISCV_ATTRIBUTES SectionType = 0x70000003
	SHT_HIPROC           SectionType = 0x7fffffff
	SHT_LOUSER           SectionType = 0x80000000
	SHT_HIUSER           SectionType = 0xffffffff
)

func (t SectionType) String() string {
	switch t {
	case SHT_NULL:
		return "SHT_NULL"
	case SHT_PROGBITS:
		return "SHT_PROGBITS"
	case SHT_SYMTAB:
		return "SHT_SYMTAB"
	case SHT_STRTAB:
		return "SHT_STRTAB"
	case SHT_RELA:
		return "SHT_RELA"
	case SHT_HASH:
		return "SHT_HASH"
	case SHT_DYNAMIC:
		return "SHT_DYNAMIC"
	case SHT_NOTE:
		return "SHT_NOTE"
	case SHT_NOBITS:
		return "SHT_NOBITS"
	case SHT_REL:
		return "SHT_REL"
	case SHT_SHLIB:
		return "SHT_SHLIB"
	case SHT_DYNSYM:
		return "SHT_DYNSYM"
	case SHT_INIT_ARRAY:
		return "SHT_INIT_ARRAY"
	case SHT_FINI_ARRAY:
		return "SHT_FINI_ARRAY"
	case SHT_PREINIT_ARRAY:
		return "SHT_PREINIT_ARRAY"
	case SHT_GROUP:
		return "SHT_GROUP"
	case SHT_SYMTAB_SHNDX:
		return "SHT_SYMTAB_SHNDX"
	case SHT_RELR:
		return "SHT_RELR"
	case SHT_GNU_ATTRIBUTES:
		return "SHT_GNU_ATTRIBUTES"
	case SHT_GNU_HASH:
		return "SHT_GNU_HASH"
	case SHT_GNU_LIBLIST:
		return "SHT_GNU_LIBLIST"
	case SHT_CHECKSUM:
		return "SHT_CHECKSUM"
	case SHT_SUNW_COMDAT:
		return "SHT_SUNW_COMDAT"
	case SHT_GNU_VERDEF:
		return "SHT_GNU_VERDEF"
	case SHT_GNU_VERNEED:
		return "SHT_GNU_VERNEED"
	case SHT_GNU_VERSYM:
		return "SHT_GNU_VERSYM"
	case SHT_RISCV_ATTRIBUTES:
		return "SHT_RISCV_ATTRIBUTES"
	default:
		if t >= SHT_LOOS && t <= SHT_HIOS {
			return "SHT_LOOS+" + strconv.FormatUint(uint64(t-SHT_LOOS), 16)
		}

		if t >= SHT_LOPROC && t <= SHT_HIPROC {
			return "SHT_LOPROC+" + strconv.FormatUint(uint64(t-SHT_LOPROC), 16)
		}

		if t >= SHT_LOUSER {
			return "SHT_LOUSER+" + strconv.FormatUint(uint64(t-SHT_LOUSER), 16)
		}

		return "SHT<" + strconv.FormatUint(uint64(t), 16) + ">"
	}
}

// SectionFlag is the section flags, sh_flags (a bit mask).
type SectionFlag uint64

const (
	SHF_WRITE            SectionFlag = 0x1
	SHF_ALLOC            SectionFlag = 0x2
	SHF_EXECINSTR        SectionFlag = 0x4
	SHF_MERGE            SectionFlag = 0x10
	SHF_STRINGS          SectionFlag = 0x20
	SHF_INFO_LINK        SectionFlag = 0x40
	SHF_LINK_ORDER       SectionFlag = 0x80
	SHF_OS_NONCONFORMING SectionFlag = 0x100
	SHF_GROUP            SectionFlag = 0x200
	SHF_TLS              SectionFlag = 0x400
	SHF_COMPRESSED       SectionFlag = 0x800
	SHF_GNU_RETAIN       SectionFlag = 0x200000
	SHF_EXCLUDE          SectionFlag = 0x80000000 // from the SHF_LOOS~GNU range
)

// String decomposes the mask into the names of known bits.
func (f SectionFlag) String() string {
	if f == 0 {
		return "0"
	}

	var parts []string
	add := func(bit SectionFlag, name string) {
		if f&bit != 0 {
			parts = append(parts, name)
		}
	}
	add(SHF_WRITE, "SHF_WRITE")
	add(SHF_ALLOC, "SHF_ALLOC")
	add(SHF_EXECINSTR, "SHF_EXECINSTR")
	add(SHF_MERGE, "SHF_MERGE")
	add(SHF_STRINGS, "SHF_STRINGS")
	add(SHF_INFO_LINK, "SHF_INFO_LINK")
	add(SHF_LINK_ORDER, "SHF_LINK_ORDER")
	add(SHF_OS_NONCONFORMING, "SHF_OS_NONCONFORMING")
	add(SHF_GROUP, "SHF_GROUP")
	add(SHF_TLS, "SHF_TLS")
	add(SHF_COMPRESSED, "SHF_COMPRESSED")
	add(SHF_GNU_RETAIN, "SHF_GNU_RETAIN")
	add(SHF_EXCLUDE, "SHF_EXCLUDE")
	s := ""
	var sSb783 strings.Builder
	for i, p := range parts {
		if i > 0 {
			sSb783.WriteString("|")
		}

		sSb783.WriteString(p)
	}

	s += sSb783.String()
	return s
}

// SHN_* are special section indices (sh_link/shndx and the like).
const (
	SHN_UNDEF     uint32 = 0
	SHN_LORESERVE uint32 = 0xff00
	SHN_LOPROC    uint32 = 0xff00
	SHN_HIPROC    uint32 = 0xff1f
	SHN_LOOS      uint32 = 0xff20
	SHN_HIOS      uint32 = 0xff3f
	SHN_ABS       uint32 = 0xfff1
	SHN_COMMON    uint32 = 0xfff2
	SHN_XINDEX    uint32 = 0xffff
)

// shndxName returns the name of a special section index, or "".
func shndxName(shndx uint32) string {
	switch shndx {
	case SHN_UNDEF:
		return "SHN_UNDEF"
	case SHN_ABS:
		return "SHN_ABS"
	case SHN_COMMON:
		return "SHN_COMMON"
	case SHN_XINDEX:
		return "SHN_XINDEX"
	}

	if shndx >= SHN_LOPROC && shndx <= SHN_HIPROC {
		return "SHN_LOPROC+" + strconv.FormatUint(uint64(shndx-SHN_LOPROC), 16)
	}

	if shndx >= SHN_LOOS && shndx <= SHN_HIOS {
		return "SHN_LOOS+" + strconv.FormatUint(uint64(shndx-SHN_LOOS), 16)
	}

	if shndx >= SHN_LORESERVE {
		return "SHN_LORESERVE+" + strconv.FormatUint(uint64(shndx-SHN_LORESERVE), 16)
	}

	return ""
}

// --- Dynamic section ---

// DynTag is the tag of a .dynamic section entry, d_tag.
type DynTag uint64

const (
	DT_NULL            DynTag = 0
	DT_NEEDED          DynTag = 1
	DT_PLTRELSZ        DynTag = 2
	DT_PLTGOT          DynTag = 3
	DT_HASH            DynTag = 4
	DT_STRTAB          DynTag = 5
	DT_SYMTAB          DynTag = 6
	DT_REAL            DynTag = 7
	DT_RELASZ          DynTag = 8
	DT_RELAENT         DynTag = 9
	DT_STRSZ           DynTag = 10
	DT_SYMENT          DynTag = 11
	DT_INIT            DynTag = 12
	DT_FINI            DynTag = 13
	DT_SONAME          DynTag = 14
	DT_RPATH           DynTag = 15
	DT_SYMBOLIC        DynTag = 16
	DT_REL             DynTag = 17
	DT_RELSZ           DynTag = 18
	DT_RELENT          DynTag = 19
	DT_PLTREL          DynTag = 20
	DT_DEBUG           DynTag = 21
	DT_TEXTREL         DynTag = 22
	DT_JMPREL          DynTag = 23
	DT_BIND_NOW        DynTag = 24
	DT_INIT_ARRAY      DynTag = 25
	DT_FINI_ARRAY      DynTag = 26
	DT_INIT_ARRAYSZ    DynTag = 27
	DT_FINI_ARRAYSZ    DynTag = 28
	DT_RUNPATH         DynTag = 29
	DT_FLAGS           DynTag = 30
	DT_PREINIT_ARRAY   DynTag = 32
	DT_PREINIT_ARRAYSZ DynTag = 33
	DT_SYMTAB_SHNDX    DynTag = 34
	DT_RELRSZ          DynTag = 35
	DT_RELR            DynTag = 36
	DT_RELRENT         DynTag = 37
	DT_LOOS            DynTag = 0x6000000d
	DT_HIOS            DynTag = 0x6ffff000
	DT_GNU_HASH        DynTag = 0x6ffffef5
	DT_VERSYM          DynTag = 0x6ffffff0
	DT_FLAGS_1         DynTag = 0x6ffffffb
	DT_VERDEF          DynTag = 0x6ffffffc
	DT_VERDEFNUM       DynTag = 0x6ffffffd
	DT_VERNEED         DynTag = 0x6ffffffe
	DT_VERNEEDNUM      DynTag = 0x6fffffff
	DT_LOPROC          DynTag = 0x70000000
	DT_HIPROC          DynTag = 0x7fffffff
)

func (t DynTag) String() string {
	switch t {
	case DT_NULL:
		return "DT_NULL"
	case DT_NEEDED:
		return "DT_NEEDED"
	case DT_PLTRELSZ:
		return "DT_PLTRELSZ"
	case DT_PLTGOT:
		return "DT_PLTGOT"
	case DT_HASH:
		return "DT_HASH"
	case DT_STRTAB:
		return "DT_STRTAB"
	case DT_SYMTAB:
		return "DT_SYMTAB"
	case DT_REAL:
		return "DT_REAL"
	case DT_RELASZ:
		return "DT_RELASZ"
	case DT_RELAENT:
		return "DT_RELAENT"
	case DT_STRSZ:
		return "DT_STRSZ"
	case DT_SYMENT:
		return "DT_SYMENT"
	case DT_INIT:
		return "DT_INIT"
	case DT_FINI:
		return "DT_FINI"
	case DT_SONAME:
		return "DT_SONAME"
	case DT_RPATH:
		return "DT_RPATH"
	case DT_SYMBOLIC:
		return "DT_SYMBOLIC"
	case DT_REL:
		return "DT_REL"
	case DT_RELSZ:
		return "DT_RELSZ"
	case DT_RELENT:
		return "DT_RELENT"
	case DT_PLTREL:
		return "DT_PLTREL"
	case DT_DEBUG:
		return "DT_DEBUG"
	case DT_TEXTREL:
		return "DT_TEXTREL"
	case DT_JMPREL:
		return "DT_JMPREL"
	case DT_BIND_NOW:
		return "DT_BIND_NOW"
	case DT_INIT_ARRAY:
		return "DT_INIT_ARRAY"
	case DT_FINI_ARRAY:
		return "DT_FINI_ARRAY"
	case DT_INIT_ARRAYSZ:
		return "DT_INIT_ARRAYSZ"
	case DT_FINI_ARRAYSZ:
		return "DT_FINI_ARRAYSZ"
	case DT_RUNPATH:
		return "DT_RUNPATH"
	case DT_FLAGS:
		return "DT_FLAGS"
	case DT_PREINIT_ARRAY:
		return "DT_PREINIT_ARRAY"
	case DT_PREINIT_ARRAYSZ:
		return "DT_PREINIT_ARRAYSZ"
	case DT_SYMTAB_SHNDX:
		return "DT_SYMTAB_SHNDX"
	case DT_RELRSZ:
		return "DT_RELRSZ"
	case DT_RELR:
		return "DT_RELR"
	case DT_RELRENT:
		return "DT_RELRENT"
	case DT_GNU_HASH:
		return "DT_GNU_HASH"
	case DT_VERSYM:
		return "DT_VERSYM"
	case DT_FLAGS_1:
		return "DT_FLAGS_1"
	case DT_VERDEF:
		return "DT_VERDEF"
	case DT_VERDEFNUM:
		return "DT_VERDEFNUM"
	case DT_VERNEED:
		return "DT_VERNEED"
	case DT_VERNEEDNUM:
		return "DT_VERNEEDNUM"
	default:
		if t >= DT_LOOS && t <= DT_HIOS {
			return "DT_LOOS+" + strconv.FormatUint(uint64(t-DT_LOOS), 16)
		}

		if t >= DT_LOPROC && t <= DT_HIPROC {
			return "DT_LOPROC+" + strconv.FormatUint(uint64(t-DT_LOPROC), 16)
		}

		return "DT<" + strconv.FormatUint(uint64(t), 16) + ">"
	}
}

// --- Symbols ---

// SymbolBind is the symbol binding class (the high nibble of st_info).
type SymbolBind uint8

const (
	STB_LOCAL      SymbolBind = 0
	STB_GLOBAL     SymbolBind = 1
	STB_WEAK       SymbolBind = 2
	STB_GNU_UNIQUE SymbolBind = 10
)

func (b SymbolBind) String() string {
	switch b {
	case STB_LOCAL:
		return "STB_LOCAL"
	case STB_GLOBAL:
		return "STB_GLOBAL"
	case STB_WEAK:
		return "STB_WEAK"
	case STB_GNU_UNIQUE:
		return "STB_GNU_UNIQUE"
	}

	return "STB<" + strconv.Itoa(int(b)) + ">"
}

// SymbolType is the symbol type (the low nibble of st_info).
type SymbolType uint8

const (
	STT_NOTYPE    SymbolType = 0
	STT_OBJECT    SymbolType = 1
	STT_FUNC      SymbolType = 2
	STT_SECTION   SymbolType = 3
	STT_FILE      SymbolType = 4
	STT_COMMON    SymbolType = 5
	STT_TLS       SymbolType = 6
	STT_GNU_IFUNC SymbolType = 10
)

func (t SymbolType) String() string {
	switch t {
	case STT_NOTYPE:
		return "STT_NOTYPE"
	case STT_OBJECT:
		return "STT_OBJECT"
	case STT_FUNC:
		return "STT_FUNC"
	case STT_SECTION:
		return "STT_SECTION"
	case STT_FILE:
		return "STT_FILE"
	case STT_COMMON:
		return "STT_COMMON"
	case STT_TLS:
		return "STT_TLS"
	case STT_GNU_IFUNC:
		return "STT_GNU_IFUNC"
	}

	return "STT<" + strconv.Itoa(int(t)) + ">"
}

// SymbolVisibility is the symbol visibility (the low 2 bits of st_other).
type SymbolVisibility uint8

const (
	STV_DEFAULT   SymbolVisibility = 0
	STV_INTERNAL  SymbolVisibility = 1
	STV_HIDDEN    SymbolVisibility = 2
	STV_PROTECTED SymbolVisibility = 3
)

func (v SymbolVisibility) String() string {
	switch v {
	case STV_DEFAULT:
		return "STV_DEFAULT"
	case STV_INTERNAL:
		return "STV_INTERNAL"
	case STV_HIDDEN:
		return "STV_HIDDEN"
	case STV_PROTECTED:
		return "STV_PROTECTED"
	}

	return "STV<" + strconv.Itoa(int(v)) + ">"
}

// --- Notes ---

// NoteType is the type of a SHT_NOTE/PT_NOTE entry. Note names differ by
// owner (name); the types below are the most common ones.
type NoteType uint32

const (
	NT_PRSTATUS            NoteType = 1
	NT_PRFPREG             NoteType = 2
	NT_PRPSINFO            NoteType = 3
	NT_TASKSTRUCT          NoteType = 4
	NT_AUXV                NoteType = 6
	NT_SIGINFO             NoteType = 0x53494749
	NT_FILE                NoteType = 0x46494c45
	NT_GNU_ABI_TAG         NoteType = 1
	NT_GNU_HWCAP           NoteType = 2
	NT_GNU_BUILD_ID        NoteType = 3
	NT_GNU_GOLD_VERSION    NoteType = 4
	NT_GNU_PROPERTY_TYPE_0 NoteType = 5
	NT_GO_BUILD_ID         NoteType = 4 // name "Go"
)

// --- AArch64 relocations (LP64; values verified against binutils/LLVM) ---

// RelocAarch64 is the relocation type for EM_AARCH64.
type RelocAarch64 uint32

const (
	R_AARCH64_NONE RelocAarch64 = 0

	R_AARCH64_P32_ABS32  RelocAarch64 = 1
	R_AARCH64_P32_ABS16  RelocAarch64 = 2
	R_AARCH64_P32_PREL32 RelocAarch64 = 3
	R_AARCH64_P32_PREL16 RelocAarch64 = 4

	R_AARCH64_ABS64               RelocAarch64 = 257
	R_AARCH64_ABS32               RelocAarch64 = 258
	R_AARCH64_ABS16               RelocAarch64 = 259
	R_AARCH64_PREL64              RelocAarch64 = 260
	R_AARCH64_PREL32              RelocAarch64 = 261
	R_AARCH64_PREL16              RelocAarch64 = 262
	R_AARCH64_MOVW_UABS_G0        RelocAarch64 = 263
	R_AARCH64_MOVW_UABS_G0_NC     RelocAarch64 = 264
	R_AARCH64_MOVW_UABS_G1        RelocAarch64 = 265
	R_AARCH64_MOVW_UABS_G1_NC     RelocAarch64 = 266
	R_AARCH64_MOVW_UABS_G2        RelocAarch64 = 267
	R_AARCH64_MOVW_UABS_G2_NC     RelocAarch64 = 268
	R_AARCH64_MOVW_UABS_G3        RelocAarch64 = 269
	R_AARCH64_MOVW_SABS_G0        RelocAarch64 = 270
	R_AARCH64_MOVW_SABS_G1        RelocAarch64 = 271
	R_AARCH64_MOVW_SABS_G2        RelocAarch64 = 272
	R_AARCH64_LD_PREL_LO19        RelocAarch64 = 273
	R_AARCH64_ADR_PREL_LO21       RelocAarch64 = 274
	R_AARCH64_ADR_PREL_PG_HI21    RelocAarch64 = 275
	R_AARCH64_ADR_PREL_PG_HI21_NC RelocAarch64 = 276
	R_AARCH64_ADD_ABS_LO12_NC     RelocAarch64 = 277
	R_AARCH64_LDST8_ABS_LO12_NC   RelocAarch64 = 278
	R_AARCH64_TSTBR14             RelocAarch64 = 279
	R_AARCH64_CONDBR19            RelocAarch64 = 280
	R_AARCH64_JUMP26              RelocAarch64 = 282
	R_AARCH64_CALL26              RelocAarch64 = 283
	R_AARCH64_LDST16_ABS_LO12_NC  RelocAarch64 = 284
	R_AARCH64_LDST32_ABS_LO12_NC  RelocAarch64 = 285
	R_AARCH64_LDST64_ABS_LO12_NC  RelocAarch64 = 286
	R_AARCH64_MOVW_PREL_G0        RelocAarch64 = 287
	R_AARCH64_MOVW_PREL_G0_NC     RelocAarch64 = 288
	R_AARCH64_MOVW_PREL_G1        RelocAarch64 = 289
	R_AARCH64_MOVW_PREL_G1_NC     RelocAarch64 = 290
	R_AARCH64_MOVW_PREL_G2        RelocAarch64 = 291
	R_AARCH64_MOVW_PREL_G2_NC     RelocAarch64 = 292
	R_AARCH64_MOVW_PREL_G3        RelocAarch64 = 293
	R_AARCH64_LDST128_ABS_LO12_NC RelocAarch64 = 299
	R_AARCH64_MOVW_GOTOFF_G0_NC   RelocAarch64 = 301
	R_AARCH64_MOVW_GOTOFF_G1      RelocAarch64 = 302
	R_AARCH64_MOVW_GOTOFF_G1_NC   RelocAarch64 = 303
	R_AARCH64_MOVW_GOTOFF_G2      RelocAarch64 = 304
	R_AARCH64_MOVW_GOTOFF_G2_NC   RelocAarch64 = 305
	R_AARCH64_GOTREL64            RelocAarch64 = 307
	R_AARCH64_GOTREL32            RelocAarch64 = 308
	R_AARCH64_GOT_LD_PREL19       RelocAarch64 = 309
	R_AARCH64_LD64_GOTOFF_LO15    RelocAarch64 = 310
	R_AARCH64_ADR_GOT_PAGE        RelocAarch64 = 311
	R_AARCH64_LD64_GOT_LO12_NC    RelocAarch64 = 312
	R_AARCH64_LD64_GOTPAGE_LO15   RelocAarch64 = 313
	R_AARCH64_PLT32               RelocAarch64 = 314

	R_AARCH64_TLSGD_ADR_PREL21             RelocAarch64 = 512
	R_AARCH64_TLSGD_ADR_PAGE21             RelocAarch64 = 513
	R_AARCH64_TLSGD_ADD_LO12_NC            RelocAarch64 = 514
	R_AARCH64_TLSGD_MOVW_G1                RelocAarch64 = 515
	R_AARCH64_TLSGD_MOVW_G0_NC             RelocAarch64 = 516
	R_AARCH64_TLSLD_ADR_PREL21             RelocAarch64 = 517
	R_AARCH64_TLSLD_ADR_PAGE21             RelocAarch64 = 518
	R_AARCH64_TLSIE_MOVW_GOTTPREL_G1       RelocAarch64 = 539
	R_AARCH64_TLSIE_MOVW_GOTTPREL_G0_NC    RelocAarch64 = 540
	R_AARCH64_TLSIE_ADR_GOTTPREL_PAGE21    RelocAarch64 = 541
	R_AARCH64_TLSIE_LD64_GOTTPREL_LO12_NC  RelocAarch64 = 542
	R_AARCH64_TLSIE_LD_GOTTPREL_PREL19     RelocAarch64 = 543
	R_AARCH64_TLSLE_MOVW_TPREL_G2          RelocAarch64 = 544
	R_AARCH64_TLSLE_MOVW_TPREL_G1          RelocAarch64 = 545
	R_AARCH64_TLSLE_MOVW_TPREL_G1_NC       RelocAarch64 = 546
	R_AARCH64_TLSLE_MOVW_TPREL_G0          RelocAarch64 = 547
	R_AARCH64_TLSLE_MOVW_TPREL_G0_NC       RelocAarch64 = 548
	R_AARCH64_TLSLE_ADD_TPREL_HI12         RelocAarch64 = 549
	R_AARCH64_TLSLE_ADD_TPREL_LO12         RelocAarch64 = 550
	R_AARCH64_TLSLE_ADD_TPREL_LO12_NC      RelocAarch64 = 551
	R_AARCH64_TLSDESC_LD_PREL19            RelocAarch64 = 560
	R_AARCH64_TLSDESC_ADR_PREL21           RelocAarch64 = 561
	R_AARCH64_TLSDESC_ADR_PAGE21           RelocAarch64 = 562
	R_AARCH64_TLSDESC_LD64_LO12_NC         RelocAarch64 = 563
	R_AARCH64_TLSDESC_ADD_LO12_NC          RelocAarch64 = 564
	R_AARCH64_TLSDESC_OFF_G1               RelocAarch64 = 565
	R_AARCH64_TLSDESC_OFF_G0_NC            RelocAarch64 = 566
	R_AARCH64_TLSDESC_LDR                  RelocAarch64 = 567
	R_AARCH64_TLSDESC_ADD                  RelocAarch64 = 568
	R_AARCH64_TLSDESC_CALL                 RelocAarch64 = 569
	R_AARCH64_TLSLE_LDST128_TPREL_LO12     RelocAarch64 = 570
	R_AARCH64_TLSLE_LDST128_TPREL_LO12_NC  RelocAarch64 = 571
	R_AARCH64_TLSLD_LDST128_DTPREL_LO12    RelocAarch64 = 572
	R_AARCH64_TLSLD_LDST128_DTPREL_LO12_NC RelocAarch64 = 573

	R_AARCH64_COPY         RelocAarch64 = 1024
	R_AARCH64_GLOB_DAT     RelocAarch64 = 1025
	R_AARCH64_JUMP_SLOT    RelocAarch64 = 1026
	R_AARCH64_RELATIVE     RelocAarch64 = 1027
	R_AARCH64_TLS_DTPMOD64 RelocAarch64 = 1028
	R_AARCH64_TLS_DTPREL64 RelocAarch64 = 1029
	R_AARCH64_TLS_TPREL64  RelocAarch64 = 1030
	R_AARCH64_TLSDESC      RelocAarch64 = 1031
	R_AARCH64_IRELATIVE    RelocAarch64 = 1032
)

var relocAarch64Names = map[RelocAarch64]string{
	R_AARCH64_NONE:                         "R_AARCH64_NONE",
	R_AARCH64_P32_ABS32:                    "R_AARCH64_P32_ABS32",
	R_AARCH64_P32_ABS16:                    "R_AARCH64_P32_ABS16",
	R_AARCH64_P32_PREL32:                   "R_AARCH64_P32_PREL32",
	R_AARCH64_P32_PREL16:                   "R_AARCH64_P32_PREL16",
	R_AARCH64_ABS64:                        "R_AARCH64_ABS64",
	R_AARCH64_ABS32:                        "R_AARCH64_ABS32",
	R_AARCH64_ABS16:                        "R_AARCH64_ABS16",
	R_AARCH64_PREL64:                       "R_AARCH64_PREL64",
	R_AARCH64_PREL32:                       "R_AARCH64_PREL32",
	R_AARCH64_PREL16:                       "R_AARCH64_PREL16",
	R_AARCH64_MOVW_UABS_G0:                 "R_AARCH64_MOVW_UABS_G0",
	R_AARCH64_MOVW_UABS_G0_NC:              "R_AARCH64_MOVW_UABS_G0_NC",
	R_AARCH64_MOVW_UABS_G1:                 "R_AARCH64_MOVW_UABS_G1",
	R_AARCH64_MOVW_UABS_G1_NC:              "R_AARCH64_MOVW_UABS_G1_NC",
	R_AARCH64_MOVW_UABS_G2:                 "R_AARCH64_MOVW_UABS_G2",
	R_AARCH64_MOVW_UABS_G2_NC:              "R_AARCH64_MOVW_UABS_G2_NC",
	R_AARCH64_MOVW_UABS_G3:                 "R_AARCH64_MOVW_UABS_G3",
	R_AARCH64_MOVW_SABS_G0:                 "R_AARCH64_MOVW_SABS_G0",
	R_AARCH64_MOVW_SABS_G1:                 "R_AARCH64_MOVW_SABS_G1",
	R_AARCH64_MOVW_SABS_G2:                 "R_AARCH64_MOVW_SABS_G2",
	R_AARCH64_LD_PREL_LO19:                 "R_AARCH64_LD_PREL_LO19",
	R_AARCH64_ADR_PREL_LO21:                "R_AARCH64_ADR_PREL_LO21",
	R_AARCH64_ADR_PREL_PG_HI21:             "R_AARCH64_ADR_PREL_PG_HI21",
	R_AARCH64_ADR_PREL_PG_HI21_NC:          "R_AARCH64_ADR_PREL_PG_HI21_NC",
	R_AARCH64_ADD_ABS_LO12_NC:              "R_AARCH64_ADD_ABS_LO12_NC",
	R_AARCH64_LDST8_ABS_LO12_NC:            "R_AARCH64_LDST8_ABS_LO12_NC",
	R_AARCH64_TSTBR14:                      "R_AARCH64_TSTBR14",
	R_AARCH64_CONDBR19:                     "R_AARCH64_CONDBR19",
	R_AARCH64_JUMP26:                       "R_AARCH64_JUMP26",
	R_AARCH64_CALL26:                       "R_AARCH64_CALL26",
	R_AARCH64_LDST16_ABS_LO12_NC:           "R_AARCH64_LDST16_ABS_LO12_NC",
	R_AARCH64_LDST32_ABS_LO12_NC:           "R_AARCH64_LDST32_ABS_LO12_NC",
	R_AARCH64_LDST64_ABS_LO12_NC:           "R_AARCH64_LDST64_ABS_LO12_NC",
	R_AARCH64_MOVW_PREL_G0:                 "R_AARCH64_MOVW_PREL_G0",
	R_AARCH64_MOVW_PREL_G0_NC:              "R_AARCH64_MOVW_PREL_G0_NC",
	R_AARCH64_MOVW_PREL_G1:                 "R_AARCH64_MOVW_PREL_G1",
	R_AARCH64_MOVW_PREL_G1_NC:              "R_AARCH64_MOVW_PREL_G1_NC",
	R_AARCH64_MOVW_PREL_G2:                 "R_AARCH64_MOVW_PREL_G2",
	R_AARCH64_MOVW_PREL_G2_NC:              "R_AARCH64_MOVW_PREL_G2_NC",
	R_AARCH64_MOVW_PREL_G3:                 "R_AARCH64_MOVW_PREL_G3",
	R_AARCH64_LDST128_ABS_LO12_NC:          "R_AARCH64_LDST128_ABS_LO12_NC",
	R_AARCH64_MOVW_GOTOFF_G0_NC:            "R_AARCH64_MOVW_GOTOFF_G0_NC",
	R_AARCH64_MOVW_GOTOFF_G1:               "R_AARCH64_MOVW_GOTOFF_G1",
	R_AARCH64_MOVW_GOTOFF_G1_NC:            "R_AARCH64_MOVW_GOTOFF_G1_NC",
	R_AARCH64_MOVW_GOTOFF_G2:               "R_AARCH64_MOVW_GOTOFF_G2",
	R_AARCH64_MOVW_GOTOFF_G2_NC:            "R_AARCH64_MOVW_GOTOFF_G2_NC",
	R_AARCH64_GOTREL64:                     "R_AARCH64_GOTREL64",
	R_AARCH64_GOTREL32:                     "R_AARCH64_GOTREL32",
	R_AARCH64_GOT_LD_PREL19:                "R_AARCH64_GOT_LD_PREL19",
	R_AARCH64_LD64_GOTOFF_LO15:             "R_AARCH64_LD64_GOTOFF_LO15",
	R_AARCH64_ADR_GOT_PAGE:                 "R_AARCH64_ADR_GOT_PAGE",
	R_AARCH64_LD64_GOT_LO12_NC:             "R_AARCH64_LD64_GOT_LO12_NC",
	R_AARCH64_LD64_GOTPAGE_LO15:            "R_AARCH64_LD64_GOTPAGE_LO15",
	R_AARCH64_PLT32:                        "R_AARCH64_PLT32",
	R_AARCH64_TLSGD_ADR_PREL21:             "R_AARCH64_TLSGD_ADR_PREL21",
	R_AARCH64_TLSGD_ADR_PAGE21:             "R_AARCH64_TLSGD_ADR_PAGE21",
	R_AARCH64_TLSGD_ADD_LO12_NC:            "R_AARCH64_TLSGD_ADD_LO12_NC",
	R_AARCH64_TLSGD_MOVW_G1:                "R_AARCH64_TLSGD_MOVW_G1",
	R_AARCH64_TLSGD_MOVW_G0_NC:             "R_AARCH64_TLSGD_MOVW_G0_NC",
	R_AARCH64_TLSLD_ADR_PREL21:             "R_AARCH64_TLSLD_ADR_PREL21",
	R_AARCH64_TLSLD_ADR_PAGE21:             "R_AARCH64_TLSLD_ADR_PAGE21",
	R_AARCH64_TLSIE_MOVW_GOTTPREL_G1:       "R_AARCH64_TLSIE_MOVW_GOTTPREL_G1",
	R_AARCH64_TLSIE_MOVW_GOTTPREL_G0_NC:    "R_AARCH64_TLSIE_MOVW_GOTTPREL_G0_NC",
	R_AARCH64_TLSIE_ADR_GOTTPREL_PAGE21:    "R_AARCH64_TLSIE_ADR_GOTTPREL_PAGE21",
	R_AARCH64_TLSIE_LD64_GOTTPREL_LO12_NC:  "R_AARCH64_TLSIE_LD64_GOTTPREL_LO12_NC",
	R_AARCH64_TLSIE_LD_GOTTPREL_PREL19:     "R_AARCH64_TLSIE_LD_GOTTPREL_PREL19",
	R_AARCH64_TLSLE_MOVW_TPREL_G2:          "R_AARCH64_TLSLE_MOVW_TPREL_G2",
	R_AARCH64_TLSLE_MOVW_TPREL_G1:          "R_AARCH64_TLSLE_MOVW_TPREL_G1",
	R_AARCH64_TLSLE_MOVW_TPREL_G1_NC:       "R_AARCH64_TLSLE_MOVW_TPREL_G1_NC",
	R_AARCH64_TLSLE_MOVW_TPREL_G0:          "R_AARCH64_TLSLE_MOVW_TPREL_G0",
	R_AARCH64_TLSLE_MOVW_TPREL_G0_NC:       "R_AARCH64_TLSLE_MOVW_TPREL_G0_NC",
	R_AARCH64_TLSLE_ADD_TPREL_HI12:         "R_AARCH64_TLSLE_ADD_TPREL_HI12",
	R_AARCH64_TLSLE_ADD_TPREL_LO12:         "R_AARCH64_TLSLE_ADD_TPREL_LO12",
	R_AARCH64_TLSLE_ADD_TPREL_LO12_NC:      "R_AARCH64_TLSLE_ADD_TPREL_LO12_NC",
	R_AARCH64_TLSDESC_LD_PREL19:            "R_AARCH64_TLSDESC_LD_PREL19",
	R_AARCH64_TLSDESC_ADR_PREL21:           "R_AARCH64_TLSDESC_ADR_PREL21",
	R_AARCH64_TLSDESC_ADR_PAGE21:           "R_AARCH64_TLSDESC_ADR_PAGE21",
	R_AARCH64_TLSDESC_LD64_LO12_NC:         "R_AARCH64_TLSDESC_LD64_LO12_NC",
	R_AARCH64_TLSDESC_ADD_LO12_NC:          "R_AARCH64_TLSDESC_ADD_LO12_NC",
	R_AARCH64_TLSDESC_OFF_G1:               "R_AARCH64_TLSDESC_OFF_G1",
	R_AARCH64_TLSDESC_OFF_G0_NC:            "R_AARCH64_TLSDESC_OFF_G0_NC",
	R_AARCH64_TLSDESC_LDR:                  "R_AARCH64_TLSDESC_LDR",
	R_AARCH64_TLSDESC_ADD:                  "R_AARCH64_TLSDESC_ADD",
	R_AARCH64_TLSDESC_CALL:                 "R_AARCH64_TLSDESC_CALL",
	R_AARCH64_TLSLE_LDST128_TPREL_LO12:     "R_AARCH64_TLSLE_LDST128_TPREL_LO12",
	R_AARCH64_TLSLE_LDST128_TPREL_LO12_NC:  "R_AARCH64_TLSLE_LDST128_TPREL_LO12_NC",
	R_AARCH64_TLSLD_LDST128_DTPREL_LO12:    "R_AARCH64_TLSLD_LDST128_DTPREL_LO12",
	R_AARCH64_TLSLD_LDST128_DTPREL_LO12_NC: "R_AARCH64_TLSLD_LDST128_DTPREL_LO12_NC",
	R_AARCH64_COPY:                         "R_AARCH64_COPY",
	R_AARCH64_GLOB_DAT:                     "R_AARCH64_GLOB_DAT",
	R_AARCH64_JUMP_SLOT:                    "R_AARCH64_JUMP_SLOT",
	R_AARCH64_RELATIVE:                     "R_AARCH64_RELATIVE",
	R_AARCH64_TLS_DTPMOD64:                 "R_AARCH64_TLS_DTPMOD64",
	R_AARCH64_TLS_DTPREL64:                 "R_AARCH64_TLS_DTPREL64",
	R_AARCH64_TLS_TPREL64:                  "R_AARCH64_TLS_TPREL64",
	R_AARCH64_TLSDESC:                      "R_AARCH64_TLSDESC",
	R_AARCH64_IRELATIVE:                    "R_AARCH64_IRELATIVE",
}

func (r RelocAarch64) String() string {
	if s, ok := relocAarch64Names[r]; ok {
		return s
	}

	return "R_AARCH64_<" + strconv.FormatUint(uint64(r), 10) + ">"
}

// --- RISC-V relocations (riscv-elf-psabi; values verified against binutils) ---

// RelocRiscv is the relocation type for EM_RISCV.
type RelocRiscv uint32

const (
	R_RISCV_NONE              RelocRiscv = 0
	R_RISCV_32                RelocRiscv = 1
	R_RISCV_64                RelocRiscv = 2
	R_RISCV_RELATIVE          RelocRiscv = 3
	R_RISCV_COPY              RelocRiscv = 4
	R_RISCV_JUMP_SLOT         RelocRiscv = 5
	R_RISCV_TLS_DTPMOD32      RelocRiscv = 6
	R_RISCV_TLS_DTPMOD64      RelocRiscv = 7
	R_RISCV_TLS_DTPREL32      RelocRiscv = 8
	R_RISCV_TLS_DTPREL64      RelocRiscv = 9
	R_RISCV_TLS_TPREL32       RelocRiscv = 10
	R_RISCV_TLS_TPREL64       RelocRiscv = 11
	R_RISCV_BRANCH            RelocRiscv = 16
	R_RISCV_JAL               RelocRiscv = 17
	R_RISCV_CALL              RelocRiscv = 18
	R_RISCV_CALL_PLT          RelocRiscv = 19
	R_RISCV_GOT_HI20          RelocRiscv = 20
	R_RISCV_TLS_GOT_HI20      RelocRiscv = 21
	R_RISCV_TLS_GD_HI20       RelocRiscv = 22
	R_RISCV_PCREL_HI20        RelocRiscv = 23
	R_RISCV_PCREL_LO12_I      RelocRiscv = 24
	R_RISCV_PCREL_LO12_S      RelocRiscv = 25
	R_RISCV_HI20              RelocRiscv = 26
	R_RISCV_LO12_I            RelocRiscv = 27
	R_RISCV_LO12_S            RelocRiscv = 28
	R_RISCV_TPREL_HI20        RelocRiscv = 29
	R_RISCV_TPREL_LO12_I      RelocRiscv = 30
	R_RISCV_TPREL_LO12_S      RelocRiscv = 31
	R_RISCV_TPREL_ADD         RelocRiscv = 32
	R_RISCV_ADD8              RelocRiscv = 33
	R_RISCV_ADD16             RelocRiscv = 34
	R_RISCV_ADD32             RelocRiscv = 35
	R_RISCV_ADD64             RelocRiscv = 36
	R_RISCV_SUB6              RelocRiscv = 52
	R_RISCV_SUB8              RelocRiscv = 37
	R_RISCV_SUB16             RelocRiscv = 38
	R_RISCV_SUB32             RelocRiscv = 39
	R_RISCV_SUB64             RelocRiscv = 40
	R_RISCV_GNU_VTINHERIT     RelocRiscv = 41
	R_RISCV_GNU_VTENTRY       RelocRiscv = 42
	R_RISCV_ALIGN             RelocRiscv = 43
	R_RISCV_RVC_BRANCH        RelocRiscv = 44
	R_RISCV_RVC_JUMP          RelocRiscv = 45
	R_RISCV_RVC_LUI           RelocRiscv = 46
	R_RISCV_GPREL_I           RelocRiscv = 47
	R_RISCV_GPREL_S           RelocRiscv = 48
	R_RISCV_TPREL_I           RelocRiscv = 49
	R_RISCV_TPREL_S           RelocRiscv = 50
	R_RISCV_RELAX             RelocRiscv = 51
	R_RISCV_SET6              RelocRiscv = 53
	R_RISCV_SET8              RelocRiscv = 54
	R_RISCV_SET16             RelocRiscv = 55
	R_RISCV_SET32             RelocRiscv = 56
	R_RISCV_32_PCREL          RelocRiscv = 57
	R_RISCV_IRELATIVE         RelocRiscv = 58
	R_RISCV_PLT32             RelocRiscv = 59
	R_RISCV_SET_ULEB128       RelocRiscv = 60
	R_RISCV_SUB_ULEB128       RelocRiscv = 61
	R_RISCV_TLSDESC_HI20      RelocRiscv = 62
	R_RISCV_TLSDESC_LOAD_LO12 RelocRiscv = 63
	R_RISCV_TLSDESC_ADD_LO12  RelocRiscv = 64
	R_RISCV_TLSDESC_CALL      RelocRiscv = 65
)

var relocRiscvNames = map[RelocRiscv]string{
	R_RISCV_NONE:              "R_RISCV_NONE",
	R_RISCV_32:                "R_RISCV_32",
	R_RISCV_64:                "R_RISCV_64",
	R_RISCV_RELATIVE:          "R_RISCV_RELATIVE",
	R_RISCV_COPY:              "R_RISCV_COPY",
	R_RISCV_JUMP_SLOT:         "R_RISCV_JUMP_SLOT",
	R_RISCV_TLS_DTPMOD32:      "R_RISCV_TLS_DTPMOD32",
	R_RISCV_TLS_DTPMOD64:      "R_RISCV_TLS_DTPMOD64",
	R_RISCV_TLS_DTPREL32:      "R_RISCV_TLS_DTPREL32",
	R_RISCV_TLS_DTPREL64:      "R_RISCV_TLS_DTPREL64",
	R_RISCV_TLS_TPREL32:       "R_RISCV_TLS_TPREL32",
	R_RISCV_TLS_TPREL64:       "R_RISCV_TLS_TPREL64",
	R_RISCV_BRANCH:            "R_RISCV_BRANCH",
	R_RISCV_JAL:               "R_RISCV_JAL",
	R_RISCV_CALL:              "R_RISCV_CALL",
	R_RISCV_CALL_PLT:          "R_RISCV_CALL_PLT",
	R_RISCV_GOT_HI20:          "R_RISCV_GOT_HI20",
	R_RISCV_TLS_GOT_HI20:      "R_RISCV_TLS_GOT_HI20",
	R_RISCV_TLS_GD_HI20:       "R_RISCV_TLS_GD_HI20",
	R_RISCV_PCREL_HI20:        "R_RISCV_PCREL_HI20",
	R_RISCV_PCREL_LO12_I:      "R_RISCV_PCREL_LO12_I",
	R_RISCV_PCREL_LO12_S:      "R_RISCV_PCREL_LO12_S",
	R_RISCV_HI20:              "R_RISCV_HI20",
	R_RISCV_LO12_I:            "R_RISCV_LO12_I",
	R_RISCV_LO12_S:            "R_RISCV_LO12_S",
	R_RISCV_TPREL_HI20:        "R_RISCV_TPREL_HI20",
	R_RISCV_TPREL_LO12_I:      "R_RISCV_TPREL_LO12_I",
	R_RISCV_TPREL_LO12_S:      "R_RISCV_TPREL_LO12_S",
	R_RISCV_TPREL_ADD:         "R_RISCV_TPREL_ADD",
	R_RISCV_ADD8:              "R_RISCV_ADD8",
	R_RISCV_ADD16:             "R_RISCV_ADD16",
	R_RISCV_ADD32:             "R_RISCV_ADD32",
	R_RISCV_ADD64:             "R_RISCV_ADD64",
	R_RISCV_SUB6:              "R_RISCV_SUB6",
	R_RISCV_SUB8:              "R_RISCV_SUB8",
	R_RISCV_SUB16:             "R_RISCV_SUB16",
	R_RISCV_SUB32:             "R_RISCV_SUB32",
	R_RISCV_SUB64:             "R_RISCV_SUB64",
	R_RISCV_GNU_VTINHERIT:     "R_RISCV_GNU_VTINHERIT",
	R_RISCV_GNU_VTENTRY:       "R_RISCV_GNU_VTENTRY",
	R_RISCV_ALIGN:             "R_RISCV_ALIGN",
	R_RISCV_RVC_BRANCH:        "R_RISCV_RVC_BRANCH",
	R_RISCV_RVC_JUMP:          "R_RISCV_RVC_JUMP",
	R_RISCV_RVC_LUI:           "R_RISCV_RVC_LUI",
	R_RISCV_GPREL_I:           "R_RISCV_GPREL_I",
	R_RISCV_GPREL_S:           "R_RISCV_GPREL_S",
	R_RISCV_TPREL_I:           "R_RISCV_TPREL_I",
	R_RISCV_TPREL_S:           "R_RISCV_TPREL_S",
	R_RISCV_RELAX:             "R_RISCV_RELAX",
	R_RISCV_SET6:              "R_RISCV_SET6",
	R_RISCV_SET8:              "R_RISCV_SET8",
	R_RISCV_SET16:             "R_RISCV_SET16",
	R_RISCV_SET32:             "R_RISCV_SET32",
	R_RISCV_32_PCREL:          "R_RISCV_32_PCREL",
	R_RISCV_IRELATIVE:         "R_RISCV_IRELATIVE",
	R_RISCV_PLT32:             "R_RISCV_PLT32",
	R_RISCV_SET_ULEB128:       "R_RISCV_SET_ULEB128",
	R_RISCV_SUB_ULEB128:       "R_RISCV_SUB_ULEB128",
	R_RISCV_TLSDESC_HI20:      "R_RISCV_TLSDESC_HI20",
	R_RISCV_TLSDESC_LOAD_LO12: "R_RISCV_TLSDESC_LOAD_LO12",
	R_RISCV_TLSDESC_ADD_LO12:  "R_RISCV_TLSDESC_ADD_LO12",
	R_RISCV_TLSDESC_CALL:      "R_RISCV_TLSDESC_CALL",
}

func (r RelocRiscv) String() string {
	if s, ok := relocRiscvNames[r]; ok {
		return s
	}

	return "R_RISCV_<" + strconv.FormatUint(uint64(r), 10) + ">"
}
