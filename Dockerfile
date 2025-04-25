# Use an official Go runtime as a parent image
FROM golang:1.24.2-bookworm

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the entire code into the container
COPY . .

# Build the Go app inside the container
RUN go build -o main .

# Expose port 8080
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
RUN apt-get update && apt-get install -y postgresql-client

ENV DB_URL=${DB_URL}
ENV PORT=8081
