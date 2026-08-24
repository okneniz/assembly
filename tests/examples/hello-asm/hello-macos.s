// hello world on direct macOS syscalls (arm64), position-independent code.
//
// How it works: libc is not called - the program talks to the kernel
// directly through the svc trap (supervisor call). The Darwin/BSD arm64
// convention differs from Linux in two things:
//
//   x16        - the system call number (on Linux that is x8!)
//   0x2000000  - the SYSCALL_CLASS_UNIX class prefix (2 << 24), which BSD
//                packs into the high bits of the number: the full number of
//                write = 0x2000004, exit = 0x2000001
//
// Because of the prefix the number does not fit into the imm16 of a movz
// instruction, so it is assembled from two halves: movz stores the high part
// (shifted), movk completes the low part - see movz/movk below.
//
//   svc #0x80 - the canonical Darwin trap immediate (dispatching goes by
//              x16, so #0 works too - but real macOS binaries write 0x80)
//
//   fd 1  = stdout, 12 = length of "hello world\n" (the kernel writes the
//   bytes verbatim)

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
