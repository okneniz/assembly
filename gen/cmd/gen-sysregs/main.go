// Command gen-sysregs builds the ARM64 system-register name table consumed by
// the "sysreg" field transform (see arch/arm64/transforms.go).
//
// It merges two data sources:
//
//   - The official ARM System Register XML (one AArch64-*.xml file per
//     register), vendored under arch/arm64/data/sysreg/. Download:
//     https://developer.arm.com/-/media/developer/products/architecture/armv8-a-architecture/2020-06/SysReg_xml_v86A-2020-06.tar.gz
//   - Apple implementation-defined registers from m1n1's apple_regs.json
//     (MIT), vendored under arch/arm64/data/apple_regs.json. The official ARM
//     XML does not describe these; they live entirely in the IMPDEF encoding
//     space (op0=3, CRn=15).
//
// Output: arch/arm64/sysregs_generated.go, a map keyed by the 15-bit value the
// "sysreg" schema field yields (instruction bits 19:5). op0's high bit is
// pinned to 1 by the MSR/MRS schema mask, so the key packs only op0's low bit:
//
//	key = (op0&1)<<14 | op1<<11 | CRn<<7 | CRm<<3 | op2
//
// Usage:
//
//	go run ./gen/cmd/gen-sysregs -i arch/arm64/data/sysreg -apple arch/arm64/data/apple_regs.json -o arch/arm64/sysregs_generated.go
//
// The value-encoding parser is adapted from golang.org/x/arch/arm64/arm64gen
// (BSD-style license, Copyright 2019 The Go Authors), which is the proven
// reference for this XML format.
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// --- XML data structures (correspond to the ARM register_page schema) ---

type registerPage struct {
	Registers struct {
		Register struct {
			RegShortName string `xml:"reg_short_name"`
			// reg_array is the canonical array declaration for some arrayed
			// registers (e.g. PMEVCNTR<n>_EL0): <reg_array><reg_array_start>
			// 0</reg_array_start><reg_array_end>30</reg_array_end></reg_array>.
			RegArray struct {
				Start string `xml:"reg_array_start"`
				End   string `xml:"reg_array_end"`
			} `xml:"reg_array"`
			// reg_variables wraps reg_variable in most other arrayed registers.
			RegVariables struct {
				RegVariable struct {
					Max string `xml:"max,attr"`
				} `xml:"reg_variable"`
			} `xml:"reg_variables"`
			AccessMechanisms struct {
				AccessMechanism []struct {
					Accessor string `xml:"accessor,attr"`
					Encoding struct {
						Enc []enc `xml:"enc"`
					} `xml:"encoding"`
				} `xml:"access_mechanism"`
			} `xml:"access_mechanisms"`
		} `xml:"register"`
	} `xml:"registers"`
}

// enc is one <enc n="op0" v="0b11"/> element. n is one of op0/op1/CRn/CRm/op2;
// v is a binary literal, possibly parameterized by the array index n:
// "0b010:n[3]", "0b1:n[1:0]", "n[3:0]", "n[2:0]".
type enc struct {
	N string `xml:"n,attr"`
	V string `xml:"v,attr"`
}

// appleReg is one entry of m1n1's apple_regs.json: enc is [op0,op1,CRn,CRm,op2].
type appleReg struct {
	Name string `json:"name"`
	Enc  []int  `json:"enc"`
}

func main() {
	xmlDir := flag.String("i", "arch/arm64/data/sysreg", "ARM sysreg XML directory")
	applePath := flag.String(
		"apple",
		"arch/arm64/data/apple_regs.json",
		"m1n1 apple_regs.json path",
	)
	outPath := flag.String("o", "arch/arm64/sysregs_generated.go", "output Go file")
	flag.Parse()

	// The ARM XML under arch/arm64/data/sysreg is gitignored (bulky, regenerable)
	// and must be fetched before generating. Fail fast with an actionable hint
	// rather than a raw "no such file or directory" from the parser.
	requireData(*xmlDir, *applePath)

	// Ordered: ARM first, then Apple (Apple wins on collision since it occupies
	// the IMPDEF space that ARM leaves unnamed anyway).
	names := map[uint32]string{}

	parseARM(*xmlDir, names)
	parseApple(*applePath, names)
	emit(*outPath, names)
}

