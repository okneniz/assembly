// Command gen-arm-instr builds the ARM A64 instruction encoding table from the
// official ARM A64 ISA XML (arch/arm64/data/instr, fetched by update-instr.sh
// from the Exploration-Tools A64 ISA release on developer.arm.com).
//
// Each per-instruction XML has an <instructionsection> with <classes> containing
// one or more <iclass>, each with a <regdiagram> of <box hibit width [name]>
// elements whose <c> children give the bit values (MSB-first). Fixed boxes (with
// 0/1 c-values, whether named or not) contribute to the match/mask; named boxes
// with no c-values are variable operand fields. This mirrors how the hand-written
// arch/arm64/schemas.go describes encodings, but sourced authoritatively.
//
// Output: arch/arm64/isa_generated.go, `var armISA = []armISAEntry{...}`.
//
// Usage:
//
//	go run ./gen/cmd/gen-arm-instr -i arch/arm64/data/instr -o arch/arm64/isa_generated.go
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/okneniz/assembly/arch/arm64"
)

type instructionSection struct {
	Classes struct {
		Iclass []struct {
			Regdiagram struct {
				Box []struct {
					Hibit string   `xml:"hibit,attr"`
					Width string   `xml:"width,attr"`
					Name  string   `xml:"name,attr"`
					C     []string `xml:"c"`
				} `xml:"box"`
			} `xml:"regdiagram"`
			Encoding []struct {
				Asmtemplate struct {
					Text []string `xml:"text"`
				} `xml:"asmtemplate"`
			} `xml:"encoding"`
		} `xml:"iclass"`
	} `xml:"classes"`
}

type entry struct {
	Name   string
	Match  uint32
	Mask   uint32
	Fields []arm64.Field
}

func main() {
	in := flag.String("i", "arch/arm64/data/instr", "A64 ISA instruction XML dir")
	out := flag.String("o", "arch/arm64/isa_generated.go", "output Go file")
	flag.Parse()

	requireData(*in)

	files, err := filepath.Glob(filepath.Join(*in, "*.xml"))
	if err != nil {
		log.Fatalf("glob %s: %v", *in, err)
	}

	sort.Strings(files)

	var entries []entry
	var noBoxes int
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}

		var s instructionSection
		if err := xml.Unmarshal(data, &s); err != nil {
			continue // skip non-instruction / unparseable files
		}

		for _, ic := range s.Classes.Iclass {
			e, ok := parseBoxes(ic.Regdiagram.Box)
			if !ok {
				noBoxes++
				continue
			}

			name, merr := mnemonic(ic.Encoding)
			if merr != nil {
				log.Printf("skip %s: %v", f, merr)
				continue
			}

			e.Name = name
			if e.Name == "" {
				continue
			}

			entries = append(entries, e)
		}
	}

	log.Printf(
		"parsed %d iclasses from %d files (%d skipped without regdiagram)",
		len(entries),
		len(files),
		noBoxes,
	)
	emit(*out, entries)
}

// decode turns a regdiagram's boxes into {match, mask, fields}: each box covers
// bit positions [hibit-width+1, hibit]; its <c> children are MSB-first bit
// values. Fixed c-values (0/1) set match/mask; a named box with NO c-values is
// a variable operand field.
func parseBoxes(boxes []struct {
	Hibit string   `xml:"hibit,attr"`
	Width string   `xml:"width,attr"`
	Name  string   `xml:"name,attr"`
	C     []string `xml:"c"`
}) (entry, bool) {
	if len(boxes) == 0 {
		return entry{}, false
	}

	var e entry
	for _, b := range boxes {
		hibit, err := strconv.Atoi(b.Hibit)
		if err != nil {
			continue
		}

		// width is optional; "0"/absent means width 1 (as before)
		w := 1
		if v, werr := strconv.Atoi(b.Width); werr == nil && v != 0 {
			w = v
		}

		// A named box = a field if all c-values are SOFT (parenthesized:
		// "(1)", "(should be 11111)" — should-be operands like Ra in
		// umulh). Bare "0"/"1" are hard opcode pins: they stay in
		// Match/Mask and are not declared a field (the size boxes of
		// and/orr are opcode traps).
		softOnly := true
		for _, c := range b.C {
			if c == "0" || c == "1" {
				softOnly = false
				break
			}
		}

		if b.Name != "" && softOnly {
			e.Fields = append(
				e.Fields,
				arm64.NewField(b.Name, uint(hibit-w+1), uint(w)),
			)
			continue
		}

		for i, c := range b.C {
			pos := uint(hibit - i)
			switch c {
			case "1":
				e.Mask |= 1 << pos
				e.Match |= 1 << pos
			case "0":
				e.Mask |= 1 << pos
			}

			// other tokens (x, (0), (1), ?) → variable/unmatched, ignored
		}
	}

	return e, true
}

