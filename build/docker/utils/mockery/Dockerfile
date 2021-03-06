ARG GO_VERSION=1.16
ARG MOCKERY_VERSION=v2.4.0

FROM vektra/mockery:${MOCKERY_VERSION} as base

FROM golang:${GO_VERSION}
COPY --from=base /usr/local/bin/mockery /usr/local/bin/mockery
ENTRYPOINT ["/usr/local/bin/mockery"]
