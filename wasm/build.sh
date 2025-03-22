docker run --platform linux/amd64 \
    -v $(pwd):/src \
    p3rmaw3b/ao:latest \
    /opt/wasi-sdk/bin/clang \
    --target=wasm32-wasi \
    -Wl,--initial-memory=67108864,--max-memory=2147483648 \
    -o /src/print_string.wasm \
    /src/print_string.c

# docker run --platform linux/amd64 \
#     -v $(pwd):/src \
#     p3rmaw3b/ao:latest-64 \
#     emcc \
#     -O3 \
#     -s WASM=1 \
#     -s STANDALONE_WASM=1 \
#     -s INITIAL_MEMORY=67108864 \
#     -s MAXIMUM_MEMORY=2147483648 \
#     -s EXPORTED_FUNCTIONS='["_print_string", "_malloc", "_free"]' \
#     -o /src/print_string-emcc.wasm \
#     /src/print_string.c

# docker run --platform linux/amd64 \
#     -v $(pwd):/src \
#     p3rmaw3b/ao:latest-64 \
#     emcc \
#     -O3 \
#     -s WASM=1 \
#     -s STANDALONE_WASM=1 \
#     -s TOTAL_MEMORY=67108864 \
#     -s TOTAL_STACK=1048576 \
#     -s ALLOW_MEMORY_GROWTH=1 \
#     -s EXPORTED_FUNCTIONS='["_print_string", "_malloc", "_free"]' \
#     -s EXPORTED_RUNTIME_METHODS='["ccall", "cwrap"]' \
#     -s ERROR_ON_UNDEFINED_SYMBOLS=0 \
#     -o /src/print_string-emcc.wasm \
#     /src/print_string.c