FROM golang:1.20 AS builder

ENV CGO_ENABLED=0

WORKDIR /example
RUN go mod init example
COPY . .
RUN go build -ldflags="-s -w"

FROM gcr.io/distroless/static-debian11
COPY --from=ghcr.io/aidansteele/serverful:main /opt/extensions/serverful /opt/extensions/serverful
COPY --from=builder /example/example /bin/example
ENTRYPOINT ["/bin/example"]
