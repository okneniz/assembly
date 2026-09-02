package objdump

// A parsec grammar for the objdump output line:
//
//	100003f64: fd 7b bf a9  stp x29, x30, [sp, #-0x10]!   (Mach-O: bytes)
//	a: 1305 addi a0,a0,1                                   (ELF: hex word)
//	1e: 02 88                                              (16-bit instruction)
//
// The code column is exactly 2 or 4 two-digit hex bytes, or a 4/8-digit hex
// word; no hex digits may follow it (3- and 5-byte "tails" and words of
// other widths are not instructions - objdump prints those as data). The
// grammar works on the normalized line (Normalize collapses whitespace), so
// the separators inside columns are single spaces.

import (
	"bufio"
	"strings"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"
)

var (
	cHex = parsecstrings.Try(parsecstrings.OneOf("hex digit",
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F'))
	cSp    = parsecstrings.Try(parsecstrings.Space("space"))
	cColon = parsecstrings.Try(parsecstrings.Eq("':'", ':'))
	cByte  = parsecstrings.Count(2, "code byte", cHex)
)

func isHexRune(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func hexVal(r rune) uint64 {
	switch {
	case r >= '0' && r <= '9':
		return uint64(r - '0')
	case r >= 'a' && r <= 'f':
		return uint64(r-'a') + 10
	default:
		return uint64(r-'A') + 10
	}
}

// cAddr - a hex address up to ':' (>=1 digit) → uint64.
var cAddr = parsecstrings.Cast(
	parsecstrings.Some(8, "address", cHex),
	func(rs []rune) (uint64, error) {
		var v uint64
		for _, r := range rs {
			v = v<<4 | hexVal(r)
		}

		return v, nil
	},
)

// noByteAfter - a guard check for the byte code field: skip spaces and make
// sure the next token is NOT a two-digit hex byte (otherwise the code column
// continues with a third/fifth byte - that is data, not an instruction).
// Peek reads character by character; the position is restored.
func noByteAfter(
	buf parsec.Buffer[rune, parsecstrings.Position],
) parsec.Error[parsecstrings.Position] {
	start := buf.Position()
	for {
		if _, err := cSp(buf); err != nil {
			break
		}
	}

	if buf.IsEOF() {
		return nil
	}

	var tok []rune
	for !buf.IsEOF() {
		r, err := buf.Read(false)
		if err != nil || r == ' ' {
			break
		}

		tok = append(tok, r)
		if _, err := buf.Read(true); err != nil {
			break
		}
	}

	if err := buf.Seek(start); err != nil {
		return parsec.NewParseError(buf.Position(), err.Error())
	}

	if len(tok) == 2 && isHexRune(tok[0]) && isHexRune(tok[1]) {
		return parsec.NewParseError(start, "unexpected third code byte")
	}

	return nil
}

// bytesField - exactly n two-digit hex bytes separated by spaces. Each byte
// is a separate token: after two hex digits there must be a space or the end
// of the line (otherwise "ad" inside the mnemonic "add" would be swallowed
// as a byte).
func bytesField(n int) parsec.Combinator[rune, parsecstrings.Position, []string] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) ([]string, parsec.Error[parsecstrings.Position]) {
		toks := make([]string, 0, n)
		for i := range n {
			if i > 0 {
				if _, err := cSp(buf); err != nil {
					return nil, err
				}
			}

			t, err := cByte(buf)
			if err != nil {
				return nil, err
			}

			if !buf.IsEOF() {
				if r, perr := buf.Read(false); perr == nil && r != ' ' {
					return nil, parsec.NewParseError(
						buf.Position(),
						"code byte is not a separate token",
					)
				}
			}

			toks = append(toks, string(t))
		}

		if err := noByteAfter(buf); err != nil {
			return nil, err
		}

		return toks, nil
	}
}

// wordField - exactly n hex digits as a single word (4 or 8).
func wordField(n int) parsec.Combinator[rune, parsecstrings.Position, []string] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) ([]string, parsec.Error[parsecstrings.Position]) {
		rs, err := parsecstrings.Count(n, "code word", cHex)(buf)
		if err != nil {
			return nil, err
		}

		if !buf.IsEOF() {
			if r, perr := buf.Read(false); perr == nil && isHexRune(r) {
				return nil, parsec.NewParseError(
					buf.Position(),
					"unexpected hex digits after code word",
				)
			}
		}

		return []string{string(rs)}, nil
	}
}

// cCodeField - the machine-code column: bytes (2/4) or a hex word (8/4
// digits).
var cCodeField = parsecstrings.Choice("code field",
	parsecstrings.Try(bytesField(4)),
	parsecstrings.Try(bytesField(2)),
	parsecstrings.Try(wordField(8)),
	parsecstrings.Try(wordField(4)),
)

// cInstrLine - the address, ':' and the code column; the tail (mnemonic and
// operands) is not consumed by the grammar. Returns addr.
var cInstrLine = func(buf parsec.Buffer[rune, parsecstrings.Position]) (uint64, parsec.Error[parsecstrings.Position]) {
	addr, err := cAddr(buf)
	if err != nil {
		return 0, err
	}

	if _, err := cColon(buf); err != nil {
		return 0, err
	}

	if _, err := parsecstrings.SkipMany(cSp, cCodeField)(buf); err != nil {
		return 0, err
	}

	return addr, nil
}

// ParseLine recognizes a normalized line as an objdump instruction line and
// returns its address.
func ParseLine(s string) (uint64, bool) {
	addr, err := parsecstrings.ParseString(s, cInstrLine)
	if err != nil {
		return 0, false
	}

	return addr, true
}

// Normalize collapses whitespace sequences into single spaces and trims the
// ends (equivalent to the former regexp \s+ → " ").
func Normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// StripComments removes objdump annotations after ';' (symbol stub, adrp
// hints), ' <symbol>' (RISC-V branch/jump target annotations) and ' // =N'
// (llvm prints the decimal value of an immediate as a comment).
func StripComments(s string) string {
	if i := strings.IndexByte(s, ';'); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}

	if i := strings.Index(s, " <"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}

	if i := strings.Index(s, " //"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}

	return s
}

// ParseByAddr parses the objdump output into a map addr→normalized line
// (instruction lines only; section headers and symbol tables are filtered out
// by the grammar).
func ParseByAddr(output string) map[uint64]string {
	out := make(map[uint64]string)
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Buffer(make([]byte, 0, 1024), 10*1024*1024)
	for scanner.Scan() {
		raw := scanner.Text()
		if len(raw) == 0 {
			continue
		}

		norm := Normalize(raw)
		addr, ok := ParseLine(norm)
		if !ok {
			continue
		}

		out[addr] = norm
	}

	return out
}
