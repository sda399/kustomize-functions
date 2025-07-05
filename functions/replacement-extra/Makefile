.PHONY: image
SHELL := /usr/bin/env bash
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
MYGOBIN = $(shell go env GOBIN)
ifeq ($(MYGOBIN),)
MYGOBIN = $(shell go env GOPATH)/bin
endif
export PATH := $(MYGOBIN):$(PATH)

NAME := replacement-extra

image: build
	./local-resource/bin/$(NAME) gen ./
	docker build -f ./Dockerfile -t ghcr.io/zprd/kustomize-functions/$(NAME):v1.0.0 image/


build:
	cd image && go build -v -o ../local-resource/bin/$(NAME) .
	chmod +x ./local-resource/bin/$(NAME)

run: build
	DEBUG=true \
		kustomize build --enable-alpha-plugins --enable-exec --stack-trace ./local-resource

test:
	cd image && go test -v .

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
		| sh -s -- -b $(MYGOBIN) v2.1.6

	cd image && $(MYGOBIN)/golangci-lint run
