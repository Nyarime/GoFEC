package raptorq

// GenerateRepair produces additional fountain repair symbols from known source
// shards without re-encoding the full block. startESI is typically K+existingRepair.
func (c *Codec) GenerateRepair(source [][]byte, startESI uint32, count int) []Symbol {
	if count <= 0 {
		return nil
	}
	out := make([]Symbol, count)
	for i := 0; i < count; i++ {
		esi := startESI + uint32(i)
		out[i] = Symbol{
			ESI:  esi,
			Data: c.ltEncode(source, esi),
		}
	}
	return out
}

// GenerateRepairFromData splits data into K source symbols and generates count
// repair symbols starting at ESI startESI.
func (c *Codec) GenerateRepairFromData(data []byte, startESI uint32, count int) []Symbol {
	K := c.sourceSymbols
	T := c.symbolSize
	source := splitData(data, K, T)
	return c.GenerateRepair(source, startESI, count)
}
