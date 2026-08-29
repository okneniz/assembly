package opcodes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTable(t *testing.T) {
	table := "" +
		"00108000 add.d                  DJK             @qemu\n" +
		"03000000 cu52i.d                DJSk12          @orig_name=lu52i.d @qemu\n" +
		"24000000 ldox4.w                DJSk14          @orig_name=ldptr.w @orig_fmt=DJSk14ps2\n" +
		"06482000 tlbclr                 EMPTY\n" +
		"\n" +
		"   \n" +
		"002a0000 break                  Ud15            @la32 @primary"

	entries, err := Parse([]rune(table))
	require.NoError(t, err)
	require.Len(t, entries, 5)

	require.Equal(t, NewEntry(0x00108000, "add.d", "DJK", "", "", []string{"qemu"}), entries[0])
	require.Equal(
		t,
		NewEntry(0x03000000, "cu52i.d", "DJSk12", "lu52i.d", "", []string{"qemu"}),
		entries[1],
	)
	require.Equal(
		t,
		NewEntry(0x24000000, "ldox4.w", "DJSk14", "ldptr.w", "DJSk14ps2", nil),
		entries[2],
	)
	require.Equal(t, NewEntry(0x06482000, "tlbclr", "EMPTY", "", "", nil), entries[3])
	require.Equal(
		t,
		NewEntry(0x002a0000, "break", "Ud15", "", "", []string{"la32", "primary"}),
		entries[4],
	)
}

func TestParseOfficialName(t *testing.T) {
	entries, err := Parse([]rune("03000000 cu52i.d   DJSk12  @orig_name=lu52i.d\n"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "lu52i.d", entries[0].OfficialName())

	entries, err = Parse([]rune("00108000 add.d   DJK  @qemu\n"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "add.d", entries[0].OfficialName())
}

func TestParseCorruptLine(t *testing.T) {
	// A non-hex word is not a line and not a blank line.
	_, err := Parse([]rune("zzzz add.d DJK\n"))
	require.ErrorContains(t, err, "corrupt line")

	// A truncated word either (7 digits, then whitespace).
	_, err = Parse([]rune("0010800 add.d DJK\n"))
	require.ErrorContains(t, err, "corrupt line")

	// A valid line followed by a corrupt one reports the second position.
	_, err = Parse([]rune("00108000 add.d DJK @qemu\n???\n"))
	require.ErrorContains(t, err, "corrupt line")
}

func TestParseLineFieldErrors(t *testing.T) {
	for _, tc := range []struct{ name, line string }{
		{"no space after the word", "00108000"},
		{"bad mnemonic", "00108000 +$$"},
		{"no space after the mnemonic", "00108000 add.d@"},
		{"bad format", "00108000 add.d -"},
		{"trailing garbage", "00108000 add.d DJK $$$"},
	} {
		_, err := Parse([]rune(tc.line + "\n"))
		require.ErrorContains(t, err, "corrupt line", tc.name)
	}
}

func TestParseNoTrailingNewline(t *testing.T) {
	entries, err := Parse([]rune("00108000 add.d DJK"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "add.d", entries[0].Name)
}

func TestCastWord(t *testing.T) {
	v, err := castWord([]rune("00108000"))
	require.NoError(t, err)
	require.Equal(t, uint32(0x00108000), v)

	_, err = castWord([]rune("108000"))
	require.ErrorContains(t, err, "want 8 digits, got 6")

	// Not reachable through the grammar (the digits are pre-validated);
	// the guard stays for direct callers.
	_, err = castWord([]rune("0000000g"))
	require.ErrorContains(t, err, "hex word")
}

func TestCutPrefix(t *testing.T) {
	v, ok := cutPrefix("orig_name=lu52i.d", "orig_name=")
	require.True(t, ok)
	require.Equal(t, "lu52i.d", v)

	_, ok = cutPrefix("qemu", "orig_name=")
	require.False(t, ok)
}
