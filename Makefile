GO ?= go
BUILD := container0

all: linux

linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD)
