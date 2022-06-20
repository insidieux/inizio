ARG GO_VERSION=1.18
FROM golang:${GO_VERSION}

RUN go install github.com/google/wire/cmd/wire@v0.5.0

ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"

ENTRYPOINT ["/go/bin/wire"]
