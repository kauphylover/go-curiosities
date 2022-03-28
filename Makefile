SHELL=/usr/local/bin/bash
.DEFAULT_GOAL=run
VERSION=0.1
GIT_COMMIT=$(shell git rev-parse HEAD)

GOLDFLAGS += -X go-curiosities/pkg/version.semver=$(VERSION)
GOLDFLAGS += -X go-curiosities/pkg/version.gitCommit=$(GIT_COMMIT)

hello:
	@echo "Hello"
	@echo $(GIT_COMMIT)

build:
	go build -o bin/main -ldflags "$(GOLDFLAGS)" main.go

run: build
	@bin/main

clean:
	@rm -rf bin/

test:
	@echo "TBD"

