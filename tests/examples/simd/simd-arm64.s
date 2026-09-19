.text
and.16b v0, v1, v2
and.8b v3, v4, v5
bic.16b v0, v1, v2
orr.16b v0, v1, v2
orr.8b v3, v4, v5
orn.16b v0, v1, v2
eor.16b v0, v1, v2
bsl.16b v0, v1, v2
bit.16b v0, v1, v2
bif.16b v0, v1, v2
add.16b v0, v1, v2
add.8b v3, v4, v5
cmeq.16b v0, v1, v2
addp.16b v0, v1, v2
sqrshl.16b v0, v1, v2
cnt.8b v0, v1
cnt.16b v0, v1
rev32.8b v0, v1
not.16b v0, v1
abs.8b v0, v1
rbit.16b v0, v1
shl.16b v0, v1, #3
shl.8b v3, v4, #5
sri.16b v0, v1, #4
ushr.16b v0, v1, #5
sshr.16b v0, v1, #5
aese v0.16b, v1.16b
aesmc v0.16b, v1.16b
dup.16b v0, w0
dup.4s v0, v1[2]
mov.s v0[1], w1
smov x0, v0.s[1]
umov w0, v0.s[1]
mov.16b v0, v1
mov.8b v2, v3
tbl.16b v0, { v1 }, v2
saddw2.8h v0, v1, v2
uaddw2.4s v0, v1, v2
usubw2.8h v0, v1, v2
mla.4s v0, v1, v2[1]
mls.4s v0, v1, v2[1]
mul.4s v0, v1, v2[1]
sqdmulh.8h v0, v1, v2[3]
sqrdmulh.4s v0, v1, v2[0]
sqrdmlah.4s v0, v1, v2[0]
sqrdmlsh.4s v0, v1, v2[0]
smlal2.4s v0, v1, v2[1]
fmla.4s v0, v1, v2[1]
fmls.2s v0, v1, v2[0]
fmul.2d v0, v1, v2[0]
fmulx.4s v0, v1, v2[2]
fcmla.4s v0, v1, v2[0], #90
ins.s v0[1], v1[2]
