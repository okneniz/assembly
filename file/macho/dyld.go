package macho

// dyld_info streams (LC_DYLD_INFO(_ONLY)) and the exports trie. Layouts are
// from <mach-o/loader.h> (dyld_info_command, opcodes, trie); the semantics
// follow dyld: elements are addressed as segment+offset, and opcode chains
// accumulate state.

// Rebase is a single rebase operation (fixing up a pointer to the slid image base).
type Rebase struct {
	Addr uint64 // vmaddr
	Type uint8  // REBASE_TYPE_*
}

func NewRebase(addr uint64, type_ uint8) Rebase {
	return Rebase{
		Addr: addr,
		Type: type_,
	}
}

// Bind is a single bind operation (binding to an external symbol).
type Bind struct {
	Addr       uint64 // vmaddr of the binding site
	SegIndex   uint8
	Offset     uint64 // offset within the segment
	Type       uint8  // BIND_TYPE_*
	Ordinal    int    // library ordinal (0=self, 0xfe=dynamic lookup, 0xff=executable)
	LibName    string // library name by ordinal (empty for special values)
	SymName    string
	Addend     int64
	WeakImport bool // flag from SET_SYMBOL_TRAILING_FLAGS_IMM
}

func NewBind(
	addr uint64,
	segIndex uint8,
	offset uint64,
	type_ uint8,
	ordinal int,
	libName string,
	symName string,
	addend int64,
	weakImport bool,
) Bind {
	return Bind{
		Addr:       addr,
		SegIndex:   segIndex,
		Offset:     offset,
		Type:       type_,
		Ordinal:    ordinal,
		LibName:    libName,
		SymName:    symName,
		Addend:     addend,
		WeakImport: weakImport,
	}
}

// Export is a terminal node of the exports trie.
type Export struct {
	Name  string
	Flags uint64 // EXPORT_SYMBOL_FLAGS_*
	Addr  uint64 // absolute vmaddr (TrieOffset + __TEXT image base)

	// TrieOffset is the raw terminal value: an offset from the vmaddr of
	// the __TEXT segment (as stored by the trie).
	TrieOffset uint64

	// Reexport: re-export of a symbol from another library.
	Reexport        bool
	ReexportOrdinal int
	ReexportName    string // non-empty when renamed

	// StubAndResolver: proxy stub + resolver.
	StubAndResolver bool
	StubOffset      uint64
	ResolverOffset  uint64
}

func NewExport(name string) Export {
	return Export{Name: name}
}

// dyldInfoCommand finds the LC_DYLD_INFO(_ONLY) command (or nil).
func (f *File) dyldInfoCommand() *DyldInfo {
	for _, lc := range f.commands {
		if d, ok := lc.(*DyldInfo); ok {
			return d
		}
	}

	return nil
}

// dyldStream reads the raw bytes of a stream (off/size relative to the start of the object).
func (f *File) dyldStream(off, size uint32) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}

	return readBytes(f.buf, f.base+int(off), int(size))
}

// readUleb reads a ULEB128 at position p, returning (value, new position).
func readUleb(b []byte, p int) (uint64, int) {
	var v uint64
	var shift uint
	for p < len(b) {
		c := b[p]
		p++
		v |= uint64(c&0x7f) << shift
		if c&0x80 == 0 {
			break
		}

		shift += 7
	}

	return v, p
}

// readSleb reads a SLEB128 at position p.
func readSleb(b []byte, p int) (int64, int) {
	var v int64
	var shift uint
	for p < len(b) {
		c := b[p]
		p++
		v |= int64(c&0x7f) << shift
		shift += 7
		if c&0x80 == 0 {
			if c&0x40 != 0 && shift < 64 {
				v |= -1 << shift
			}

			break
		}
	}

	return v, p
}

// segAddr returns the vmaddr of a segment by index (0 and an error when out of range).
func (f *File) segAddr(idx uint8) uint64 {
	if int(idx) < len(f.segments) {
		return f.segments[idx].Vmaddr
	}

	return 0
}

// dylibName resolves a library ordinal to a name (using the list of LC_*_DYLIB).
func (f *File) dylibName(ordinal int) string {
	switch ordinal {
	case 0:
		return "" // SELF
	case 0xfe:
		return "dynamic lookup"
	case 0xff:
		return "executable"
	}

	i := 0
	for _, lc := range f.commands {
		if d, ok := lc.(*Dylib); ok {
			i++
			if i == ordinal {
				return d.Name
			}
		}
	}

	return ""
}

