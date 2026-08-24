//go:build amd64 || arm64

package xor

// Bytes XORs src into dst: dst[i] ^= src[i]
// Implemented in xor_amd64.s and xor_arm64.s; xor_generic.go covers every
// other architecture.
func Bytes(dst, src []byte)
