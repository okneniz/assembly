package loong64

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAllMnemonicsCountError - every shape rejects a wrong operand count.
func TestAllMnemonicsCountError(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }

	for _, m := range Mnemonics() {
		_, err := BuildInstr(m, []Op{r(1), r(2), r(3), r(4), r(5)})
		require.ErrorContains(t, err, "operands", m)
	}
}

// TestShapeKindErrors - the wrong-operand-kind branch of every factory,
// via one representative mnemonic per factory.
func TestShapeKindErrors(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }
	n := OpNum

	for _, tc := range []struct {
		mnem  string
		ops   []Op
		wantE string
	}{
		// regs3: a number in the third register slot.
		{"add.w", []Op{r(1), r(2), n(3)}, "want a register"},
		// regs2: a number in the second slot.
		{"clz.d", []Op{r(1), n(2)}, "want a register"},
		// regI20: a register in the value slot.
		{"lu12i.w", []Op{r(1), r(2)}, "want an immediate"},
		// code15: a register instead of the code.
		{"break", []Op{r(1)}, "want an immediate"},
		// target: a register instead of the target.
		{"b", []Op{r(1)}, "want an immediate"},
		// regTarget: a register as the target.
		{"beqz", []Op{r(1), r(2)}, "want an immediate"},
		// regs2Target: a register in the third slot.
		{"beq", []Op{r(1), r(2), r(3)}, "want an immediate"},
		// regs2Off (jirl): a register as the offset.
		{"jirl", []Op{r(1), r(2), r(3)}, "want an immediate"},
		// regregRole (addi.w): a register as the immediate.
		{"addi.w", []Op{r(1), r(2), r(3)}, "want an immediate"},
		// regs3Role (bytepick.w): a register as the index.
		{"bytepick.w", []Op{r(1), r(2), r(3), r(4)}, "want an immediate"},
		// regRegU5U5: a register as lsb.
		{"bstrins.w", []Op{r(1), r(2), n(5), r(3)}, "want an immediate"},
		// u5RegI12 (preld): a register as the hint.
		{"preld", []Op{r(1), r(2), n(8)}, "want an immediate"},
		// u5Regs2 (preldx): a register as the hint.
		{"preldx", []Op{r(1), r(2), r(3)}, "want an immediate"},
		// regU8 (ldpte): a register as the level.
		{"ldpte", []Op{r(1), r(2)}, "want an immediate"},
		// regU14 (csrrd): a register as the csr.
		{"csrrd", []Op{r(1), r(2)}, "want an immediate"},
	} {
		_, err := BuildInstr(tc.mnem, tc.ops)
		require.ErrorContains(t, err, tc.wantE, tc.mnem)
	}
}

// TestAPIWrappers - MarshalDTO and WriteWord (the pseudo layer's pair
// primitives).
func TestAPIWrappers(t *testing.T) {
	in := newUnknown(newBase(0x10, 0x1c000000))

	b, err := MarshalDTO(in.base, ".word", "<unknown>", "", nil)
	require.NoError(t, err)

	var dto map[string]any
	require.NoError(t, json.Unmarshal(b, &dto))
	require.Equal(t, ".word", dto["mnemonic"])

	var buf bytes.Buffer
	written, err := WriteWord(&buf, 0x001039ac)
	require.NoError(t, err)
	require.Equal(t, int64(4), written)
	require.Equal(t, []byte{0xac, 0x39, 0x10, 0x00}, buf.Bytes())
}

// TestNarrowRoles - the bytepick role bounds (2 and 3 bits).
func TestNarrowRoles(t *testing.T) {
	v, err := New().UImm2(3)
	require.NoError(t, err)
	require.Equal(t, int64(3), v.Val())

	_, err = New().UImm2(4)
	require.ErrorContains(t, err, "outside")

	v3, err := New().UImm3(7)
	require.NoError(t, err)
	require.Equal(t, int64(7), v3.Val())

	_, err = New().UImm3(8)
	require.ErrorContains(t, err, "outside")
}

