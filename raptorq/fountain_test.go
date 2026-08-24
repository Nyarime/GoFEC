package raptorq

import (
	"bytes"
	"testing"
)

func TestGenerateRepairFountain(t *testing.T) {
	K, T := 8, 64
	data := make([]byte, K*T)
	for i := range data {
		data[i] = byte(i)
	}
	c := New(K, T)
	base := c.Encode(data, 2)
	extra := c.GenerateRepairFromData(data, uint32(K+2), 4)

	received := append(base[:K], base[K:]...)
	received = append(received, extra...)
	// Drop source symbol 2
	filtered := make([]Symbol, 0, len(received)-1)
	for _, s := range received {
		if s.ESI != 2 {
			filtered = append(filtered, s)
		}
	}
	got, err := c.Decode(filtered, len(data))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("fountain repair decode mismatch")
	}
}
