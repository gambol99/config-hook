#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-01-30 11:23:43 +0000 (Fri, 30 Jan 2015)
#
#  vim:ts=2:sw=2:et
#
NAME=config-hook
AUTHOR=gambol99
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')

build:
	godep go build -o stage/${NAME}

docker: build
	docker build -t ${AUTHOR}/${NAME} .

clean:
	rm -f ./stage/${NAME}
	go clean

changelog:
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog

all: clean changelog build docker

release:
	rm -rf release
	mkdir release
	GOOS=linux godep go build -o release/$(NAME)
	cd release && tar -zcf $(NAME)_$(VERSION)_linux_$(HARDWARE).tgz $(NAME)
	go chanagelog
	cp changelog release/changelog
	rm release/$(NAME)

.PHONY: build release changelog

