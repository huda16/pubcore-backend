FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3 && \
    swag init && \
    CGO_ENABLED=0 GOOS=linux go build -o pubcore-server .

# Runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/pubcore-server .
COPY --from=builder /app/.env.example .env.example

EXPOSE 8080

CMD ["./pubcore-server"]
