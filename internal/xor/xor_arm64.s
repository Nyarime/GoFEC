#include "textflag.h"

// func Bytes(dst, src []byte)
// ARM64 NEON: 16字节/次 EOR
TEXT ·Bytes(SB), NOSPLIT, $0-48
	MOVD	dst_base+0(FP), R0
	MOVD	dst_len+8(FP), R2
	MOVD	src_base+24(FP), R1
	MOVD	src_len+32(FP), R3

	// min(len)
	CMP	R3, R2
	CSEL	LT, R2, R3, R2

	// NEON: 16字节/次
neon16:
	CMP	$16, R2
	BLT	word8

	VLD1	(R1), [V0.B16]
	VLD1	(R0), [V1.B16]
	VEOR	V0.B16, V1.B16, V1.B16
	VST1	[V1.B16], (R0)
	ADD	$16, R0
	ADD	$16, R1
	SUB	$16, R2
	B	neon16

word8:
	CMP	$8, R2
	BLT	tail

	MOVD	(R1), R4
	MOVD	(R0), R5
	EOR	R4, R5, R5
	MOVD	R5, (R0)
	ADD	$8, R0
	ADD	$8, R1
	SUB	$8, R2
	B	word8

tail:
	CBZ	R2, done
	MOVBU	(R1), R4
	MOVBU	(R0), R5
	EOR	R4, R5, R5
	MOVBU	R5, (R0)
	ADD	$1, R0
	ADD	$1, R1
	SUB	$1, R2
	B	tail

done:
	RET
