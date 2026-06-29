.PHONY: build build-dev build-all clean install test fmt vet docs docs-build

BINARY_NAME=vibe
INSTALL_PATH=$(HOME)/.oh-my-zsh/custom/plugins/vibe
DOCS_DIR=docs
DOCS_PORT?=1313

# VERSION is derived from git so `make install` stamps a real release version by
# default. If HEAD is exactly on a tag, that tag is used (e.g. v0.3.6); otherwise
# the most recent tag with commit suffix is used (e.g. v0.3.6-1-g68df413). Falls
# back to "dev" only when git/tags are unavailable. Override explicitly with
# `make install VERSION=v1.2.3`.
VERSION ?= $(shell git describe --tags 2>/dev/null || echo dev)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) main.go

# build-dev produces a build stamped as "dev" (the updater never self-updates a
# dev build). Use this for local development to silence update prompts.
build-dev:
	go build -ldflags "-X main.version=dev" -o $(BINARY_NAME) main.go

build-all:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-linux-arm64 main.go

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/
	cp vibe.plugin.zsh $(INSTALL_PATH)/
	cp _vibe $(INSTALL_PATH)/
	@echo "Installed to $(INSTALL_PATH)"
	@echo "Add 'vibe' to your plugins list in ~/.zshrc and reload your shell"

clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# docs serves the Hugo documentation site locally for review at
# http://localhost:$(DOCS_PORT). Includes drafts and future-dated content.
# Override the port with `make docs DOCS_PORT=8080`.
docs:
	@command -v hugo >/dev/null 2>&1 || { echo "hugo not found. Install it: brew install hugo"; exit 1; }
	cd $(DOCS_DIR) && hugo server --buildDrafts --buildFuture --port $(DOCS_PORT)

# docs-build performs a production build of the docs into docs/public.
docs-build:
	@command -v hugo >/dev/null 2>&1 || { echo "hugo not found. Install it: brew install hugo"; exit 1; }
	cd $(DOCS_DIR) && hugo --minify
