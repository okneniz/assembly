.text
fadd d0, d1, d2
fadd s0, s1, s2
fsub d0, d1, d2
fsub s0, s1, s2
fmul d0, d1, d2
fmul s0, s1, s2
fdiv d0, d1, d2
fdiv s0, s1, s2
fmax d0, d1, d2
fmax s0, s1, s2
fmin d0, d1, d2
fmin s0, s1, s2
fneg d3, d4
fneg s3, s4
fmov d5, d6
fmov s5, s6
fmov d7, x0
fmov x1, d8
fmov s9, w2
fmov w3, s10
fmov d11, #0.5
fmov s12, #0.5
fmov d13, #1.5
fmov s14, #1.5
fcmp d15, d16
fcmp s15, s16
fcmp d17, #0.0
fcmp s17, #0.0
fcvt s18, d19
fcvt d18, s19
scvtf d20, w21
scvtf d20, x21
scvtf s20, w21
scvtf s20, x21
ucvtf d22, w21
ucvtf d22, x21
ucvtf s22, w21
ucvtf s22, x21
fcvtzs w23, d24
fcvtzs x23, d24
fcvtzs w23, s24
fcvtzs x23, s24
fcvtzu w25, d26
fcvtzu x25, d26
fcvtzu w25, s26
fcvtzu x25, s26
fmadd d27, d28, d29, d30
fmadd s27, s28, s29, s30
fnmsub d27, d28, d29, d30
fnmsub s27, s28, s29, s30
ldr d0, [x1, #8]
ldr s0, [x1, #4]
str d2, [x1, #16]
str s2, [x1, #8]
