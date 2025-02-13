FROM golang:1.23.5 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o stress-test-command main.go

FROM alpine:latest

# Copy binary from builder
COPY --from=builder /app/stress-test-command /stress-test-command

# Run
ENTRYPOINT ["/stress-test-command"]