# Assembly [![Powered by GLM](https://img.shields.io/badge/Powered_by-GLM_(Z.ai)-black?logo=zdotai)](https://z.ai)

Low-level toolkit to work at the lowest level.

Data driven - the authoritative decode data is generated, not hardcoded: match/mask tables come from each architecture's machine-readable source (the ARM A64 ISA XML, Spike's `encoding.h`, the `loongarch-opcodes` tables) via `gen/cmd/*` generators; what cannot be machine-generated (operand formats, aliases) is hand-written and joins the generated tables by mnemonic.

## Motivation

- interested in low-level stuff or just system programming
- interested in AI / agentic development (all code has been written by GLM (z.ai))
- interested in spec-driven development (A64 ISA XML + RISC-V header file + loongarch-opcodes tables)
- want to build something complex straight to machine code - without C or asm

## Architectures

Three ISA backends, fully separated (no shared line between them):

- **ARM64** — decode + assemble; match/mask from the official A64 ISA XML
- **RISC-V** (RV64GC) — decode + assemble with RVC compression; encodings from Spike's `encoding.h`
- **LoongArch** (LA64) — decode + assemble, the full scalar integer set (248 instructions); encodings from the `loongarch-opcodes` tables

## Status

- **99.97% byte-for-byte match** with `llvm-objdump` on a real Go Mach-O binary (155,576 of 155,625 instructions); RISC-V example matches 100%
- **LoongArch** — the third backend: the full scalar integer set (248 instructions), 100% against `llvm-objdump` in a per-word differential over the whole table
- **Assembler with round-trip oracle** — GNU-as compatible syntax subset plus assembly's own `Text()` output:
  - ARM64: **99.86%** byte-exact round-trip over the 155k-instruction example (the rest are decode-equivalent encodings of ambiguous aliases)
  - RISC-V: 89/93 byte-exact + 4 FP-rounding (lossy in objdump-style text) on the example; 400/400 on a synthetic corpus
  - LoongArch: byte-exact on the example; `li.w`/`li.d` ladders and `la` pairs match `llvm-mc` byte-for-byte
  - RVC compressor: 32-bit forms compress to 16-bit exactly like GNU as
- **Data-driven schemas (ARM64)** — each ARM64 instruction is a `Schema{Mask, Value, Fields, Formatter}` with string-keyed transforms/formatters (no closures); RISC-V and LoongArch join generated `{match, mask}` tables to hand-written decode tables
- **Self-contained binary parsers** — full ELF and Mach-O container parsers in `file/elf` and `file/macho`, fully separated like `arch/arm64`, `arch/riscv` and `arch/loong64` (not a single shared line), all via parsec:
  - **ELF**, 32/64 LE/BE: program headers, section headers (incl. extended numbering), symtab/dynsym, RELA/REL with all `R_AARCH64_*`/`R_RISCV_*`, dynamic + `DT_*`, notes/build-id, verneed/verdef, `.gnu.hash`/`.hash`
  - **Mach-O**, incl. FAT/Universal: all ~50 load commands, nlist + stabs, dysymtab + indirect, section relocations, dyld_info rebase/bind/lazy/weak opcode streams, exports trie, chained fixups, function starts, data-in-code, thread states, code-signature superblob
  - zero runtime dependency on `debug/macho` or `debug/elf` (stdlib only as the diff oracle in tests)
- **Interactive web UI** — drag-and-drop disassembly with syntax highlighting and virtual scrolling, plus an assembler panel (POST /api/v1/asm)
- **Diff tool** — `assembly-diff` compares output against objdump for iterative coverage development

## Economics

The build statistics live in [ECONOMICS.md](ECONOMICS.md) — a log with
an entry per big feature (the base build, then each architecture that
landed on top).

## Quick start

