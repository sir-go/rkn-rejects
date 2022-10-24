#!/usr/bin/env bash
docker run -it --rm \
  -v /var/log/checks:/var/log/checks \
  --dns 195.208.4.1 \
  --dns 195.208.5.1 \
  z-client \
  -w 25 \
  -t 20s \
  -rh 172.17.0.1 \
  -lt 10s \
  -d /var/log/checks
