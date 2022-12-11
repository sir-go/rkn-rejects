# Check
[![Go](https://github.com/sir-go/rkn-rejects/actions/workflows/go-check.yml/badge.svg)](https://github.com/sir-go/rkn-rejects/actions/workflows/go-check.yml)

## What it does
- get targets for checking from the redis set
- run bunch of workers
- wait for all workers are done

Each worker is an HTTP client that tries to get data from the target
resource and log the result to the logfile if the target is accessible.

## Tests
```bash
go test -v ./cmd/check/...
gosec ./cmd/check/...
```

## Docker
```bash
docker build . -t check

docker run -it --rm \
  -v /tmp/checks:/var/log/checks \
  --dns 195.208.4.1 \
  --dns 195.208.5.1 \
  check  -w 25  -t 20s  -rh 172.17.0.1  -lt 10s \
  -d /var/log/checks
```

## Build
```bash
go mod download
go build -o check ./cmd/check;
```
If the check will run on the same host that the firewall does,
it should run from the docker container.

## Run
```bash
check -w 25 t 20s  -lt 10s  -d /tmp/checks
```

### Flags

| key  | default   | description                     |
|------|-----------|---------------------------------|
| -rh  | localhost | redis host                      |
| -rp  | 6379      | redis port                      |
| -ra  |           | redis password                  |
| -rd  | 0         | redis DB num                    |
| -rk  | check     | redis set name                  |
| -rtc | 15s       | redis connection timeout        |
| -rtr | 15s       | redis read timeout              |
| -log | info      | log level                       |
| -w   | 10        | workers amount (max 75)         |
| -s   | 5ms       | workers delay between requests  |
| -m   | -1(inf)   | checks amount limit             |
| -t   | 3s        | check TCP timeout               |
| -lt  | 10s       | log polling interval            |
| -d   | /tmp      | path to the logs for each check |
