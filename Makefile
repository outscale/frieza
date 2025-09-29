all: help

.PHONY: help
help:
	@echo "help:"
	@echo "- build   : build frieza"
	@echo "- install : install frieza"
	@echo "- test    : run all tests"
	@echo "- release : will generate artefacts locally"

.PHONY: test
test: test-reuse test-go-fmt build
	@echo all tests OK

.PHONY: test-reuse
test-reuse:
	@echo test reuse:
	docker run --rm --volume $(PWD):/data fsfe/reuse:latest lint

.PHONY: test-go-fmt
test-go-fmt:
	@echo test go fmt:
	test -z $(gofmt -l .)

.PHONY: build
build:
	@echo building:
	cd cmd/frieza && go build -ldflags "-X github.com/outscale/frieza/internal/common.version=`cat version` -X github.com/outscale/frieza/internal/common.commit=`git rev-list -1 HEAD`"

.PHONY: install
install:
	@echo installing:
	cd cmd/frieza && go install -ldflags "-X github.com/outscale/frieza/internal/common.version=`cat version` -X github.com/outscale/frieza/internal/common.commit=`git rev-list -1 HEAD`"

.PHONY: release
release:
	goreleaser release --snapshot --rm-dist
