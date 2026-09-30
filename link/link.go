// Package link - the in-memory linker: the .S sources of a program
// assemble into ONE unit output and resolve together at its phase
// boundary. Every source's text and data sections deposit first, the
// zero-fill reserves of all of them after - the program's bss is the
// unit's single tail, as in ld. The label namespace is the link
// interface: a .L local stays behind its file key, every other name
// (and a .global-promoted one) is shared by all the sources, a
// duplicate definition or an unresolved reference surfaces at resolve.
package link

import (
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/unit"
)

// Source is one .S input of a link: its text and the file name that
// keys its .L locals and its line-map origin.
type Source struct {
	File string
	Src  string
}

// Deps is the arch injection of the driver: the assembler entry that
// parses one source (alias.ParseSourceUnit on arm64). The driver itself
// knows no architecture.
type Deps struct {
	Parse func(file, src string) (*asm.SourceUnit, []asm.AsmError)
}

// Link assembles the sources into one unit and resolves it: parse
// errors are reported (assembly continues past a failed source, every
// error surfaces together), the sections of all sources deposit before
// the bss reserves of all of them, entry ("" - none) is the unit's
// entry label, place is the placement policy of the resolve. The result
// is the program: encodable records, the symbol table, the line map.
func Link(deps Deps, sources []Source, entry string, place unit.Place) *unit.Fixed {
	u := unit.New()

	var parsed []*asm.SourceUnit
	var errs []error

	for _, s := range sources {
		su, asmErrs := deps.Parse(s.File, s.Src)
		for _, e := range asmErrs {
			errs = append(errs, e)
		}

		if su != nil {
			parsed = append(parsed, su)
		}
	}

	for _, su := range parsed {
		su.Deposit(u)
	}

	for _, su := range parsed {
		su.DepositBss(u)
	}

	if entry != "" {
		u.Entry(entry)
	}

	f := u.Resolve(place)
	f.Errs = append(errs, f.Errs...)
	return f
}

// EntryOf picks the entry symbol of a linked program: the explicit
// name, else the start/_start convention. The name still has to exist -
// the caller reports a program with no entry of its own.
func EntryOf(f *unit.Fixed, explicit string) (string, bool) {
	if explicit != "" {
		return explicit, true
	}

	for _, name := range []string{"start", "_start"} {
		if _, ok := f.Syms[name]; ok {
			return name, true
		}
	}

	return "", false
}
