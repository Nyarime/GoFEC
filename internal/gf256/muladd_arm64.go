//go:build arm64

package gf256

//go:noescape
func mulAddNEON(dst, src []byte, loTable, hiTable *[16]byte)

// MulAddRegion dst ^= src * coeff
// 自动选择: SVE(可变长) > NEON(16B/op) > scalar
func MulAddRegion(dst, src []byte, coeff byte) {
	if coeff == 0 { return }

	n := len(dst)
	if len(src) < n { n = len(src) }

	if coeff == 1 {
		for i := 0; i < n; i++ { dst[i] ^= src[i] }
		return
	}

	// TODO: hasSVE时走SVE路径(更宽向量)

	// NEON VTBL: 16字节/次
	aligned := n &^ 15
	if aligned > 0 {
		mulAddNEON(dst[:aligned], src[:aligned], &mulLo[coeff], &mulHi[coeff])
	}
	for i := aligned; i < n; i++ {
		dst[i] ^= Mul(src[i], coeff)
	}
}
