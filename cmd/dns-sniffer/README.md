# DNS sniffer
[![Go Test cmd/dns](https://github.com/sir-go/rkn-rejects/actions/workflows/go-dns.yml/badge.svg)](https://github.com/sir-go/rkn-rejects/actions/workflows/go-dns.yml)

## What it does
- read all DNS answer packets from `nf_qeueue`
- check a hostname in the answer
- add IPs from A-record to the allowed nftables set with TTL if all 
  IP addresses in the answer and
  the hostname are not found in the denied lists

## Tests
```bash
go test -v ./cmd/dns-sniffer/...
gosec ./cmd/dns-sniffer/...
```

## Build
```bash
go mod download
go build -o dns ./cmd/dns-sniffer;
```

## Run
### Flags
| key    | default       | description                          |
|--------|---------------|--------------------------------------|
| -nfq   | 100-103       | nf queue num range                   |
| -mdone | 1             | nf mark value                        |
| -nfql  | 0xFF          | max nf queue capacity                |
| -nft   | rkn           | nftables table name                  |
| -nfs   | allow_sniffed | nftables allowed set name            |
| -rh    | localhost     | redis host                           |
| -rp    | 6379          | redis port                           |
| -ra    |               | redis password                       |
| -rd    | 0             | redis DB num                         |
| -rk    | domains       | redis domains set name               |
| -rtc   | 15s           | redis connection timeout             |
| -rtr   | 15s           | redis read timeout                   |
| -log   | info          | log level                            |
| -dry   | false         | print the config without any actions |

