package leopard

import (
	"errors"
	"fmt"
)

const (
	// MaxShards is the maximum data+parity shard count (65536).
	MaxShards = kOrder
	// DefaultShardSize is the default even shard size in bytes.
	DefaultShardSize = 64 * 1024
)

// Codec is a high-level Leopard-RS encoder/decoder over a single payload.
type Codec struct {
	enc          *Encoder
	dataShards   int
	parityShards int
	shardSize    int
}

// NewCodec creates a Leopard-RS codec.
// dataShards + parityShards must be <= MaxShards.
func NewCodec(dataShards, parityShards int) (*Codec, error) {
	enc, err := New(dataShards, parityShards)
	if err != nil {
		return nil, err
	}
	return &Codec{
		enc:          enc,
		dataShards:   dataShards,
		parityShards: parityShards,
	}, nil
}

// PlanShards picks data/parity shard counts for a payload length and redundancy
// percentage. shardSize must be even; when 0, DefaultShardSize is used.
func PlanShards(dataLen, percent, shardSize int) (dataShards, parityShards int, err error) {
	if dataLen < 1 {
		return 0, 0, errors.New("leopard: empty payload")
	}
	if percent < 1 {
		percent = 1
	}
	if shardSize <= 0 {
		shardSize = DefaultShardSize
	}
	if shardSize%2 != 0 {
		shardSize++
	}
	dataShards = (dataLen + shardSize - 1) / shardSize
	parityShards = dataShards * percent / 100
	if parityShards < 1 {
		parityShards = 1
	}
	if dataShards+parityShards > MaxShards {
		return 0, 0, fmt.Errorf("leopard: %d+%d shards exceeds %d", dataShards, parityShards, MaxShards)
	}
	return dataShards, parityShards, nil
}

// NewCodecForPayload builds a codec sized to dataLen with percent redundancy.
func NewCodecForPayload(dataLen, percent, shardSize int) (*Codec, int, error) {
	d, p, err := PlanShards(dataLen, percent, shardSize)
	if err != nil {
		return nil, 0, err
	}
	c, err := NewCodec(d, p)
	if err != nil {
		return nil, 0, err
	}
	if shardSize <= 0 {
		shardSize = DefaultShardSize
	}
	if shardSize%2 != 0 {
		shardSize++
	}
	need := ((dataLen + d - 1) / d)
	if need%2 != 0 {
		need++
	}
	if need > shardSize {
		shardSize = need
	}
	c.shardSize = shardSize
	return c, shardSize, nil
}

// DataShards returns the number of data shards.
func (c *Codec) DataShards() int { return c.dataShards }

// ParityShards returns the number of parity shards.
func (c *Codec) ParityShards() int { return c.parityShards }

// ShardSize returns the even shard size used by the last Encode call, or 0.
func (c *Codec) ShardSize() int { return c.shardSize }

// Encode splits data into data shards and generates parity shards.
func (c *Codec) Encode(data []byte) ([][]byte, error) {
	shardSize := c.shardSize
	if shardSize <= 0 {
		shardSize = ((len(data) + c.dataShards - 1) / c.dataShards)
		if shardSize%2 != 0 {
			shardSize++
		}
		if shardSize < 2 {
			shardSize = 2
		}
	}
	c.shardSize = shardSize

	total := c.dataShards + c.parityShards
	shards := make([][]byte, total)
	paddedLen := c.dataShards * shardSize
	padded := make([]byte, paddedLen)
	copy(padded, data)

	for i := 0; i < c.dataShards; i++ {
		shards[i] = make([]byte, shardSize)
		copy(shards[i], padded[i*shardSize:(i+1)*shardSize])
	}
	for i := c.dataShards; i < total; i++ {
		shards[i] = make([]byte, shardSize)
	}

	if err := c.enc.Encode(shards); err != nil {
		return nil, err
	}
	return shards, nil
}

// EncodeParity returns only parity shards (systematic: data stays in the payload).
func (c *Codec) EncodeParity(data []byte) ([][]byte, error) {
	all, err := c.Encode(data)
	if err != nil {
		return nil, err
	}
	return all[c.dataShards:], nil
}

// Decode recovers missing data shards. shards[i] may be nil when missing;
// present[i] must be true when shards[i] is non-nil.
func (c *Codec) Decode(shards [][]byte, present []bool, dataLen int) ([]byte, error) {
	if err := c.enc.Decode(shards, present); err != nil {
		return nil, err
	}
	out := make([]byte, 0, dataLen)
	for i := 0; i < c.dataShards; i++ {
		if shards[i] == nil {
			return nil, errors.New("leopard: data shard not recovered")
		}
		out = append(out, shards[i]...)
	}
	if len(out) > dataLen {
		out = out[:dataLen]
	}
	return out, nil
}

// Type returns the codec name.
func (c *Codec) Type() string { return "leopard" }
