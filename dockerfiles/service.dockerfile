# Use a more recent version of the Go base image
FROM golang:1.21-alpine AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files separately to leverage Docker layer caching
COPY ../go.mod .
COPY ../go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the source code
COPY .. .

# Build the application
RUN go build -o reptr ./service/cmd/server/main.go

# Use a smaller base image for the final runtime image
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=build /app/reptr .

# Expose the application port
EXPOSE 8080

# Set the command to run the application
CMD ["./reptr"]
