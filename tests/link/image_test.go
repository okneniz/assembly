package linktests

// The image law: the two placement engines must build ONE image. The
// resolve side places the streams (MachoPlaceStreams inside the link),
// the writer side lays the sections (MachoPlaceSections inside MachoOf)
// - and the reader, a neutral third party, reads the emitted image
// back. Every symbol the resolve computed must sit inside a section the
// image describes, and the sections must describe EXACTLY the streams:
// the same bases, the same sizes, nothing invented, nothing dropped.
// The bss alignment hole and the section-count desync of the specimen
// round both die here, without a single process spawn.

import (
	"os"
	"path/filepath"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/file"
	macho "github.com/okneniz/assembly/file/macho"
	"github.com/okneniz/assembly/link"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
	unitarm64 "github.com/okneniz/assembly/unit/arm64"
)

func TestLinkImageDescribesProgram(t *testing.T) {
	arb := linkarb.NewArchProgramArb(seedRnd(t), linkarb.Arm64)
	ohsnap.Check(t, 5000, arb, func(a linkarb.ArchProgram) bool {
		srcs := sources(a)

		f := link.Link(arm64Deps(), srcs, "", file.MachoPlaceStreams)
		if len(f.Errs) != 0 {
			t.Logf("resolve: %v", f.Errs)
			return false
		}

		img, err := unitarm64.MachoOf(f, "_start")
		if err != nil {
			t.Logf("image: %v", err)
			return false
		}

		view, err := openImage(t, img.Bytes())
		if err != nil {
			t.Fatalf("read back: %v", err)
		}

		text, terr := f.EncodeText()
		data, derr := f.EncodeData()
		if terr != nil || derr != nil {
			return false
		}

		// the streams as the image must describe them (a zero-size one
		// is absent, as MachoOf builds it)
		want := map[string]uint64{
			"__text": uint64(len(text)),
			"__data": uint64(len(data)),
			"__bss":  uint64(f.DataMem - len(data)),
		}

		sections := view.Sections()
		seen := map[string]bool{}

		for _, s := range sections {
			size, known := want[s.SectName]
			if !known {
				t.Logf("section %s: invented", s.SectName)
				return false
			}

			if size != s.Size {
				t.Logf("section %s: size %d, want %d", s.SectName, s.Size, size)
				return false
			}

			seen[s.SectName] = true
		}

		for name, size := range want {
			if size == 0 && seen[name] {
				t.Logf("section %s: present but empty", name)
				return false
			}
		}

		// the bases: the writer's placement must be the resolve's one -
		// the text stream at its base, the data stream TILED by __data
		// and __bss from its base, contiguously, up to DataMem
		textAddr, dataAddr := file.MachoPlaceStreams(len(text), len(data), f.DataMem)
		for _, s := range sections {
			base := textAddr
			if s.SectName == "__data" {
				base = dataAddr
			}

			if s.SectName == "__bss" {
				base = dataAddr + uint64(len(data))
			}

			if s.Addr != base {
				t.Logf("section %s: addr %#x, want %#x", s.SectName, s.Addr, base)
				return false
			}
		}

		// every symbol of the resolve sits inside a section of the
		// image (containment, or exactly its end - the _end shape)
		for name, addr := range f.Syms {
			if !sectionHolds(sections, addr) {
				t.Logf("symbol %s at %#x sits in no section of the image", name, addr)
				return false
			}
		}

		return true
	})
}

// openImage reads an emitted image back through the neutral parser.
func openImage(t *testing.T, blob []byte) (*macho.File, error) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "image")
	if err := os.WriteFile(path, blob, 0o644); err != nil {
		return nil, err
	}

	return macho.Open(path)
}

// sectionHolds is the containment rule of machoSectOf: an address lies
// inside a section, or exactly at its end.
func sectionHolds(sections []*macho.Section, addr uint64) bool {
	for _, s := range sections {
		if addr >= s.Addr && addr < s.Addr+s.Size {
			return true
		}

		if addr == s.Addr+s.Size {
			return true
		}
	}

	return false
}
