package arm64

import (
	"encoding/binary"
	"testing"

	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// sysWord assembles an MRS/MSR instruction word from a system-register encoding
// (op0,op1,CRn,CRm,op2) and a destination/source register Rt. base is 0xD5300000
// for MRS (read, L=1) and 0xD5100000 for MSR (write, L=0). This mirrors how the
// ARM encoding packs the sysreg selector into bits 19:5.
func sysWord(base, op0, op1, crn, crm, op2, rt uint32) uint32 {
	return base | op0<<19 | op1<<16 | crn<<12 | crm<<8 | op2<<5 | rt
}

// disasmOne decodes a single word and returns the instruction text
// (mnemonic + operands) via ObjDump - without the addr\tbytes prefix.
func disasmOne(t *testing.T, word uint32) string {
	t.Helper()
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], word)
	insts, err := Parse(0)(bytes.Buffer(buf[:]))
	require.NoError(t, err)
	require.Len(t, insts, 1, "expected 1 instruction")
	return insts[0].ObjDump(disasm.DefaultViewCtx())
}

func TestSysRegNameTransform(t *testing.T) {
	cases := []struct {
		key  uint32
		name string
	}{
		{
			0x4080,
			"SCTLR_EL1",
		}, // op0=3,op1=0,CRn=1,CRm=0,op2=0
		{
			0x5807,
			"DCZID_EL0",
		}, // op0=3,op1=3,CRn=0,CRm=0,op2=7
		{
			0x5a10,
			"NZCV",
		}, // op0=3,op1=3,CRn=4,CRm=2,op2=0
		{
			0x5a20,
			"FPCR",
		}, // op0=3,op1=3,CRn=4,CRm=4,op2=0
		{
			0x5a21,
			"FPSR",
		}, // op0=3,op1=3,CRn=4,CRm=4,op2=1
		{
			0x5e82,
			"TPIDR_EL0",
		}, // op0=3,op1=3,CRn=13,CRm=0,op2=2
		{
			0x5e83,
			"TPIDRRO_EL0",
		}, // op0=3,op1=3,CRn=13,CRm=0,op2=3
		{
			0x4780,
			"HID0_EL1",
		}, // Apple IMPDEF: op0=3,op1=0,CRn=15,CRm=0,op2=0
	}
	for _, c := range cases {
		got := sysRegName(c.key)
		require.Equal(t, c.name, got, "sysRegName(%#x)", c.key)
	}
}

func TestSysRegNameFallback(t *testing.T) {
	// An encoding absent from the table: op0=2,op1=7,CRn=15,CRm=15,op2=7.
	key := uint32(0x3fff)
	require.NotContains(
		t,
		sysregNames,
		key,
		"test key %#x unexpectedly named; pick another absent encoding",
		key,
	)
	got := sysRegName(key)
	want := "S2_7_C15_C15_7"
	require.Equal(t, want, got, "fallback")
}

func TestSysRegRoundTrip(t *testing.T) {
	cases := []struct {
		word uint32
		want string
	}{
		// ARM architectural registers, read (MRS) and write (MSR).
		{
			sysWord(0xD5300000, 3, 3, 0, 0, 7, 0),
			"mrs x0, DCZID_EL0",
		},
		{
			sysWord(0xD5300000, 3, 3, 4, 2, 0, 1),
			"mrs x1, NZCV",
		},
		{
			sysWord(0xD5100000, 3, 3, 4, 2, 0, 2),
			"msr NZCV, x2",
		},
		{
			sysWord(0xD5100000, 3, 0, 1, 0, 0, 3),
			"msr SCTLR_EL1, x3",
		},
		// Apple M1 implementation-defined register.
		{
			sysWord(0xD5300000, 3, 0, 15, 0, 0, 4),
			"mrs x4, HID0_EL1",
		},
		// Unknown encoding falls through to the S<op0>_<op1>_C<CRn>_C<CRm>_<op2> form.
		{
			sysWord(0xD5300000, 2, 7, 15, 15, 7, 5),
			"mrs x5, S2_7_C15_C15_7",
		},
	}
	for _, c := range cases {
		got := disasmOne(t, c.word)
		require.Equal(t, c.want, got, "word %#08x", c.word)
	}
}
