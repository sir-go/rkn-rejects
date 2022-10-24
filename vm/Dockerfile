FROM alpine:latest

WORKDIR /opt
ADD bin/check .
RUN apk update && \
    apk add --no-cache tzdata
ENV TZ="Europe/Moscow"

RUN mkdir /var/log/checks

ENTRYPOINT ["/opt/check"]
