package dap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/debug/session"
)

func TestResolveLine(t *testing.T) {
	lines := []session.Line{
		session.NewLine("hello.s", 2, 0x1000, 4),
		session.NewLine("hello.s", 2, 0x1004, 4), // one line, two instructions
		session.NewLine("hello.s", 4, 0x1008, 4),
		session.NewLine("/abs/other.s", 6, 0x1010, 4),
	}

	cases := []struct {
		name  string
		file  string
		line  int
		addr  uint64
		found bool
	}{
		{"direct hit", "hello.s", 4, 0x1008, true},
		{"earliest instruction of the line", "hello.s", 2, 0x1000, true},
		{"base name fallback", "other.s", 6, 0x1010, true},
		{"no such line", "hello.s", 3, 0, false},
		{"no such file", "missing.s", 2, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			addr, ok := resolveLine(lines, c.file, c.line)
			require.Equal(t, c.found, ok)
			require.Equal(t, c.addr, addr)
		})
	}
}

func TestBreakpointSet(t *testing.T) {
	bps := newBreakpointSet()
	require.Empty(t, bps.takeFile("a.s"))

	bps.putFile("a.s", []uint64{0x1000})
	require.Equal(t, []uint64{0x1000}, bps.takeFile("a.s"))
	require.Empty(t, bps.takeFile("a.s"))

	bps.putLabels([]uint64{0x2000, 0x2008})
	require.Equal(t, []uint64{0x2000, 0x2008}, bps.takeLabels())
	require.Empty(t, bps.takeLabels())

	bps.putAddrs([]uint64{0x3000})
	require.Equal(t, []uint64{0x3000}, bps.takeAddrs())
	require.Empty(t, bps.takeAddrs())
}

func TestSortByAddr(t *testing.T) {
	lines := []session.Line{
		session.NewLine("a.s", 3, 0x2000, 4),
		session.NewLine("a.s", 1, 0x1000, 4),
	}

	got := sortByAddr(lines)
	require.Equal(t, uint64(0x1000), got[0].Addr)
	require.Equal(t, uint64(0x2000), got[1].Addr)
	require.Equal(t, uint64(0x2000), lines[0].Addr) // the input stays untouched
}
