# Dockerfile for Rate-Limited API Server

# Use the official Golang image as the base image
FROM golang:1.17

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application inside the container
RUN go build -o rate-limited-api-server

# Expose the port on which the server will listen
EXPOSE 8080

# Command to run the server when the container starts
CMD ["./rate-limited-api-server"]