```bash
# Build everything
go build ./...

# Run the full test suite (requires objdump and qemu-system on PATH)
make tests

# Start the interactive server
go run ./cmd/assembly-server
# → open http://127.0.0.1:8080, drop a binary on the page

# Check coverage against objdump
go run ./tests/cmd/assembly-diff tests/examples/hello-world/hello-world

# Full assembler cycle: assemble → execute → disassemble
make -C tests/examples/hello-asm          # prints "hello world"
make -C tests/examples/hello-asm disasm   # disassembles hello.bin back

# Cross-architecture VMs (qemu-system): riscv64 + arm64 + loong64 bare-metal
make -C tests/examples/hello-asm vm       # all three VMs print "hello world"
make -C tests/examples/hello-asm compare  # same program, three ISAs side by side
```

## Usage

```bash
# Assemble (raw binary out, hex dump to stdout, symbol table)
go run ./cmd/assembly-asm -arch riscv64 -base 0x1000 -o prog.bin --sym prog.sym --hex prog.s

# Assemble into a minimal static ELF64 — runs natively on Linux arm64/riscv64/loong64,
# under qemu on any host, no system assembler/linker needed
go run ./cmd/assembly-asm -arch arm64 -base 0x40100000 -format elf -o prog.elf prog.s

# Disassemble a raw binary (the reverse direction of the same tool;
# our own ELF layout is recognized automatically)
go run ./cmd/assembly-asm -arch arm64 --disasm prog.bin

# Interactive web UI: disassembly viewer + assembler panel
go run ./cmd/assembly-server        # → http://127.0.0.1:8080
#   POST /api/v1/disasm  — multipart binary → disassembly JSON
#   POST /api/v1/asm     — {arch, source, baseAddr} → sections/symbols/errors

# Compare local disassembly against objdump (coverage tool)
go run ./tests/cmd/assembly-diff tests/examples/hello-world/hello-world

# The demo folder does all of it: assemble → execute (native + VMs) → disasm
make -C tests/examples/hello-asm            # native: "hello world"
make -C tests/examples/hello-asm vm         # qemu-system riscv64 + arm64 + loong64 VMs
make -C tests/examples/hello-asm compare    # same program across three ISAs
```

Syntax: GNU-as compatible subset plus assembly's own `Text()` output
(objdump-style absolute branch targets, aliases) — the latter gives the
round-trip oracle `assemble(Text(instr)) == instr`. v1 limits: no
relocations/.o output, no `.macro`, no relaxation (documented in
`.agent/README.md`, kept locally).

## How it's tested

One command runs every gate: `make tests`. What it actually checks:

- **Unit + property tests** — `go test ./...`:
  - round-trip laws on generated instruction samples: `enc∘dec∘enc == enc` in a fixed context
  - any counterexample shrinks to a minimal one
  - race-detector clean
- **Differential gates vs objdump** — every instruction diffed against objdump, address-keyed, on real binaries:
  - ARM64: 99.97% on a Go Mach-O build (155,576 of 155,625)
  - RISC-V: 100% on the example
- **Round-trip fidelity on real binaries**:
  - the project's own CLIs, cross-built for linux/arm64, linux/riscv64 and linux/loong64
  - each binary: disasm → listing → assemble, 3 iterations
  - sha256 gate on `.text` + the whole ELF rebuilt by our own writer
  - 9/9 binaries pass (3 CLIs × 3 architectures), millions of instructions each
- **Behavioral VM matrix**:
  - original vs rebuilt binaries run side by side
  - isolated `qemu-system` VMs, aarch64 + riscv64 + loong64
  - the code we emit must actually execute — and behave identically
- **Decode is proven, not sampled**:
  - property oracle equivalent to the linear rule list
  - tree depth bounded by the instruction word width

On top of that, `make lint` (gofmt, go vet, errcheck, golangci-lint) runs on every change.

## Dependencies

