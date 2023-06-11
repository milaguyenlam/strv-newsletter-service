# Start from a Debian based image with Go installed
FROM golang:1.20

# Create a directory for the app
WORKDIR /go/src/app

# Copy over the go mod and sum files
COPY go.mod go.sum ./

# Download any necessary dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]