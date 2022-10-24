export GO111MODULE=on
GOPROXY ?= https://proxy.golang.org,direct
export GOPROXY

CURRENT_DIR = $(shell pwd)
BUILD_DIR = ${CURRENT_DIR}/bin

.PHONY: default

check_musl:
	CC=musl-gcc CGO_ENABLED=1 \
	go build -o ${BUILD_DIR}/check rejects-go/cmd/check;

check_gcc:
	go build -o ${BUILD_DIR}/check rejects-go/cmd/check;

get_rkn:
	go build -o ${BUILD_DIR}/get_rkn rejects-go/cmd/check;

dns:
	go build -o ${BUILD_DIR}/dns rejects-go/cmd/dns-sniffer;

dpi:
	go build -o ${BUILD_DIR}/dpi rejects-go/cmd/dpi-sniffer;

default: all
all: check_musl get_rkn dns dpi
