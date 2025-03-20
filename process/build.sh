chmod +x /Users/webbergao/work/src/arweave/aos-sqlite/container/src/emcc-lua
docker run --platform linux/amd64 \
    -e DEBUG=1 \
    -v .:/src \
    -v  /Users/webbergao/work/src/arweave/aos-sqlite/container/src/emcc-lua:/usr/local/bin/emcc-lua \
    -v  /Users/webbergao/work/src/arweave/aos-sqlite/container/src/clang-lua:/usr/local/bin/clang-lua \
    p3rmaw3b/ao:latest clang-lua

cp process.wasm ../wasm/