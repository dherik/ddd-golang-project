# Use the official Go image as a base
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# Create a minimal runtime image
# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/main .

# Expose a port if your application listens on a specific port
# EXPOSE 8080

# Run the application
CMD ["/main"]
