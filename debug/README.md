# debug/

The debugger of the toolkit: it runs a program under qemu, frozen at
start, and lets you drive it — breakpoints at labels and source lines,
stepping, register and memory inspection, disassembly, the source line
of every stop. One static binary, no gdb, no DWARF: the debug
information is the toolkit's own — the symbol table and the line map
come from the assembly step itself.

One engine, three frontends:

- the console REPL — `assembly-debug`
- the editors — VSCode and Zed, both driven by the same DAP adapter
  (`assembly-debug-dap`)
- Go programs — a debugging script is an ordinary program importing
  `debug/session`

The rest of this file: using it (below), then the package layout for
the developers of the debugger.

## Running it

A full session over a fixture (`help` lists the commands; the machine
boots frozen, stops at the entry):

```console
$ cd tests/examples/hello-asm
$ assembly-debug -arch arm64 hello-arm-vm.s
assembly-debug: qemu-system-aarch64 hello-arm-vm.s
pc 0x40100000  hello-arm-vm.s:10  mov   x0, #0x09000000 // PL011 UART base
40100000:	00 20 a1 d2	mov x0, #0x9000000
type help for the commands
(asmdb) b done
breakpoint at 0x40100024
(asmdb) c
hello world
pc 0x40100024  hello-arm-vm.s:21  movz  x0, #0x8400, lsl #16 // PSCI SYSTEM_OFF (0x84000008): qemu virt
40100024:	00 80 b0 d2	mov x0, #0x84000000
(asmdb) syms
0x40100024 done
0x4010000c loop
0x40100030 msg
0x40100000 start
(asmdb) x 0x40100030 12
40100030: 68 65 6c 6c 6f 20 77 6f 72 6c 64 0a
(asmdb) dis 3
40100024:	00 80 b0 d2	mov x0, #0x84000000
40100028:	00 01 80 f2	movk x0, #0x8
4010002c:	02 00 00 d4	hvc #0x6a0
(asmdb) list
   hello-arm-vm.s:19      b     loop
   hello-arm-vm.s:20  done:
-> hello-arm-vm.s:21      movz  x0, #0x8400, lsl #16 // PSCI SYSTEM_OFF (0x84000008): qemu virt
   hello-arm-vm.s:22      movk  x0, #0x8             // implements PSCI via HVC without firmware
   hello-arm-vm.s:23      hvc   #0
(asmdb) s
pc 0x40100028  hello-arm-vm.s:22  movk  x0, #0x8             // implements PSCI via HVC without firmware
40100028:	00 01 80 f2	movk x0, #0x8
(asmdb) q
```

The same drive without the prompt, one batch command (the gates run
this way):

```bash
assembly-debug -arch arm64 -e "b done; c; regs" tests/examples/hello-asm/hello-arm-vm.s
```

A prebuilt image with its symbol sidecar:

```bash
assembly-debug -arch arm64 -bin hello.elf -sym hello.sym
```

## Editor sessions

One DAP adapter (`assembly-debug-dap`) drives both editors; the launch
attributes are the same everywhere: `arch`, `source` (or `bin` + `sym`),
`base`, `icount`.

### VSCode

Install the extension and its adapter with one command (from the
repository root; restart VSCode afterwards — extensions are scanned at
startup):

```bash
make vscode
```

The target builds the adapter, copies the extension into
`~/.vscode/extensions`, and bundles the adapter binary with it — a
dock-launched VSCode runs its children with the minimal system PATH, so
the extension carries the absolute adapter path in its manifest and the
debugger resolves qemu outside PATH on its own.

A `.vscode/launch.json` configuration:

```json
{
    "type": "assembly-debug",
    "request": "launch",
    "name": "Debug hello (arm64)",
    "arch": "arm64",
    "source": "${workspaceFolder}/hello.s"
}
```

F5 boots the program frozen under the arch's qemu and stops at the
entry: breakpoints at source lines and labels, stepping, the register
scope in the variables panel, the disassembly and memory views, the
program's output in the debug console.

### Zed

Zed speaks DAP, but a launch configuration naming an adapter it does
not know never reaches the picker — only extension-registered names
do. The working route is a name override: the configuration claims the
built-in name `CodeLLDB`, and the `dap` setting says which binary that
name actually launches — ours.

`.zed/settings.json` (the absolute path of the adapter built by
`make cli`):

```json
{
  "dap": {
    "CodeLLDB": {
      "binary": "/path/to/assembly/bin/assembly-debug-dap"
    }
  }
}
```

