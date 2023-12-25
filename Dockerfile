FROM golang:1.21.5-bookworm as builder
WORKDIR /app
COPY . .
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN apt-get update && apt-get install -y protobuf-compiler
RUN make init && make build

FROM scratch
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
COPY --from=builder /app/.mysql.env .
COPY --from=builder /app/.rabbit.env .
CMD ["./server"]
