SHELL:=/bin/bash

.DEFAULT_GOAL := help
# put here commands, that have the same name as files in dir
.PHONY: run clean generate build docker_build docker_push

BUILD_TAG=$(shell git describe --tags --abbrev=0 HEAD)
BUILD_HASH=$(shell git rev-parse --short HEAD)
BUILD_BRANCH=$(shell git symbolic-ref HEAD |cut -d / -f 3)
BUILD_VERSION="${BUILD_TAG}-${BUILD_HASH}"
BUILD_TIME=$(shell date --utc +%F-%H:%m:%SZ)
PACKAGE=publications
LDFLAGS=-extldflags=-static -w -s -X ${PACKAGE}/internal/version.Version=${BUILD_VERSION} -X ${PACKAGE}/internal/version.BuildTime=${BUILD_TIME}
CONTAINER_IMAGE_REGISTRY=local/publications

help:
	@echo "build, build-images, deps, build-api, build-api-image, generate-api, build-sql-migrations-image"

version:
	@echo "${BUILD_VERSION}"

# ONLY TABS IN THE START OF COMMAND, NO SPACES!
build: deps build-api

build-images: build-api-image build-sql-migrations-image

clean:
	@echo "[INFO] Cleaning build files"
	rm -f build/*

deps:
	@echo "[INFO] Downloading and installing dependencies"
	go mod download

build-api: deps generate-api
	@echo "[INFO] Building API Server binary"
	go build -ldflags "${LDFLAGS}" -o build/publications-api ./cmd/publications-api
	@echo "[INFO] Build successful"

generate-api:
	@echo "[INFO] Running code generations for API"
	go generate cmd/publications-api/main.go

build-api-image:
	@echo "[INFO] Building API container image"
	docker build -t ${CONTAINER_IMAGE_REGISTRY}/publications-api:${BUILD_BRANCH}-${BUILD_HASH} \
	-t ${CONTAINER_IMAGE_REGISTRY}/publications-api:${BUILD_VERSION} \
	--build-arg BUILD_VERSION=${BUILD_VERSION} -f cmd/publications-api/Dockerfile .

build-sql-migrations-image:
	@echo "[INFO] Building SQL migrations image"
	docker build -t ${CONTAINER_IMAGE_REGISTRY}/publications-sql-migrations:${BUILD_BRANCH}-${BUILD_HASH} \
	-t ${CONTAINER_IMAGE_REGISTRY}/publications-sql-migrations:${BUILD_VERSION} \
	-f migrations/Dockerfile .
