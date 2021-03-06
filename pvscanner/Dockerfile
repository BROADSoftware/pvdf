# Build the manager binary
FROM golang:1.15.6-alpine3.12 as builder
RUN apk add alpine-sdk

WORKDIR /workspace
# Copy the Go Modules manifests
COPY pvscanner/go.mod pvscanner/go.mod
COPY pvscanner/go.sum pvscanner/go.sum
COPY shared/go.mod shared/go.mod
COPY shared/go.sum shared/go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN cd pvscanner; go mod download

# Copy the go source
COPY pvscanner/main.go pvscanner/main.go
COPY pvscanner/pkg/ pvscanner/pkg/

COPY shared/common/ shared/common/
COPY shared/pkg/ shared/pkg/

WORKDIR /workspace/pvscanner
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o pvscanner main.go

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y install --no-install-recommends \
        file \
        xfsprogs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=builder /workspace/pvscanner/pvscanner /pvscanner

CMD ["/pvscanner"]
