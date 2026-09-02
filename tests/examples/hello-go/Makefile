# The pure-Go asm examples: the same programs as the .s family of
# ../hello-asm, written on the prog/<arch> chains (no .s source, no
# expression evaluator - everything is computed by Go).
#
#   make run-macos  - build hello-go-macos and execute it via the runner
#                     harness (the raw code image; a macho wrapper of the
#                     same code is in file.MachoBuilder)
#   make run-loong  - build hello-go-loong and run it in the QEMU virt
#                     machine (bare-metal, flash at 0x1c000000)
#
# Run from this directory (or make -C tests/examples/hello-go from the
# repository root).

ASM := go run ./macos
LOONG := go run ./loong
RUNNER := cc -O2 -o runner ../hello-asm/runner.c && ./runner
QEMU_LOONG ?= qemu-system-loongarch64

all: run-macos

hello-go-macos.bin:
	$(ASM) -o hello-go-macos.bin

run-macos: hello-go-macos.bin
	@echo "--- executing hello-go-macos.bin (native, via the runner) ---"
	@$(RUNNER) hello-go-macos.bin

hello-go-loong.bin:
	$(LOONG) -o hello-go-loong.bin

# The machine has no poweroff register: the program idles, the gate
# stops qemu once the line is out.
run-loong: hello-go-loong.bin
	@echo "--- executing hello-go-loong.bin (qemu virt, bare-metal) ---"
	@$(QEMU_LOONG) -machine virt -nographic \
		-device loader,file=hello-go-loong.bin,addr=0x1c000000,cpu-num=0 \
		> loong.out 2>&1 & QEMU=$$!; \
	for i in $$(seq 1 50); do grep -q hello loong.out && break; sleep 0.2; done; \
	kill $$QEMU 2>/dev/null; cat loong.out; rm -f loong.out

clean:
	rm -f hello-go-macos.bin hello-go-loong.bin runner loong.out

.PHONY: all run-macos run-loong clean
