package raptorq

import "unsafe"

func xorFast(dst, src []byte) {
	n := len(dst)
	if len(src) < n { n = len(src) }
	words := n / 8
	if words > 0 {
		dw := unsafe.Slice((*uint64)(unsafe.Pointer(&dst[0])), words)
		sw := unsafe.Slice((*uint64)(unsafe.Pointer(&src[0])), words)
		for i := range dw { dw[i] ^= sw[i] }
	}
	for i := words * 8; i < n; i++ { dst[i] ^= src[i] }
}
