# Build stage
FROM golang:1.23.4-bookworm AS builder

WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with correct path
RUN CGO_ENABLED=1 go build -o server-info-api main.go

# Final stage
FROM alpine:latest AS final

# Install required dependencies for SQLite
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Create directory for SQLite database
RUN mkdir -p /app/data

# Copy the binary from the builder stage to the correct location
COPY --from=builder /build/server-info-api /app/
COPY --from=builder /build/data.csv /app/data

# ensure executable is ...well .. executable. 
RUN chmod +x /app/server-info-api

# Expose the application port
EXPOSE 8080

# Set environment variable for the port
ENV PORT=8080

# Run the executable from the correct path
CMD ["/app/server-info-api"]