// requireData verifies the vendored source data is present and non-empty. On a
// fresh clone, arch/arm64/data/sysreg is gitignored and must be fetched first;
// this prints a friendly, actionable hint instead of letting the parser fail
// with an opaque read error. Exits non-zero if anything is missing.
func requireData(xmlDir, applePath string) {
	missing := false

	hasXML := false
	if entries, err := os.ReadDir(xmlDir); err == nil {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "AArch64-") && strings.HasSuffix(e.Name(), ".xml") {
				hasXML = true
				break
			}
		}
	}

	if !hasXML {
		fmt.Fprintf(os.Stderr, "\nerror: no ARM System Register XML found in %s\n", xmlDir)
		fmt.Fprintf(
			os.Stderr,
			"       That directory is gitignored (the 490 XML files are bulky and\n",
		)
		fmt.Fprintf(os.Stderr, "       regenerable) and must be downloaded before generating.\n\n")
		fmt.Fprintf(os.Stderr, "       Fix:    make update-sysreg-data\n")
		fmt.Fprintf(os.Stderr, "       Then:  make gen-sysregs\n\n")
		missing = true
	}

	if _, err := os.Stat(applePath); err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: Apple register data not found at %s\n", applePath)
		fmt.Fprintf(os.Stderr, "       Fix:    make update-sysreg-data\n\n")
		missing = true
	}

	if missing {
		os.Exit(1)
	}
}

// parseARM walks the vendored ARM XML and fills names.
func parseARM(xmlDir string, names map[uint32]string) {
	entries, err := os.ReadDir(xmlDir)
	if err != nil {
		log.Fatalf("read xml dir %s: %v", xmlDir, err)
	}

	// Sort for deterministic output regardless of filesystem walk order.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var added, skipped int
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "AArch64-") || !strings.HasSuffix(e.Name(), ".xml") {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(xmlDir, e.Name()))
		if err != nil {
			log.Printf("%s: %v", e.Name(), err)
			continue
		}

		var page registerPage
		if err := xml.Unmarshal(raw, &page); err != nil {
			log.Printf("%s: parse error: %v", e.Name(), err)
			continue
		}

		r := page.Registers.Register
		name := r.RegShortName
		if name == "" {
			continue
		}

		// Reserved/unallocated encoding placeholder.
		if strings.Contains(name, "<op1>_<Cn>_<Cm>_<op2>") {
			skipped++
			continue
		}

		ams := r.AccessMechanisms.AccessMechanism
		if len(ams) == 0 {
			skipped++
			continue
		}

		// Use the first access mechanism's encoding (the canonical one; the
		// XML lists MRS and MSR variants that share it).
		encs := ams[0].Encoding.Enc
		if !strings.Contains(ams[0].Accessor, "MRS") && !strings.Contains(ams[0].Accessor, "MSR") {
			skipped++ // system instructions (TLBI/DC/IC), not MRS/MSR registers
			continue
		}

		// Expand arrayed registers (e.g. DBGBCR<n>_EL1) into one entry per
		// instance; the enc v-values are parameterized by the instance. The
		// instance count comes from reg_array (canonical) or reg_variable max.
		if strings.Contains(name, "<n>") {
			maxStr := r.RegArray.End
			if maxStr == "" {
				maxStr = r.RegVariables.RegVariable.Max
			}

			maxN, err := strconv.Atoi(maxStr)
			if err != nil {
				log.Printf(
					"%s: arrayed register %s missing instance count: %v",
					e.Name(),
					name,
					err,
				)
				skipped++
				continue
			}

			for n := 0; n <= maxN; n++ {
				key, ok := encKey(encs, n)
				if !ok {
					continue
				}

				instance := strings.ReplaceAll(name, "<n>", strconv.Itoa(n))
				add(names, key, instance, &added)
			}

			continue
		}

		key, ok := encKey(encs, 0)
		if !ok {
			skipped++
			continue
		}

		add(names, key, name, &added)
	}

	log.Printf("ARM: %d registers added, %d files skipped", added, skipped)
}

