.PHONY: default
default: bin/cs bin/servedir

VERSION := $(shell git describe --tags 2> /dev/null || echo "unreleased")
V_DIRTY := $(shell git describe --exact-match HEAD 2> /dev/null > /dev/null || echo "-unreleased")

bin/cs: $(shell find  vendor/mcquay.me/cs -type f)
	@echo cs
	@GOOS=linux go build -v -o bin/cs ./vendor/mcquay.me/cs

bin/servedir: $(shell find  vendor/mcquay.me/servedir -type f)
	@echo servedir
	@GOOS=linux go build -v -o bin/servedir ./vendor/mcquay.me/servedir

.PHONY: docker-build
docker-build: bin/cs bin/servedir
	docker build -f Dockerfile  . -t smcquay/alpine:$(VERSION)$(V_DIRTY)
