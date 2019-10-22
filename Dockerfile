# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13 AS build_base
WORKDIR /go/src/github.com/dev10/fantom-asset-management

# Add Maintainer Info
LABEL maintainer="Peter Badenhorst <peter.badenhorst@fantom.foundation>"

# Copy go mod and sum files
COPY go.mod go.sum ./
# Force the go compile to use modules
ENV GO111MODULE=on
# use Docker layer caching system to only get changed dependencies
RUN go mod download

# This image builds the actual code
FROM build_base AS fam_builder

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
#RUN make install
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -mod=readonly -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/famcli
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -mod=readonly -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/famd

# Fresh minimal image
FROM alpine AS fam
COPY --from=fam_builder /go/bin/famd /bin/famd
COPY --from=fam_builder /go/bin/famcli /bin/famcli

# Expose port(s) to the outside world
EXPOSE 26656

# Command to run the executable
CMD ["/bin/famd", "start", "--trace"]
