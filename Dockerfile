# Build stage
FROM golang:1.25 as builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o http-proxy .

# Final stage
FROM alpine:3.22.2

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/http-proxy .

# Expose default port
EXPOSE 8080

# Run the proxy
ENTRYPOINT ["./http-proxy"]
CMD ["-listen", ":8080"]
