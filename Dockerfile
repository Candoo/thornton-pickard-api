# -----------------------------------
# Stage 1: Build the Go Application
# -----------------------------------
FROM golang:1.24-alpine AS builder 

WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Install the Swag tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate the API documentation (creates the 'docs' folder)
RUN swag init --generalInfo ./cmd/api/main.go --output ./docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/api


# -----------------------------------
# Stage 2: Create a Minimal Runtime Image
# -----------------------------------
FROM alpine:latest

# Set necessary timezone and CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /main .

# Expose the port (Gin default)
EXPOSE 8080

# Run the application
CMD ["./main"]