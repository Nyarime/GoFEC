// Package raptorqgo provides pure Go implementations of
// RaptorQ (RFC 6330) and LDPC erasure codes.
//
// FEC Types:
//   - RaptorQ: rateless fountain code, unlimited repair symbols
//   - LDPC: low-density parity-check, high efficiency for large blocks
//
// Usage:
//   codec := raptorqgo.NewRaptorQ(numSource)
//   encoded := codec.Encode(data)
//   decoded := codec.Decode(partialBlocks)
package gofec

// Codec FEC编解码接口
type Codec interface {
	// Encode 编码数据为多个块(source + repair)
	Encode(data []byte) []Block
	// Decode 从部分块恢复数据
	Decode(blocks []Block, dataLen int) ([]byte, error)
	// Type 返回编码类型
	Type() string
}

// Block 编码块
type Block struct {
	ID   int64  // 块序号(source=0..N-1, repair=N..)
	Data []byte // 块数据
}
