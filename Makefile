WAMR_VERSION = 2.2.0
WAMR_DIR = _build/wamr


WAMR_FLAGS = -DWAMR_ENABLE_LOG=1 -DWAMR_BUILD_DUMP_CALL_STACK=1 -DCMAKE_BUILD_TYPE=Debug

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
    WAMR_BUILD_PLATFORM = darwin
    ifeq ($(UNAME_M),arm64)
        WAMR_BUILD_TARGET = AARCH64
    else
        WAMR_BUILD_TARGET = X86_64
    endif
else
    WAMR_BUILD_PLATFORM = linux
    WAMR_BUILD_TARGET = X86_64
endif

.PHONY: all clean debug debug-clean utils

all: $(WAMR_DIR)/lib/libvmlib.a utils

debug: debug-clean $(WAMR_DIR)
	HB_DEBUG=1 make $(WAMR_DIR)/lib/libvmlib.a

debug-clean:
	rm -rf $(WAMR_DIR)

clean:
	rm -rf $(WAMR_DIR)

# Clone the WAMR repository at our target release
$(WAMR_DIR):
	git clone \
		https://github.com/bytecodealliance/wasm-micro-runtime.git \
		$(WAMR_DIR) \
		-b WAMR-$(WAMR_VERSION) \
		--single-branch

$(WAMR_DIR)/lib/libvmlib.a: $(WAMR_DIR)
	sed -i '742a tbl_inst->is_table64 = 1;' ./_build/wamr/core/iwasm/aot/aot_runtime.c; \
	cmake \
		$(WAMR_FLAGS) \
		-S $(WAMR_DIR) \
		-B $(WAMR_DIR)/lib \
		-DWAMR_BUILD_TARGET=$(WAMR_BUILD_TARGET) \
		-DWAMR_BUILD_PLATFORM=$(WAMR_BUILD_PLATFORM) \
		-DWAMR_BUILD_MEMORY64=1 \
		-DWAMR_DISABLE_HW_BOUND_CHECK=1 \
		-DWAMR_BUILD_EXCE_HANDLING=1 \
		-DWAMR_BUILD_SHARED_MEMORY=0 \
		-DWAMR_BUILD_AOT=1 \
		-DWAMR_BUILD_LIBC_WASI=1 \
		-DWAMR_BUILD_FAST_INTERP=0 \
		-DWAMR_BUILD_INTERP=1 \
		-DWAMR_BUILD_JIT=0 \
		-DWAMR_BUILD_FAST_JIT=0 \
		-DWAMR_BUILD_DEBUG_AOT=1 \
		-DWAMR_BUILD_TAIL_CALL=1 \
		-DWAMR_BUILD_AOT_STACK_FRAME=1 \
		-DWAMR_BUILD_MEMORY_PROFILING=1 \
		-DWAMR_BUILD_LIB_PTHREAD=1 \
		-DWAMR_BUILD_LIBC_BUILTIN=1 \
		-DWAMR_BUILD_DUMP_CALL_STACK=1
	make -C $(WAMR_DIR)/lib -j8

	# Copy the libvmlib.a file to the go directory
	cp $(WAMR_DIR)/lib/libvmlib.a $(WAMR_DIR)/language-bindings/go/wamr/packaged/lib/darwin-aarch64/

	# Copy the include directory to the go directory
	cp -r $(WAMR_DIR)/core/iwasm/include/*.h $(WAMR_DIR)/language-bindings/go/wamr/packaged/include/

	
utils:
	# Compile libs
	gcc -c ./libs/wamr_imports.c \
		-o ./libs/wamr_imports.o \
		-I./_build/wamr/language-bindings/go/wamr/packaged/include
	ar rcs ./libs/libwamr_imports.a ./libs/wamr_imports.o
	rm ./libs/wamr_imports.o

	cp ./libs/libwamr_imports.a ./_build/wamr/language-bindings/go/wamr/libwamr_imports.a

	# Copy modifyed wamr go lib
	cp ./libs/utils.go ./_build/wamr/language-bindings/go/wamr/utils.go
	cp ./libs/wamr_imports.h ./_build/wamr/language-bindings/go/wamr/wamr_imports.h
	

# Print the library path
print-lib-path:
	@echo $(CURDIR)/$(WAMR_DIR)/lib/libvmlib.a