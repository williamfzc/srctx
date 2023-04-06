# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

default:
	go build ./cmd/srctx

# linux
build_linux_amd64:
	GOOS=linux GOARCH=amd64 ${GOCMD} build -o srctx_linux_amd64 ./cmd/srctx

# windows
build_windows_amd64:
	GOOS=windows GOARCH=amd64 ${GOCMD} build -o srctx_windows_amd64.exe ./cmd/srctx

# mac
build_macos_amd64:
	GOOS=darwin GOARCH=amd64 ${GOCMD} build -o srctx_macos_amd64 ./cmd/srctx
build_macos_arm64:
	GOOS=darwin GOARCH=arm64 ${GOCMD} build -o srctx_macos_arm64 ./cmd/srctx

test:
	$(GOTEST) ./...
