#!/bin/bash
#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-03-20 14:13:45 +0000 (Fri, 20 Mar 2015)
#
#  vim:ts=2:sw=2:et
#

ETCD_PACKAGE="https://github.com/coreos/etcd/releases/download"
ETCD_VERSION="v2.0.0"
ETCD_URL="${ETCD_PACKAGE}/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz"
TEMPDIR=$(mktemp -d)

say() {
  [ -n "$1" ] && echo "*** $@"
}

failed() {
  say "$1"
  exit 1
}

extract() {
  filename="${1}"
  directory="${2}"
  [ -z "${filename}"  ] && failed "you have not specified a filename to extract"
  [ -f "${filename}"  ] || failed "the file: ${filename} does not exist"
  [ -d "${directory}" ] || failed "the directory: ${directory} does not exist"

  if [[ "${filename}" =~ ^.*tar.gz$ ]]; then
    say "extracting the file: ${filename} to: ${directory}"
    tar -zxf "${filename}" -C "${directory}" || failed "unable to extract the file: ${filename}"
  fi
}

download() {
  url="${1}"
  dest="${2}"
  say "attempting to download package: ${url} to: ${dest}"
  curl -skL ${url} > ${dest} || failed "unable to download the package"
}

run() {
  command="${1}"
  [ -n "${command}" ] || failed "you have not specified anything to run"
  say "executing command: ${command}"
  eval "${command}" || failed "failed to execute: ${command}"
}

setup_etcd() {
  download "${ETCD_URL}" "${TEMPDIR}/etcd.tar.gz"
  extract "${TEMPDIR}/etcd.tar.gz" "${TEMPDIR}"
  run "${TEMPDIR}/etcd*/etcd -data-dir=/tmp -bind-addr=127.0.0.1:4001 >/dev/null 2>&1 &"
}

main() {
  setup_etcd
  make test
}

main "$@"
