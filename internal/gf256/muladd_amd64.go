//go:build amd64

package gf256

// mulAddAVX2 AVX2 VPSHUFB加速的GF(256) MulAdd
//go:noescape
func mulAddAVX2(dst, src []byte, loTable, hiTable *[16]byte)

// MulAddRegion dst ^= src * coeff (AVX2加速)
func MulAddRegion(dst, src []byte, coeff byte) {
	if coeff == 0 { return }
	if coeff == 1 {
		// 纯XOR
		for i := range dst {
			if i < len(src) { dst[i] ^= src[i] }
		}
		return
	}

	// 预计算split tables
	var loTable, hiTable [16]byte
	for i := byte(0); i < 16; i++ {
		loTable[i] = Mul(i, coeff)
		hiTable[i] = Mul(i<<4, coeff)
	}

	// AVX2主循环
	n := len(dst)
	if len(src) < n { n = len(src) }
	
	aligned := n &^ 31 // 32字节对齐
	if aligned > 0 {
		mulAddAVX2(dst[:aligned], src[:aligned], &loTable, &hiTable)
	}

	// 尾部标量
	for i := aligned; i < n; i++ {
		dst[i] ^= Mul(src[i], coeff)
	}
}
