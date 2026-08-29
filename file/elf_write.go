package file

// WriteELF assembles an executable ELF64 from neutral Sections
// (the emitter is elf.Write).

import "github.com/okneniz/assembly/file/elf"

// e_machine of supported architectures (see <linux/elf.h> EM_*); the subpackage
// constants are aliased so that callers of file.WriteELF do not have to import it.
const (
	EM_AARCH64   uint16 = uint16(elf.EM_AARCH64)
	EM_RISCV     uint16 = uint16(elf.EM_RISCV)
	EM_LOONGARCH uint16 = uint16(elf.EM_LOONGARCH)
)

// WriteELF assembles an executable ELF64 LE from sections. machine is e_machine
// (EM_AARCH64/EM_RISCV/EM_LOONGARCH), flags is e_flags (processor-specific:
// EM_LOONGARCH needs the ABI bits, 0x43 = base+double-float); entry is the
// absolute address of the entry point
// (usually base + the offset of the start symbol). Sections are placed into the
// file back to back, data only (names and addresses are not preserved in this
// format). s.Size > len(s.Data) is a zero NOBITS tail in memory
// (p_memsz > p_filesz), allowed only for the last section.
func WriteELF(
	machine uint16,
	flags uint32,
	base, entry uint64,
	sections []Section,
) ([]byte, error) {
	blobs := make([]elf.Blob, len(sections))
	for i, s := range sections {
		blobs[i] = elf.NewNobitsBlob(s.Addr, s.Data, int(s.Size))
	}

	return elf.Write(elf.Machine(machine), flags, base, entry, blobs)
}
