# FROM emscripten/emsdk:3.1.59
# LABEL maintainer "tom wilson <tom@hyper.io>"

FROM --platform=linux/amd64 ubuntu:22.04

# The working directory used by the base image is /src, so we can mount volumes to there
# to expose files on the host to the ao container
# 
# https://github.com/emscripten-core/emsdk/blob/9b0db91883452051aca8deddc932363aab29060b/docker/Dockerfile#L120

RUN apt-get update -qq -y
RUN apt-get install -y curl vim make gcc libreadline6-dev libssl-dev zlib1g-dev zip unzip \
    python3 \
    python3-pip

# install clang
RUN apt-get update && apt-get install -y \
  clang \
  lld \
  llvm

ENV LUA_VERSION 5.3.4
ENV LUAROCKS_VERSION 2.4.4
ENV PYTHON_VERSION 3.6.6


# Intall yaml
RUN pip3 install pyyaml


# Download lua source
RUN cd / && \
  curl -L http://www.lua.org/ftp/lua-${LUA_VERSION}.tar.gz | tar xzf -


# Install NodeJS
ENV NODE_VERSION=18.17.0
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
ENV NVM_DIR=/root/.nvm
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"
RUN node --version
RUN npm --version


# 安装 wasi-sdk
RUN curl -L https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-22/wasi-sdk-22.0-linux.tar.gz | \
    tar xz -C /opt && \
    ln -s /opt/wasi-sdk-22.0 /opt/wasi-sdk


ENV CC '/opt/wasi-sdk/bin/clang --target=wasm32-wasi'
ENV NM '/opt/wasi-sdk/bin/llvm-nm'
ENV AR '/opt/wasi-sdk/bin/llvm-ar'

RUN cd /lua-${LUA_VERSION} && \
    make clean && \
    make generic \
    CC="/opt/wasi-sdk/bin/clang \
        --target=wasm32-wasi \
        --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
        -isystem /opt/wasi-sdk/share/wasi-sysroot/include/wasm32-wasip1 \
        #-isystem ${WASIX_DIR}/include \
        -O3 \
        -mllvm -wasm-enable-sjlj \
        -D_WASI_EMULATED_SIGNAL \
        -D_WASI_EMULATED_PROCESS_CLOCKS \
        -D_WASI_EMULATED_MMAN \
        -D_WASI_EMULATED_GETPID \
        # -D__wasilibc_unmodified_upstream \
        -DL_tmpnam=20" \
    AR="/opt/wasi-sdk/bin/llvm-ar rcu" \
    RANLIB="/opt/wasi-sdk/bin/llvm-ranlib" \
    LDFLAGS="-Wl,--export-all \
             #-L${WASIX_DIR} -lwasix \
             -Wl,--no-entry \
             -Wl,--allow-undefined \
             -Wl,--stack-first \
             -Wl,--import-memory \
             -Wl,--max-memory=4294967296"


#############################
##### Install Commands ######
#############################

# NOTE: Some commands do not require additional work
# when building the image, but are still commented
# here for posterity

###############
#### Init #####
###############
# The init command currently does not require running
# within the ao container, and so requires no additional setup

###############
#### Repl #####
###############
# The repl command simply fires up lua repl, which is installed
# as part of lua runtime above, and so requires no additional setup

###############
#### Build ####
###############
COPY ./src/emcc-lua /usr/local/bin/emcc-lua
COPY ./src/clang-lua /usr/local/bin/clang-lua
COPY ./src/emcc_lua_lib /usr/local/emcc-lua/emcc_lua_lib

COPY ./src/definition.yml /opt/definition.yml
COPY ./src/loader.lua /opt/loader.lua

RUN mkdir -p /opt/src
COPY ./src/json.lua /opt/src/json.lua
COPY ./src/ao.lua /opt/src/ao.lua
# COPY ./src/pack.lua /opt/pack.lua
COPY ./src/main.c /opt/main.c
# COPY ./src/linit.c /opt/linit.c
COPY ./src/lsqlite3.c /opt/lsqlite3.c
COPY ./src/lsqlite3.h /opt/lsqlite3.h
COPY ./src/sqlite3.c /opt/sqlite3.c
COPY ./src/sqlite3.h /opt/sqlite3.h
COPY ./src/main.lua /opt/main.lua


RUN chmod +x /usr/local/bin/emcc-lua
RUN chmod +x /usr/local/bin/clang-lua


###########################################
### Publish, Bundler, Process, Contract ###
###########################################
# We first create a directory for the node impls to be placed
# and dependencies installed
# 
# By running npm link, we allow any commands exposed by
# the node module to be ran globally within the container
RUN mkdir -p /opt/node
COPY ./src/node /opt/node
RUN cd /opt/node && \
  npm install --omit="dev" && \
  npm link