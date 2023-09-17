FROM golang:1.21-bullseye as builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOFLAGS="-buildmode=pie"

COPY go.mod go.sum /app

RUN go mod download

COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./ui /app/ui

RUN go build -ldflags "-s -w" \
    -gcflags=all="-trimpath $(pwd)" \
    -asmflags=all="-trimpath $(pwd)" \
    -o /app/server \
    ./cmd/server



FROM gcr.io/distroless/base-debian11:nonroot
WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8080
EXPOSE 8081

USER 65534

ENTRYPOINT [ "/app/server" ]
