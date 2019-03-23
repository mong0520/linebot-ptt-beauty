# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

build:
	docker build -t mong0520/linebot-ptt-beauty .

release:
	docker push mong0520/linebot-ptt-beauty

debug:
	docker run -it --rm -v $(PWD)/ssl:/ssl mong0520/linebot-ptt-beauty /bin/bash
