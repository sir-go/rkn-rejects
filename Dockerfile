FROM golang:1.18-alpine3.16 as builder
WORKDIR /go/src/app
COPY go.* ./
COPY cmd/check .
RUN go mod download
RUN CC=musl-gcc CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags '-w -s -extldflags "-static"' -o /app .

FROM alpine:3.16
COPY --from=builder /app /app
RUN apk update && \
    apk add --no-cache tzdata
ENV TZ="Europe/Moscow"
RUN mkdir /var/log/checks
ENTRYPOINT ["/app"]
