package macho

// LoadCommand is a typed load command. Each structure below mirrors the
// layout from <mach-o/loader.h>; Generic keeps the raw bytes of unknown and
// historical commands (LC_FVMFILE, LC_PREPAGE, ...), so they are not lost
// from the output.

// LoadCommand is the interface of a typed load command.
type LoadCommand interface {
	Cmd() Cmd
	Cmdsize() uint32
}

// Generic is an unknown/historical command: raw bytes after the header.
type Generic struct {
	cmd  Cmd
	size uint32
	Raw  []byte // payload after cmd/cmdsize
}

func NewGeneric(cmd Cmd, size uint32, raw []byte) *Generic {
	return &Generic{
		cmd:  cmd,
		size: size,
		Raw:  raw,
	}
}

func (g *Generic) Cmd() Cmd {
	return g.cmd
}
func (g *Generic) Cmdsize() uint32 {
	return g.size
}

// --- Commands with a string (lc_str: offset from the start of the command) ---

// Dylib — LC_LOAD_DYLIB / LC_ID_DYLIB / LC_LOAD_WEAK_DYLIB /
// LC_REEXPORT_DYLIB / LC_LAZY_LOAD_DYLIB / LC_LOAD_UPWARD_DYLIB.
type Dylib struct {
	cmd                  Cmd
	size                 uint32
	Name                 string
	Timestamp            uint32
	CurrentVersion       uint32 // packed: mm.mm.mmhh.pp
	CompatibilityVersion uint32
}

func NewDylib(
	cmd Cmd,
	size uint32,
	name string,
	timestamp uint32,
	currentVersion uint32,
	compatibilityVersion uint32,
) *Dylib {
	return &Dylib{
		cmd:                  cmd,
		size:                 size,
		Name:                 name,
		Timestamp:            timestamp,
		CurrentVersion:       currentVersion,
		CompatibilityVersion: compatibilityVersion,
	}
}

func (c *Dylib) Cmd() Cmd {
	return c.cmd
}
func (c *Dylib) Cmdsize() uint32 {
	return c.size
}

// StringCmd — LC_LOAD_DYLINKER / LC_ID_DYLINKER / LC_DYLD_ENVIRONMENT /
// LC_RPATH / LC_SUB_FRAMEWORK / LC_SUB_UMBRELLA / LC_SUB_CLIENT /
// LC_SUB_LIBRARY.
type StringCmd struct {
	cmd  Cmd
	size uint32
	Str  string
}

func NewStringCmd(cmd Cmd, size uint32, str string) *StringCmd {
	return &StringCmd{
		cmd:  cmd,
		size: size,
		Str:  str,
	}
}

func (c *StringCmd) Cmd() Cmd {
	return c.cmd
}
func (c *StringCmd) Cmdsize() uint32 {
	return c.size
}

// --- Symbol tables ---

// Symtab is LC_SYMTAB: the location of the nlist table and strings.
type Symtab struct {
	size    uint32
	Symoff  uint32
	Nsyms   uint32
	Stroff  uint32
	Strsize uint32
}

func NewSymtab(size uint32, symoff uint32, nsyms uint32, stroff uint32, strsize uint32) *Symtab {
	return &Symtab{
		size:    size,
		Symoff:  symoff,
		Nsyms:   nsyms,
		Stroff:  stroff,
		Strsize: strsize,
	}
}

func (c *Symtab) Cmd() Cmd {
	return LC_SYMTAB
}
func (c *Symtab) Cmdsize() uint32 {
	return c.size
}

// Dysymtab is LC_DYSYMTAB: ranges and auxiliary tables.
type Dysymtab struct {
	size uint32

	ToLocalSymbols      uint32 // start index of locals
	NLocalSymbols       uint32
	ToExtDefinedSymbols uint32 // start index of externally defined
	NExtDefinedSymbols  uint32
	ToUndefSymbols      uint32 // start index of undefined
	NUndefSymbols       uint32

	TOCOffset      uint32
	NToc           uint32
	ModTabOffset   uint32
	NModTab        uint32
	ExtRefSymOff   uint32
	NExtRefSyms    uint32
	IndirectSymOff uint32
	NIndirectSyms  uint32

	ExtRelOffset uint32
	NExtRel      uint32
	LocRelOffset uint32
	NLocRel      uint32
}

