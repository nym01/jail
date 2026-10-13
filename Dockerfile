FROM golang:1.26 AS builder

WORKDIR /src
RUN apt-get update \
    && apt-get install -y --no-install-recommends bzip2 \
    && rm -rf /var/lib/apt/lists/*
RUN ARCH=$(uname -m) \
    && URL=https://storage.googleapis.com/gvisor/releases/release/latest/${ARCH} \
    && curl -fsSLO ${URL}/gvisor.tar.bz2 \
    && curl -fsSLO ${URL}/gvisor.tar.bz2.sha512 \
    && sha512sum -c gvisor.tar.bz2.sha512 \
    && mkdir /gvisor \
    && tar -xjf gvisor.tar.bz2 -C /gvisor
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o jail .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3 \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /rootfs/usr/bin /rootfs/proc /rootfs/dev /rootfs/tmp \
    && cp -a /usr/bin/python3* /rootfs/usr/bin/ \
    && cp -a /usr/lib /usr/lib64 /rootfs/usr/ \
    && ln -s usr/lib /rootfs/lib \
    && ln -s usr/lib64 /rootfs/lib64

COPY --from=builder /src/jail /usr/local/bin/jail
COPY --from=builder /gvisor/runsc /usr/local/bin/runsc
COPY --from=builder /gvisor/gvisor-bin /usr/local/bin/gvisor-bin

CMD ["jail"]
