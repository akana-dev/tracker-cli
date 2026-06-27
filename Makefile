BINARY_NAME=tracker
BUILD_DIR=build

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X tracker/internal/version.Version=$(VERSION) \
	-X tracker/internal/version.Commit=$(COMMIT) \
	-X tracker/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: build build-all install test clean release checksums

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/tracker

build-all: clean
	@mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/tracker
	GOOS=linux   GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/tracker
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/tracker
	GOOS=darwin  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/tracker
	GOOS=darwin  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/tracker
	@$(MAKE) checksums

checksums:
	@cd $(BUILD_DIR) && \
	if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum $(BINARY_NAME)-* > checksums.txt; \
	elif command -v shasum >/dev/null 2>&1; then \
		shasum -a 256 $(BINARY_NAME)-* > checksums.txt; \
	else \
		echo "Ошибка: не найден sha256sum или shasum"; \
		exit 1; \
	fi
	@echo "Checksums:"
	@cat $(BUILD_DIR)/checksums.txt

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

release: build-all
	@echo ""
	@echo "Релиз $(VERSION) готов в $(BUILD_DIR)/"
	@echo "Создайте тег: git tag -a $(VERSION) -m 'Release $(VERSION)'"
	@echo "И запушьте:    git push origin $(VERSION)"