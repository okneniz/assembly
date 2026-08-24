package elf

import (
	"encoding/binary"
	"errors"
)

// ErrNoSymbols is returned by Symbols and DynamicSymbols when the file has no
// corresponding section (a stripped binary) - the contract matches debug/elf.
var ErrNoSymbols = errors.New("no symbol section")

// Symbol is an Elf{32,64}_Sym with normalized fields. st_info packs
// bind/type (the Bind/Type methods), st_other is the visibility (Visibility),
// and st_shndx is the index of the defining section (SHN_XINDEX is resolved
// through SHT_SYMTAB_SHNDX).
//
// Elf64_Sym (24 bytes): st_name(4)@0 st_info(1)@4 st_other(1)@5 st_shndx(2)@6
//
//	st_value(8)@8 st_size(8)@16
//
// Elf32_Sym (16 bytes): st_name(4)@0 st_value(4)@4 st_size(4)@8
//
//	st_info(1)@12 st_other(1)@13 st_shndx(2)@14
type Symbol struct {
	Name  string
	Value uint64
	Size  uint64
	Info  uint8
	Other uint8
	Shndx uint32
}

// Bind returns the symbol binding class (the high nibble of st_info).
func (s Symbol) Bind() SymbolBind {
	return SymbolBind(s.Info >> 4)
}

// Type returns the symbol type (the low nibble of st_info).
func (s Symbol) Type() SymbolType {
	return SymbolType(s.Info & 0xf)
}

// Visibility returns the symbol visibility (the low 2 bits of st_other).
func (s Symbol) Visibility() SymbolVisibility {
	return SymbolVisibility(s.Other & 3)
}

// SectionIndexName describes a special section index; "" for a regular one.
func (s Symbol) SectionIndexName() string {
	return shndxName(s.Shndx)
}

const (
	sym64Size  = 24
	sym32Size  = 16
	sym64Info  = 4
	sym64Other = 5
	sym64Shndx = 6
	sym64Value = 8
	sym64SizeF = 16
	sym32Value = 4
	sym32SizeF = 8
	sym32Info  = 12
	sym32Other = 13
	sym32Shndx = 14
)

// Symbols returns the .symtab table (linker symbols); ErrNoSymbols for
// stripped binaries. The result is cached.
func (f *File) Symbols() ([]Symbol, error) {
	if !f.symtabDone {
		f.symtab, f.symtabErr = f.parseSymbols(SHT_SYMTAB)
		f.symtabDone = true
	}

	return f.symtab, f.symtabErr
}

// DynamicSymbols returns the .dynsym table (runtime symbols). The result is
// cached.
func (f *File) DynamicSymbols() ([]Symbol, error) {
	if !f.dynsymDone {
		f.dynsym, f.dynsymErr = f.parseSymbols(SHT_DYNSYM)
		f.dynsymDone = true
	}

	return f.dynsym, f.dynsymErr
}

// parseSymbols finds the section of the given type, its string table (sh_link),
// and resolves SHT_SYMTAB_SHNDX for symbols with st_shndx == SHN_XINDEX.
func (f *File) parseSymbols(typ SectionType) ([]Symbol, error) {
	var symsec *Section
	for _, s := range f.sections {
		if s.Type == typ {
			symsec = s
			break
		}
	}

	if symsec == nil {
		return nil, ErrNoSymbols
	}

	if symsec.Link >= uint32(len(f.sections)) {
		return nil, errf(
			"symbol section %s: sh_link=%d is outside the section table",
			symsec.Name,
			symsec.Link,
		)
	}

	strsec := f.sections[symsec.Link]
	if strsec.Type != SHT_STRTAB {
		return nil, errf(
			"symbol section %s: sh_link points to %s, want SHT_STRTAB",
			symsec.Name,
			strsec.Type,
		)
	}

	strtab, err := strsec.Data()
	if err != nil {
		return nil, err
	}

	data, err := symsec.Data()
	if err != nil {
		return nil, err
	}

	size := sym32Size
	if f.class == CLASS64 {
		size = sym64Size
	}

	if symsec.Entsize != 0 && symsec.Entsize < uint64(size) {
		return nil, errf(
			"symbol section %s: entsize=%d is less than the minimum %d",
			symsec.Name,
			symsec.Entsize,
			size,
		)
	}

	entsize := symsec.Entsize
	if entsize == 0 {
		entsize = uint64(size)
	}

	// SHT_SYMTAB_SHNDX: an array of uint32 extended indices for XINDEX symbols.
	var xndx []byte
	for _, s := range f.sections {
		if s.Type == SHT_SYMTAB_SHNDX && s.Link == uint32(symsec.index) {
			if xndx, err = s.Data(); err != nil {
				return nil, err
			}

			break
		}
	}

	// The first entry of the table is the reserved zero symbol (all zeros);
	// like debug/elf, we skip it.
	out := make([]Symbol, 0, len(data)/int(entsize))
	for off := int(entsize); off+size <= len(data); off += int(entsize) {
		var sym Symbol
		nameOff := u32(data[off:], f.order)
		if f.class == CLASS64 {
			sym.Info = data[off+sym64Info]
			sym.Other = data[off+sym64Other]
			sym.Shndx = uint32(u16(data[off+sym64Shndx:], f.order))
			sym.Value = u64(data[off+sym64Value:], f.order)
			sym.Size = u64(data[off+sym64SizeF:], f.order)
		} else {
			sym.Value = uint64(u32(data[off+sym32Value:], f.order))
			sym.Size = uint64(u32(data[off+sym32SizeF:], f.order))
			sym.Info = data[off+sym32Info]
			sym.Other = data[off+sym32Other]
			sym.Shndx = uint32(u16(data[off+sym32Shndx:], f.order))
		}

		if nameOff < uint32(len(strtab)) {
			sym.Name = cstr(strtab[nameOff:])
		}

		if sym.Shndx == SHN_XINDEX {
			i := off / int(entsize)
			if xndx != nil && (i+1)*4 <= len(xndx) {
				sym.Shndx = u32(xndx[i*4:], f.order)
			}
		}

		out = append(out, sym)
	}

	return out, nil
}

// u16/u32/u64 read little/big-endian values from a byte slice.
func u16(b []byte, order binary.ByteOrder) uint16 {
	return order.Uint16(b)
}
func u32(b []byte, order binary.ByteOrder) uint32 {
	return order.Uint32(b)
}
func u64(b []byte, order binary.ByteOrder) uint64 {
	return order.Uint64(b)
}
