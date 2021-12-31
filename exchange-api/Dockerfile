FROM golang:1.17-buster as builder

ARG APP_NAME

WORKDIR /exchange-api

COPY go.* ./

RUN go mod download

COPY . ./

RUN go build -v -o executable ./transporters/${APP_NAME}

FROM debian:buster-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /exchange-api/executable /executable

CMD ["./executable"]
