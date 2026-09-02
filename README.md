# Assembly [![Powered by GLM](https://img.shields.io/badge/Powered_by-GLM_(Z.ai)-black?logo=zdotai)](https://z.ai)

Low-level toolkit to work at the lowest level.

Data driven - the authoritative decode data is generated, not hardcoded: match/mask tables come from each architecture's machine-readable source (the ARM A64 ISA XML, Spike's `encoding.h`, the `loongarch-opcodes` tables) via `gen/cmd/*` generators; what cannot be machine-generated (operand formats, aliases) is hand-written and joins the generated tables by mnemonic.

## Install 

Install the CLI (Go puts the `assembly` binary on PATH):

```bash
go install github.com/okneniz/assembly/cmd/assembly@latest
```

Or as package:

```bash
go get github.com/okneniz/assembly/assembly@latest
```

## How to use

Just build executable files!

### Assembler

```
cat tests/examples/hello-asm/hello-macos.s
start:
    mov     x0, #1                  // write(fd=stdout, ...)
    adr     x1, msg                 // ... buf - string address (pc-relative)
    mov     x2, #12                 // ... len = 12
    movz    x16, #0x200, lsl #16    // x16 = 0x2000000 | ...
    movk    x16, #0x4               // ... 0x4 = write
    svc     #0x80

    mov     x0, #0                  // exit(return code = 0)
    movz    x16, #0x200, lsl #16    // x16 = 0x2000000 | ...
    movk    x16, #0x1               // ... 0x1 = exit
    svc     #0x80

msg:
    .ascii "hello world\n"
```

Assemble it into an executable and run it — from the repository root, on any arm64 Mac, with nothing but the CLI:

```bash
assembly -arch arm64 -format macho -o hello tests/examples/hello-asm/hello-macos.s
./hello                  # prints "hello world"
```

### DSL for low level programming

The `prog` package is the Go twin of an `.s` file: every chain method is a source line, a macro is an ordinary Go function returning `*Program`, labels and branch targets resolve at assembly time. This is `tests/examples/hello-go/macos/main.go` (imports: `prog/arm64`, `arch/arm64`):

```go
p := prog.New().
    Label("start").
    Mov(prog.X0, fdStdout).                     // write(fd=stdout, ...)
    Adr(prog.X1, "msg").                        // ... buf = the string
    Mov(prog.X2, int64(len(msg))).              // ... len
    Movz(prog.X16, sysClassUnix>>16, arch.Hw1). // x16 = 0x2000000 | ...
    Movk(prog.X16, sysWrite, arch.Hw0).         // ... sysWrite
    Svc(trapMach).
    Mov(prog.X0, 0). // exit(0)
    Movz(prog.X16, sysClassUnix>>16, arch.Hw1).
    Movk(prog.X16, sysExit, arch.Hw0).
    Svc(trapMach).
    Label("msg").
    Ascii(msg).
    Entry("start")

bin, _ := p.Build()
code, syms, errs := bin.Assemble(0x1000) // == hello.bin, byte for byte
```

No parser and no expression evaluator: the program is a chain of Go calls, and everything around it (constants, lengths, immediate math) is ordinary Go.

### Disassembler

`assembly -arch arm64 --disasm hello` reads the executable back (the tool recognizes its own container):

```
100000000: 20 00 80 d2  mov x0, #0x1
100000004: 21 01 00 10  adr x1, #36
100000008: 82 01 80 d2  mov x2, #0xc
10000000c: 10 40 a0 d2  mov x16, #0x2000000
100000010: 90 00 80 f2  movk x16, #0x4
100000014: 01 10 00 d4  svc #0x80
100000018: 00 00 80 d2  mov x0, #0x0
10000001c: 10 40 a0 d2  mov x16, #0x2000000
100000020: 30 00 80 f2  movk x16, #0x1
100000024: 01 10 00 d4  svc #0x80
100000028: 68 65 6c 6c  ldnp opc=0x1 imm7=0x58 Rt2=0x19 Rn=0xb Rt=0x8
10000002c: 6f 20 77 6f  umlal2.4s v15, v3, v7[3]
100000030: 72 6c 64 0a  bic w18, w3, w4, lsr #27
```

