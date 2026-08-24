package macho

import "encoding/binary"

// Misc __LINKEDIT metadata and entry points.

// Dice is a data-in-code entry: a range in __text containing non-code.
type Dice struct {
	Offset uint32 // offset from the start of __TEXT
	Length uint16
	Kind   DiceKind
}

func NewDice(offset uint32, length uint16, kind DiceKind) Dice {
	return Dice{
		Offset: offset,
		Length: length,
		Kind:   kind,
	}
}

// DataInCode returns the LC_DATA_IN_CODE table.
func (f *File) DataInCode() ([]Dice, error) {
	data, err := f.linkeditBlob(LC_DATA_IN_CODE)
	if err != nil || data == nil {
		return nil, err
	}

	var out []Dice
	for off := 0; off+8 <= len(data); off += 8 {
		out = append(
			out,
			NewDice(
				u32(data[off:], f.order),
				readU16(data[off+4:], f.order),
				DiceKind(readU16(data[off+6:], f.order)),
			),
		)
	}

	return out, nil
}

// FunctionStarts returns function start addresses (LC_FUNCTION_STARTS):
// a stream of ULEB128 deltas from the __TEXT segment vmaddr.
func (f *File) FunctionStarts() ([]uint64, error) {
	data, err := f.linkeditBlob(LC_FUNCTION_STARTS)
	if err != nil || data == nil {
		return nil, err
	}

	base := f.textVmaddr()

	var out []uint64
	addr := base
	p := 0
	for p < len(data) {
		delta, np := readUleb(data, p)
		p = np
		if delta == 0 {
			continue // alignment terminator
		}

		addr += delta
		out = append(out, addr)
	}

	return out, nil
}

// textVmaddr returns the vmaddr of the __TEXT segment (0 if absent).
func (f *File) textVmaddr() uint64 {
	for _, s := range f.segments {
		if s.SegName == "__TEXT" {
			return s.Vmaddr
		}
	}

	if len(f.segments) > 0 {
		return f.segments[0].Vmaddr
	}

	return 0
}

// UUIDCommand returns the file's UUID (nil if there is no command).
func (f *File) UUIDCommand() *[16]byte {
	for _, lc := range f.commands {
		if u, ok := lc.(*UUID); ok {
			id := u.ID
			return &id
		}
	}

	return nil
}

// BuildVersionCommand returns LC_BUILD_VERSION (nil if absent).
func (f *File) BuildVersionCommand() *BuildVersion {
	for _, lc := range f.commands {
		if b, ok := lc.(*BuildVersion); ok {
			return b
		}
	}

	return nil
}

// Entry returns the virtual address of the entry point: from LC_MAIN
// (entryoff + __TEXT vmaddr) or from the register context of LC_UNIXTHREAD
// (pc for arm64/x86_64). false means the entry point was not found.
func (f *File) Entry() (uint64, bool) {
	for _, lc := range f.commands {
		switch c := lc.(type) {
		case *Main:
			return f.textVmaddr() + c.EntryOff, true
		case *Thread:
			if c.Cmd() != LC_UNIXTHREAD {
				continue
			}

			switch f.hdr.CpuType {
			case CPU_TYPE_ARM64, CPU_TYPE_ARM64_32:
				// ARM_THREAD_STATE64: x0..x28, fp, lr, sp, pc (u64), cpsr.
				// pc is a u64 at offset 64 bytes from the start of the state.
				if len(c.State) >= 34 {
					return uint64(c.State[64]) | uint64(c.State[65])<<32, true
				}
			case CPU_TYPE_X86_64:
				// x86_THREAD_STATE64: rax..r15 (16 x u64 = 128 bytes), then rip -
				// a u64 at offset 128 (<mach/i386/thread_status.h>).
				if len(c.State) >= 34 {
					return uint64(c.State[128>>2]) | uint64(c.State[132>>2])<<32, true
				}
			}
		}
	}

	return 0, false
}

// CodeSignatureBlob is an element of the code signature superblob index.
type CodeSignatureBlob struct {
	Type   uint32 // CSMACHO types: 0=code directory, 2=requirements, ...
	Offset uint32 // from the start of the superblob
}

func NewCodeSignatureBlob(type_ uint32, offset uint32) CodeSignatureBlob {
	return CodeSignatureBlob{
		Type:   type_,
		Offset: offset,
	}
}

// CodeSignature parses the superblob envelope of LC_CODE_SIGNATURE
// (magic 0xfade0cc0, a count, and an index array). Deep parsing of the
// signature itself (hashes, CMS) is out of scope for the package.
func (f *File) CodeSignature() ([]CodeSignatureBlob, error) {
	data, err := f.linkeditBlob(LC_CODE_SIGNATURE)
	if err != nil || data == nil {
		return nil, err
	}

	// The code signature format is big-endian, regardless of the file's byte order.
	const csMagicSuperblob = 0xfade0cc0
	if len(data) < 12 {
		return nil, errf("code signature: truncated (%d bytes)", len(data))
	}

	if magic := binary.BigEndian.Uint32(data); magic != csMagicSuperblob {
		return nil, errf("code signature: magic %#x, want superblob %#x", magic, csMagicSuperblob)
	}

	count := binary.BigEndian.Uint32(data[4:])
	var out []CodeSignatureBlob
	for i := uint32(0); i < count && 12+int(i)*8+8 <= len(data); i++ {
		out = append(
			out,
			NewCodeSignatureBlob(
				binary.BigEndian.Uint32(data[12+i*8:]),
				binary.BigEndian.Uint32(data[12+i*8+4:]),
			),
		)
	}

	return out, nil
}

// linkeditBlob reads the linkedit_data_command blob by command type.
func (f *File) linkeditBlob(cmd Cmd) ([]byte, error) {
	for _, lc := range f.commands {
		if ld, ok := lc.(*LinkeditData); ok && ld.Cmd() == cmd {
			if ld.DataSize == 0 {
				return nil, nil
			}

			return f.dyldStream(ld.DataOff, ld.DataSize)
		}
	}

	return nil, nil
}
