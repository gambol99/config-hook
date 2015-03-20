#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-01-30 11:23:43 +0000 (Fri, 30 Jan 2015)
#
#  vim:ts=2:sw=2:et
#
NAME=config-hook
AUTHOR=gambol99
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')
PWD=$(shell pwd)

.PHONY: build release changelog clean docker test units

build:
	go get github.com/tools/godep
	godep go build -o stage/${NAME}

docker: build
	docker build -t ${AUTHOR}/${NAME} .

clean:
	rm -f ./stage/${NAME}

changelog:
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog

test: build
	go get github.com/stretchr/testify
	godep go test -v ./...

units:
	docker run --rm -v "${PWD}":/go/src/github.com/${AUTHOR}/${NAME} \
  -w /go/src/github.com/${AUTHOR}/${NAME} -e GOOS=linux golang:1.3.3 /bin/bash tests/unit_testing.sh

release:
	rm -rf release
	mkdir release
	GOOS=linux godep go build -o release/$(NAME)
	cd release && tar -zcf $(NAME)_$(VERSION)_linux_$(HARDWARE).tgz $(NAME)
	go chanagelog
	cp changelog release/changelog
	rm release/$(NAME)


