package riscv

// EncOpts - the encoding modes of a computed instruction. A value, not
// an environment: there is no symbolic knowledge here - the NoRVC flag
// combines .option norvc and "the slot was symbolic" (a decision of the
// asm/riscv syntax layer: symbolic targets are not compressed so that
// sizes agree across assembler passes).

// EncOpts - encoding modes: NoRVC disables RVC compression.
type EncOpts struct {
	NoRVC bool
}