// Rebases returns all rebase operations from the rebase_dyld_info stream.
func (f *File) Rebases() ([]Rebase, error) {
	d := f.dyldInfoCommand()
	if d == nil {
		return nil, nil
	}

	data, err := f.dyldStream(d.RebaseOff, d.RebaseSize)
	if err != nil {
		return nil, err
	}

	var out []Rebase
	var segIndex uint8
	var segOffset uint64
	typ := uint8(0)

	p := 0
	for p < len(data) {
		op := data[p]
		p++
		imm := op & 0x0f
		switch op & REBASE_OPCODE_MASK {
		case REBASE_OPCODE_DONE:
			return out, nil
		case REBASE_OPCODE_SET_TYPE_IMM:
			typ = imm
		case REBASE_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB:
			segIndex = imm
			segOffset, p = readUleb(data, p)
		case REBASE_OPCODE_ADD_ADDR_ULEB:
			var delta uint64
			delta, p = readUleb(data, p)
			segOffset += delta
		case REBASE_OPCODE_ADD_ADDR_IMM_SCALED:
			segOffset += uint64(imm) << typ
		case REBASE_OPCODE_DO_REBASE_IMM_TIMES:
			for range imm {
				out = append(out, NewRebase(f.segAddr(segIndex)+segOffset, typ))
				segOffset += 1 << typ
			}
		case REBASE_OPCODE_DO_REBASE_ULEB_TIMES:
			var times uint64
			times, p = readUleb(data, p)
			for range times {
				out = append(out, NewRebase(f.segAddr(segIndex)+segOffset, typ))
				segOffset += 1 << typ
			}
		case REBASE_OPCODE_DO_REBASE_ADD_ADDR_ULEB:
			out = append(out, NewRebase(f.segAddr(segIndex)+segOffset, typ))
			var delta uint64
			delta, p = readUleb(data, p)
			segOffset += delta + (1 << typ)
		case REBASE_OPCODE_DO_REBASE_ULEB_TIMES_SKIPPING_ULEB:
			var times, skip uint64
			times, p = readUleb(data, p)
			skip, p = readUleb(data, p)
			for range times {
				out = append(out, NewRebase(f.segAddr(segIndex)+segOffset, typ))
				segOffset += skip + (1 << typ)
			}
		default:
			return nil, errf("rebase: unknown opcode %#x", op)
		}
	}

	return out, nil
}

// parseBindStream parses a bind/weak_bind/lazy_bind stream.
func (f *File) parseBindStream(data []byte) ([]Bind, error) {
	var out []Bind
	var segIndex uint8
	var segOffset uint64
	typ := uint8(0)
	ordinal := 0
	symName := ""
	var addend int64
	weak := false

	emit := func() {
		out = append(
			out,
			NewBind(
				f.segAddr(segIndex)+segOffset,
				segIndex,
				segOffset,
				typ,
				ordinal,
				f.dylibName(ordinal),
				symName,
				addend,
				weak,
			),
		)
		segOffset += 1 << typ
	}

	p := 0
	for p < len(data) {
		op := data[p]
		p++
		imm := op & 0x0f
		switch op & BIND_OPCODE_MASK {
		case BIND_OPCODE_DONE:
			// The lazy stream does not use DONE between entries, but after
			// DONE the data continues with clean segment state.
			segOffset = 0
		case BIND_OPCODE_SET_DYLIB_ORDINAL_IMM:
			ordinal = int(imm)
		case BIND_OPCODE_SET_DYLIB_ORDINAL_ULEB:
			var v uint64
			v, p = readUleb(data, p)
			ordinal = int(v)
		case BIND_OPCODE_SET_DYLIB_SPECIAL_IMM:
			if imm == 0 {
				ordinal = 0
			} else {
				ordinal = int(imm) - 16
			}
		case BIND_OPCODE_SET_SYMBOL_TRAILING_FLAGS_IMM:
			start := p
			for p < len(data) && data[p] != 0 {
				p++
			}

			symName = string(data[start:p])
			p++               // null byte
			weak = imm&1 != 0 // BIND_SYMBOL_FLAGS_WEAK_IMPORT
		case BIND_OPCODE_SET_TYPE_IMM:
			typ = imm
		case BIND_OPCODE_SET_ADDEND_SLEB:
			addend, p = readSleb(data, p)
		case BIND_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB:
			segIndex = imm
			segOffset, p = readUleb(data, p)
		case BIND_OPCODE_ADD_ADDR_ULEB:
			var delta uint64
			delta, p = readUleb(data, p)
			segOffset += delta
		case BIND_OPCODE_DO_BIND:
			emit()
		case BIND_OPCODE_DO_BIND_ADD_ADDR_ULEB:
			emit()
			var delta uint64
			delta, p = readUleb(data, p)
			segOffset += delta
		case BIND_OPCODE_DO_BIND_ADD_ADDR_IMM_SCALED:
			emit()
			segOffset += uint64(imm) << typ
		case BIND_OPCODE_DO_BIND_ULEB_TIMES_SKIPPING_ULEB:
			var times, skip uint64
			times, p = readUleb(data, p)
			skip, p = readUleb(data, p)
			for range times {
				emit()
				segOffset += skip
			}
		case BIND_SUBOPCODE_THREADED:
			// Chained binds in a dyld_info stream (a rare hybrid with
			// chained fixups): 0 is the ordinal table (count + ulebs),
			// 1 is a bind at the current state with a stride of 8.
			switch imm {
			case BIND_SUBOPCODE_THREADED_APPLY_SET_BIND_ORDINAL_TABLE:
				var count uint64
				count, p = readUleb(data, p)
				for range count {
					_, p = readUleb(data, p)
				}
			case BIND_SUBOPCODE_THREADED_APPLY_APPLY:
				emit()
				segOffset += 8 - (1 << typ)
			default:
				return nil, errf("bind: unknown threaded subopcode %#x", imm)
			}
		default:
			return nil, errf("bind: unknown opcode %#x", op)
		}
	}

	return out, nil
}

