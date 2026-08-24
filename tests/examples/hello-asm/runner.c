// A minimal loader for assembly-asm binaries: mmaps the whole file and
// transfers control to its start. The code in hello-*.s is
// position-independent (adr pc-relative), so any address will do.
//
// macOS (Apple Silicon): exec pages are signed, so anonymous mapped memory
// requires MAP_JIT + icache invalidation.
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

#ifdef __APPLE__
#include <libkern/OSCacheControl.h>
#include <pthread.h>
#define JIT_FLAGS MAP_JIT
#else
#define JIT_FLAGS 0
#endif

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "usage: runner FILE.bin\n");
        return 2;
    }

    FILE *f = fopen(argv[1], "rb");
    if (!f) {
        perror("open");
        return 1;
    }
    fseek(f, 0, SEEK_END);
    long sz = ftell(f);
    fseek(f, 0, SEEK_SET);

    void *mem = mmap(NULL, sz, PROT_READ | PROT_WRITE | PROT_EXEC,
                     MAP_PRIVATE | MAP_ANONYMOUS | JIT_FLAGS, -1, 0);
    if (mem == MAP_FAILED) {
        perror("mmap");
        return 1;
    }
#ifdef __APPLE__
    // Apple Silicon: MAP_JIT pages are protected W^X per-thread - writing
    // is enabled by a toggle; before execution we turn the protection back
    // on.
    pthread_jit_write_protect_np(0);
#endif
    if (fread(mem, 1, sz, f) != (size_t)sz) {
        perror("read");
        return 1;
    }
    fclose(f);

#ifdef __APPLE__
    pthread_jit_write_protect_np(1);
    sys_icache_invalidate(mem, sz);
#else
    __builtin___clear_cache((char *)mem, (char *)mem + sz);
#endif

    ((void (*)(void))mem)();
    return 0; // unreachable: the program terminates with the exit syscall
}
