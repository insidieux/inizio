override APP_NAME=inizio
override GO_VERSION=1.18
override GOLANGCI_LINT_VERSION=v1.46.2
override SECUREGO_GOSEC_VERSION=2.12.0
override HADOLINT_VERSION=v2.10.0
override WIRE_VERSION=v0.5.0
override PROTOC_VERSION=3.3.0
override MOCKERY_VERSION=v2.13.1
override CHANGELOG_GENERATOR_VERSION=1.15.2

GOOS?=$(shell go env GOOS || echo linux)
GOARCH?=$(shell go env GOARCH || echo amd64)
CGO_ENABLED?=0

DOCKER_REGISTRY?=docker.io
DOCKER_IMAGE?=${DOCKER_REGISTRY}/insidieux/${APP_NAME}
DOCKER_TAG?=latest
CHANGELOG_GITHUB_TOKEN?=

ifeq (, $(shell which docker))
$(error "Binary docker not found in $(PATH)")
endif

.PHONY: all
all: cleanup vendor lint test build

# --- [ CI helpers ] ---------------------------------------------------------------------------------------------------

.PHONY: cleanup
cleanup:
	@rm ${PWD}/bin/${APP_NAME} || true
	@rm ${PWD}/coverage.out || true
	@rm -r ${PWD}/vendor || true

.PHONY: tidy
tidy:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go mod tidy

.PHONY: vendor
vendor:
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

.PHONY: wire
wire:
	@docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg WIRE_VERSION=${WIRE_VERSION} \
		-f ${PWD}/build/docker/utils/wire/Dockerfile \
		-t wire:custom \
			build/docker/utils/wire
	@find ${PWD} -type f -name "wire_gen.go" -delete
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		wire:custom \
			/project/...

.PHONY: protoc
protoc: $(shell find api/protobuf -type f -name "*.proto")
	@find ${PWD}/internal -type f -name "*.pb.go" -delete
	for file in $^ ; do \
		docker run --rm \
			-v ${PWD}:/project \
			-w /project \
			rvolosatovs/protoc:${PROTOC_VERSION} \
				--proto_path /project \
				--go_out /project/pkg \
				--go_opt=paths=source_relative \
				--go-grpc_out=/project/pkg \
				--go-grpc_opt=paths=source_relative \
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
	@docker run \
		--rm \
		-v ${PWD}:/project \
		-w /project \
		vektra/mockery:${MOCKERY_VERSION} \
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

.PHONY: generate-changelog
generate-changelog:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		-e CHANGELOG_GITHUB_TOKEN=${CHANGELOG_GITHUB_TOKEN} \
		ferrarimarco/github-changelog-generator:${CHANGELOG_GENERATOR_VERSION} \
			--user insidieux \
			--project inizio \
			--no-unreleased
