ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.16

### Build binary
FROM golang:${GO_VERSION} as build-binary
COPY . /project
WORKDIR /project
RUN go mod vendor \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build \
        -v \
        -mod vendor \
        -o /project/bin/inizio \
            /project/cmd/inizio

### Image
FROM alpine:${ALPINE_VERSION} as image
COPY --from=build-binary /project/bin/inizio /usr/local/bin/inizio
COPY configs/inizio/plugins.yaml /etc/inizio/plugins.yaml
RUN chmod +x /usr/local/bin/inizio
VOLUME ["/etc/inizio"]
VOLUME ["/usr/local/bin/inizio-plugins"]
ENTRYPOINT ["inizio"]
