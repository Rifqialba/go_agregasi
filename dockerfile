# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o aggregation-api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/aggregation-api .

COPY config.yaml .
COPY validation_schema.json .

EXPOSE 8080

CMD ["./aggregation-api"]