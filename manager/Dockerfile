ARG BUILDER_IMAGE=golang:1.24.1-alpine

FROM ${BUILDER_IMAGE} AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd && CGO_ENABLED=1 go build -o manager main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/manager .

RUN apk add --no-cache libc6-compat sqlite-libs

EXPOSE ${MAIN_SERVER_PORT}
EXPOSE ${PROBE_SERVER_PORT}

CMD ["./manager"]