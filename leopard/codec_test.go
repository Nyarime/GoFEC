package leopard

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestCodecPayloadRoundtrip(t *testing.T) {
	data := make([]byte, 512*1024)
	rand.Read(data)

	c, _, err := NewCodecForPayload(len(data), 25, 0)
	if err != nil {
		t.Fatal(err)
	}
	shards, err := c.Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	present := make([]bool, len(shards))
	for i := range shards {
		present[i] = true
	}
	// Erase two data shards
	shards[1] = nil
	shards[3] = nil
	present[1] = false
	present[3] = false

	got, err := c.Decode(shards, present, len(data))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("decode mismatch")
	}
}

func TestEncodeParityOnly(t *testing.T) {
	data := []byte("leopard parity only test payload")
	c, err := NewCodec(4, 2)
	if err != nil {
		t.Fatal(err)
	}
	parity, err := c.EncodeParity(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(parity) != 2 {
		t.Fatalf("parity shards = %d", len(parity))
	}
}
