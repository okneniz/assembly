package loong64

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// errWriter - a writer that always fails (the Encode error path).
type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestUnknown(t *testing.T) {
	in := newUnknown(newBase(0x90000010, 0xdeadbeef))
	require.Equal(t, uint64(0x90000010), in.Addr())
	require.Equal(t, 4, in.Len())
	require.Equal(t, "<unknown>", in.ObjDump(disasm.DefaultViewCtx()))

	var buf bytes.Buffer
	n, err := in.Encode(&buf, 0x90000010)
	require.NoError(t, err)
	require.Equal(t, int64(4), n)
	require.Equal(t, []byte{0xef, 0xbe, 0xad, 0xde}, buf.Bytes())
}

func TestUnknownEncodeError(t *testing.T) {
	in := newUnknown(newBase(0, 0x1c000000))

	_, err := in.Encode(errWriter{}, 0)
	require.ErrorContains(t, err, "write failed")
}

func TestUnknownMarshalJSON(t *testing.T) {
	in := newUnknown(newBase(0x90000010, 0x1c000000))

	b, err := in.MarshalJSON()
	require.NoError(t, err)

	var dto map[string]any
	require.NoError(t, json.Unmarshal(b, &dto))
	require.Equal(t, ".word", dto["mnemonic"])
	require.Equal(t, "0x90000010", dto["addr"])
	require.Equal(t, "0x1c000000", dto["raw"])
}

// writeWord - the LE byte order and the error path.
func TestWriteWord(t *testing.T) {
	var buf bytes.Buffer
	n, err := writeWord(&buf, 0x00108000)
	require.NoError(t, err)
	require.Equal(t, int64(4), n)
	require.Equal(t, binary.LittleEndian.AppendUint32(nil, 0x00108000), buf.Bytes())

	_, err = writeWord(errWriter{}, 0)
	require.ErrorContains(t, err, "write failed")
}
