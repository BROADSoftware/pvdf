# Build the manager binary
FROM golang:1.15.6-alpine3.12 as builder
RUN apk add alpine-sdk

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o /workspace/pvscanner main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless and https://github.com/GoogleContainerTools/distroless/tree/master/base for more details
# (Tag :debug for debug friendly image)
#FROM gcr.io/distroless/static:debug
FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /workspace/pvscanner /

CMD ["/pvscanner"]
