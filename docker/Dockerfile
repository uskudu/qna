FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/app

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/app .
# Copy goose binary for migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Copy migrations directory for migrate service
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

CMD ["./app"]
