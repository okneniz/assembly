package elf

import (
	"encoding/binary"
	"encoding/hex"
)

// Note is an entry from a SHT_NOTE section or PT_NOTE segment. Layout (gABI,
// 4-byte field alignment): namesz(4) descsz(4) type(4) name (padded to 4)
// desc (padded to 4). The name designates the "owner" (GNU, Go, LINUX, ...)
// and takes part in interpreting Type.
type Note struct {
	Name string
	Type NoteType
	Desc []byte
}

func NewNote(name string, type_ NoteType, desc []byte) Note {
	return Note{
		Name: name,
		Type: type_,
		Desc: desc,
	}
}

// Notes returns all entries: from every PT_NOTE segment and every SHT_NOTE
// section (duplicates are possible if a note is visible in both places).
func (f *File) Notes() ([]Note, error) {
	var out []Note

	for _, p := range f.progs {
		if p.Type != PT_NOTE || p.Filesz == 0 {
			continue
		}

		data, err := readBytes(f.buf, int(p.Off), int(p.Filesz))
		if err != nil {
			return nil, err
		}

		notes, err := parseNotes(data, f.order)
		if err != nil {
			return nil, errf("PT_NOTE segment: %v", err)
		}

		out = append(out, notes...)
	}

	for _, s := range f.sections {
		if s.Type != SHT_NOTE || s.Size == 0 {
			continue
		}

		data, err := s.Data()
		if err != nil {
			return nil, err
		}

		notes, err := parseNotes(data, f.order)
		if err != nil {
			return nil, errf("section %s: %v", s.Name, err)
		}

		out = append(out, notes...)
	}

	return out, nil
}

// parseNotes parses a stream of notes with 4-byte field alignment.
func parseNotes(data []byte, order binary.ByteOrder) ([]Note, error) {
	var out []Note
	for off := 0; off+12 <= len(data); {
		namesz := u32(data[off:], order)
		descsz := u32(data[off+4:], order)
		typ := u32(data[off+8:], order)

		nameStart := off + 12
		nameEnd := nameStart + int(namesz)
		if nameEnd > len(data) {
			return nil, errf("note: name (offset %d, size %d) is out of bounds", nameStart, namesz)
		}

		descStart := align4(nameEnd)
		descEnd := descStart + int(descsz)
		if descEnd > len(data) {
			return nil, errf(
				"note: desc (offset %d, size %d) is out of bounds",
				descStart,
				descsz,
			)
		}

		name := cstr(data[nameStart:nameEnd])
		out = append(out, NewNote(name, NoteType(typ), data[descStart:descEnd:descEnd]))
		off = align4(descEnd)
	}

	return out, nil
}

// align4 rounds up to a multiple of 4.
func align4(n int) int {
	return (n + 3) &^ 3
}

// BuildID returns the build id (hex) - the contents of the GNU note
// NT_GNU_BUILD_ID; for Go binaries, if it is absent, the Go note
// (name "Go", NT_GO_BUILD_ID).
func (f *File) BuildID() (string, error) {
	notes, err := f.Notes()
	if err != nil {
		return "", err
	}

	var goID string
	for _, n := range notes {
		switch {
		case n.Name == "GNU" && n.Type == NT_GNU_BUILD_ID:
			return hex.EncodeToString(n.Desc), nil
		case n.Name == "Go" && n.Type == NT_GO_BUILD_ID:
			goID = string(n.Desc)
		}
	}

	if goID != "" {
		return goID, nil
	}

	return "", errf("build id is missing")
}
