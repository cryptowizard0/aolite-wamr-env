package core

import "errors"

// 返回 WASM 模块中导出项的数量
func (c *Context) GetExportCount() (int32, error) {
	if !c.Initialized || c.Module == nil {
		return 0, errors.New("module not initialized")
	}
	return c.Module.GetExportCount(), nil
}

// 返回指定索引的导出项类型信息
func (c *Context) GetExportType(index int32) (*ExportType, error) {
	if !c.Initialized || c.Module == nil {
		return nil, errors.New("module not initialized")
	}

	exportType := c.Module.GetExportType(index)
	return &ExportType{
		Kind: exportType.Kind,
		Name: exportType.Name,
	}, nil
}

// ... existing code ...

// uint64ToBytes 将 uint64 转换为字节数组
func uint64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(val)
	b[1] = byte(val >> 8)
	b[2] = byte(val >> 16)
	b[3] = byte(val >> 24)
	b[4] = byte(val >> 32)
	b[5] = byte(val >> 40)
	b[6] = byte(val >> 48)
	b[7] = byte(val >> 56)
	return b
}
