package arm64

import (
	"testing"
	"time"
)

// benchSink protects decoding from being eliminated by the compiler.
var benchSink bool

// BenchmarkDecodeOne — production decodeOne (dtree decision trees) over the
// full corpus coverage; one iteration = the whole corpus once.
func BenchmarkDecodeOne(b *testing.B) {
	words := benchWordsAll(b)

	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()
	for range b.N {
		for _, w := range words {
			benchSink = decodeOne(w, 0) != nil
		}
	}

	elapsed := time.Since(start).Seconds()

	b.ReportMetric(float64(len(words))*float64(b.N)/elapsed, "instr/s")
}
