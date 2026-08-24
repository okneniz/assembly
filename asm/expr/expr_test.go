package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpressions(t *testing.T) {
	cases := []struct {
		src  string
		want int64
	}{
		{
			"42",
			42,
		},
		{
			"0x2a",
			0x2a,
		},
		{
			"0X2A",
			0x2A,
		},
		{
			"0b101",
			5,
		},
		{
			"010",
			8,
		}, // octal
		{
			"0",
			0,
		},
		{
			"1+2*3",
			7,
		},
		{
			"(1+2)*3",
			9,
		},
		{
			"1<<4",
			16,
		},
		{
			"0xff&0x0f",
			0x0f,
		},
		{
			"6/2+1",
			4,
		},
		{
			"7%3",
			1,
		},
		{
			"-5+2",
			-3,
		},
		{
			"~0",
			-1,
		},
		{
			"1|2^3&4",
			1 | (2 ^ (3 & 4)),
		},
		{
			"'a'",
			97,
		},
		{
			"'\\n'",
			10,
		},
	}
	for _, c := range cases {
		e, err := ParseExpr(c.src)
		require.NoError(t, err, "ParseExpr(%q)", c.src)
		v, err := e.Eval(nil)
		require.NoError(t, err, "Eval(%q)", c.src)
		require.Equal(t, c.want, v, "Eval(%q)", c.src)
	}
}

func TestBadExpressions(t *testing.T) {
	for _, src := range []string{"", "1+", "(1", "0x", "abc(", "1 2"} {
		_, err := ParseExpr(src)
		require.Error(t, err, "ParseExpr(%q) should fail", src)
	}
}
