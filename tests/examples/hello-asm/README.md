# hello-asm - a full assembly cycle demo

An assembly source is built by `assembly-asm`, executed, and disassembled
back - all three directions with a single utility.

```bash
cd tests/examples/hello-asm

make          # build hello.bin + runner, execute → "hello world"
make disasm   # disassemble hello.bin back
make clean
```

## Virtual machines (qemu-system)

`make vm` runs the bare-metal variants in full qemu virtual machines
(`brew install qemu`) - hypervisor-level isolation, safe for untrusted code:

```bash
make vm-riscv   # qemu-system-riscv64 -machine virt: UART + sifive_test-exit
make vm-arm     # qemu-system-aarch64 -machine virt: PL011 + PSCI SYSTEM_OFF
make compare    # the same program on three instruction sets, side by side
```

- `hello-riscv.s` / `hello-arm-vm.s` - bare-metal: byte-by-byte output to the
  UART (MMIO), termination via the sifive_test test device (riscv) / PSCI
  through HVC (arm).
- addresses: riscv - RAM 0x80000000; arm - 0x40100000 (above the DTB area).
- `TestHelloVM` in the repository root runs both VMs and requires "hello
  world".

For true isolation of untrusted arm64/Linux code on macOS there is also
Virtualization.framework (UTM, Lima) - a VM of the native architecture; riscv
is possible only through emulation (qemu). See `asm.EmitELF` - the ELF works
in both.

## What is what here

- `hello-macos.s` / `hello-linux.s` - hello world on direct syscalls
  (write/exit), position-independent (the string address is `adr`
  pc-relative), the message embedded into the code via `.ascii`. The variant
  is selected by `uname`.
- `runner.c` - a minimal loader: mmaps the binary (on macOS with `MAP_JIT`
  and the W^X toggle plus icache invalidation) and jumps to the start of the
  code.
- `hello.sym` - the symbol table (`start`, `msg`) from `--sym`.

The disassembled tail is the message data: the decoder honestly shows the
`.word` fallback and random "instructions", just like objdump on data inside
.text.

## All Makefile targets

| Target | What it does |
|---|---|
| `make` / `make run` | build `hello.bin` (the syscall variant selected by uname) + runner → execute |
| `make disasm` | disassemble `hello.bin` |
| `make vm-riscv` | build `hello-riscv.elf` → run in qemu-system-riscv64 |
| `make vm-arm` | build `hello-arm-vm.elf` → run in qemu-system-aarch64 |
| `make vm` | both VMs in a row |
| `make compare` | one program on three ISAs: sizes + full disassembly side by side |
| `make clean` | remove artifacts (bin/sym/elf/runner) |

## Troubleshooting

- **`qemu-system-*: command not found`** - `brew install qemu` (you need
  qemu-system; it is not part of qemu-user in brew).
- **Strange output / the wrong thing runs** - `make clean && make`: the CLI
  does not write `-o` on build errors, but old artifacts can survive a
  rebuild incrementally.
- **Linux host**: `make run` takes `hello-linux.s` and executes the ELF
  directly (`-format elf`), no runner needed.

## Manual run

```bash
go run ../../../cmd/assembly-asm -arch arm64 -o hello.bin --sym hello.sym hello-macos.s
cc -O2 -o runner runner.c
./runner hello.bin
go run ../../../cmd/assembly-asm -arch arm64 --disasm hello.bin
```

Regression test: `TestHelloAsmExample` in the repository root checks that
both sources assemble and their binaries pass a byte-exact round-trip.
