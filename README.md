# RaptorQ-Go

纯Go实现的FEC纠删码库。

## 编码类型

| 类型 | 状态 | 说明 |
|------|------|------|
| LDPC | ✅ v1.0 | 低密度校验码，迭代BP解码 |
| RaptorQ | 🔧 开发中 | RFC 6330喷泉码 |

## 使用

```go
import "github.com/nyarime/raptorq-go/ldpc"

codec := ldpc.New(10, 4, 0.3) // 10数据+4校验, 密度0.3
encoded := codec.Encode(dataShards)
// 丢失部分块后
codec.Decode(partialShards) // 迭代恢复
```

## 被引用

- [NRUP](https://github.com/Nyarime/NRUP) - UDP传输协议的FEC层
