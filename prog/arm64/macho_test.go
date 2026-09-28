package arm64

// Binary.MachO over a two-stream program: a label may sit exactly at
// the end of its stream (a compiler's fixup table closes on one - the
// linker's _end is the same shape); the image builds with the boundary
// label bound to its own section.

import (
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/unit"
)

func TestMachOEndOfStreamLabel(t *testing.T) {
	p := New(unit.New()).Entry("start")

	p.Label("start")
	p.Movz(X0, 42, arch.Hw0)
	p.Svc(0x80)

	p.Data()
	p.Label("value").Quad(7)
	p.Label("end") // one past the last data byte

	bin, errs := p.Build()
	require.Empty(t, errs)

	img, err := bin.MachO("start")
	require.NoError(t, err)
	require.NotEmpty(t, img.Bytes())
}