func NewDysymtab(size uint32) *Dysymtab {
	return &Dysymtab{size: size}
}

func (c *Dysymtab) Cmd() Cmd {
	return LC_DYSYMTAB
}
func (c *Dysymtab) Cmdsize() uint32 {
	return c.size
}

// --- Metadata ---

// UUID — LC_UUID.
type UUID struct {
	size uint32
	ID   [16]byte
}

func NewUUID(size uint32) *UUID {
	return &UUID{size: size}
}

func (c *UUID) Cmd() Cmd {
	return LC_UUID
}
func (c *UUID) Cmdsize() uint32 {
	return c.size
}

// LinkeditData is linkedit_data_command: a blob in __LINKEDIT. It covers
// LC_CODE_SIGNATURE, LC_SEGMENT_SPLIT_INFO, LC_FUNCTION_STARTS,
// LC_DATA_IN_CODE, LC_DYLIB_CODE_SIGN_DRS, LC_LINKER_OPTIMIZATION_HINT,
// LC_DYLD_EXPORTS_TRIE, LC_DYLD_CHAINED_FIXUPS, LC_ATOM_INFO.
type LinkeditData struct {
	cmd      Cmd
	size     uint32
	DataOff  uint32
	DataSize uint32
}

func NewLinkeditData(cmd Cmd, size uint32, dataOff uint32, dataSize uint32) *LinkeditData {
	return &LinkeditData{
		cmd:      cmd,
		size:     size,
		DataOff:  dataOff,
		DataSize: dataSize,
	}
}

func (c *LinkeditData) Cmd() Cmd {
	return c.cmd
}
func (c *LinkeditData) Cmdsize() uint32 {
	return c.size
}

// EncryptionInfo — LC_ENCRYPTION_INFO / LC_ENCRYPTION_INFO_64.
type EncryptionInfo struct {
	cmd       Cmd
	size      uint32
	CryptOff  uint32
	CryptSize uint32
	CryptID   uint32
	Pad       uint32 // _64 only
}

func NewEncryptionInfo(
	cmd Cmd,
	size uint32,
	cryptOff uint32,
	cryptSize uint32,
	cryptID uint32,
) *EncryptionInfo {
	return &EncryptionInfo{
		cmd:       cmd,
		size:      size,
		CryptOff:  cryptOff,
		CryptSize: cryptSize,
		CryptID:   cryptID,
	}
}

func (c *EncryptionInfo) Cmd() Cmd {
	return c.cmd
}
func (c *EncryptionInfo) Cmdsize() uint32 {
	return c.size
}

// VersionMin is version_min_command (LC_VERSION_MIN_MACOSX etc.).
type VersionMin struct {
	cmd     Cmd
	size    uint32
	Version uint32 // packed: x.y.z
	Sdk     uint32
}

func NewVersionMin(cmd Cmd, size uint32, version uint32, sdk uint32) *VersionMin {
	return &VersionMin{
		cmd:     cmd,
		size:    size,
		Version: version,
		Sdk:     sdk,
	}
}

func (c *VersionMin) Cmd() Cmd {
	return c.cmd
}
func (c *VersionMin) Cmdsize() uint32 {
	return c.size
}

// BuildTool is an element of the tools array of the LC_BUILD_VERSION command.
type BuildTool struct {
	Tool    Tool
	Version uint32
}

func NewBuildTool(tool Tool, version uint32) BuildTool {
	return BuildTool{
		Tool:    tool,
		Version: version,
	}
}

// BuildVersion — LC_BUILD_VERSION.
type BuildVersion struct {
	size     uint32
	Platform Platform
	Minos    uint32
	Sdk      uint32
	Ntools   uint32
	Tools    []BuildTool
}

func NewBuildVersion(
	size uint32,
	platform Platform,
	minos uint32,
	sdk uint32,
	ntools uint32,
) *BuildVersion {
	return &BuildVersion{
		size:     size,
		Platform: platform,
		Minos:    minos,
		Sdk:      sdk,
		Ntools:   ntools,
	}
}

func (c *BuildVersion) Cmd() Cmd {
	return LC_BUILD_VERSION
}
func (c *BuildVersion) Cmdsize() uint32 {
	return c.size
}

