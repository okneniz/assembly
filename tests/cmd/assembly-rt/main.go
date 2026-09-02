// Command assembly-rt - round-trip fidelity of a prebuilt binary: the .text
// section is disassembled by our decoder into a listing, the listing is
// assembled back by our driver, and so on for N iterations in a row (N > 2 is
// mandatory): the bytes of every iteration must match the original up to
// sha256, and the listing must be stable across iterations. Lines that do not
// assemble or assemble into different bytes are replaced with a
// .word/.half filler: the image bytes are preserved, and the failure lands in
// the report taxonomy.
//
// Usage:
//
//	assembly-rt -bin prog [-iterations 3] [-rebuild out.elf] [-json]
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/file"
)

// report - the outcome of a run over one binary.
type report struct {
	Bin      string       `json:"bin"`
	Format   string       `json:"format"`
	Arch     string       `json:"arch"`
	TextAddr uint64       `json:"text_addr"`
	TextSize int          `json:"text_size"`
	OrigSHA  string       `json:"orig_sha256"`
	FinalSHA string       `json:"final_sha256"`
	Rebuilt  string       `json:"rebuilt,omitempty"`
	Pass     bool         `json:"pass"`
	Iters    []iterReport `json:"iterations"`
}

func newReport(
	bin string,
	format string,
	arch string,
	textAddr uint64,
	textSize int,
	origSHA string,
) *report {
	return &report{
		Bin:      bin,
		Format:   format,
		Arch:     arch,
		TextAddr: textAddr,
		TextSize: textSize,
		OrigSHA:  origSHA,
	}
}

// iterReport - one iteration of "disasm → listing → assembly".
type iterReport struct {
	N           int            `json:"n"`
	InstrLines  int            `json:"instr_lines"`
	FillerLines int            `json:"filler_lines"`
	Reasons     map[string]int `json:"filler_reasons"`
	Samples     []string       `json:"samples,omitempty"`
	AsmErrs     int            `json:"asm_errors"`
	ErrSamples  []string       `json:"asm_error_samples,omitempty"`
	ListingSHA  string         `json:"listing_sha256"`
	TextSHA     string         `json:"text_sha256"`
	Fail        string         `json:"fail,omitempty"`
}

// iterResult - iterReport plus the reassembled bytes (not included in JSON).
type iterResult struct {
	rep  iterReport
	text []byte
}

