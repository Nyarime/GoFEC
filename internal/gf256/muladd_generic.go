//go:build !amd64

package gf256

func MulAddRegion(dst, src []byte, coeff byte) {
	MulAdd(dst, src, coeff)
}
