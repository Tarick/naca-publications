SHELL:=/bin/bash

.DEFAULT_GOAL := help
# put here commands, that have the same name as files in dir
.PHONY: run clean generate build docker_build docker_push deploy-to-local-k8s build-and-deploy

BUILD_TAG=$(shell git describe --tags --abbrev=0 HEAD)
BUILD_HASH=$(shell git rev-parse --short HEAD)
BUILD_BRANCH=$(shell git symbolic-ref HEAD |cut -d / -f 3)
BUILD_VERSION="${BUILD_TAG}-${BUILD_HASH}"
BUILD_TIME=$(shell date --utc +%F-%H:%m:%SZ)
PACKAGE=publications
LDFLAGS=-extldflags=-static -w -s -X ${PACKAGE}/internal/version.Version=${BUILD_VERSION} -X ${PACKAGE}/internal/version.BuildTime=${BUILD_TIME}
# This CONTAINER_REGISTRY must be sourced from environment and it must be FQDN,
# containerd registry plugin doesn't give a shit about short names even if they're present locally, appends docker.io to it
CONTAINER_IMAGE_REGISTRY=${CONTAINER_REGISTRY_FQDN}/publications

help:
	@echo "build, build-images, deps, build-api, build-api-image, generate-api, build-sql-migrations-image, build-importer, build-and-deploy, deploy-to-local-k8s"


version:
	@echo "${BUILD_VERSION}"

# ONLY TABS IN THE START OF COMMAND, NO SPACES!
build: deps build-api build-importer

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
	buildctl build --frontend dockerfile.v0 --opt build-arg:BUILD_VERSION=${BUILD_VERSION} \
	--local context=. --local dockerfile=cmd/publications-api  \
	--output type=image,\"name=${CONTAINER_IMAGE_REGISTRY}/publications-api:${BUILD_BRANCH}-${BUILD_HASH},${CONTAINER_IMAGE_REGISTRY}/publications-api:${BUILD_VERSION}\"
	@echo "[INFO] Image built successfully"

build-importer: deps
	@echo "[INFO] Building Publications importer binary"
	go build -ldflags "${LDFLAGS}" -o build/publications-importer ./cmd/publications-importer
	@echo "[INFO] Build successful"

build-sql-migrations-image:
	@echo "[INFO] Building SQL migrations image"
	buildctl build --frontend dockerfile.v0 --opt build-arg:BUILD_VERSION=${BUILD_VERSION} \
	--local context=. --local dockerfile=migrations  \
	--output type=image,\"name=${CONTAINER_IMAGE_REGISTRY}/publications-sql-migrations:${BUILD_BRANCH}-${BUILD_HASH},${CONTAINER_IMAGE_REGISTRY}/publications-sql-migrations:${BUILD_VERSION}\"
	@echo "[INFO] Image built successfully"

build-and-deploy: build-images deploy-to-local-k8s

deploy-to-local-k8s:
	@echo "[INFO] Deploying current Publications to local k8s service"
	@echo "[INFO] Deleting old SQL migrations"
	helmfile --environment local --selector app_name=publications-sql-migrations -f ../naca-ops-config/helm/helmfile.yaml destroy
	# Ugly workaround for helm 3.4
	helmfile --environment local --selector app_name=publications-api -f ../naca-ops-config/helm/helmfile.yaml destroy
	@echo "[INFO] Deploying rss-feeds images with tag ${BUILD_VERSION}"
	PUBLICATIONS_TAG=${BUILD_VERSION} helmfile --environment local --selector tier=naca-publications -f ../naca-ops-config/helm/helmfile.yaml sync --skip-deps
