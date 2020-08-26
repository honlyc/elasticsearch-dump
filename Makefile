
SHELL = /bin/bash

GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=es-dump
BINARY_UNIX=$(BINARY_NAME)-unix

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./$(BINARY_UNIX) -v ./es-dump.go