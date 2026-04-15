package raptorq

import (
	"fmt"
	"testing"
)

func TestPerformanceMatrix(t *testing.T) {
	configs := []struct{ K, T, repair int }{
		{8, 64, 4},       // 512B, 小数据
		{16, 128, 8},     // 2KB
		{32, 256, 16},    // 8KB
		{64, 512, 32},    // 32KB
		{128, 1024, 64},  // 128KB
	}

	for _, c := range configs {
		data := make([]byte, c.K*c.T)
		for i := range data { data[i] = byte(i) }

		codec := New(c.K, c.T)
		symbols := codec.Encode(data, c.repair)

		// 全部收到
		decoded, err := codec.Decode(symbols, len(data))
		status := "✅"
		if err != nil || len(decoded) != len(data) {
			status = "❌"
		}

		t.Logf("%s K=%d T=%d (%dKB) → %d符号 repair=%d",
			status, c.K, c.T, c.K*c.T/1024, len(symbols), c.repair)
	}
}

func TestLossRecoveryMatrix(t *testing.T) {
	codec := New(16, 128) // 2KB
	data := make([]byte, 16*128)
	for i := range data { data[i] = byte(i % 256) }

	for loss := 1; loss <= 6; loss++ {
		symbols := codec.Encode(data, loss*3) // 修复=丢失*2

		// 丢失前loss个源符号
		received := []Symbol{}
		for i, s := range symbols {
			if i >= loss { // 跳过前loss个
				received = append(received, s)
			}
		}

		decoded, err := codec.Decode(received, len(data))
		if err != nil {
			t.Logf("⚠️ 丢失%d块: 解码失败(%v)", loss, err)
		} else {
			match := true
			for i := range data {
				if i < len(decoded) && data[i] != decoded[i] { match = false; break }
			}
			if match {
				t.Logf("✅ 丢失%d块: 恢复成功", loss)
			} else {
				t.Logf("❌ 丢失%d块: 数据不匹配", loss)
			}
		}
	}
}

func BenchmarkRaptorQSizes(b *testing.B) {
	sizes := []struct{ name string; K, T int }{
		{"512B", 8, 64},
		{"2KB", 16, 128},
		{"8KB", 32, 256},
		{"32KB", 64, 512},
	}

	for _, s := range sizes {
		data := make([]byte, s.K*s.T)
		codec := New(s.K, s.T)

		b.Run(fmt.Sprintf("Encode_%s", s.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				codec.Encode(data, s.K/2)
			}
		})
	}
}
