// Package file is a format-neutral dispatcher for binary file formats:
// ArchKind, a neutral Section, the FileFormat
// interface, and Detect. Parsing itself belongs entirely to the self-contained
// subpackages file/elf and file/macho: the formats
// share not a single line of code, and this package merely wraps the result in a
// minimal contract for cmd/server/tests.
package file

import (
	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"

	"github.com/okneniz/assembly/file/elf"
	"github.com/okneniz/assembly/file/macho"
)

// ArchKind identifies the CPU architecture extracted from binary metadata
// (ELF e_machine / Mach-O cputype). Consumers switch on it to pick a decoder.
type ArchKind int

const (
	ArchUnknown ArchKind = iota
	ArchARM64
	ArchRISCV64
	ArchLOONGARCH64
)

// FileFormat — a binary file format (ELF, Mach-O), implemented in this package.
type FileFormat interface {
	Name() string
	ArchKind() ArchKind
	Sections() ([]Section, error)
	Section(name string) (*Section, error)
	CodeSection() (*Section, error)
}

// Detect determines the file format from its first 4 bytes (read by a parsec
// combinator) and returns a ready-to-use FileFormat, or an error.
func Detect(path string) (FileFormat, error) {
	magic, err := readMagic(path)
	if err != nil {
		return nil, err
	}

	if elf.IsMagic(magic) {
		f, err := elf.Open(path)
		if err != nil {
			return nil, err
		}

		return elfFormat{f: f}, nil
	}

	if macho.IsMagic(magic) {
		f, err := macho.Open(path)
		if err != nil {
			return nil, err
		}

		return machoFormat{f: f}, nil
	}

	return nil, errUnsupported(magic)
}

// readMagic reads the first 4 bytes of a file with a combinator (one buffer per file).
func readMagic(path string) ([4]byte, error) {
	buf, err := bytes.BufferFromFile(path)
	if err != nil {
		return [4]byte{}, err
	}

	raw, err := common.Count(4, "file: expected 4 magic bytes", bytes.Any())(buf)
	if err != nil {
		return [4]byte{}, err
	}

	var magic [4]byte
	copy(magic[:], raw)
	return magic, nil
}
