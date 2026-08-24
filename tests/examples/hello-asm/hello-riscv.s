// hello world for RISC-V bare-metal on qemu virt (qemu-system-riscv64).
// No OS: output is byte writes to the UART (ns16550 @ 0x10000000, MMIO),
// termination is writing 0x5555 (FINISHER_PASS) to the sifive_test test
// device (@ 0x100000), after which qemu powers off the virtual machine.
//
// Build: assembly-asm -arch riscv64 -base 0x80000000 -format elf
// Run: qemu-system-riscv64 -machine virt -bios none -nographic -kernel ...

start:
    li    a0, 0x10000000      # UART base
    la    a1, msg             # cursor
    la    a2, end             # end
loop:
    bgeu  a1, a2, done        # while cursor < end
    lb    a5, 0(a1)           # string byte
    sb    a5, 0(a0)           # → UART
    addi  a1, a1, 1
    j     loop
done:
    li    a5, 0x100000        # sifive_test
    li    a6, 0x5555          # FINISHER_PASS → exit(0)
    sw    a6, 0(a5)
hang:
    j     hang

msg:
    .ascii "hello world\n"
end:
