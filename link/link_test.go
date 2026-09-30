package link

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/unit"
)

// flatPlace - the streams at fixed bases: the addresses under test are
// layout facts, not placement ones.
func flatPlace(textBase, dataBase uint64) unit.Place {
	return func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return textBase, dataBase
	}
}

// deps - the arm64 assembler as the driver's test injection.
func deps() Deps {
	return Deps{Parse: alias.ParseSourceUnit}
}

// linkCase - one table row of the driver: the sources, the entry, and
// the expected resolve - a subset of the symbol table (the global
// interface), the memory size of the data stream, or an error.
type linkCase struct {
	name     string
	sources  []Source
	entry    string
	wantSyms map[string]uint64
	wantMem  int
	wantErr  string
}

func TestLink(t *testing.T) {
	cases := []linkCase{
		{
			name: "sections merge per stream, bss aggregates last",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
a:
  movz x0, #1
.data
da:
  .word 0xAA
.bss
ba:
  .zero 8
`,
				},
				{
					File: "b.S",
					Src: `
.text
b:
  movz x1, #2
.data
db:
  .word 0xBB
.bss
bb:
  .zero 4
`,
				},
			},
			wantSyms: map[string]uint64{
				"a": 0x1000, "b": 0x1004,
				"da": 0x8000, "db": 0x8004,
				"ba": 0x8008, "bb": 0x8010,
			},
			wantMem: 20,
		},
		{
			name: "cross-file reference through the shared namespace",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
fa:
  movz x0, #1
  b fb
`,
				},
				{
					File: "b.S",
					Src: `
.text
fb:
  ret
`,
				},
			},
			wantSyms: map[string]uint64{"fa": 0x1000, "fb": 0x1008},
		},
		{
			name: "same .L locals of different sources stay isolated",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
.Lx:
  movz x0, #1
  b .Lx
`,
				},
				{
					File: "b.S",
					Src: `
.text
.Lx:
  movz x1, #2
  b .Lx
`,
				},
			},
			wantSyms: map[string]uint64{
				"a.S:.Lx": 0x1000,
				"b.S:.Lx": 0x1008,
			},
		},
		{
			name: ".global promotes a .L name into the shared namespace",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.global .Lshare
.text
.Lshare:
  ret
`,
				},
				{
					File: "b.S",
					Src: `
.global .Lshare
.text
entry:
  bl .Lshare
`,
				},
			},
			wantSyms: map[string]uint64{".Lshare": 0x1000, "entry": 0x1004},
		},
		{
			name: "a duplicate global is a link error",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
dup:
  ret
`,
				},
				{
					File: "b.S",
					Src: `
.text
dup:
  ret
`,
				},
			},
			wantErr: `label "dup" redefined`,
		},
		{
			name: "an unresolved reference is a link error",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
go:
  bl missing
`,
				},
			},
			wantErr: `undefined symbol "missing"`,
		},
		{
			name: "the entry label must exist",
			sources: []Source{
				{
					File: "a.S",
					Src: `
.text
start:
  ret
`,
				},
			},
			entry:   "gone",
			wantErr: `entry: undefined label "gone"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := Link(deps(), tc.sources, tc.entry, flatPlace(0x1000, 0x8000))

			if tc.wantErr != "" {
				require.NotEmpty(t, f.Errs)
				require.ErrorContains(t, errors.Join(f.Errs...), tc.wantErr)
				return
			}

			require.Empty(t, f.Errs)

			for name, addr := range tc.wantSyms {
				require.Equalf(t, addr, f.Syms[name], "symbol %q", name)
			}

			if tc.wantMem > 0 {
				require.Equal(t, tc.wantMem, f.DataMem)
			}
		})
	}
}

func TestLinkParseErrorsSurface(t *testing.T) {
	// a source that does not assemble: its errors join the resolve ones,
	// the healthy source still links (every error at once)
	f := Link(
		deps(),
		[]Source{
			{File: "bad.S", Src: ".text\n  frobnicate x0\n"},
			{File: "good.S", Src: ".text\nok:\n  ret\n"},
		},
		"",
		flatPlace(0x1000, 0x8000),
	)

	require.NotEmpty(t, f.Errs)
	require.ErrorContains(t, errors.Join(f.Errs...), "frobnicate")
	require.Equal(t, uint64(0x1000), f.Syms["ok"])
}

func TestEntryOf(t *testing.T) {
	f := &unit.Fixed{Syms: map[string]uint64{"start": 1}}

	name, ok := EntryOf(f, "")
	require.True(t, ok)
	require.Equal(t, "start", name)

	name, ok = EntryOf(f, "explicit")
	require.True(t, ok)
	require.Equal(t, "explicit", name)

	_, ok = EntryOf(&unit.Fixed{Syms: map[string]uint64{}}, "")
	require.False(t, ok)
}
