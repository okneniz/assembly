package file

import "fmt"

// unsupportedFormatError is returned by Detect when none of the detectors
// (ELF, Mach-O) recognized the file format.
type unsupportedFormatError struct {
	magic [4]byte
}

func newUnsupportedFormatError(magic [4]byte) unsupportedFormatError {
	return unsupportedFormatError{magic: magic}
}

func (e unsupportedFormatError) Error() string {
	return fmt.Sprintf("assembly/format: unsupported file format (magic: % x)", e.magic[:])
}

func errUnsupported(magic [4]byte) error {
	return newUnsupportedFormatError(magic)
}
