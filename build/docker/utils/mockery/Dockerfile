ARG GO_VERSION=1.16
ARG MOCKERY_VERSION=v2.5.1

### Build binary
FROM vektra/mockery:${MOCKERY_VERSION} as build-binary

### Image
FROM golang:${GO_VERSION}
COPY --from=build-binary /usr/local/bin/mockery /usr/local/bin/mockery
ENTRYPOINT ["/usr/local/bin/mockery"]