// DyldInfo is LC_DYLD_INFO / LC_DYLD_INFO_ONLY: offsets of the opcode
// streams and the exports trie (the streams are parsed by File methods in
// dyld.go).
type DyldInfo struct {
	cmd          Cmd
	size         uint32
	RebaseOff    uint32
	RebaseSize   uint32
	BindOff      uint32
	BindSize     uint32
	WeakBindOff  uint32
	WeakBindSize uint32
	LazyBindOff  uint32
	LazyBindSize uint32
	ExportOff    uint32
	ExportSize   uint32
}

func NewDyldInfo(cmd Cmd, size uint32) *DyldInfo {
	return &DyldInfo{
		cmd:  cmd,
		size: size,
	}
}

func (c *DyldInfo) Cmd() Cmd {
	return c.cmd
}
func (c *DyldInfo) Cmdsize() uint32 {
	return c.size
}

// Main is LC_MAIN: the entry point as an offset into __TEXT from the base address.
type Main struct {
	size      uint32
	EntryOff  uint64
	StackSize uint64
}

func NewMain(size uint32, entryOff uint64, stackSize uint64) *Main {
	return &Main{
		size:      size,
		EntryOff:  entryOff,
		StackSize: stackSize,
	}
}

func (c *Main) Cmd() Cmd {
	return LC_MAIN
}
func (c *Main) Cmdsize() uint32 {
	return c.size
}

// SourceVersion — LC_SOURCE_VERSION: packed A.B.C.D.E.
type SourceVersion struct {
	size    uint32
	Version uint64
}

func NewSourceVersion(size uint32, version uint64) *SourceVersion {
	return &SourceVersion{
		size:    size,
		Version: version,
	}
}

func (c *SourceVersion) Cmd() Cmd {
	return LC_SOURCE_VERSION
}
func (c *SourceVersion) Cmdsize() uint32 {
	return c.size
}

// LinkerOption is LC_LINKER_OPTION: -foo bar linker option pairs.
type LinkerOption struct {
	size    uint32
	Options []string
}

func NewLinkerOption(size uint32) *LinkerOption {
	return &LinkerOption{size: size}
}

func (c *LinkerOption) Cmd() Cmd {
	return LC_LINKER_OPTION
}
func (c *LinkerOption) Cmdsize() uint32 {
	return c.size
}

// Note is LC_NOTE: arbitrary named data.
type Note struct {
	size      uint32
	DataOwner string // 16 bytes
	Offset    uint64
	Size      uint64
}

func NewNote(size uint32, dataOwner string, offset uint64, size_ uint64) *Note {
	return &Note{
		size:      size,
		DataOwner: dataOwner,
		Offset:    offset,
		Size:      size_,
	}
}

func (c *Note) Cmd() Cmd {
	return LC_NOTE
}
func (c *Note) Cmdsize() uint32 {
	return c.size
}

// ThreadState is a single thread state block of LC_THREAD/LC_UNIXTHREAD.
type ThreadState struct {
	Flavor uint32
	Count  uint32 // number of uint32 words in State
	State  []uint32
}

// Thread is LC_THREAD / LC_UNIXTHREAD: a chain of state blocks (each
// architecture has its own flavors and layouts).
type Thread struct {
	cmd    Cmd
	size   uint32
	Flavor uint32
	Count  uint32
	State  []uint32
}

func NewThread(cmd Cmd, size uint32, flavor uint32, count uint32) *Thread {
	return &Thread{
		cmd:    cmd,
		size:   size,
		Flavor: flavor,
		Count:  count,
	}
}

func (c *Thread) Cmd() Cmd {
	return c.cmd
}
func (c *Thread) Cmdsize() uint32 {
	return c.size
}

// FilesetEntry is LC_FILESET_ENTRY: a nested Mach-O (dyld shared cache).
type FilesetEntry struct {
	size     uint32
	Vmaddr   uint64
	FileOff  uint64
	EntryID  string
	Reserved uint32
}

func NewFilesetEntry(
	size uint32,
	vmaddr uint64,
	fileOff uint64,
	entryID string,
	reserved uint32,
) *FilesetEntry {
	return &FilesetEntry{
		size:     size,
		Vmaddr:   vmaddr,
		FileOff:  fileOff,
		EntryID:  entryID,
		Reserved: reserved,
	}
}

func (c *FilesetEntry) Cmd() Cmd {
	return LC_FILESET_ENTRY
}
func (c *FilesetEntry) Cmdsize() uint32 {
	return c.size
}

