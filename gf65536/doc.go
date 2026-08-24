// Package gf65536 implements GF(2^16) finite field arithmetic with
// log/exp tables and SIMD-accelerated region multiply on amd64.
//
// This field uses irreducible polynomial 0x1002B. Leopard-RS in this module
// uses a separate Cantor-basis field (0x1002D) internally.
package gf65536
