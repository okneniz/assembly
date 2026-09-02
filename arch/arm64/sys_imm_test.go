package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSvcBrkCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"svc #0x80",
			New().Svc(imm16(t, 0x80)),
			0xd4001001,
		},
		{
			"brk #1",
			New().Brk(imm16(t, 1)),
			0xd4200020,
		},
		{
			"brk #0",
			New().Brk(imm16(t, 0)),
			0xd4200000,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	svc := New().Svc(imm16(t, 1))
	_, ok := svc.(sysImm)
	require.True(t, ok, "type = %T, want sysImm", svc)
	brk := New().Brk(imm16(t, 1))
	_, ok = brk.(sysImm)
	require.True(t, ok, "type = %T, want sysImm", brk)
}
