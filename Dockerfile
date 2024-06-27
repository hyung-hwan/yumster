FROM alpine:3.20

RUN apk add createrepo_c --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing
RUN mkdir -p /repo
ADD ./yumster.yml /
ADD ./yumster /

EXPOSE 8080

CMD ["/yumster"]
