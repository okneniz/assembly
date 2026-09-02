# The toolchain container of the test gates: Go + clang/LLVM with all the
# targets (RISC-V and LoongArch included). Everything the LLVM-oracle tests
# need lives here; QEMU is deliberately NOT installed - the VM gates stay on
# the host (make -C tests rt-vm).
#
# make -C tests docker builds and runs it; the engine is colima.

FROM golang:1.26-trixie

# clang - the assembler and fixtures compiler (one multi-target binary);
# llvm  - llvm-mc, llvm-objdump, llvm-readelf (the test oracles).
RUN apt-get update \
	&& apt-get install -y --no-install-recommends clang llvm \
	&& rm -rf /var/lib/apt/lists/*
