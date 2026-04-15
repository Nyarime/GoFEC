// Package xor provides hardware-accelerated XOR operations.
// AVX2: 32 bytes/op, SSE2: 16 bytes/op, fallback: 8 bytes/op
package xor

// Bytes XORs src into dst: dst[i] ^= src[i]
// Uses AVX2 when available, falls back to uint64.
func Bytes(dst, src []byte)
