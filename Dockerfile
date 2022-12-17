FROM alpine:3.17.0

COPY bin/linux/amd64/myhttpserver /myhttpserver
EXPOSE 8888
ENTRYPOINT /myhttpserver
