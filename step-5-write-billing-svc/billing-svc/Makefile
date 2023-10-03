VERSION ?= $(shell git rev-parse --short HEAD)
SERVICE = billing-svc
ARCH ?= $(shell uname -m)

GO = CGO_ENABLED=$(CGO_ENABLED) GOFLAGS=-mod=vendor go
CGO_ENABLED ?= 0
GO_BUILD_FLAGS = -ldflags "-X main.version=${VERSION}"

# Utility functions
check_defined = \
	$(strip $(foreach 1,$1, \
		$(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = $(if $(value $1),, \
	$(error undefined '$1' variable: $2))

# Pattern #1 example: "example : description = Description for example target"
# Pattern #2 example: "### Example separator text
help: HELP_SCRIPT = \
	if (/^([a-zA-Z0-9-\.\/]+).*?: description\s*=\s*(.+)/) { \
		printf "\033[34m%-40s\033[0m %s\n", $$1, $$2 \
	} elsif(/^\#\#\#\s*(.+)/) { \
		printf "\033[33m>> %s\033[0m\n", $$1 \
	}

.PHONY: help
help:
	@perl -ne '$(HELP_SCRIPT)' $(MAKEFILE_LIST)

### Dev

.PHONY: run
run: description = Run $(SERVICE)
run:
	$(GO) run `ls -1 *.go | grep -v _test.go` -d

.PHONY: start/deps
start/deps: description = Start dependencies
start/deps:
	docker-compose up -d rabbitmq natsjs

### Build

.PHONY: build
build: description = Build $(SERVICE)
build: clean build/linux-amd64

.PHONY: build/linux-amd64
build/linux-amd64: description = Build $(SERVICE) for linux
build/linux-amd64: clean
	GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o ./build/$(SERVICE)-linux-amd64

.PHONY: clean
clean: description = Remove existing build artifacts
clean:
	$(RM) ./build/$(SERVICE)-*

### Docker

.PHONY: docker/build
docker/build: description = Build docker image
docker/build:
	docker build -t dselans/$(SERVICE) --build-arg TARGETOS=linux --build-arg TARGETARCH=amd64 \
	-f ./Dockerfile .

.PHONY: docker/push
docker/push: description = Push docker image
docker/push:
	docker push dselans/$(SERVICE)
