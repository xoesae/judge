FROM debian:bookworm-slim AS base

RUN apt-get -y update && apt-get install -y \
    libc6 \
    libstdc++6 \
    libprotobuf32 \
    libnl-route-3-200 \
    python3 \
    python3-pip \
 && rm -rf /var/lib/apt/lists/*

FROM base AS build

RUN apt-get update && apt-get install -y \
    autoconf \
    bison \
    flex \
    gcc \
    g++ \
    git \
    libprotobuf-dev \
    libnl-route-3-dev \
    libtool \
    make \
    pkg-config \
    protobuf-compiler \
 && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/google/nsjail.git

RUN set -ex && cd /nsjail && make V=1

FROM base AS run

COPY --from=build /nsjail/nsjail /bin/nsjail
