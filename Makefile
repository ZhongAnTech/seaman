INSTALL_DIR ?= ./seaman
GO_TAGS ?=
CGO_ENABLED ?= 0
GO_PACKAGES ?= ./...
GO_BUILD = go build -v -tags="${GO_TAGS}"

export CGO_ENABLED := ${CGO_ENABLED}

all: build

build:
	${GO_BUILD} ./cmd/seaman

install: build
	install -d ${INSTALL_DIR}
	install -d ${INSTALL_DIR}/log
	cp -Rp seaman configs ${INSTALL_DIR}

clean:
	rm -f seaman
	go clean -cache ${GO_PACKAGES}

docker:
	docker build -t seaman .
