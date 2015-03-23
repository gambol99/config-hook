#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-01-30 14:51:56 +0000 (Fri, 30 Jan 2015)
#
#  vim:ts=2:sw=2:et
#
FROM progrium/busybox
MAINTAINER Rohith <gambol99@gmail.com>

ADD bin/cfgctl /bin/cfgctl
ADD bin/cfhook /bin/cfgctl

RUN chmod +x /bin/cfgctl && \
    chmod +x /bin/cfhook && \
    mkdir -p /data/config

VOLUME "/data/config"

ENTRYPOINT [ "/bin/cfgctl" ]
