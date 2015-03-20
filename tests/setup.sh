#!/bin/bash
#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-03-20 15:56:51 +0000 (Fri, 20 Mar 2015)
#
#  vim:ts=2:sw=2:et
#
ETCD_VERSION="v2.0.0"
ETCD_DIR="/opt/etcd*"
WORKDIR="/go/src/github.com/gambol99/config-hook"

say() {
    [ -n "$1" ] || echo "** $1"
}

failed() {
    say "[ERROR] $1"
    exit 1
}

say "Starting Etcd Service"
cd /opt/etcd-${ETCD_VERSION}-linux-amd64 || failed "Unable to move into etcd directory"
./etcd >/dev/null 2&>1 &
[ $? -ne 0 ] && failed "Unable to start the Etcd Service"

cd $WORKDIR
say "Running Unit tests"
godep go test -v ./... $@
