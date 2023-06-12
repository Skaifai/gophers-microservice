# Use the official Golang image as the base image
FROM golang:1.16-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product-service ./cmd/product-service

# Create a new stage for the final production image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the previous stage
COPY --from=builder /app/product-service .

# Expose the port on which the microservice will run
EXPOSE 8080

# Run the microservice binary
CMD ["./product-service"]
