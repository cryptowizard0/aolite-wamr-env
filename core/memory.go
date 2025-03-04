package core

import (
	"errors"
	"unsafe"
)

// ReadMemory 从WASM内存中读取数据
func (c *Context) ReadMemory(offset uint64, size uint64) ([]byte, error) {
	if !c.Initialized {
		return nil, errors.New("runtime not initialized")
	}

	// 验证地址
	if !c.Instance.ValidateAppAddr(offset, size) {
		return nil, errors.New("invalid memory address or size")
	}

	// 获取本地内存地址
	nativeAddr := c.Instance.AddrAppToNative(offset)
	if nativeAddr == nil {
		return nil, errors.New("failed to get native address")
	}

	// 复制内存数据
	buffer := make([]byte, size)
	copy(buffer, (*[1 << 30]byte)(unsafe.Pointer(nativeAddr))[:size:size])

	return buffer, nil
}

// WriteMemory 写入数据到WASM内存
func (c *Context) WriteMemory(offset uint64, data []byte) error {
	if !c.Initialized {
		return errors.New("runtime not initialized")
	}

	size := uint64(len(data))
	if size == 0 {
		return nil
	}

	// 验证地址
	if !c.Instance.ValidateAppAddr(offset, size) {
		return errors.New("invalid memory address or size")
	}

	// 获取本地内存地址
	nativeAddr := c.Instance.AddrAppToNative(offset)
	if nativeAddr == nil {
		return errors.New("failed to get native address")
	}

	// 写入内存数据
	copy((*[1 << 30]byte)(unsafe.Pointer(nativeAddr))[:size:size], data)

	return nil
}

// GetMemorySize 获取WASM内存大小
func (c *Context) GetMemorySize() uint64 {
	if !c.Initialized {
		return 0
	}

	// Get memory range to determine size
	_, start, end := c.Instance.GetAppAddrRange(0)
	return end - start
}
