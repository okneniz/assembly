# Assembly Debug — the VSCode extension

A declarative extension: no JavaScript. It declares the `assembly-debug`
debug type and launches the bundled adapter binary (the `make vscode`
target installs the folder together with a freshly built adapter and
writes its absolute path into the manifest — a dock-launched VSCode runs
its children with the minimal system PATH, so the adapter must not be
resolved through it).

## Install

From the repository root (restart VSCode afterwards — extensions are
scanned at startup):

```bash
make vscode
```

## Use

Create `.vscode/launch.json` in the workspace:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "type": "assembly-debug",
            "request": "launch",
            "name": "Debug hello (arm64)",
            "arch": "arm64",
            "source": "${workspaceFolder}/hello.s"
        }
    ]
}
```

Press F5. The program boots frozen under the arch's qemu
(`brew install qemu` first), stops at the entry, and the editor
debugger drives it:

- breakpoints at source lines and at labels (the BREAKPOINTS panel
  section "Function Breakpoints")
- step, continue, pause — the register scope in the VARIABLES panel
- the program's own printing in the DEBUG CONSOLE
- the DISASSEMBLY view (right-click → Open Disassembly View) and the
  MEMORY viewer (a register's context menu → View Memory), both
  through the toolkit's own decoders

## Attributes

| Attribute | Meaning |
|---|---|
| `arch` | `arm64`, `riscv64`, or `loong64` |
| `source` | the `.s` file: assembled in-process, labels and lines come along |
| `bin` | a prebuilt image instead of a source |
| `sym` | the symbol sidecar of `bin` (the CLI's `-sym` output) |
| `base` | load address override (hex or dec; the arch default otherwise) |
| `icount` | deterministic virtual clock |

Exactly one of `source` or `bin` runs the program.