// Routines is LC_ROUTINES / LC_ROUTINES_64: the address of init routines (historical).
type Routines struct {
	cmd         Cmd
	size        uint32
	InitAddress uint64
	InitModule  uint64
}

func NewRoutines(cmd Cmd, size uint32, initAddress uint64, initModule uint64) *Routines {
	return &Routines{
		cmd:         cmd,
		size:        size,
		InitAddress: initAddress,
		InitModule:  initModule,
	}
}

func (c *Routines) Cmd() Cmd {
	return c.cmd
}
func (c *Routines) Cmdsize() uint32 {
	return c.size
}

// TwoLevelHints is LC_TWOLEVEL_HINTS (historical).
type TwoLevelHints struct {
	size   uint32
	Offset uint32
	Nhints uint32
}

func NewTwoLevelHints(size uint32, offset uint32, nhints uint32) *TwoLevelHints {
	return &TwoLevelHints{
		size:   size,
		Offset: offset,
		Nhints: nhints,
	}
}

func (c *TwoLevelHints) Cmd() Cmd {
	return LC_TWOLEVEL_HINTS
}
func (c *TwoLevelHints) Cmdsize() uint32 {
	return c.size
}

// --- Parsing ---

// LoadCommands returns all load commands of the file in file order.
func (f *File) LoadCommands() []LoadCommand {
	return f.commands
}

// parseLoadCommands walks the ncmds commands, validates cmdsize, and builds
// the typed structures; segments are additionally laid out into
// f.segments/f.sections.
func (f *File) parseLoadCommands() error {
	hdrSize := headerSize32
	if f.is64 {
		hdrSize = headerSize64
	}

	pos := f.base + hdrSize

	for i := range f.hdr.Ncmds {
		cmdRaw, err := readU32At(f.buf, f.order, pos)
		if err != nil {
			return err
		}

		cmdsize, err := readU32At(f.buf, f.order, pos+4)
		if err != nil {
			return err
		}

		if cmdsize < 8 {
			return errf(
				"load command %d (%s): cmdsize=%d is less than the minimum of 8",
				i,
				Cmd(cmdRaw),
				cmdsize,
			)
		}

		align := uint32(4)
		if f.is64 {
			align = 8
		}

		if cmdsize%align != 0 {
			return errf(
				"load command %d (%s): cmdsize=%d is not a multiple of %d",
				i,
				Cmd(cmdRaw),
				cmdsize,
				align,
			)
		}

		lc, err := f.parseCommand(pos, Cmd(cmdRaw), cmdsize)
		if err != nil {
			return errf("load command %d (%s): %v", i, Cmd(cmdRaw), err)
		}

		if lc != nil {
			f.commands = append(f.commands, lc)
		}

		pos += int(cmdsize)
	}

	return nil
}

