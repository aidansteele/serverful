FROM golang:1.20 AS builder

ENV CGO_ENABLED=0

WORKDIR /serverful
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w"

FROM scratch AS ext
COPY --from=builder /serverful/serverful /opt/extensions/serverful
