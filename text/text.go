// Package text formats the machine-code column of instructions
// (Mach-O and ELF styles).
package text

import (
	"encoding/binary"
	"fmt"
)

// CodeStyle controls how the machine-code column of an instruction is displayed.
// This is a format-dependent (Mach-O vs ELF), not arch-dependent, setting,
// which is why it lives here and not in the arch packages.
type CodeStyle int

const (
	// CodeBytes is bytes separated by spaces ("fd 23 00 d1"), the Mach-O/ARM64 style.
	CodeBytes CodeStyle = iota
	// CodeWord is a 32-bit word as a single hex value ("ff010113"), the ELF style.
	CodeWord
)

// FormatCode returns the machine-code column in the given style. n is the number
// of instruction bytes (2 for compressed RISC-V, 4 for 32-bit ones; 0 => 4, a fixed width).
func FormatCode(raw uint32, n int, style CodeStyle) string {
	if style == CodeWord {
		if n == 2 {
			return fmt.Sprintf("%04x", raw&0xffff)
		}

		return fmt.Sprintf("%08x", raw)
	}

	return formatBytes(raw, n)
}

// StyleFor picks the code-column style by the file-format name (file.FileFormat.
// Name): ELF shows an instruction as a single hex word, Mach-O as
// space-separated bytes, the way objdump does.
func StyleFor(format string) CodeStyle {
	if format == "ELF" {
		return CodeWord
	}

	return CodeBytes
}

// formatBytes returns the machine code as n little-endian bytes separated by
// spaces ("fd 23 00 d1" for 4 bytes, "02 88" for 2 bytes). n=0 means the
// default width (4).
func formatBytes(raw uint32, n int) string {
	if n == 0 {
		n = 4
	}

	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], raw)
	switch n {
	case 2:
		return fmt.Sprintf("%02x %02x", b[0], b[1])
	case 4:
		return fmt.Sprintf("%02x %02x %02x %02x", b[0], b[1], b[2], b[3])
	default:
		out := make([]byte, 0, n*3)
		for i := 0; i < n && i < 4; i++ {
			if i > 0 {
				out = append(out, ' ')
			}

			out = append(out, hexByte(b[i])...)
		}

		return string(out)
	}
}

// hexByte returns the two hex digits of a byte without involving fmt.
func hexByte(b byte) string {
	const digits = "0123456789abcdef"
	return string([]byte{digits[b>>4], digits[b&0xf]})
}
