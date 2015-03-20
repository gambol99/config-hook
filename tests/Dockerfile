#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-03-20 15:56:51 +0000 (Fri, 20 Mar 2015)
#
#  vim:ts=2:sw=2:et
#
FROM golang:1.3.3
MAINTAINER Rohith <gambol99@gmail.com>

ADD ./setup.sh /bin/setup.sh
RUN chmod +x /bin/setup.sh
RUN curl -skL https://github.com/coreos/etcd/releases/download/v2.0.0/etcd-v2.0.0-linux-amd64.tar.gz > /opt/etcd.tar.gz && \
    tar zxvf /opt/etcd.tar.gz -C /opt && \
    go get github.com/tools/godep && \
    go get github.com/stretchr/testify

WORKDIR [ "/go/src/github.com/gambol99/config-hook" ]
ENTRYPOINT [ "/bin/setup.sh" ]
