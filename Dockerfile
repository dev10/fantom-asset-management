# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13

# Add Maintainer Info
LABEL maintainer="Peter Badenhorst <peter.badenhorst@fantom.foundation>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
ENV GOBIN=/app
RUN make install

# Expose port 8080 to the outside world
EXPOSE 26656

# Command to run the executable
CMD ["./famd", "start", "--trace"]
