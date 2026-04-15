#include "textflag.h"

// func mulAddAVX2(dst, src []byte, loTable, hiTable *[16]byte)
// AVX2 VPSHUFB: 32字节/次 GF(256) MulAdd
TEXT ·mulAddAVX2(SB), NOSPLIT, $0-56
	MOVQ	dst_base+0(FP), DI
	MOVQ	dst_len+8(FP), CX
	MOVQ	src_base+24(FP), SI
	MOVQ	src_len+32(FP), DX
	MOVQ	loTable+48(FP), R8
	MOVQ	hiTable+56(FP), R9

	// min(len)
	CMPQ	CX, DX
	CMOVQGT	DX, CX

	// 加载lo/hi table到YMM寄存器
	VMOVDQU	(R8), X14        // lo table (16字节)
	VINSERTI128 $1, X14, Y14, Y14  // broadcast到32字节
	VMOVDQU	(R9), X15
	VINSERTI128 $1, X15, Y15, Y15

	// nibble mask
	VPCMPEQB Y13, Y13, Y13   // all 1s
	VPSRLW	$4, Y13, Y13     // 0x0F0F...
	VPAND	Y13, Y13, Y13    // mask = 0x0F

loop32:
	CMPQ	CX, $32
	JL	tail

	VMOVDQU	(SI), Y0          // src 32B
	VPAND	Y0, Y13, Y1      // lo nibble
	VPSRLQ	$4, Y0, Y2
	VPAND	Y2, Y13, Y2      // hi nibble

	VPSHUFB	Y1, Y14, Y1      // lo lookup
	VPSHUFB	Y2, Y15, Y2      // hi lookup
	VPXOR	Y1, Y2, Y1       // gf_mul result

	VMOVDQU	(DI), Y3          // dst
	VPXOR	Y1, Y3, Y3       // dst ^= result
	VMOVDQU	Y3, (DI)

	ADDQ	$32, SI
	ADDQ	$32, DI
	SUBQ	$32, CX
	JMP	loop32

tail:
	VZEROUPPER
	// 剩余用标量处理
	TESTQ	CX, CX
	JZ	done
	RET  // Go层处理tail

done:
	RET
