# Use golang as the base image
FROM golang:1.22.5-alpine

# Set the working directory
WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Copy the Go module files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
#COPY . .

# Command to run the application
CMD ["air", "-c", ".air.toml"]