// parseCommand parses a single command at offset off (absolute).
func (f *File) parseCommand(off int, cmd Cmd, cmdsize uint32) (LoadCommand, error) {
	str := func(lcStrOff uint32) (string, error) {
		if int(lcStrOff)+1 > int(cmdsize) {
			return "", errf("lc_str offset %d outside the command (cmdsize %d)", lcStrOff, cmdsize)
		}

		raw, err := readBytes(f.buf, off+int(lcStrOff), int(cmdsize-lcStrOff))
		if err != nil {
			return "", err
		}

		return cstr(raw), nil
	}
	u32 := func(o int) (uint32, error) {
		return readU32At(f.buf, f.order, off+o)
	}
	u64 := func(o int) (uint64, error) {
		return readU64At(f.buf, f.order, off+o)
	}

	switch cmd {
	case LC_SEGMENT, LC_SEGMENT_64:
		return f.parseSegment(off, cmd, cmdsize)
	case LC_SYMTAB:
		symoff, err := u32(8)
		if err != nil {
			return nil, err
		}

		nsyms, err := u32(12)
		if err != nil {
			return nil, err
		}

		stroff, err := u32(16)
		if err != nil {
			return nil, err
		}

		strsize, err := u32(20)
		if err != nil {
			return nil, err
		}

		return NewSymtab(cmdsize, symoff, nsyms, stroff, strsize), nil
	case LC_DYSYMTAB:
		d := NewDysymtab(cmdsize)
		fields := []*uint32{
			&d.ToLocalSymbols, &d.NLocalSymbols,
			&d.ToExtDefinedSymbols, &d.NExtDefinedSymbols,
			&d.ToUndefSymbols, &d.NUndefSymbols,
			&d.TOCOffset, &d.NToc,
			&d.ModTabOffset, &d.NModTab,
			&d.ExtRefSymOff, &d.NExtRefSyms,
			&d.IndirectSymOff, &d.NIndirectSyms,
			&d.ExtRelOffset, &d.NExtRel,
			&d.LocRelOffset, &d.NLocRel,
		}
		for i, fp := range fields {
			v, err := u32(8 + i*4)
			if err != nil {
				return nil, err
			}

			*fp = v
		}

		return d, nil
	case LC_LOAD_DYLIB,
		LC_ID_DYLIB,
		LC_LOAD_WEAK_DYLIB,
		LC_REEXPORT_DYLIB,
		LC_LAZY_LOAD_DYLIB,
		LC_LOAD_UPWARD_DYLIB:
		nameOff, err := u32(8)
		if err != nil {
			return nil, err
		}

		ts, err := u32(12)
		if err != nil {
			return nil, err
		}

		cur, err := u32(16)
		if err != nil {
			return nil, err
		}

		compat, err := u32(20)
		if err != nil {
			return nil, err
		}

		name, err := str(nameOff)
		if err != nil {
			return nil, err
		}

		return NewDylib(cmd, cmdsize, name, ts, cur, compat), nil
	case LC_LOAD_DYLINKER, LC_ID_DYLINKER, LC_DYLD_ENVIRONMENT, LC_RPATH,
		LC_SUB_FRAMEWORK, LC_SUB_UMBRELLA, LC_SUB_CLIENT, LC_SUB_LIBRARY:
		nameOff, err := u32(8)
		if err != nil {
			return nil, err
		}

		s, err := str(nameOff)
		if err != nil {
			return nil, err
		}

		return NewStringCmd(cmd, cmdsize, s), nil
	case LC_UUID:
		raw, err := readBytes(f.buf, off+8, 16)
		if err != nil {
			return nil, err
		}

		u := NewUUID(cmdsize)
		copy(u.ID[:], raw)
		return u, nil
	case LC_CODE_SIGNATURE, LC_SEGMENT_SPLIT_INFO, LC_FUNCTION_STARTS,
		LC_DATA_IN_CODE, LC_DYLIB_CODE_SIGN_DRS, LC_LINKER_OPTIMIZATION_HINT,
		LC_DYLD_EXPORTS_TRIE, LC_DYLD_CHAINED_FIXUPS, LC_ATOM_INFO:
		dataOff, err := u32(8)
		if err != nil {
			return nil, err
		}

		dataSize, err := u32(12)
		if err != nil {
			return nil, err
		}

		return NewLinkeditData(cmd, cmdsize, dataOff, dataSize), nil
	case LC_ENCRYPTION_INFO, LC_ENCRYPTION_INFO_64:
		cryptOff, err := u32(8)
		if err != nil {
			return nil, err
		}

		cryptSize, err := u32(12)
		if err != nil {
			return nil, err
		}

		cryptID, err := u32(16)
		if err != nil {
			return nil, err
		}

		e := NewEncryptionInfo(cmd, cmdsize, cryptOff, cryptSize, cryptID)
		if cmd == LC_ENCRYPTION_INFO_64 {
			if e.Pad, err = u32(20); err != nil {
				return nil, err
			}
		}

		return e, nil
	case LC_VERSION_MIN_MACOSX,
		LC_VERSION_MIN_IPHONEOS,
		LC_VERSION_MIN_TVOS,
		LC_VERSION_MIN_WATCHOS:
		version, err := u32(8)
		if err != nil {
			return nil, err
		}

		sdk, err := u32(12)
		if err != nil {
			return nil, err
		}

		return NewVersionMin(cmd, cmdsize, version, sdk), nil
	case LC_BUILD_VERSION:
		platform, err := u32(8)
		if err != nil {
			return nil, err
		}

		minos, err := u32(12)
		if err != nil {
			return nil, err
		}

		sdk, err := u32(16)
		if err != nil {
			return nil, err
		}

		ntools, err := u32(20)
		if err != nil {
			return nil, err
		}

		b := NewBuildVersion(cmdsize, Platform(platform), minos, sdk, ntools)
		for i := range ntools {
			tool, err := u32(24 + int(i)*8)
			if err != nil {
				return nil, err
			}

			ver, err := u32(28 + int(i)*8)
			if err != nil {
				return nil, err
			}

			b.Tools = append(b.Tools, NewBuildTool(Tool(tool), ver))
		}

		return b, nil
	case LC_DYLD_INFO, LC_DYLD_INFO_ONLY:
		d := NewDyldInfo(cmd, cmdsize)
		fields := []*uint32{
			&d.RebaseOff, &d.RebaseSize,
			&d.BindOff, &d.BindSize,
			&d.WeakBindOff, &d.WeakBindSize,
			&d.LazyBindOff, &d.LazyBindSize,
			&d.ExportOff, &d.ExportSize,
		}
		for i, fp := range fields {
			v, err := u32(8 + i*4)
			if err != nil {
				return nil, err
			}

			*fp = v
		}

		return d, nil
	case LC_MAIN:
		entryOff, err := u64(8)
		if err != nil {
			return nil, err
		}

		stackSize, err := u64(16)
		if err != nil {
			return nil, err
		}

		return NewMain(cmdsize, entryOff, stackSize), nil
	case LC_SOURCE_VERSION:
		v, err := u64(8)
		if err != nil {
			return nil, err
		}

		return NewSourceVersion(cmdsize, v), nil
	case LC_LINKER_OPTION:
		count, err := u32(8)
		if err != nil {
			return nil, err
		}

		raw, err := readBytes(f.buf, off+12, int(cmdsize-12))
		if err != nil {
			return nil, err
		}

		l := NewLinkerOption(cmdsize)
		parts := splitNul(raw)
		for i := uint32(0); i < count && int(i) < len(parts); i++ {
			l.Options = append(l.Options, parts[i])
		}

		return l, nil
	case LC_NOTE:
		raw, err := readBytes(f.buf, off+8, 16)
		if err != nil {
			return nil, err
		}

		o, err := u64(24)
		if err != nil {
			return nil, err
		}

		s, err := u64(32)
		if err != nil {
			return nil, err
		}

		return NewNote(cmdsize, cstr(raw), o, s), nil
	case LC_THREAD, LC_UNIXTHREAD:
		flavor, err := u32(8)
		if err != nil {
			return nil, err
		}

		count, err := u32(12)
		if err != nil {
			return nil, err
		}

		if int(16+count*4) > int(cmdsize) {
			return nil, errf("thread state (%d words) does not fit into cmdsize %d", count, cmdsize)
		}

		th := NewThread(cmd, cmdsize, flavor, count)
		for i := range count {
			w, err := u32(16 + int(i)*4)
			if err != nil {
				return nil, err
			}

			th.State = append(th.State, w)
		}

		return th, nil
	case LC_FILESET_ENTRY:
		vmaddr, err := u64(8)
		if err != nil {
			return nil, err
		}

		fileoff, err := u64(16)
		if err != nil {
			return nil, err
		}

		idOff, err := u32(24)
		if err != nil {
			return nil, err
		}

		reserved, err := u32(28)
		if err != nil {
			return nil, err
		}

		id, err := str(idOff)
		if err != nil {
			return nil, err
		}

		return NewFilesetEntry(cmdsize, vmaddr, fileoff, id, reserved), nil
	case LC_ROUTINES, LC_ROUTINES_64:
		init, err := u64(8)
		if err != nil {
			return nil, err
		}

		mod, err := u64(16)
		if err != nil {
			return nil, err
		}

		return NewRoutines(cmd, cmdsize, init, mod), nil
	case LC_TWOLEVEL_HINTS:
		o, err := u32(8)
		if err != nil {
			return nil, err
		}

		n, err := u32(12)
		if err != nil {
			return nil, err
		}

		return NewTwoLevelHints(cmdsize, o, n), nil
	default:
		// Unknown/historical command: raw bytes.
		raw, err := readBytes(f.buf, off+8, int(cmdsize-8))
		if err != nil {
			return nil, err
		}

		return NewGeneric(cmd, cmdsize, raw), nil
	}
}

// splitNul splits null-terminated strings out of a block of bytes.
func splitNul(b []byte) []string {
	var out []string
	start := 0
	for i, c := range b {
		if c == 0 {
			out = append(out, string(b[start:i]))
			start = i + 1
		}
	}

	if start < len(b) {
		out = append(out, string(b[start:]))
	}

	return out
}
