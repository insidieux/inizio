ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}

RUN GO111MODULE=on go get -u github.com/gobuffalo/packr/v2/packr2

ENTRYPOINT ["/go/bin/packr2"]
