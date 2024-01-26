# Start from the official Go image to build our application.
FROM golang:1.21-alpine AS builder

# Set the working directory outside GOPATH to enable the support for modules.
WORKDIR /app

# Copy go.mod and go.sum to download all dependencies.
COPY go.* ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the application.
RUN go build -o main .

# Use a smaller base image to create a final production image.
FROM alpine:latest  
WORKDIR /root/

# Copy the pre-built binary from the builder stage.
COPY --from=builder /app/main .

# Command to run the executable.
CMD ["./main"]