func main() {
	binFlag := flag.String("bin", "", "source binary (ELF or Mach-O)")
	iterFlag := flag.Int("iterations", 3, "round-trip iterations; must be > 2")
	rebuildFlag := flag.String(
		"rebuild",
		"",
		"write a full rebuild via our ELF writer to this path (ELF originals only)",
	)
	jsonFlag := flag.Bool("json", false, "print machine-readable JSON report")
	flag.Usage = func() {
		fmt.Fprintln(
			os.Stderr,
			"usage: assembly-rt -bin FILE [-iterations N] [-rebuild OUT.elf] [-json]",
		)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *binFlag == "" || *iterFlag <= 2 {
		fmt.Fprintln(os.Stderr, "rt: -bin is required and -iterations must be > 2")
		os.Exit(2)
	}

	rep, err := run(*binFlag, *iterFlag, *rebuildFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rt:", err)
		os.Exit(1)
	}

	if *jsonFlag {
		// compact single-line JSON - the jsonl format for make rt-run.
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(rep); err != nil {
			fmt.Fprintln(os.Stderr, "encode:", err)
			os.Exit(1)
		}
	} else {
		if err := printReport(os.Stdout, rep); err != nil {
			fmt.Fprintln(os.Stderr, "print:", err)
			os.Exit(1)
		}
	}

	if !rep.Pass {
		os.Exit(1)
	}
}

// run performs the full run: iterations round-trip iterations over .text
// and, on success, a rebuild of the whole binary with our ELF writer.
func run(binPath string, iterations int, rebuildPath string) (*report, error) {
	ff, err := file.Detect(binPath)
	if err != nil {
		return nil, err
	}

	sec, err := ff.CodeSection()
	if err != nil {
		return nil, err
	}

	arch := ff.ArchKind()
	assemble, err := backendFor(arch)
	if err != nil {
		return nil, err
	}

	if rebuildPath != "" && ff.Name() != "ELF" {
		return nil, fmt.Errorf("rebuild is only supported for ELF originals, got %s", ff.Name())
	}

	rep := newReport(binPath, ff.Name(), archName(arch), sec.Addr, len(sec.Data), shaHex(sec.Data))

	cur := append([]byte(nil), sec.Data...)
	prevListing := ""
	for i := 1; i <= iterations; i++ {
		it := iterate(cur, sec.Addr, arch, assemble, i, prevListing)
		prevListing = it.rep.ListingSHA
		rep.Iters = append(rep.Iters, it.rep)
		if it.rep.Fail != "" {
			break // a divergence is recorded; running further is pointless
		}

		cur = it.text
	}

	rep.FinalSHA = shaHex(cur)
	rep.Pass = len(rep.Iters) == iterations
	for _, it := range rep.Iters {
		if it.Fail != "" {
			rep.Pass = false
		}
	}

	rep.Pass = rep.Pass && rep.FinalSHA == rep.OrigSHA

	if rebuildPath != "" && rep.Pass {
		if err := rebuildELF(binPath, cur, rebuildPath); err != nil {
			return rep, fmt.Errorf("rebuild: %w", err)
		}

		rep.Rebuilt = rebuildPath
	}

	return rep, nil
}

// backendFor returns the assembler backend of the architecture (aliases/
// pseudo - same as in assembly).
func backendFor(arch file.ArchKind) (func(string, uint64) (*asm.Result, []asm.AsmError), error) {
	if arch == file.ArchARM64 {
		return alias.Assemble, nil
	}

	if arch == file.ArchRISCV64 {
		return pseudo.Assemble, nil
	}

	if arch == file.ArchLOONGARCH64 {
		return lpseudo.Assemble, nil
	}

	return nil, fmt.Errorf("unsupported architecture %d", arch)
}

// archName - a human-readable architecture name for the report.
func archName(arch file.ArchKind) string {
	if arch == file.ArchARM64 {
		return "arm64"
	}

	if arch == file.ArchRISCV64 {
		return "riscv64"
	}

	if arch == file.ArchLOONGARCH64 {
		return "loong64"
	}

	return "unknown"
}

// shaHex - the sha256 hash of a buffer in hex.
func shaHex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// printReport prints the human-readable outcome; a write error goes up.
func printReport(w io.Writer, rep *report) error {
	if _, err := fmt.Fprintf(w, "bin: %s [%s %s] .text @ %#x (%d bytes)\n",
		rep.Bin, rep.Arch, rep.Format, rep.TextAddr, rep.TextSize); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "orig sha256: %s\n", rep.OrigSHA); err != nil {
		return err
	}

	for _, it := range rep.Iters {
		if _, err := fmt.Fprintf(w, "iteration %d: %d instr, %d filler %v, asm-errors %d",
			it.N, it.InstrLines, it.FillerLines, it.Reasons, it.AsmErrs); err != nil {
			return err
		}

		if it.Fail == "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, " — FAIL: %s\n", it.Fail); err != nil {
				return err
			}
		}

		for _, s := range it.Samples {
			if _, err := fmt.Fprintf(w, "  filler: %s\n", s); err != nil {
				return err
			}
		}

		for _, s := range it.ErrSamples {
			if _, err := fmt.Fprintf(w, "  asm: %s\n", s); err != nil {
				return err
			}
		}
	}

	verdict := "FAIL"
	if rep.Pass {
		verdict = "PASS"
	}

	if _, err := fmt.Fprintf(w, "gate: %s (final sha256 %s)\n", verdict, rep.FinalSHA); err != nil {
		return err
	}

	if rep.Rebuilt != "" {
		if _, err := fmt.Fprintf(w, "rebuilt: %s\n", rep.Rebuilt); err != nil {
			return err
		}
	}

	return nil
}
