// hello world for ARM64 bare-metal on qemu virt (qemu-system-aarch64).
// No OS: output is byte writes to the PL011 UART (@ 0x09000000, MMIO),
// termination is PSCI SYSTEM_OFF via HVC (qemu virt provides PSCI without firmware).
//
// Build: assembly -arch arm64 -base 0x40100000 -format elf
// Run: qemu-system-aarch64 -machine virt -cpu cortex-a53 -nographic
//         -device loader,file=...,cpu-num=0

start:
    mov   x0, #0x09000000 // PL011 UART base
    adr   x1, msg // cursor
    add   x2, x1, #12 // end = msg + 12
loop:
    cmp   x1, x2
    b.ge  done
    ldrb  w3, [x1] // string byte
    strb  w3, [x0] // → UART
    add   x1, x1, #1
    b     loop
done:
    movz  x0, #0x8400, lsl #16 // PSCI SYSTEM_OFF (0x84000008): qemu virt
    movk  x0, #0x8             // implements PSCI via HVC without firmware
    hvc   #0

msg:
    .ascii "hello world\n"