// mnemonic extracts the first alphabetic token of the asmtemplate(s) ("ADD",
// "NOP", "BL"), lowercased. asmtemplates sometimes carry a leading variant
// index like "0". The tokenizer is a parsec grammar (mnemonic.go).
func mnemonic(encs []struct {
	Asmtemplate struct {
		Text []string `xml:"text"`
	} `xml:"asmtemplate"`
}) (string, error) {
	for _, enc := range encs {
		m, err := firstAlphaRun(enc.Asmtemplate.Text)
		if err != nil {
			return "", err
		}

		if m != "" {
			return m, nil
		}
	}

	return "", nil
}

func requireData(dir string) {
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "\nerror: %s not found (or not a dir)\n", dir)
		fmt.Fprintf(
			os.Stderr,
			"       It is gitignored (canonical, regenerable) and must be downloaded.\n\n",
		)
		fmt.Fprintf(os.Stderr, "       Fix:    make update-arm-instr-data\n")
		fmt.Fprintf(os.Stderr, "       Then:  make gen-arm-instr\n\n")
		os.Exit(1)
	}
}

func emit(outPath string, entries []entry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Name != entries[j].Name {
			return entries[i].Name < entries[j].Name
		}

		return entries[i].Match < entries[j].Match
	})

	var b strings.Builder
	fmt.Fprintln(&b, "// Code generated by gen/cmd/gen-arm-instr; DO NOT EDIT.")
	fmt.Fprintln(&b, "//")
	fmt.Fprintln(&b, "// ARM A64 instruction encodings, parsed from the official A64 ISA XML")
	fmt.Fprintln(&b, "// (arch/arm64/data/instr). Each entry is one <iclass> regdiagram: fixed")
	fmt.Fprintln(&b, "// boxes → Match/Mask; named variable boxes → Fields. Authoritative")
	fmt.Fprintln(&b, "// encoding data straight from the ARM release.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "package arm64")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "type armISAEntry struct {")
	fmt.Fprintln(&b, "\tName   string")
	fmt.Fprintln(&b, "\tMatch  uint32")
	fmt.Fprintln(&b, "\tMask   uint32")
	fmt.Fprintln(&b, "\tFields []Field")
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	fmt.Fprintln(
		&b,
		"func newArmISAEntry(name string, match uint32, mask uint32, fields []Field) armISAEntry {",
	)
	fmt.Fprintln(&b, "\treturn armISAEntry{")
	fmt.Fprintln(&b, "\t\tName:   name,")
	fmt.Fprintln(&b, "\t\tMatch:  match,")
	fmt.Fprintln(&b, "\t\tMask:   mask,")
	fmt.Fprintln(&b, "\t\tFields: fields,")
	fmt.Fprintln(&b, "\t}")
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "var armISA = []armISAEntry{")
	for _, e := range entries {
		if len(e.Fields) == 0 {
			fmt.Fprintf(&b, "\tnewArmISAEntry(%q, 0x%08x, 0x%08x, nil),\n", e.Name, e.Match, e.Mask)
			continue
		}

		fmt.Fprintf(
			&b,
			"\tnewArmISAEntry(%q, 0x%08x, 0x%08x, []Field{\n",
			e.Name,
			e.Match,
			e.Mask,
		)
		for _, f := range e.Fields {
			fmt.Fprintf(&b, "\t\tNewField(%q, %d, %d),\n", f.Name, f.Offset, f.Width)
		}

		fmt.Fprintln(&b, "\t}),")
	}

	fmt.Fprintln(&b, "}")
	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		log.Fatalf("write %s: %v", outPath, err)
	}

	log.Printf("wrote %d entries to %s", len(entries), outPath)
}
