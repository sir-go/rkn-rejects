# DPI sniffer
[![Go Test cmd/dpi](https://github.com/sir-go/rkn-rejects/actions/workflows/go-dpi.yml/badge.svg)](https://github.com/sir-go/rkn-rejects/actions/workflows/go-dpi.yml)

## What it does

- read all packets from `nf_qeueue`
- check TLS SNI extension, HTTP headers, and the payload of packet
- if IP address is denied or found a denied hostname then the packet marks as "bad"
  end returns to the firewall (it will be rejected)

## Tests
```bash
go test -v ./cmd/dpi-sniffer/...
gosec ./cmd/dpi-sniffer/...
```

### Build
```bash
go mod download
go build -o dni ./cmd/dni-sniffer;
```

### Flags

| key    | default       | description                          |
|--------|---------------|--------------------------------------|
| -nfq   | 200-203       | nf queue num range                   |
| -mdone | 1             | nf mark value for checked packets    |
| -mbad  | 3             | nf mark value for "bad" packets      |
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