// TestShapeFirstSlotErrors - a wrong operand in the earlier slots (each
// factory's remaining branches).
func TestShapeFirstSlotErrors(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }
	n := OpNum

	for _, tc := range []struct {
		mnem string
		ops  []Op
	}{
		{"add.w", []Op{n(1), r(2), r(3)}},            // regs3 slot 0
		{"add.w", []Op{r(1), n(2), r(3)}},            // regs3 slot 1
		{"clz.d", []Op{n(1), r(2)}},                  // regs2 slot 0
		{"lu12i.w", []Op{n(1), n(2)}},                // regI20 slot 0
		{"beqz", []Op{n(1), n(2)}},                   // regTarget slot 0
		{"beq", []Op{n(1), r(2), n(3)}},              // regs2Target slot 0
		{"beq", []Op{r(1), n(2), n(3)}},              // regs2Target slot 1
		{"jirl", []Op{n(1), r(2), n(3)}},             // regs2Off slot 0
		{"jirl", []Op{r(1), n(2), n(3)}},             // regs2Off slot 1
		{"addi.w", []Op{n(1), r(2), n(3)}},           // regregRole slot 0
		{"addi.w", []Op{r(1), n(2), n(3)}},           // regregRole slot 1
		{"bytepick.w", []Op{n(1), r(2), r(3), n(4)}}, // regs3Role slot 0
		{"bytepick.w", []Op{r(1), n(2), r(3), n(4)}}, // regs3Role slot 1
		{"bytepick.w", []Op{r(1), r(2), n(3), n(4)}}, // regs3Role slot 2
		{"bstrins.w", []Op{n(1), r(2), n(5), n(3)}},  // regRegU5U5 slot 0
		{"bstrins.w", []Op{r(1), n(2), n(5), n(3)}},  // regRegU5U5 slot 1
		{"bstrins.w", []Op{r(1), r(2), n(1), n(3)}},  // regRegU5U5 msb<lsb
		{"bstrins.d", []Op{n(1), r(2), n(63), n(0)}}, // regRegU6U6 slot 0
		{"bstrins.d", []Op{r(1), n(2), n(63), n(0)}}, // regRegU6U6 slot 1
		{"bstrins.d", []Op{r(1), r(2), n(3), n(63)}}, // regRegU6U6 msb<lsb
		{"preld", []Op{n(5), n(13), n(8)}},           // u5RegI12 reg slot
		{"preldx", []Op{n(5), n(13), r(14)}},         // u5Regs2 rj slot
		{"ldpte", []Op{n(1), n(1)}},                  // regU8 slot 0
		{"csrrd", []Op{n(1), n(5)}},                  // regU14 slot 0
	} {
		_, err := BuildInstr(tc.mnem, tc.ops)
		require.Error(t, err, tc.mnem)
	}
}

// TestShapeRoleErrors - role-range failures inside multi-role shapes.
func TestShapeRoleErrors(t *testing.T) {
	r := func(n int) Op { return OpReg(uint8(n)) }
	n := OpNum

	for _, tc := range []struct {
		mnem string
		ops  []Op
	}{
		{"bstrins.w", []Op{r(1), r(2), n(32), n(3)}},  // msb out of range
		{"bstrins.w", []Op{r(1), r(2), n(5), n(32)}},  // lsb out of range
		{"bstrins.w", []Op{r(1), r(2), r(5), n(3)}},   // msb is a register
		{"bstrins.d", []Op{r(1), r(2), n(64), n(0)}},  // msb out of range
		{"bstrins.d", []Op{r(1), r(2), n(63), n(64)}}, // lsb out of range
		{"bstrins.d", []Op{r(1), r(2), n(63), r(0)}},  // lsb is a register
		{"preld", []Op{n(32), r(13), n(8)}},           // hint out of range
		{"preldx", []Op{n(32), r(13), r(14)}},         // hint out of range
		{"preld", []Op{n(5), r(13), n(2048)}},         // si12 out of range
		{"preldx", []Op{n(5), r(13), n(14)}},          // rk is a number
	} {
		_, err := BuildInstr(tc.mnem, tc.ops)
		require.Error(t, err, tc.mnem)
	}
}
