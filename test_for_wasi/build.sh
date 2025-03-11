emcc test_for_wasi/test.c -o test_for_wasi/test666.wasm \
    -s STANDALONE_WASM=1 \
    -s EXPORTED_FUNCTIONS='["_main"]' \
    -s WASM=1 \
    -s NO_FILESYSTEM=0 


docker run --platform linux/amd64 -v .:/src p3rmaw3b/ao:latest emcc test_for_wasi/test.c  -o test_for_wasi/test123.wasm  \
-s STANDALONE_WASM=1 \
    -s EXPORTED_FUNCTIONS='["_main"]' \
    -s WASM=1 \
   -s NO_FILESYSTEM=0