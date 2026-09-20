package unit

// Sym is a deferred record: it knows its byte size up front (the layout
// counts labels and stream offsets before any symbol value is known) and
// resolves into one or more purely encodable records once its address
// and the program symbols are known. A branch to a label, an address
// pair, an inline-asm fragment, a symbolic data word - every hole in the
// stream is a Sym; the resolve phase fills them all in one walk.
type Sym interface {
	// Size is the reserved byte count (stable from parse time: the
	// resolve phase encodes exactly this many bytes or fails).
	Size() int

	// Resolve evaluates the record at its final address.
	Resolve(ctx Ctx) ([]Resolved, error)
}
