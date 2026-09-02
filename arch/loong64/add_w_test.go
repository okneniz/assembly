package loong64

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAddWCtor(t *testing.T) {
	// llvm-mc-verified: add.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001039ac),
		ctorWord(t, New().AddW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AddW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AddW)
	require.True(t, ok, "type = %T, want AddW", in)
}

func TestAddWDecodeEncode(t *testing.T) {
	in := decodeOne(0x001039ac, 0x90000000)

	addw, ok := in.(AddW)
	require.True(t, ok, "type = %T, want AddW", in)
	require.Equal(t, "add.w $t0, $t1, $t2", addw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), addw.Addr())
	require.Equal(t, 4, addw.Len())
	require.Equal(t, uint32(0x001039ac), ctorWord(t, addw))
}

func TestAddWObjDumpNoAlias(t *testing.T) {
	// The bare register numbering prints as $rN when there is no ABI name
	// (llvm-mc-verified word: add.w $r21, $fp, $s0).
	in := decodeOne(0x00105ed5, 0)
	require.Equal(t, "add.w $r21, $fp, $s0", in.ObjDump(disasm.DefaultViewCtx()))
}

func TestAddWEncodeError(t *testing.T) {
	in := New().AddW(lreg(t, 0), lreg(t, 1), lreg(t, 2))

	_, err := in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "write failed")
}

func TestAddWMarshalJSON(t *testing.T) {
	in := New().AddW(lreg(t, 12), lreg(t, 13), lreg(t, 14))

	b, err := in.MarshalJSON()
	require.NoError(t, err)

	var dto map[string]any
	require.NoError(t, json.Unmarshal(b, &dto))
	require.Equal(t, "add.w", dto["mnemonic"])
	require.Equal(t, "add.w $t0, $t1, $t2", dto["operands"])
	require.Equal(t, "LA64", dto["group"])

	fields, ok := dto["fields"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "$t0", fields["rd"])
	require.Equal(t, "$t1", fields["rj"])
	require.Equal(t, "$t2", fields["rk"])
}
