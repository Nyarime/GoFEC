// Package ldpc implements LDPC (Low-Density Parity-Check) erasure codes.
//
// LDPC codes use sparse parity-check matrices for efficient
// encoding and iterative belief-propagation decoding.
package ldpc

import (
	"errors"
)

// Codec LDPC编解码器
type Codec struct {
	numData   int     // 数据块数
	numParity int     // 校验块数
	density   float64 // 稀疏矩阵密度(0.0-1.0)
	matrix    [][]bool // 校验矩阵
}

// New 创建LDPC编解码器
// numData: 数据块数, numParity: 校验块数, density: 0.1-0.3推荐
func New(numData, numParity int, density float64) *Codec {
	c := &Codec{
		numData:   numData,
		numParity: numParity,
		density:   density,
	}
	c.generateMatrix()
	return c
}

// generateMatrix 生成随机稀疏校验矩阵
func (c *Codec) generateMatrix() {
	// v1.1: PEG构造(最大化girth，BP收敛更快)
	dv := int(c.density * float64(c.numParity)) + 1
	if dv < 2 { dv = 2 }
	if dv > c.numParity { dv = c.numParity }
	c.matrix = BuildPEGMatrix(c.numData, c.numParity, dv)
}

// Encode 编码: 数据块→数据块+校验块
func (c *Codec) Encode(shards [][]byte) [][]byte {
	if len(shards) != c.numData { return nil }
	
	shardSize := len(shards[0])
	result := make([][]byte, c.numData+c.numParity)
	
	// 复制数据块
	for i := 0; i < c.numData; i++ {
		result[i] = make([]byte, shardSize)
		copy(result[i], shards[i])
	}
	
	// 生成校验块(XOR)
	for i := 0; i < c.numParity; i++ {
		result[c.numData+i] = make([]byte, shardSize)
		for j := 0; j < c.numData; j++ {
			if c.matrix[i][j] {
				xorBytes(result[c.numData+i], shards[j])
			}
		}
	}
	
	return result
}

// Decode 迭代信念传播解码
func (c *Codec) Decode(shards [][]byte) error {
	if len(shards) != c.numData+c.numParity {
		return errors.New("wrong shard count")
	}
	
	// 找丢失的块
	missing := []int{}
	for i, s := range shards {
		if s == nil { missing = append(missing, i) }
	}
	if len(missing) == 0 { return nil }
	if len(missing) > c.numParity { return errors.New("too many missing") }
	
	// 迭代BP解码(简化版)
	for iter := 0; iter < 50; iter++ {
		recovered := false
		for _, pi := range missing {
			if pi >= c.numData { continue } // 只恢复数据块
			// 找包含此数据块的校验方程
			for ci := 0; ci < c.numParity; ci++ {
				if !c.matrix[ci][pi] { continue }
				// 检查该方程是否只缺这一个
				unknowns := 0
				for _, mi := range missing {
					if c.matrix[ci][mi] { unknowns++ }
				}
				if unknowns == 1 && shards[c.numData+ci] != nil {
					// 可以恢复！
					shards[pi] = make([]byte, len(shards[c.numData+ci]))
					copy(shards[pi], shards[c.numData+ci])
					for j := 0; j < c.numData; j++ {
						if j != pi && c.matrix[ci][j] && shards[j] != nil {
							xorBytes(shards[pi], shards[j])
						}
					}
					recovered = true
				}
			}
		}
		if !recovered { break }
		// 更新missing
		missing = missing[:0]
		for i, s := range shards {
			if s == nil { missing = append(missing, i) }
		}
		if len(missing) == 0 { return nil }
	}
	
	if len(missing) > 0 { return errors.New("decode failed") }
	return nil
}

func xorBytes(dst, src []byte) {
	for i := range dst {
		if i < len(src) { dst[i] ^= src[i] }
	}
}

func (c *Codec) Type() string { return "ldpc" }
