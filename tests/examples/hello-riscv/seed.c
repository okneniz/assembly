// Freestanding snippets exercising RV64IMAFD (integer + mul/div + atomics +
// float). Compiled with -march=rv64imafd (no compressed) so every instruction
// is 32-bit.
int sum(int n) {
    int s = 0;
    for (int i = 1; i <= n; i++) s += i;
    return s;
}
int mulop(int a, int b) { return a * b; }
int divop(int a, int b) { return b ? a / b : 0; }
long mulhop(long a, long b) { return a * b; }
float fadd(float a, float b) { return a + b; }
float fmul(float a, float b) { return a * b; }
double dadd(double a, double b) { return a + b; }
double dmul(double a, double b) { return a * b; }
int atomicop(int *p, int v) {
    int old;
    __asm__ volatile("amoadd.w %0, %1, (%2)" : "=r"(old) : "r"(v), "r"(p) : "memory");
    return old;
}
static int arr[64];
int loop(int n) {
    int s = 0;
    for (int i = 0; i < n; i++) { arr[i] = i * (i + 1); s += arr[i] & 0x7f; }
    return s;
}
int recurse(int n) { if (n <= 1) return 1; return recurse(n - 1) + recurse(n - 2); }
int stacky(int a, int b, int c, int d, int e, int f) {
    int w = a + 1, x = b + 2, y = c + 3, z = d + 4, p = e + 5, q = f + 6;
    return (w * x - y * z) + (p * q) - (w + q);
}
unsigned int read_fcsr(void) {
    unsigned int v;
    __asm__ volatile("csrr %0, fcsr" : "=r"(v));
    return v;
}
void write_fflags(unsigned int v) {
    __asm__ volatile("csrw fflags, %0" :: "r"(v));
}