// parseApple reads m1n1's apple_regs.json and overlays it on names.
func parseApple(path string, names map[uint32]string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read apple json %s: %v", path, err)
	}

	var regs []appleReg
	if err := json.Unmarshal(raw, &regs); err != nil {
		log.Fatalf("parse apple json: %v", err)
	}

	var added int
	for _, r := range regs {
		if len(r.Enc) != 5 {
			continue
		}

		key := sysregKey(
			uint32(r.Enc[0]),
			uint32(r.Enc[1]),
			uint32(r.Enc[2]),
			uint32(r.Enc[3]),
			uint32(r.Enc[4]),
		)
		// Apple overlay takes precedence over any ARM entry.
		names[key] = r.Name
		added++
	}

	log.Printf("Apple: %d registers added (overlay)", added)
}

// encKey resolves the five enc elements (by their n attribute) for a given
// array instance and packs them into the 15-bit sysreg field key. ok=false if
// the encoding is malformed or incomplete. The v-attribute grammar lives in
// parse_enc.go (parsec).
func encKey(encs []enc, instance int) (uint32, bool) {
	byName := map[string]string{}
	for _, e := range encs {
		byName[e.N] = e.V
	}

	values := make([]uint64, 0, 5)
	for _, name := range []string{"op0", "op1", "CRn", "CRm", "op2"} {
		v, err := parseEncValue(byName[name])
		if err != nil {
			return 0, false
		}

		values = append(values, v.resolve(instance))
	}

	return sysregKey(
		uint32(values[0]),
		uint32(values[1]),
		uint32(values[2]),
		uint32(values[3]),
		uint32(values[4]),
	), true
}

// sysregKey packs the encoding into the 15-bit value the "sysreg" field
// produces. op0's high bit is fixed by the MSR/MRS mask, so only its low bit
// is encoded here.
func sysregKey(op0, op1, crn, crm, op2 uint32) uint32 {
	return (op0&1)<<14 | op1<<11 | crn<<7 | crm<<3 | op2
}

// add inserts a name unless the key is already taken (first-wins for ARM;
// Apple calls parseApple set directly to override).
func add(names map[uint32]string, key uint32, name string, count *int) {
	if _, exists := names[key]; exists {
		return
	}

	names[key] = name
	*count++
}

// emit writes the generated Go map, sorted by key for stable output.
func emit(outPath string, names map[uint32]string) {
	keys := make([]uint32, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	var b strings.Builder
	fmt.Fprintln(&b, "// Code generated by gen/cmd/gen-sysregs; DO NOT EDIT.")
	fmt.Fprintln(&b, "//")
	fmt.Fprintln(&b, "// ARM64 system register names, sourced from the official ARM System")
	fmt.Fprintln(&b, "// Register XML (arch/arm64/data/sysreg) and Apple implementation-defined")
	fmt.Fprintln(&b, "// registers from m1n1's apple_regs.json (arch/arm64/data/apple_regs.json).")
	fmt.Fprintln(&b, "//")
	fmt.Fprintln(&b, "// Keyed by the 15-bit value the \"sysreg\" schema field yields (instruction")
	fmt.Fprintln(&b, "// bits 19:5): (op0&1)<<14 | op1<<11 | CRn<<7 | CRm<<3 | op2. op0's high bit")
	fmt.Fprintln(
		&b,
		"// is pinned to 1 by the MSR/MRS schema mask, so only its low bit is encoded.",
	)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "package arm64")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "var sysregNames = map[uint32]string{")
	for _, k := range keys {
		fmt.Fprintf(&b, "\t0x%04x: %q,\n", k, names[k])
	}

	fmt.Fprintln(&b, "}")

	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		log.Fatalf("write %s: %v", outPath, err)
	}

	log.Printf("wrote %d entries to %s", len(names), outPath)
}
