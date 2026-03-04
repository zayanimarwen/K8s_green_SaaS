# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git ca-certificates

ENV GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod GOPATH=/go

COPY . .
RUN go mod download -x 2>&1 | tail -5 || true
ARG CMD=api
RUN CGO_ENABLED=0 GOOS=linux GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod \
    go build -ldflags="-w -s" -o app ./cmd/${CMD}

# Production stage
FROM golang:1.22-alpine AS prod
RUN adduser -D -u 1001 appuser
WORKDIR /app
COPY --from=builder /build/app .
USER appuser
EXPOSE 9000
ENTRYPOINT ["./app"]
