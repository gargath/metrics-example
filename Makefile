include .env

BINARY := server
VERSION := $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")
ECHO := echo


.NOTPARALLEL:

.PHONY: all
all: test build

.PHONY: check
check: fmt lint vet
	$(GINKGO) -p ./pkg/... ./cmd/... || find . -name \*_test.db -type f -delete
	find . -name \*_test.db -type f -delete

.PHONY: build
build: clean $(BINARY)

.PHONY: clean
clean:
	rm -f $(BINARY)
	find . -name \*_test.db -type f -delete

.PHONY: distclean
distclean: clean
	rm -f .env
	rm -f anaximander.db

# Run go fmt against code
.PHONY: fmt
fmt:
	$(GO) fmt ./pkg/... ./cmd/...

# Run go vet against code
.PHONY: vet
vet:
	$(GO) vet -tags dev -composites=false ./pkg/... ./cmd/...

.PHONY: lint
lint:
	@ $(ECHO) "\033[36mLinting code\033[0m"
	$(LINTER) run --disable-all \
                --exclude-use-default=false \
                --enable=govet \
                --enable=ineffassign \
                --enable=deadcode \
                --enable=golint \
                --enable=goconst \
                --enable=gofmt \
                --enable=goimports \
                --skip-dirs=pkg/client/ \
                --deadline=120s \
                --tests ./...
	@ $(ECHO)

.PHONY: test
test:
	@ $(ECHO) "\033[36mRunning test suite in Ginkgo\033[0m"
	$(GINKGO) -v -p -race -randomizeAllSpecs ./pkg/... ./cmd/... || find . -name \*_test.db -type f -delete
	@ $(ECHO)
	find . -name \*_test.db -type f -delete

# Build binary
$(BINARY): fmt vet
	GO111MODULE=on CGO_ENABLED=1 $(GO) build -o $(BINARY) -ldflags="-X main.VERSION=${VERSION}" github.com/gargath/metrics-example/cmd
