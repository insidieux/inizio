ARG GO_VERSION=1.16

### Image
FROM golang:${GO_VERSION}
ARG WIRE_VERSION=v0.5.0
RUN go get -u -t github.com/google/wire/cmd/wire@${WIRE_VERSION}
ENTRYPOINT ["/go/bin/wire"]
