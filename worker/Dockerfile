ARG BUILDER_IMAGE=golang:1.24.1-alpine

FROM ${BUILDER_IMAGE} AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd && go build -o worker main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/worker .

EXPOSE ${SERVER_PORT}

CMD ["./worker"]