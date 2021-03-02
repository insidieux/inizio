ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}

RUN go get -u -t github.com/google/wire/cmd/wire

ENTRYPOINT ["/go/bin/wire"]
