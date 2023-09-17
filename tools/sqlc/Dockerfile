FROM golang:1.21-bullseye as builder
WORKDIR /app
RUN go env -w GOBIN=/app
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.21.0


FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/sqlc /app/sqlc
CMD [ "./sqlc", "generate" ]
