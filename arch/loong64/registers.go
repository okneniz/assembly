package loong64

// laRegNames - register names in LoongArch ABI notation, as llvm-objdump
// prints them: zero, ra, tp, sp, a0-a7, t0-t8, then r21 (no ABI alias),
// fp, s0-s8.
var laRegNames = [32]string{
	"zero", "ra", "tp", "sp", "a0", "a1", "a2", "a3",
	"a4", "a5", "a6", "a7", "t0", "t1", "t2", "t3",
	"t4", "t5", "t6", "t7", "t8", "r21", "fp", "s0",
	"s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8",
}

// laRegName - the printed form of a register number ("$t0", "$r21", ...).
func laRegName(num uint8) string {
	return "$" + laRegNames[num]
}
