// Package gofec provides pure Go forward error correction codecs:
// RaptorQ fountain codes, LDPC block codes and Leopard-RS (GF(2^16), up to
// 65536 shards). See subpackages for typed APIs.
package gofec

import (
	"github.com/nyarime/gofec/leopard"
	"github.com/nyarime/gofec/ldpc"
	"github.com/nyarime/gofec/raptorq"
)

// Codec is a generic FEC encoder/decoder over one payload buffer.
type Codec interface {
	Type() string
}

// NewRaptorQ creates a RaptorQ codec with K source symbols of T bytes each.
func NewRaptorQ(sourceSymbols, symbolSize int) *raptorq.Codec {
	return raptorq.New(sourceSymbols, symbolSize)
}

// NewLDPC creates an LDPC codec.
func NewLDPC(numData, numParity int, density float64) *ldpc.Codec {
	return ldpc.New(numData, numParity, density)
}

// NewLeopard creates a Leopard-RS codec for large payloads.
func NewLeopard(dataShards, parityShards int) (*leopard.Codec, error) {
	return leopard.NewCodec(dataShards, parityShards)
}

// NewLeopardForPayload sizes a Leopard codec to fit dataLen at percent redundancy.
func NewLeopardForPayload(dataLen, percent, shardSize int) (*leopard.Codec, error) {
	c, _, err := leopard.NewCodecForPayload(dataLen, percent, shardSize)
	return c, err
}
