FROM debian:bookworm-slim AS base

RUN apt-get -y update && apt-get install -y \
    libc6 \
    libstdc++6 \
    libprotobuf32 \
    libnl-route-3-200 \
    python3 \
    python3-pip \
    debootstrap \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /rootfs && \
    debootstrap --variant=minbase bookworm /rootfs http://deb.debian.org/debian && \
    chroot /rootfs apt-get update && \
    chroot /rootfs apt-get install -y python3 && \
    chroot /rootfs apt-get clean && \
    rm -rf /rootfs/var/lib/apt/lists/*

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
    bash  \
    coreutils \
 && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/google/nsjail.git

RUN set -ex && cd /nsjail && make V=1

# ms golang
FROM golang:1.24 AS go-builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main cmd/api/api.go

FROM base AS run

COPY --from=build /nsjail/nsjail /bin/nsjail
COPY --from=go-builder /app/main .
COPY --from=base /rootfs /tmp/myroot

EXPOSE 8080

CMD ["./main"]
