// Package debug is the arch-neutral core of the debugger: the target
// description vocabulary. The executor (debug/qemu), the protocol
// (debug/rsp), and the engine (debug/session) are all arch-blind; the
// per-arch packages (debug/arm64, ...) provide the Target.
//
// The interface carries data and rendered text, never instructions -
// the codebase rule of no shared instruction model holds here too: the
// arch packages decode and render with their own concrete types and
// hand over listing lines.
package debug

// Target is the arch-specific surface of debugging: the register
// layout, the qemu executor configuration, and the disassembly glue.
// One implementation per arch package; the implementations are values
// (no state) - NewTarget ctors in debug/arm64 and friends.
type Target interface {
	// Arch is the architecture name ("arm64").
	Arch() string

	// QemuBinary is the executor executable (qemu-system-aarch64).
	QemuBinary() string

	// Registers is the ordered register set (the 'g' block layout).
	Registers() []Reg

	// PCNum and SPNum are the RSP register numbers of pc and sp.
	PCNum() int
	SPNum() int

	// QemuArgs is the arch part of the executor command line for the
	// image path (machine, cpu, loader); the executor adds the
	// common tail (-nographic, -S, the gdbstub) itself.
	QemuArgs(imgPath string) []string

	// Disasm renders the listing lines of the code buffer located at
	// addr (its own decoders; one line per instruction).
	Disasm(code []byte, addr uint64) []string

	// InstrLen is the instruction length at the head of the code
	// buffer (4, or 2 for a compressed riscv instruction); a nil or
	// short buffer yields the arch constant - the breakpoint kind.
	InstrLen(code []byte) int
}
