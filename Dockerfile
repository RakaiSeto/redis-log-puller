# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o redis-log-puller .

# Stage 2: Run
FROM slim:latest

WORKDIR /app

# Install ca-certificates in case we need to make HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary from the previous stage
COPY --from=builder /app/redis-log-puller .

# Command to run the executable
CMD ["./redis-log-puller"]
