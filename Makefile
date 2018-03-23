.PHONY: default
default: bin/cs bin/servedir

bin/cs: $(shell find  vendor/mcquay.me/cs -type f)
	@echo cs
	@GOOS=linux go build -v -o bin/cs ./vendor/mcquay.me/cs

bin/servedir: $(shell find  vendor/mcquay.me/servedir -type f)
	@echo servedir
	@GOOS=linux go build -v -o bin/servedir ./vendor/mcquay.me/servedir
