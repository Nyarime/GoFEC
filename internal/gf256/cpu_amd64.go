//go:build amd64

package gf256

import "golang.org/x/sys/cpu"

var hasGFNI = cpu.X86.HasAVX512 && cpu.X86.HasAVX512VL

// GFNI affine变换矩阵: 每个GF(2^8)系数对应一个8x8位矩阵
var gfniMatrix [256]uint64

func init() {
	for c := 0; c < 256; c++ {
		var m uint64
		for bit := 0; bit < 8; bit++ {
			row := Mul(byte(c), 1<<bit)
			m |= uint64(row) << (bit * 8)
		}
		gfniMatrix[c] = m
	}
}
