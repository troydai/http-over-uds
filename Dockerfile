FROM golang:1.19-alpine3.16 AS builder

RUN apk add make

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY ./ ./
RUN mkdir -p /artifacts
RUN make build

FROM alpine:3.16 AS server

WORKDIR /app
RUN mkdir -p /app/socket
COPY --from=builder /src/artifacts/bin/server ./server

ENTRYPOINT [ "/app/server" ]

FROM alpine:3.16 AS client

RUN apk add busybox curl

CMD exec /bin/sh -c "trap : TERM INT; sleep 9999999999d & wait"
