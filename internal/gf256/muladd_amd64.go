//go:build amd64

package gf256

//go:noescape
func mulAddAVX2(dst, src []byte, loTable, hiTable *[16]byte)

//go:noescape
func mulAddGFNI(dst, src []byte, matrix uint64)

// MulAddRegion dst ^= src * coeff
// 自动选择: GFNI(64B/op) > AVX2(32B/op) > scalar
func MulAddRegion(dst, src []byte, coeff byte) {
	if coeff == 0 { return }

	n := len(dst)
	if len(src) < n { n = len(src) }

	if coeff == 1 {
		// 纯XOR, 走xor包更快
		for i := 0; i < n; i++ { dst[i] ^= src[i] }
		return
	}

	if hasGFNI {
		// AVX512 GFNI: 64字节/次, 单指令GF乘法
		aligned := n &^ 63
		if aligned > 0 {
			mulAddGFNI(dst[:aligned], src[:aligned], gfniMatrix[coeff])
		}
		for i := aligned; i < n; i++ {
			dst[i] ^= Mul(src[i], coeff)
		}
		return
	}

	// AVX2 VPSHUFB: 32字节/次
	aligned := n &^ 31
	if aligned > 0 {
		mulAddAVX2(dst[:aligned], src[:aligned], &mulLo[coeff], &mulHi[coeff])
	}
	for i := aligned; i < n; i++ {
		dst[i] ^= Mul(src[i], coeff)
	}
}
