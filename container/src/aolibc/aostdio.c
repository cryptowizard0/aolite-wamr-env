#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
// #include <emscripten.h>
#include <wasi/api.h> 

#if 0
#define AO_LOG(...) fprintf(stderr, __VA_ARGS__)
#else
#define AO_LOG(...)
#endif

// WeaveDrive async wrapper functions. These allow us to call the WeaveDrive
// async JS code from C.
// EM_ASYNC_JS(int, weavedrive_open, (const char* c_filename, const char* mode), {
//     const filename = UTF8ToString(Number(c_filename));
//     if (!Module.WeaveDrive) {
//         return Promise.resolve(null)
//     }

//     const drive = Module.WeaveDrive(Module, FS);
//     return await drive.open(filename);
// });

// EM_ASYNC_JS(int, weavedrive_read, (int fd, int *dst_ptr, size_t length), {
//     const drive = Module.WeaveDrive(Module, FS);
//     return Promise.resolve(await drive.read(fd, dst_ptr, length));
// });

// EM_ASYNC_JS(int, weavedrive_close, (int fd), {
//   const drive = Module.WeaveDrive(Module, FS);
//   return drive.close(fd);
// });


FILE* fopen(const char* filename, const char* mode) {
    AO_LOG("AO: Called fopen: %s, %s\n", filename, mode);
    
    __wasi_oflags_t oflags = 0;
    __wasi_rights_t fs_rights_base = 0;
    
    if (strchr(mode, 'r')) {
        fs_rights_base |= __WASI_RIGHTS_FD_READ | __WASI_RIGHTS_FD_SEEK;
    }
    if (strchr(mode, 'w')) {
        oflags |= __WASI_OFLAGS_CREAT | __WASI_OFLAGS_TRUNC;
        fs_rights_base |= __WASI_RIGHTS_FD_WRITE | __WASI_RIGHTS_FD_SEEK;
    }
    if (strchr(mode, 'a')) {
        oflags |= __WASI_OFLAGS_CREAT;
        fs_rights_base |= __WASI_RIGHTS_FD_WRITE | __WASI_RIGHTS_FD_SEEK;
    }
    
    __wasi_fd_t fd;
    __wasi_errno_t ret = __wasi_path_open(
        3,              // AT_FDCWD
        0,              // dirflags
        filename,       // path
        oflags,         // oflags
        fs_rights_base, // fs_rights_base
        fs_rights_base, // fs_rights_inheriting
        0,              // fs_flags
        &fd             // 输出的文件描述符
    );
    
    if (ret != 0 || fd < 0) {
        return NULL;
    }
    
    return fdopen(fd, mode);
}

size_t fread(void* ptr, size_t size, size_t nmemb, FILE* stream) {
    int fd = fileno(stream);
    size_t bytes_to_read = size * nmemb;
    size_t bytes_read = 0;
    
    __wasi_iovec_t iov = {
        .buf = ptr,
        .buf_len = bytes_to_read
    };
    
    __wasi_errno_t ret = __wasi_fd_read(fd, &iov, 1, &bytes_read);
    if (ret != 0) {
        return 0;
    }
    return bytes_read / size;
}

int fclose(FILE* stream) {
    AO_LOG("AO: fclose called\n");
    int fd = fileno(stream);
    fflush(stream);
    __wasi_errno_t ret = __wasi_fd_close(fd);
    return ret == 0 ? 0 : EOF;
}

// void* realloc(void* ptr, size_t size) {
//     void* new_ptr = memalign(16, size);
//     memcpy(new_ptr, ptr, size);
//     free(ptr);
//     //AO_LOG("DBG: Realloc called: %p -> %p, size: %zu\n", ptr, new_ptr, size);
//     return new_ptr;
// }

// Emscripten malloc does not align to 16 bytes correctly, which causes some 
// programs that use aligned memory (for example, those that use SIMD...) to
// crash. So we need to use the aligned allocator.
// void* malloc(size_t size) {
//     return memalign(16, size);
// }

int madvise(void* addr, size_t length, int advice) {
    AO_LOG("AO: madvise called with addr: %p, length: %zu, advice: %d\n", addr, length, advice);
    return 0;
}

// void* mmap(void* addr, size_t length, int prot, int flags, int fd, off_t offset) {
//     AO_LOG("AO: mmap called with addr: %p, length: %zu, prot: %d, flags: %d, fd: %d, offset: %d\n", addr, length, prot, flags, fd, offset);
//     // Allocate a buffer that fits with emscripten's normal allignments 
//     void* buffer = memalign(65536, length);
//     AO_LOG("AO: mmap: Reading from arweave to: %p, length: %zu\n", buffer, length);
//     weavedrive_read(fd, buffer, length);
//     AO_LOG("AO: mmap returned: %p\n", buffer);
//     return buffer;
// }

/*
int munmap(void* addr, size_t length) {
    AO_LOG("AO: munmap called with addr: %p, length: %zu\n", addr, length);
    return 0;
}
*/
