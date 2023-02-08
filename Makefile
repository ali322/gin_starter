VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT := $(shell basename "$(PWD)")

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOOS := "linux"
GOARCH := "amd64"

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(Build)"

mirror:
	@go env -w GOPROXY=https://goproxy.cn,direct

install:
	@go get -u

build:
	@echo ">  Building binary"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(PROJECT) *.go
upload:
	@ssh -p 22022 root@1.2.3.4 "sudo systemctl stop api"
	@scp -P 22022 -C api-starter config.yml root@1.2.3.4:/root/app/
	@ssh -p 22022 root@1.2.3.4 "sudo systemctl start api"

.PHONY: install build mirror upload