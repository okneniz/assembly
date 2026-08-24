package riscv

import (
	"encoding/binary"
	mrnd "math/rand/v2"
	"reflect"
	"testing"
	"time"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/disasm"
)

// benchCtor - the fingerprint of a constructed instruction: type + ObjDump text.
// Matching fingerprints of decodeOne(w) and the entry is the criterion that
// "the word was decoded by this entry", not an earlier one.
type benchCtor struct {
	kind reflect.Type
	text string
}

func benchCtorOf(in Instr) benchCtor {
	return benchCtor{
		kind: reflect.TypeOf(in),
		text: in.ObjDump(disasm.DefaultViewCtx()),
	}
}

// benchWordFor picks a word decodable by the required table entry:
// the candidate is match (free bits = 0), then up to 32 attempts with random
// free bits (deterministic rnd). 0 - the entry is unreachable.
func benchWordFor(
	rnd *mrnd.Rand,
	match, mask uint32,
	ctor func(word uint32, addr uint64) Instr,
) uint32 {
	ref := func(w uint32) benchCtor { return benchCtorOf(ctor(w, 0)) }

	if benchCtorOf(decodeOne(match, 0)) == ref(match) {
		return match
	}

	for range 32 {
		w := match | (rnd.Uint32() &^ mask)
		if benchCtorOf(decodeOne(w, 0)) == ref(w) {
			return w
		}
	}

	return 0
}

// benchWordsAll - one word per decodeTable entry (the 32-bit branches are
// the analog of the arm64 schemas and material for the future tree).
// Compressed RVC are hand-written switches by quadrant, not a table:
// representative encodings go into the stream
// (c.nop; 0x0000 - reserved c.addi4spn -> Unknown, worst case), variable
// length 2/4 is present.
func benchWordsAll(b *testing.B) []uint32 {
	b.Helper()

	rnd := mrnd.New(mrnd.NewPCG(0x21DEC0D5, 0xB4A5E4C4))

	out := make([]uint32, 0, len(decodeTable))
	skipped := 0

	for _, e := range decodeTable {
		mm, ok := riscvEncodings[encName(e.name)]
		if !ok {
			continue // decodeOne skips such entries
		}

		if w := benchWordFor(rnd, mm[0], mm[1], e.ctor); w != 0 {
			out = append(out, w)
		} else {
			skipped++
		}
	}

	b.Logf("words: %d (entries skipped: %d)", len(out), skipped)

	return out
}

// benchData - repeats repetitions: 32-bit words (LE) + c.nop + reserved.
func benchData(words []uint32, repeats int) []byte {
	data := make([]byte, 0, repeats*(len(words)*4+4))
	for range repeats {
		for _, w := range words {
			data = binary.LittleEndian.AppendUint32(data, w)
		}

		data = binary.LittleEndian.AppendUint16(data, 0x0001) // c.nop
		data = binary.LittleEndian.AppendUint16(data, 0x0000) // reserved → Unknown
	}

	return data
}

// BenchmarkParse decodes the whole buffer per iteration (combinator assembly
// included - that is how consumers call Parse). The instr/s metric is the
// decoder throughput; the reference point is 1M instr/s.
func BenchmarkParse(b *testing.B) {
	const repeats = 8

	words := benchWordsAll(b)
	data := benchData(words, repeats)

	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()
	for range b.N {
		if _, err := Parse(0)(bytes.Buffer(data)); err != nil {
			b.Fatal(err)
		}
	}

	elapsed := time.Since(start).Seconds()

	b.ReportMetric(float64(repeats*(len(words)+2))*float64(b.N)/elapsed, "instr/s")
}
