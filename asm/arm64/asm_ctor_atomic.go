package arm64

// Assembler constructors for atomics (ldar/stlr) and exclusives
// (stlxr/stxrb). Addressing — only [rn] without an offset.

import (
	"errors"
	"fmt"
)

// newLdarArm — ldar rt, [rn] (x/w).
func newLdarArm(ops []vOp) (Instr, error) {
	rt, rn, err := rtMem(ops, "ldar")
	if err != nil {
		return nil, err
	}

	enc := uint32(0xC8DFFC00)
	if rt[0] == 'w' {
		enc = 0x88DFFC00
	}

	return LdarOf(rt, rn, enc), nil
}

// newStlrArm — stlr rt, [rn].
func newStlrArm(ops []vOp) (Instr, error) {
	rt, rn, err := rtMem(ops, "stlr")
	if err != nil {
		return nil, err
	}

	enc := uint32(0xC89FFC00)
	if rt[0] == 'w' {
		enc = 0x889FFC00
	}

	return StlrOf(rt, rn, enc), nil
}

// rtMem — rt, [rn] without an offset.
func rtMem(ops []vOp, name string) (string, string, error) {
	if len(ops) != 2 || !ops[1].IsMem() {
		return "", "", fmt.Errorf("%s: want rt, [rn]", name)
	}

	rt, err := wantAReg(ops[0], name)
	if err != nil {
		return "", "", err
	}

	return rt, ops[1].Mem().Base(), nil
}

// newStlxrArm — stlxr rs, rt, [rn].
func newStlxrArm(ops []vOp) (Instr, error) {
	return makeExclCtor(ops, 0xC800FC00)
}

// newStxrbArm — stxrb rs, rt, [rn].
func newStxrbArm(ops []vOp) (Instr, error) {
	return makeExclCtor(ops, 0x08000000)
}

func makeExclCtor(ops []vOp, enc uint32) (Instr, error) {
	if len(ops) != 3 || !ops[2].IsMem() {
		return nil, errors.New("want rs, rt, [rn]")
	}

	rs, err := wantAReg(ops[0], "stlxr")
	if err != nil {
		return nil, err
	}

	rt, err := wantAReg(ops[1], "stlxr")
	if err != nil {
		return nil, err
	}

	if enc == 0xC800FC00 && rt[0] == 'w' {
		enc = 0x8800FC00
	}

	if enc == 0x08000000 {
		return StxrbOf(rs, rt, ops[2].Mem().Base(), enc), nil
	}

	return StlxrOf(rs, rt, ops[2].Mem().Base(), enc), nil
}
