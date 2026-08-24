package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrrImm — orr rd, rn, #bitmask.
type OrrImm struct {
	base
	logImm
}

const (
	orrImmX uint32 = 0xB2000000
	orrImmW uint32 = 0x32000000
)

func decodeOrrImm(w uint32, addr uint64) Instr {
	return OrrImm{
		newBase(addr, w),
		decodeLogImm(w),
	}
}

// movzRep/movnRep - whether the pattern is representable by a single
// MOVZ/MOVN (set bits / in the MOVN case zero bits fit into a single
// 16-bit aligned window).
func movzRep(v uint64, bits int) bool {
	for s := uint(0); s < uint(bits); s += 16 {
		if v>>s&0xffff != 0 && v&^(0xffff<<s) == 0 {
			return true
		}
	}

	return false
}

func movnRep(v uint64, bits int) bool {
	all := ^uint64(0)
	if bits == 32 {
		all = 0xffffffff
	}

	return movzRep(^v&all, bits)
}

// ObjDump - orr rd, rn, #bitmask. With Rn = ZR LLVM objdump prints the mov
// alias, but ONLY if the immediate is not representable by a single
// MOVZ/MOVN (otherwise the text is ambiguous on reassembly:
// mov x, #0x20 canonically = the MOVZ encoding).
func (i OrrImm) ObjDump(_ disasm.ViewCtx) string {
	bits := 32
	if i.is64 {
		bits = 64
	}

	if (i.rn == "xzr" || i.rn == "wzr") &&
		!movzRep(i.mask(), bits) && !movnRep(i.mask(), bits) {
		return fmt.Sprintf("mov %s, %s", i.rd, i.immText())
	}

	return fmt.Sprintf("orr %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i OrrImm) Encode(w io.Writer, pc uint64) (int64, error) {
	match := orrImmX
	if !i.is64 {
		match = orrImmW
	}

	if i.n {
		match |= 1 << 22
	}

	rd, rn, err := i.bits()
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|i.imms<<10|i.immr<<16)
}

func (i OrrImm) MarshalJSON() ([]byte, error) {
	return i.marshal("orr", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}

// immText - the immediate in objdump style: 64-bit patterns with the top
// bit set are printed signed (#-0x80000000 instead of #0xffffffff80000000).
func (i OrrImm) immText() string {
	m := i.mask()
	if i.is64 && m>>63 == 1 {
		return fmt.Sprintf("#-0x%x", ^m+1)
	}

	return fmt.Sprintf("#0x%x", m)
}
