# tests/cmd/

This directory holds commands and wrappers needed only for tests and quality
gates - everything that is not a product CLI (`cmd/assembly-asm`,
`cmd/assembly-server`).

- `assembly-diff/` - differential of our disassembler against objdump
  (coverage gate; also run in the `make rt-vm` VM matrix).
- `assembly-rt/` - round-trip fidelity: decode → asm → bytes with a sha256
  gate (`make rt-run`).
- `objdump/` - wrapper over the external GNU/llvm objdump: runs it with
  candidate probing (Apple's binary cannot handle RISC-V), Mach-O/ELF
  arguments, and a parsec grammar for its output lines; the oracle of the
  differential tests and `-diff`. Any new wrapper over an external utility
  goes here as well.
