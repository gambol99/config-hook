#!/bin/sh
#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-01-30 14:54:23 +0000 (Fri, 30 Jan 2015)
#
#  vim:ts=2:sw=2:et
#
NAME="Config Hook Service"
STORE=${STORE:-etcd://localhost:4001}
HOOK_PREFIX=${HOOK_PREFIX:-_HOOK_}
DOCKER=${DOCKER:-/var/run/docker.sock}
VERBOSE=${VERBOSE:-3}

annonce() {
    [ -n "$1" ] && {
        echo "** $@";
    }
}

annonce "Starting the ${NAME}"

/bin/config-hook -logtostderr=true -v=${VERBOSE} -prefix=${HOOK_PREFIX} -docker=${DOCKER} -store=${STORE}
