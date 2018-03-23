FROM scratch
MAINTAINER Tony Allen <cyril0allen@gmail.com>

ADD mysql-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/mysql-operator"]
