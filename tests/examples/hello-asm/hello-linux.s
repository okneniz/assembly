// hello world on direct Linux syscalls (arm64), position-independent code.
//
// How it works: libc is not called - the program talks to the kernel
// directly through the svc trap (supervisor call). The Linux arm64
// convention:
//
//   x0-x5   - call arguments (like a regular function)
//   x8      - the system call number: the kernel dispatcher looks here
//   svc #0  - switch to kernel mode; the result comes back in x0
//
// The "magic" constants are just numbers from the kernel tables:
//   fd 1  = stdout (0=stdin, 2=stderr; the descriptors are opened by the loader already)
//   64    = __NR_write, 93 = __NR_exit  (asm-generic/unistd.h)
//   12    = length of "hello world\n": the kernel writes the bytes verbatim,
//           without NUL termination - the length is passed explicitly

start:
    mov     x0, #1                  // write(fd=stdout, ...)
    adr     x1, msg                 // ... buf - string address (pc-relative)
    mov     x2, #12                 // ... len = 12
    mov     x8, #64                 // __NR_write
    svc     #0

    mov     x0, #0                  // exit(return code = 0)
    mov     x8, #93                 // __NR_exit
    svc     #0

msg:
    .ascii "hello world\n"
