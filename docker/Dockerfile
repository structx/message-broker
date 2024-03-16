FROM golang:1.22-bookworm as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/bin/server ./cmd/server

FROM gcr.io/distroless/static-debian11

COPY --from=builder /usr/bin/ /app/bin/

VOLUME [ "/var/lib/message_broker/kv", "/var/log/message_broker" ]

EXPOSE 2112 8080

ENTRYPOINT [ "/app/bin/server" ]