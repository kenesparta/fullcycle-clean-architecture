FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y protobuf-compiler
RUN make init
RUN make build

FROM scratch
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
COPY --from=builder /app/.mysql.env .
COPY --from=builder /app/.rabbit.env .
CMD ["./server"]
