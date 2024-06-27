FROM alpine:3.20

RUN apk add createrepo_c --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing
RUN mkdir -p /repo
ADD ./yumster.yml /
ADD ./yumster /

##RUN grep -q -E '^import socket$' /usr/lib/python2.7/site-packages/createrepo/__init__.py || sed -r -i  -e "/^import subprocess$/a import socket" /usr/lib/python2.7/site-packages/createrepo/__init__.py
##RUN sed -r -i  -e "s|self.olddir = '.olddata'|self.olddir = '.olddata-' + socket.gethostname()|g"  /usr/lib/python2.7/site-packages/createrepo/__init__.py

EXPOSE 8080

CMD ["/yumster"]
