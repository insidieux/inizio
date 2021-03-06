override APP_NAME=inizio
override GO_VERSION=1.16
override PROTOC_VERSION=3.1.32
override MOCKERY_VERSION=v2.5.1
override GOLANGCI_LINT_VERSION=v1.38.0
override SECUREGO_GOSEC_VERSION=v2.7.0
override HADOLINT_VERSION=v1.23.0

GOOS?=$(shell go env GOOS || echo linux)
GOARCH?=$(shell go env GOARCH || echo amd64)
CGO_ENABLED?=0

DOCKER_REGISTRY?=docker.io
DOCKER_IMAGE?=${DOCKER_REGISTRY}/insidieux/${APP_NAME}
DOCKER_TAG?=latest

ifeq (, $(shell which docker))
$(error "Binary docker not found in $(PATH)")
endif

.PHONY: all
all: cleanup wire vendor lint test build

# --- [ CI helpers ] ---------------------------------------------------------------------------------------------------

.PHONY: cleanup
cleanup:
	@rm ${PWD}/bin/${APP_NAME} || true
	@rm ${PWD}/coverage.out || true
	@find ${PWD} -type f -name "wire_gen.go" -delete
	@rm -r ${PWD}/vendor || true

.PHONY: wire
wire:
	@docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		-f ${PWD}/build/docker/utils/wire/Dockerfile \
		-t wire:custom \
			build/docker/utils/wire
	@find ${PWD} -type f -name "wire_gen.go" -delete
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		wire:custom \
			/project/...

.PHONY: vendor
vendor:
	@rm -r ${PWD}/vendor || true
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go mod tidy
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go mod vendor

.PHONY: lint-golangci-lint
lint-golangci-lint:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golangci/golangci-lint:${GOLANGCI_LINT_VERSION} \
			golangci-lint run -v

.PHONY: lint-golint
lint-golint:
	@docker run --rm \
		-v ${PWD}:/project \
		cytopia/golint \
			--set_exit_status \
			/project/cmd/...\
			/project/internal/...\
			/project/pkg/...

.PHONY: lint-gosec
lint-gosec:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		securego/gosec:${SECUREGO_GOSEC_VERSION} \
			/project/...

.PHONY: lint-dockerfile
lint-dockerfile:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		hadolint/hadolint:${HADOLINT_VERSION} \
			hadolint \
				/project/build/docker/cmd/inizio/Dockerfile

.PHONY: lint
lint:
	@make lint-golangci-lint
	@make lint-golint
	@make lint-gosec
	@make lint-dockerfile

.PHONY: test
test:
	@rm -r ${PWD}/coverage.out || true
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go test \
				-race \
				-mod vendor \
				-covermode=atomic \
				-coverprofile=/project/coverage.out \
					/project/...

.PHONY: build
build:
	@rm ${PWD}/bin/${APP_NAME} || true
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		-e GOOS=${GOOS} \
		-e GOARCH=${GOARCH} \
		-e CGO_ENABLED=${CGO_ENABLED} \
		-e GO111MODULE=on \
		golang:${GO_VERSION} \
			go build \
				-mod vendor \
				-o /project/bin/${APP_NAME} \
				-v /project/cmd/${APP_NAME}

# --- [ Local helpers ] ------------------------------------------------------------------------------------------------

.PHONY: protoc
protoc: $(shell find api/protobuf -type f -name "*.proto")
	@find ${PWD}/internal -type f -name "*.pb.go" -delete
	for file in $^ ; do \
		docker run --rm \
			-v ${PWD}:/project \
			-w /project \
			thethingsindustries/protoc:${PROTOC_VERSION} \
				--proto_path /project \
				--go_out=paths=source_relative,plugins=grpc:/project/pkg \
					$${file}; \
	done

.PHONY: mockery
mockery:
ifndef MOCKERY_SOURCE_DIR
	$(error MOCKERY_SOURCE_DIR is not set)
endif
ifndef MOCKERY_INTERFACE
	$(error MOCKERY_INTERFACE is not set)
endif
	@find ${PWD} -type f -name "mock_*_test.go" -delete
	@docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		-f ${PWD}/build/docker/utils/mockery/Dockerfile \
		-t mockery:${MOCKERY_VERSION} \
			build/docker/utils/mockery
	@docker run \
		--rm \
		-v ${PWD}:/project \
		-w /project \
		mockery:${MOCKERY_VERSION} \
			--testonly \
			--inpackage \
			--case snake \
			--log-level trace \
			--output /project/${MOCKERY_SOURCE_DIR} \
			--dir /project/${MOCKERY_SOURCE_DIR} \
			--name=${MOCKERY_INTERFACE}

.PHONY: docker-image-build
docker-image-build:
ifndef DOCKER_IMAGE
	$(error DOCKER_IMAGE is not set)
endif
ifndef DOCKER_TAG
	$(error DOCKER_TAG is not set)
endif
	@docker rmi ${DOCKER_IMAGE}:${DOCKER_TAG} || true
	@docker build \
		-f ${PWD}/build/docker/cmd/inizio/Dockerfile \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
			.
