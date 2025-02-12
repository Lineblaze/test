FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o avito_test ./cmd/app

FROM alpine:latest

COPY --from=builder /app/avito_test /avito_test

EXPOSE 8080

CMD ["/avito_test"]
