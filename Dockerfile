# Multi-stage build for Go backend
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# NEW: instalar también wget para que el HEALTHCHECK funcione
RUN apk --no-cache add ca-certificates wget   # NEW

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/main .
COPY .env.dev .env.dev

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
