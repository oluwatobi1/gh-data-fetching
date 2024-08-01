export PATH := $(shell go env GOPATH)/bin:$(PATH)

fmt:
	go fmt ./...
