
CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTAINER_DIR="$(cd "${CURRENT_DIR}/../container" && pwd)"

chmod +x ${CONTAINER_DIR}/src/emcc-lua
chmod +x ${CONTAINER_DIR}/src/clang-lua
rm process.wasm
rm ../wasm/process.wasm
docker run --platform linux/amd64 \
    -e DEBUG=1 \
    -v .:/src \
    -v  ${CONTAINER_DIR}/src/emcc-lua:/usr/local/bin/emcc-lua \
    -v  ${CONTAINER_DIR}/src/clang-lua:/usr/local/bin/clang-lua \
    p3rmaw3b/ao:latest clang-lua

cp process.wasm ../wasm/