`.zed/debug.json` — the launch configurations (the adapter name is the
lie, the attributes are ours; the adapter ignores everything it does
not read):

```json
[
  {
    "label": "hello arm64 (source)",
    "adapter": "CodeLLDB",
    "request": "launch",
    "arch": "arm64",
    "source": "$ZED_WORKTREE_ROOT/tests/examples/hello-asm/hello-arm-vm.s"
  }
]
```

F4 opens the session list; the rest is the same engine: the entry stop,
breakpoints at labels and source lines, stepping, the registers.

Both editor routes are verified against real qemu sessions: VSCode (F5,
the extension above) and Zed (F4, the name override above) drive the
same `assembly-debug-dap` adapter.

### prog programs

A program written on the prog chains debugs the same way, one
indirection deeper: the program serves the session itself. The launch
configuration carries a `command` instead of a source — the adapter
spawns it and relays the conversation byte for byte (`go run` and qemu
resolve outside PATH on their own):

```json
{
    "type": "assembly-debug",
    "request": "launch",
    "name": "Debug hello-go (riscv64)",
    "command": "go run ${workspaceFolder}/tests/examples/hello-go/riscv -debug"
}
```

The line map points at the Go source, so breakpoints land on the
chain-call lines and every stop carries the Go file:line — the editor
highlights the very line of the executing instruction. The serving
side is three pieces (see the riscv example's `-debug` mode): a
`WithPos` resolver (the line map needs one), `dap.NewImageLauncher`
over the assembled image, `dap.NewServer` on stdio.

## The layers

For the developers of the debugger: everything here is arch-blind
except the three thin target packages, and the whole thing speaks to
the machine through one wire protocol. Dependency order, bottom-up;
one abstraction per package:

- **target.go** — the arch vocabulary: `Target` carries the register layout
  (the `g`-block order), the qemu binary and its arch command line, the pc/sp
  register numbers, the disassembly rendering, and the instruction length
  (also the breakpoint kind). The interface moves data and rendered text,
  never instructions — no shared instruction model crosses the arch boundary.
- **rsp** — the GDB Remote Serial Protocol client: packet framing with
  checksums, acks and retransmits, the typed command set (`g`/`p`/`m`/`Z0`/`c`/`s`),
  stop-reply parsing. One command in flight at a time — all the debugging
  flow needs.
- **qemu** — the executor: a frozen (`-S`) `qemu-system-*` per run with its
  gdbstub on an ephemeral port; `Start` writes the image to a temp file,
  boots, and dials the stub; `Machine.Close` tears the process, the
  connection, and the temp file down.
- **session** — the engine: `New` negotiates and syncs with the halted
  machine; then breakpoints by symbol or address (`BreakAt`), continue/step/
  interrupt, register dumps, memory read/write, disassembly windows through
  the target's own decoders, and the address→line lookup (`LineAt`) over the
  normalized line map. Breakpoints resolve through the symbol table first.
- **arm64 / riscv / loong64** — one stateless `Target` each: register sets,
  qemu command lines, and the decode glue into the arch packages' own
  decoders (riscv `InstrLen` reads compressed 2-byte instructions).
- **load** — the input path shared by the frontends: an `.s` source is
  assembled in-process (symbols, line map, and source text come along; the
  image is an ELF, or the raw section concatenation for the loong flash
  loader), or a prebuilt binary runs as-is with its `-sym` sidecar.
- **dap** — the editor protocol: one DAP conversation over a byte transport
  (stdio for the adapter binary, pipes in tests). `Content-Length` framing,
  the camelCase wire subset, and the dispatch that translates onto the
  session — continue/step run in goroutines and report back as
  `stopped`/`terminated` events, source-line breakpoints resolve through the
  line map, label breakpoints through the symbol table, registers become the
  variables scope, the disassembly and memory windows go through the same
  decoders as the REPL.

No DWARF anywhere: the debug information is the toolkit's own — the symbol
table and the line map both come from the assembly step. `asm` records
source lines in pass 2; `prog` resolves them through the injected `WithPos`
position resolver (there is no built-in caller detection).

## Testing

The engine tests (`session`, `dap`) run against scripted RSP stubs — full
conversations with expected requests and canned replies, no qemu needed. The
end-to-end gates (`tests/debug_test.go`, `tests/dap_test.go`) boot real
machines: they skip with a hint when the arch's `qemu-system-*` is not on
PATH. The rsp packet codec carries its own property suite (round-trip and
single-flip detection at 12000 iterations).
