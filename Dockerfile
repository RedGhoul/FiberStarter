# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fiberstarter .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/fiberstarter .
COPY --from=builder /app/views ./views
COPY --from=builder /app/public ./public
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./fiberstarter"]
