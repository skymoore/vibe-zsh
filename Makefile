.PHONY: build clean install test

BINARY_NAME=vibe
INSTALL_PATH=$(HOME)/.oh-my-zsh/custom/plugins/vibe

build:
	go build -o $(BINARY_NAME) main.go

build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 main.go

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/
	cp vibe.plugin.zsh $(INSTALL_PATH)/
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
