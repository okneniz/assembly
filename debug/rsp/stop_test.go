package rsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStopReply(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		want    StopReply
		wantErr bool
	}{
		{
			name:    "signal only",
			payload: "S05",
			want:    NewStopReply('S', 5, map[string]string{}),
		},
		{
			name:    "stopped with fields",
			payload: "T05thread:p1.1;core:0;",
			want: NewStopReply('T', 5, map[string]string{
				"thread": "p1.1",
				"core":   "0",
			}),
		},
		{
			name:    "watchpoint report",
			payload: "T05watch:8021a4f0;",
			want: NewStopReply('T', 5, map[string]string{
				"watch": "8021a4f0",
			}),
		},
		{
			name:    "process exit",
			payload: "W2a",
			want:    NewStopReply('W', 42, map[string]string{}),
		},
		{
			name:    "terminated",
			payload: "X06",
			want:    NewStopReply('X', 6, map[string]string{}),
		},
		{name: "too short", payload: "T0", wantErr: true},
		{name: "bad hex", payload: "Tzz", wantErr: true},
		{name: "unknown kind", payload: "Q05", wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseStopReply(c.payload)
			if c.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.want, got)
		})
	}
}

func TestStopReplyExited(t *testing.T) {
	require.True(t, NewStopReply('W', 0, nil).Exited())
	require.True(t, NewStopReply('X', 6, nil).Exited())
	require.False(t, NewStopReply('T', 5, nil).Exited())
	require.False(t, NewStopReply('S', 5, nil).Exited())
}
