# WAMR Go Environment for aolite
> [!NOTE]
> This is a customized WAMR environment for [aolite](https://github.com/everVision/aolite)

This project provides a customized WebAssembly Micro Runtime (WAMR) environment specifically designed for aolite, a blockchain platform. The environment is optimized for Go applications with particular focus on blockchain-related requirements.

## Features

Custom WAMR build with:
- Memory64 support enabled
- AOT (Ahead of Time) compilation enabled
- WASI support
- Hardware bound check disabled
- Exception handling enabled
- Tail call optimization enabled
- Memory profiling enabled

## Prerequisites

- Go 1.16 or later
- CMake 3.12 or later
- Make
- GCC/Clang

## Building

1. Clone the repository:
git clone <repository-url>
cd <repository-name>

2. Build WAMR:
```
// Regular build
make

// Debug build
make debug

// Clean build files:
make clean
```

3. Go environment setup
```
go mod tidy
```
4. Run tests:
```
go test -v
```

## Configuration

The WAMR runtime is configured with the following settings:

- WAMR_BUILD_MEMORY64=1: Enables 64-bit memory support
- WAMR_DISABLE_HW_BOUND_CHECK=1: Disables hardware boundary checks
- WAMR_BUILD_EXCE_HANDLING=1: Enables exception handling
- WAMR_BUILD_AOT=1: Enables AOT compilation
- WAMR_BUILD_INTERP=1: Enables interpreter mode
- WAMR_BUILD_JIT=0: Disables JIT compilation
- WAMR_BUILD_TAIL_CALL=1: Enables tail call optimization
- WAMR_BUILD_MEMORY_PROFILING=1: Enables memory profiling
- WAMR_BUILD_LIBC_WASI=1: Enables WASI support

## Development

To modify the WAMR configuration, edit the Makefile and adjust the CMake flags as needed.

## License

MIT License

Copyright (c) 2024 everVision

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.