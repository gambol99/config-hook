#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-01-30 14:51:56 +0000 (Fri, 30 Jan 2015)
#
#  vim:ts=2:sw=2:et
#
FROM progrium/busybox
MAINTAINER Rohith <gambol99@gmail.com>

ADD stage/config-hook /bin/config-hook
ADD stage/startup.sh /startup.sh
RUN chmod +x /bin/config-hook && chmod +x /startup.sh

ENRTYPOINT [ "/startup.sh" ]
