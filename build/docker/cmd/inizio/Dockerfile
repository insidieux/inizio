ARG GO_VERSION=1.16
ARG ALPINE_VERSION=3.12

### Wire
FROM golang:${GO_VERSION} as wire
ARG WIRE_VERSION=v0.5.0
RUN GO111MODULE=on go get -u -t github.com/google/wire/cmd/wire@${WIRE_VERSION}
COPY . /project
WORKDIR /project
RUN find . -not -path "./vendor/*" -type f -name "wire_gen.go" -delete \
    && /go/bin/wire ./... \
    && mkdir -p /usr/local/src/wire \
    && find . -not -path "./vendor/*" -type f -name "wire_gen.go" -exec cp --parents {} /usr/local/src/wire \;

### Vendor
FROM golang:${GO_VERSION} as vendor
COPY . /project
WORKDIR /project
RUN go mod tidy && \
    go mod vendor

### Build binary
FROM golang:${GO_VERSION} as build-binary
COPY . /project
COPY --from=wire /usr/local/src/wire /project
COPY --from=vendor /project/vendor /project/vendor
WORKDIR /project
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build \
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
