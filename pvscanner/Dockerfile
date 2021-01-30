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


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless and https://github.com/GoogleContainerTools/distroless/tree/master/base for more details
FROM centos:7
#FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /workspace/pvscanner/pvscanner /pvscanner


CMD ["/pvscanner"]
