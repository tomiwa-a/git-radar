VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: build build-all install clean

build:
	go build $(LDFLAGS) -o build/git-radar ./cmd

build-all: clean
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/git-radar-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/git-radar-darwin-arm64 ./cmd
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/git-radar-linux-amd64 ./cmd
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/git-radar-linux-arm64 ./cmd
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/git-radar-windows-amd64.exe ./cmd

install: build
	sudo cp build/git-radar /usr/local/bin/
	@if [ "$(shell uname)" = "Darwin" ]; then \
		sudo codesign -s - /usr/local/bin/git-radar; \
	fi
	@echo "git-radar installed to /usr/local/bin/"

clean:
	rm -rf build/ dist/
