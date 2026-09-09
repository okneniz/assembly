package loong64

import (
	"encoding/binary"
	"testing"

	parsecbytes "github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"
)

func TestParseInstrs(t *testing.T) {
	data := binary.LittleEndian.AppendUint32(nil, 0x001039ac) // add.w $t0, $t1, $t2
	data = binary.LittleEndian.AppendUint32(data, 0xffffffff) // no encoding

	instrs, err := MakeDecoder()(parsecbytes.Buffer(data))
	require.NoError(t, err)
	require.Len(t, instrs, 2)

	// Both kinds of lines carry their length.
	require.Equal(t, 4, instrs[0].Len())
	require.Equal(t, 4, instrs[1].Len())
}

func TestParseTruncatedTail(t *testing.T) {
	// A full word plus a 3-byte tail: the tail is dropped, no error.
	data := binary.LittleEndian.AppendUint32(nil, 0x1c000000)
	data = append(data, 0x11, 0x22, 0x33)

	instrs, err := MakeDecoder()(parsecbytes.Buffer(data))
	require.NoError(t, err)
	require.Len(t, instrs, 1)

	// Only the tail.
	instrs, err = MakeDecoder()(parsecbytes.Buffer([]byte{1, 2, 3}))
	require.NoError(t, err)
	require.Empty(t, instrs)

	// Nothing at all.
	instrs, err = MakeDecoder()(parsecbytes.Buffer(nil))
	require.NoError(t, err)
	require.Empty(t, instrs)
}
