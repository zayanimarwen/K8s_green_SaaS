# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git ca-certificates

COPY go.mod ./
# Generer go.sum depuis go.mod (resout le probleme go.sum vide)
RUN go mod tidy 2>/dev/null || true
RUN go mod download

COPY . .
ARG CMD=api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app ./cmd/${CMD}

# Production stage
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 1001 appuser
WORKDIR /app
COPY --from=builder /build/app .
USER appuser
EXPOSE 9000
ENTRYPOINT ["./app"]

# Dev stage
FROM golang:1.22-alpine AS dev
WORKDIR /app
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest
COPY go.mod ./
RUN go mod tidy 2>/dev/null || true
RUN go mod download
CMD ["air"]
