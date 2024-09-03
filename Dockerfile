# Stage 1: Build the application
FROM golang:1.21 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o bin/api cmd/api/main.go

# Stage 2: Create the final image
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./api"]
