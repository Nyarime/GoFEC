package raptorq

// Bitset 简易位图(比map[int]bool快10x+)
type Bitset struct {
	bits []uint64
	size int
}

func NewBitset(n int) *Bitset {
	return &Bitset{bits: make([]uint64, (n+63)/64), size: n}
}

func (b *Bitset) Set(i int)       { b.bits[i/64] |= 1 << uint(i%64) }
func (b *Bitset) Clear(i int)     { b.bits[i/64] &^= 1 << uint(i%64) }
func (b *Bitset) Has(i int) bool  { return b.bits[i/64]&(1<<uint(i%64)) != 0 }

func (b *Bitset) Count() int {
	n := 0
	for _, w := range b.bits {
		for w != 0 { n++; w &= w - 1 }
	}
	return n
}

func (b *Bitset) First() int {
	for i, w := range b.bits {
		if w != 0 {
			for bit := 0; bit < 64; bit++ {
				if w&(1<<uint(bit)) != 0 { return i*64 + bit }
			}
		}
	}
	return -1
}

func (b *Bitset) ForEach(f func(int)) {
	for i, w := range b.bits {
		for w != 0 {
			bit := 0
			for ; bit < 64 && w&(1<<uint(bit)) == 0; bit++ {}
			if bit < 64 {
				f(i*64 + bit)
				w &^= 1 << uint(bit)
			}
		}
	}
}

// XOR performs bitwise XOR of two bitsets (modifies dst)
func (b *Bitset) XOR(other *Bitset) {
	for i := 0; i < len(b.bits) && i < len(other.bits); i++ {
		b.bits[i] ^= other.bits[i]
	}
}