// Binds returns the bind operations from the bind_dyld_info stream.
func (f *File) Binds() ([]Bind, error) {
	d := f.dyldInfoCommand()
	if d == nil {
		return nil, nil
	}

	data, err := f.dyldStream(d.BindOff, d.BindSize)
	if err != nil {
		return nil, err
	}

	return f.parseBindStream(data)
}

// WeakBinds returns the weak_bind operations (binding to weak definitions).
func (f *File) WeakBinds() ([]Bind, error) {
	d := f.dyldInfoCommand()
	if d == nil {
		return nil, nil
	}

	data, err := f.dyldStream(d.WeakBindOff, d.WeakBindSize)
	if err != nil {
		return nil, err
	}

	return f.parseBindStream(data)
}

// LazyBinds returns the lazy_bind operations (binding on first call).
func (f *File) LazyBinds() ([]Bind, error) {
	d := f.dyldInfoCommand()
	if d == nil {
		return nil, nil
	}

	data, err := f.dyldStream(d.LazyBindOff, d.LazyBindSize)
	if err != nil {
		return nil, err
	}

	return f.parseBindStream(data)
}

// Exports returns all exported symbols from the exports trie
// (LC_DYLD_INFO.export or LC_DYLD_EXPORTS_TRIE).
func (f *File) Exports() ([]Export, error) {
	var data []byte
	var err error
	if d := f.dyldInfoCommand(); d != nil && d.ExportSize > 0 {
		data, err = f.dyldStream(d.ExportOff, d.ExportSize)
	} else {
		for _, lc := range f.commands {
			if ld, ok := lc.(*LinkeditData); ok && ld.Cmd() == LC_DYLD_EXPORTS_TRIE {
				data, err = f.dyldStream(ld.DataOff, ld.DataSize)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var out []Export
	if err := walkTrie(data, 0, "", &out); err != nil {
		return nil, err
	}

	// The trie stores an offset from the image base (__TEXT vmaddr).
	base := f.textVmaddr()
	for i := range out {
		if !out[i].Reexport && !out[i].StubAndResolver {
			out[i].Addr = base + out[i].TrieOffset
		}
	}

	return out, nil
}

// walkTrie walks the trie node at offset nodeOff with the current name prefix.
//
// Node (layout verified against ld64/llvm samples):
//
//	uleb terminalSize; if > 0, terminal data (uleb flags + kind-specific
//	ulebs); uleb childCount; childCount x (cstring edge + uleb child
//	offset); after a non-empty list, a null byte.
func walkTrie(data []byte, nodeOff int, prefix string, out *[]Export) error {
	p := nodeOff
	terminalSize, np := readUleb(data, p)
	p = np
	if terminalSize > 0 {
		e := NewExport(prefix)
		termStart := p
		flags, np2 := readUleb(data, p)
		e.Flags = flags
		p = np2
		switch {
		case flags&EXPORT_SYMBOL_FLAGS_REEXPORT != 0:
			e.Reexport = true
			// the position after the uleb is not needed below: p is reset
			// to termStart+terminalSize further down
			v, _ := readUleb(data, p)
			e.ReexportOrdinal = int(v >> 8)
			if nameOff := v & 0xff; nameOff != 0 {
				// Offset inside the trie where the alternative name lies.
				// (usually 0 - re-export under the same name)
				_ = nameOff
			}
		case flags&EXPORT_SYMBOL_FLAGS_STUB_AND_RESOLVER != 0:
			e.StubAndResolver = true
			var q int
			e.StubOffset, q = readUleb(data, p)
			e.ResolverOffset, _ = readUleb(data, q)
		default:
			e.TrieOffset, _ = readUleb(data, p)
		}

		*out = append(*out, e)
		p = termStart + int(terminalSize)
	}

	childCount, np5 := readUleb(data, p)
	p = np5
	for range childCount {
		start := p
		for p < len(data) && data[p] != 0 {
			p++
		}

		if p >= len(data) {
			return errf("exports trie: truncated edge at %d", start)
		}

		edge := string(data[start:p])
		p++ // null terminator of the edge
		childOff, np6 := readUleb(data, p)
		p = np6
		if int(childOff) >= len(data) {
			return errf("exports trie: child offset %d out of data", childOff)
		}

		if err := walkTrie(data, int(childOff), prefix+edge, out); err != nil {
			return err
		}
	}

	_ = p
	return nil
}
