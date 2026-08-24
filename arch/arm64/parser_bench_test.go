package arm64

import (
	"encoding/binary"
	mrnd "math/rand/v2"
	"reflect"
	"testing"
	"time"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/disasm"
)

// benchCtor - fingerprint of a construct: type + ObjDump text. Matching
// fingerprints of decodeOne(w) and an entry is the criterion that "the word
// was decoded by this entry", not an earlier one.
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

// benchWordFor picks a word decoded by the required entry: the candidate is
// match (free bits = 0), then up to 32 attempts with random free bits
// (deterministic rnd). 0 - the entry is unreachable.
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

// benchWordsAll - one word per reachable decoder entry: the curated schemas
// (with a ctor), the generated isaTail tail (best effort - part of the tail
// is intercepted by the schemas and unreachable by design), a garbage word
// (Unknown - a full pass, worst case). Full branch coverage - so that when
// the structure changes (list -> tree) the benchmark runs on the same
// material.
func benchWordsAll(b *testing.B) []uint32 {
	b.Helper()

	rnd := mrnd.New(mrnd.NewPCG(0xA64DEC0D, 0xB4A5E4C4))

	out := make([]uint32, 0, len(getSchemas())+16)
	skipped := 0

	for i, s := range getSchemas() {
		match, mask := s.Value, s.Mask
		if e := schemaISAEntry(i); e != nil {
			match, mask = e.Match, e.Mask
		}

		if s.ctor == nil {
			// barrier: the word matched an entry without a ctor - it goes
			// to isaTail (the branch is mandatory in the material)
			out = append(out, match)
			continue
		}

		if w := benchWordFor(rnd, match, mask, s.ctor); w != 0 {
			out = append(out, w)
		} else {
			skipped++
		}
	}

	tail, tailMissed := 0, 0
	for k := 0; isaTailEntry(k) != nil; k++ {
		tail++

		e := isaTailEntry(k)
		ctor := func(word uint32, addr uint64) Instr { return decodeGeneric(e, word, addr) }

		if w := benchWordFor(rnd, e.Match, e.Mask, ctor); w != 0 {
			out = append(out, w)
		} else {
			tailMissed++
		}
	}

	// Unknown: a word that does not match any entry.
	unknowns := make([]uint32, 0, 9)
	unknowns = append(unknowns, 0x00000000)
	for range 8 {
		unknowns = append(unknowns, rnd.Uint32())
	}

	for _, w := range unknowns {
		if _, ok := decodeOne(w, 0).(Unknown); ok {
			out = append(out, w)
			break
		}
	}

	b.Logf("words: %d (schemas with ctor missed: %d, isaTail: %d, of which missed: %d)",
		len(out), skipped, tail, tailMissed)

	return out
}

// benchData - repeats repetitions of the set (LE).
func benchData(words []uint32, repeats int) []byte {
	data := make([]byte, 0, repeats*len(words)*4)
	for range repeats {
		for _, w := range words {
			data = binary.LittleEndian.AppendUint32(data, w)
		}
	}

	return data
}

// BenchmarkParse decodes the whole buffer per iteration (combinator
// assembly included - that is how the consumer calls Parse). The instr/s
// metric is the decoder throughput; the reference point is 1M instr/s.
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

	b.ReportMetric(float64(len(data)/4)*float64(b.N)/elapsed, "instr/s")
}
