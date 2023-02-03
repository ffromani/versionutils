RUNTIME ?= podman
REPOOWNER ?= ffromani
IMAGENAME ?= versionutils
IMAGETAG ?= latest

all: dist

.PHONY: build
build: dist

.PHONY: dist
dist: binaries

outdir:
	mkdir -p _output || :

.PHONY: deps-update
deps-update:
	go mod tidy && go mod vendor

.PHONY: deps-clean
deps-clean:
	rm -rf vendor

.PHONY: binaries
binaries: outdir deps-update
	# go flags are set in here
	./hack/build-binaries.sh

.PHONY: clean
clean:
	rm -rf _output

.PHONY: test-e2e
test-e2e: binaries
	go test test/e2e/...