Read it against the source: the same instructions in the same order.

The independent check - the same executable through **llvm-objdump** (Apple LLVM 17) - agrees instruction for instruction, and even resolves `_main` from the symbol table the writer emitted:

```console
$ objdump -d hello
00000001000002b8 <_main>:
1000002b8: d2800020    	mov  x0, #0x1                ; =1
1000002bc: 10000121    	adr  x1, 0x1000002e0 <_main+0x28>
1000002c0: d2800182    	mov  x2, #0xc                ; =12
1000002c4: d2a04010    	mov  x16, #0x2000000         ; =33554432
1000002c8: f2800090    	movk x16, #0x4
1000002cc: d4001001    	svc  #0x80
1000002d0: d2800000    	mov  x0, #0x0                ; =0
1000002d4: d2a04010    	mov  x16, #0x2000000         ; =33554432
1000002d8: f2800030    	movk x16, #0x1
1000002dc: d4001001    	svc  #0x80
1000002e0: 6c6c6568    	ldnp d8, d25, [x11, #-0x140]
1000002e4: 6f77206f    	umlal2.4s v15, v3, v7[3]
1000002e8: 0a646c72    	bic  w18, w3, w4, lsr #27
```

## Architectures

Three ISA backends, fully separated (no shared line between them):

- **ARM64** — decode + assemble; match/mask from the official A64 ISA XML
- **RISC-V** (RV64GC) — decode + assemble with RVC compression; encodings from Spike's `encoding.h`
- **LoongArch** (LA64) — decode + assemble, the full scalar integer set (248 instructions); encodings from the `loongarch-opcodes` tables

## Motivation

- interested in low-level stuff or just system programming
- interested in AI / agentic development (all code has been written by GLM (z.ai))
- interested in spec-driven development (A64 ISA XML + RISC-V header file + loongarch-opcodes tables)
- want to build something complex straight to machine code - without C or asm

## Status

- **Disassembler** — ARM64, RISC-V (RV64GC + RVC), LoongArch; diffed against llvm-objdump on real binaries
- **Assembler** — GNU-as compatible syntax; assembles back into the same bytes
- **Executable writer** — minimal ELF64 (Linux, qemu) and Mach-O (arm64 macOS, ad-hoc signed)
- **Container parsers** — ELF and Mach-O, self-contained, no `debug/elf`/`debug/macho`
- **Generated decode tables** — from the ARM A64 XML, Spike's `encoding.h`, `loongarch-opcodes`
- **Web UI** — disassembly viewer + assembler panel
- **assembly-diff** — objdump coverage gate

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
# The CLI is assembly - install it from a clone with go install ./cmd/assembly,
# globally with
#   go install github.com/okneniz/assembly/cmd/assembly@latest
# or as a pinned tool of your own module with
#   go get -tool github.com/okneniz/assembly/cmd/assembly   # then: go tool assembly

# Assemble (raw binary out, hex dump to stdout, symbol table)
assembly -arch riscv64 -base 0x1000 -o prog.bin --sym prog.sym --hex prog.s

# Assemble into a minimal static ELF64 — runs natively on Linux arm64/riscv64/loong64,
# under qemu on any host, no system assembler/linker needed
assembly -arch arm64 -base 0x40100000 -format elf -o prog.elf prog.s

# Assemble into a Mach-O executable — runs as-is on arm64 macOS,
# no linker, no codesign (the writer ad-hoc signs the image itself)
assembly -arch arm64 -format macho -o prog prog.s

# Disassemble a raw binary (the reverse direction of the same tool;
# our own ELF layout is recognized automatically)
assembly -arch arm64 --disasm prog.bin

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
