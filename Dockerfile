# Build Stage
FROM golang:1.22-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining project files
COPY . .

# Build the Go app and output it as 'main'
RUN go build -o main ./cmd

# Deployment Stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary has execute permissions
RUN chmod +x main

# Command to run the executable
CMD ["./main"]
