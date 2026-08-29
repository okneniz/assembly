# A full assembly cycle demo: assembly → execution → disassembly.
#
#   make          - build hello.bin (assembly-asm) + runner, print hello world
#   make disasm   - disassemble hello.bin back (assembly-asm -disasm)
#   make clean
#
# Run from this directory (or make -C tests/examples/hello-asm from the repository root).

UNAME := $(shell uname)
ifeq ($(UNAME),Darwin)
HELLO_SRC := hello-macos.s
else
HELLO_SRC := hello-linux.s
endif

CC ?= cc
ASM := go run ../../../cmd/assembly-asm

# qemu-system-* from brew: full virt machines (VMs) for riscv64/arm64.
QEMU_RISCV ?= qemu-system-riscv64
QEMU_ARM   ?= qemu-system-aarch64
QEMU_LOONG ?= qemu-system-loongarch64

all: run

hello.bin: $(HELLO_SRC)
	$(ASM) -arch arm64 -o hello.bin --sym hello.sym $(HELLO_SRC)

runner: runner.c
	$(CC) -O2 -o runner runner.c

run: hello.bin runner
	@echo "--- executing hello.bin ---"
	./runner hello.bin

disasm: hello.bin
	@echo "--- disassembly of hello.bin ---"
	$(ASM) -arch arm64 --disasm hello.bin

# --- Virtual machines (qemu-system, brew install qemu) ---

hello-riscv.elf: hello-riscv.s
	$(ASM) -arch riscv64 -base 0x80000000 -format elf -o hello-riscv.elf hello-riscv.s

vm-riscv: hello-riscv.elf
	@echo "--- riscv64 VM (qemu virt, bare-metal) ---"
	$(QEMU_RISCV) -machine virt -bios none -nographic -kernel hello-riscv.elf

hello-arm-vm.elf: hello-arm-vm.s
	$(ASM) -arch arm64 -base 0x40100000 -format elf -o hello-arm-vm.elf hello-arm-vm.s

vm-arm: hello-arm-vm.elf
	@echo "--- arm64 VM (qemu virt, bare-metal) ---"
	$(QEMU_ARM) -machine virt -cpu cortex-a53 -nographic -device loader,file=hello-arm-vm.elf,cpu-num=0

hello-loongarch.bin: hello-loongarch.s
	$(ASM) -arch loong64 -base 0x1c000000 -format bin -o hello-loongarch.bin hello-loongarch.s

vm-loong: hello-loongarch.bin
	@echo "--- loong64 VM (qemu virt, bare-metal) ---"
	$(QEMU_LOONG) -machine virt -nographic -device loader,file=hello-loongarch.bin,addr=0x1c000000,cpu-num=0

vm: vm-riscv vm-arm vm-loong

# Architecture comparison: the same program assembled for riscv64/arm64.
compare: hello.bin hello-riscv.elf hello-arm-vm.elf
	@echo "=== ARM64 (macOS syscall variant) - $$(( $$(stat -f%z hello.bin) )) bytes ==="
	$(ASM) -arch arm64 --disasm hello.bin
	@echo ""
	@echo "=== ARM64 bare-metal (qemu virt) - $$(( $$(stat -f%z hello-arm-vm.elf) - 4096 )) bytes ==="
	$(ASM) -arch arm64 -base 0x40100000 --disasm hello-arm-vm.elf
	@echo ""
	@echo "=== RISC-V bare-metal (qemu virt) - $$(( $$(stat -f%z hello-riscv.elf) - 4096 )) bytes ==="
	$(ASM) -arch riscv64 -base 0x80000000 --disasm hello-riscv.elf

clean:
	rm -f hello.bin hello.sym runner hello-riscv.elf hello-arm-vm.elf

.PHONY: all run disasm vm vm-riscv vm-arm compare clean
