package macho

import "encoding/binary"

// Symbol is nlist_64/nlist: an entry of the LC_SYMTAB symbol table. The
// Type field is the n_type byte (N_EXT/N_TYPE/N_PEXT/N_STAB), Sect is the
// section index (for N_SECT), Desc is the n_desc bits (weak, library
// ordinal, alt entry...), Value is the address (or the offset for N_UNDF).
//
// nlist_64 (16 bytes): n_strx(4) n_type(1) n_sect(1) n_desc(2) n_value(8)
// nlist (12 bytes):    n_strx(4) n_type(1) n_sect(1) n_desc(2) n_value(4).
type Symbol struct {
	Name  string
	Type  uint8
	Sect  uint8
	Desc  uint16
	Value uint64
}

// IsExternal: the N_EXT bit is set.
func (s Symbol) IsExternal() bool {
	return s.Type&N_EXT != 0
}

// IsPrivateExternal: the N_PEXT bit is set.
func (s Symbol) IsPrivateExternal() bool {
	return s.Type&N_PEXT != 0
}

// IsDebug: this is a stab symbol (N_STAB set - N_SO/N_FUN/...).
func (s Symbol) IsDebug() bool {
	return s.Type&N_STAB != 0
}

// Stab returns the stab symbol type (meaningful when IsDebug).
func (s Symbol) Stab() Stab {
	return Stab(s.Type)
}

// Ntype returns the raw N_TYPE value (masking the bits).
func (s Symbol) Ntype() uint8 {
	return s.Type & N_TYPE
}

// Ordinal returns the library ordinal from n_desc (for undefined symbols).
func (s Symbol) Ordinal() int {
	return LibraryOrdinal(s.Desc)
}

// IsWeakRef / IsWeakDef: weak reference/definition.
func (s Symbol) IsWeakRef() bool {
	return s.Desc&N_WEAK_REF != 0
}
func (s Symbol) IsWeakDef() bool {
	return s.Desc&N_WEAK_DEF != 0
}

// IsAltEntry: an alternate entry point of a function.
func (s Symbol) IsAltEntry() bool {
	return s.Desc&N_ALT_ENTRY != 0
}

// IsThumbDef: a Thumb function definition (32-bit ARM).
func (s Symbol) IsThumbDef() bool {
	return s.Desc&N_ARM_THUMB_DEF != 0
}

// SectionName returns the symbol's section name (for N_SECT) or "".
func (s Symbol) SectionName(f *File) string {
	if s.Ntype() != N_SECT || s.Sect == 0 || int(s.Sect) > len(f.sections) {
		return ""
	}

	return f.sections[s.Sect-1].SectName
}

const (
	nlist64Size = 16
	nlist32Size = 12
)

// Symbols returns the LC_SYMTAB symbol table (including stab symbols).
// The result is cached.
func (f *File) Symbols() ([]Symbol, error) {
	if f.symtabDone {
		return f.symtab, f.symtabErr
	}

	f.symtabDone = true

	var cmd *Symtab
	for _, lc := range f.commands {
		if st, ok := lc.(*Symtab); ok {
			cmd = st
			break
		}
	}

	if cmd == nil {
		f.symtabErr = errf("no LC_SYMTAB")
		return nil, f.symtabErr
	}

	strtab, err := readBytes(f.buf, f.base+int(cmd.Stroff), int(cmd.Strsize))
	if err != nil {
		f.symtabErr = err
		return nil, err
	}

	size := nlist32Size
	if f.is64 {
		size = nlist64Size
	}

	data, err := readBytes(f.buf, f.base+int(cmd.Symoff), int(cmd.Nsyms)*size)
	if err != nil {
		f.symtabErr = err
		return nil, err
	}

	out := make([]Symbol, 0, cmd.Nsyms)
	for i := range cmd.Nsyms {
		off := int(i) * size
		var s Symbol
		strx := u32(data[off:], f.order)
		s.Type = data[off+4]
		s.Sect = data[off+5]
		s.Desc = readU16(data[off+6:], f.order)
		if f.is64 {
			s.Value = readU64(data[off+8:], f.order)
		} else {
			s.Value = uint64(u32(data[off+8:], f.order))
		}

		if strx < uint32(len(strtab)) {
			s.Name = cstr(strtab[strx:])
		}

		out = append(out, s)
	}

	f.symtab = out
	return out, nil
}

// u32/readU16/readU64 read integers from a byte slice in the file's byte order.
func u32(b []byte, order binary.ByteOrder) uint32 {
	return order.Uint32(b)
}

// readU16 reads a uint16 from the start of b.
func readU16(b []byte, order binary.ByteOrder) uint16 {
	return order.Uint16(b)
}

// readU64 reads a uint64 from the start of b.
func readU64(b []byte, order binary.ByteOrder) uint64 {
	return order.Uint64(b)
}
