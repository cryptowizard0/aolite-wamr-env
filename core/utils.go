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
