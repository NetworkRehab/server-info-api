
# Build stage
FROM golang:1.23.4-bookwork AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o server-info-api main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/server-info-api .

# Expose the application port
EXPOSE 8080

# Set environment variable for the port
ENV PORT=8080

# Run the executable
CMD ["./server-info-api"]