//go:build arm64

package gf256

// MulAddRegion ARM64版(标量, NEON GF mul太复杂先用标量)
func MulAddRegion(dst, src []byte, coeff byte) {
	MulAdd(dst, src, coeff)
}
