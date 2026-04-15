# GoFEC 可行性研究

## 1. GF(2^16) 大有限域支持

### 现状
- GoFEC当前: GF(2^8) = 256个元素, 最多255个数据分片
- klauspost: GF(2^8) + Leopard GF(2^16), 最多65536分片

### GF(2^16)优势
- 支持 **65536个分片** (vs GF(2^8)的255)
- 大规模分布式存储必需 (如IPFS, Filecoin)
- 更大的编码空间

### 实现难度
- **中等**: 需要16位查表 (65536×65536太大) → 用split-table或log/exp
- Log/Exp表: 131KB (vs GF(2^8)的512B)
- 乘法: 需要GF(2^16)的primitive polynomial (如x^16+x^12+x^3+x+1)
- SIMD: AVX2可以处理16位元素 (VPMULLW)

### 建议
- **v2.0.0**: 作为新的codec类型 `gf16`
- 参考: klauspost的leopard16实现, malaire/reed-solomon-16 (Rust)

---

## 2. AVX512 GFNI 加速

### 什么是GFNI?
- **Galois Field New Instructions** (Intel Ice Lake+)
- 原生GF(2^8)乘法指令: `GF2P8MULB`
- 不再需要VPSHUFB查表 → 直接硬件乘法

### 性能预期
- klauspost实测: **GFNI比AVX2快3倍** (3x)
- GF(256) MulAdd: 11.5 GB/s → **~35 GB/s**
- 整体RaptorQ编码提速: ~2-3倍

### CPU支持
- Intel: Ice Lake (2019), Tiger Lake, Alder Lake, Sapphire Rapids
- AMD: Zen 4 (2022), Zen 5
- 2026年大多数新服务器都支持

### 实现难度
- **低**: 只需在gf256/mulAdd_amd64.s加一个GFNI路径
- 运行时CPU特性检测 (CPUID)
- klauspost的实现可直接参考

### 代码框架
```asm
// AVX512 GFNI路径
TEXT ·mulAddGFNI(SB), NOSPLIT, $0
  VPBROADCASTB  coeff, Z15
  // 循环
loop:
  VMOVDQU64  (SI), Z0
  VGF2P8MULB Z15, Z0, Z1   // ← 单指令GF(2^8)乘法!
  VPXORQ     (DI), Z1, Z1
  VMOVDQU64  Z1, (DI)
  // ...
```

### 建议
- **v1.x.0**: 添加GFNI路径 (运行时检测)
- 需要Go 1.21+(支持AVX512汇编)

---

## 3. Leopard-RS 高分片支持

### 什么是Leopard?
- O(n log n)复杂度的Reed-Solomon变体
- 基于FFT(快速沃尔什-哈达马变换, FWHT)
- klauspost在v1.12.0引入

### 两种模式
| 模式 | 有限域 | 最大分片 | 复杂度 |
|------|--------|----------|--------|
| Leopard8 | GF(2^8) | 256 | O(n log n) |
| Leopard16 | GF(2^16) | 65536 | O(n log n) |

### vs 标准RS/RaptorQ
| 特性 | 标准RS | Leopard | RaptorQ |
|------|--------|---------|---------|
| 复杂度 | O(n²) | O(n log n) | O(n) |
| 最大分片 | 255 | 65536 | 无限 |
| 系统码 | ✅ | ✅ | ✅ |
| 无速率 | ❌ | ❌ | ✅ |

### 实现难度
- **高**: FFT/FWHT核心算法复杂
- 需要GF(2^16)运算
- klauspost的leopard.go: ~2000行

### 建议
- **v2.0.0**: 作为高级特性
- 先做好GF(2^16) → 再做FWHT → 再做Leopard
- 或者直接集成klauspost的leopard作为可选后端

---

## 路线图

| 版本 | 特性 | 优先级 |
|------|------|--------|
| v1.0 | LDPC + RaptorQ + AVX2 + NEON | ✅ 已完成 |
| v1.x | AVX512 GFNI | 高(3倍GF提速) |
| v2.0 | GF(2^16) + Leopard16 | 中(65536分片) |
| v2.x | FWHT加速 | 低 |

## 参考
- klauspost/reedsolomon #327: GF(2^16) Mul
- klauspost/reedsolomon #320: AVX512 GFNI
- klauspost/reedsolomon #272: leopard16优化
- malaire/reed-solomon-16: Rust GF(2^16)
- FOSDEM 2026: "Demystifying Erasure Coding"
