# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

build:
	docker run --rm -v $(GOPATH):/go -w /go/src/gitlab.kkinternal.com/neilwei/project-lambda-template golang:1.10 go build

debug:
	docker run -it --rm -v $(GOPATH):/go -w /go/src/github.com/mong0520/linebot-ptt-beauty golang:1.10 /bin/bash
