FROM golang:1.13.0-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -ldflags '-linkmode external -w -extldflags "-static"' -o main .

FROM alpine:3.10

WORKDIR /app


COPY --from=builder /app/main /app/main

CMD ["/app/main"]

