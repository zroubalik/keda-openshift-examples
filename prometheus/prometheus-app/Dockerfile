FROM golang:1.17 as builder

WORKDIR /workspace

COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

COPY vendor/ vendor/
COPY main.go main.go

RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
COPY --from=builder /workspace/prometheusapp .
# 65532 is numeric for nonroot
USER 65532:65532

ENTRYPOINT ["/prometheusapp"]