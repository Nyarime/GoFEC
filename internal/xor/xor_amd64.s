#include "textflag.h"

// func Bytes(dst, src []byte)
TEXT ·Bytes(SB), NOSPLIT, $0-48
	MOVQ	dst_base+0(FP), DI    // dst pointer
	MOVQ	dst_len+8(FP), CX     // dst length
	MOVQ	src_base+24(FP), SI   // src pointer
	MOVQ	src_len+32(FP), DX    // src length

	// min(len(dst), len(src))
	CMPQ	CX, DX
	CMOVQGT	DX, CX

	// AVX2: 32字节/次
	CMPQ	CX, $32
	JL	sse2

avx2:
	VMOVDQU	(SI), Y0
	VPXOR	(DI), Y0, Y0
	VMOVDQU	Y0, (DI)
	ADDQ	$32, SI
	ADDQ	$32, DI
	SUBQ	$32, CX
	CMPQ	CX, $32
	JGE	avx2
	VZEROUPPER

sse2:
	// SSE2: 16字节/次
	CMPQ	CX, $16
	JL	word

	MOVOU	(SI), X0
	MOVOU	(DI), X1
	PXOR	X0, X1
	MOVOU	X1, (DI)
	ADDQ	$16, SI
	ADDQ	$16, DI
	SUBQ	$16, CX
	JMP	sse2

word:
	// uint64: 8字节/次
	CMPQ	CX, $8
	JL	tail

	MOVQ	(SI), AX
	XORQ	(DI), AX
	MOVQ	AX, (DI)
	ADDQ	$8, SI
	ADDQ	$8, DI
	SUBQ	$8, CX
	JMP	word

tail:
	// 逐字节
	TESTQ	CX, CX
	JZ	done
	MOVB	(SI), AL
	XORB	(DI), AL
	MOVB	AL, (DI)
	INCQ	SI
	INCQ	DI
	DECQ	CX
	JMP	tail

done:
	RET
