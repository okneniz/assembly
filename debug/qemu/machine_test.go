package qemu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullArgs(t *testing.T) {
	cases := []struct {
		name     string
		archArgs []string
		port     int
		icount   bool
		want     []string
	}{
		{
			name:     "plain",
			archArgs: []string{"-machine", "virt"},
			port:     1234,
			want: []string{
				"-machine", "virt",
				"-display", "none", "-monitor", "none", "-serial", "stdio",
				"-S", "-gdb", "tcp::1234",
			},
		},
		{
			name:     "custom port and icount",
			archArgs: []string{"-device", "loader,file=x"},
			port:     5555,
			icount:   true,
			want: []string{
				"-device", "loader,file=x",
				"-display", "none", "-monitor", "none", "-serial", "stdio",
				"-S", "-gdb", "tcp::5555",
				"-icount", "shift=auto",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, fullArgs(c.archArgs, c.port, c.icount))
		})
	}
}
