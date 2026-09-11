package load

// Input is what to load: a source file (assembled in-process) or a
// binary with an optional symbol sidecar; Base overrides the arch's
// default section base (hex or dec, "" keeps it).
type Input struct {
	SrcPath string
	BinPath string
	SymPath string
	Base    string
}