- [okneniz/parsec](https://github.com/okneniz/parsec) — parser combinators
- [okneniz/oh-snap](https://github.com/okneniz/oh-snap) — for simplest property tests
- Go standard library only (net/http, encoding/binary, embed, etc.)

## Design decisions

Conventions and trade-offs established for this codebase (ARM64, RISC-V and LoongArch).

**Data-driven schemas (ARM64).** Every ARM64 instruction is `Schema{Mask, Value, Fields, Formatter}` matched by `(word & Mask) == Value`. `Transform`/`Formatter` are string keys into registries (not closures), so a schema stays plain, serializable data.

**Decode data is generated; formatting is hand-written.** Match/mask tables come from each arch's authoritative source (Spike's `encoding.h`, ARM sysreg XML + m1n1, the loongarch-opcodes tables) via `gen/cmd/*`; operand layout and alias mnemonics are hand-written because they can't be machine-generated. RISC-V joins the two by mnemonic: `decodeTable` holds the format config, `riscvEncodings` (generated) holds the match/mask.

**No shared instruction type — instructions represent themselves.** Each arch owns its decode result with its own `Parse` and rendering; duplication between the arches is accepted as the lesser evil compared to a bad abstraction (a shared "universal instruction" model was tried and reverted). What crosses package boundaries are only small capability interfaces (`disasm.ObjDump`, `asm.Instr`); consumers `switch` on `file.ArchKind` and call the arch directly — no registry, no shared `Architecture` interface.

## System register data

MSR/MRS operands are named via a generated table (`arch/arm64/sysregs_generated.go`, ~800 entries) rather than a hand-coded switch. It merges two sources, vendored under `arch/arm64/data/`:

- **ARM System Register XML** (`data/sysreg/AArch64-*.xml`) — the official machine-readable spec, covering all architectural A-profile registers across EL1/EL2/EL3.
- **`data/apple_regs.json`** — Apple M1 implementation-defined registers (the IMPDEF space, `op0=3 CRn=15`) reverse-engineered by the [m1n1](https://github.com/AsahiLinux/m1n1) project. These are absent from the ARM XML.

The `sysreg` field on the MSR/MRS schemas uses a `"sysreg"` transform (same pattern as `cond`/`shiftName`) that looks up the 15-bit system-register selector against this table, falling back to the objdump-style `S<op0>_<op1>_C<CRn>_C<CRm>_<op2>` form when unknown.

### Refreshing the data

```bash
# Re-download the vendored sources (ARM XML + m1n1 apple_regs.json)
make update-sysreg-data

# Regenerate the Go table (also formats the tree and runs vet)
make gen-sysregs
```

Pin a different ARM release or m1n1 revision via env vars (see `arch/arm64/data/update.sh`):

```bash
ARM_SYSREG_URL='https://.../SysReg_xml_v9x-YYYY-MM.tar.gz' make update-sysreg-data
```

The vendored data carries its own licenses — see [`arch/arm64/data/LICENSE`](arch/arm64/data/LICENSE). The ARM XML is under Arm's click-through terms (each file retains its header); `apple_regs.json` is MIT (© The Asahi Linux Contributors).

## LoongArch opcode data

The `arch/loong64` encodings come from the community
[loongarch-opcodes](https://github.com/loongson-community/loongarch-opcodes)
tables (12 scalar-integer subsets, CC-BY 4.0), vendored under
`arch/loong64/data/`:

- the official mnemonics are restored from the upstream `@orig_name` renames;
  `csrrd`/`csrwr` — merged upstream — are derived back from `csrxchg`
- `gen-loongarch-instr` emits the 248-entry `{match, mask}` table; the mask
  is derived from the operand-format notation, not from trailing zero bits

```bash
make update-loong-data   # re-download the vendored tables
make gen-loongarch-instr # regenerate the Go table
```

## License

MIT.

The vendored system-register data under `arch/arm64/data/` is third-party and
carries its own licenses (Arm click-through for the XML, MIT for `apple_regs.json`);
see [`arch/arm64/data/LICENSE`](arch/arm64/data/LICENSE). The loongarch-opcodes
tables under `arch/loong64/data/` are CC-BY 4.0; see
[`arch/loong64/data/LICENSE`](arch/loong64/data/LICENSE).
