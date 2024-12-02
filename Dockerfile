FROM golang:1.23.3 AS build

WORKDIR /build

# Include project files
COPY . .

# Download dependencies (this layer will be cached if go.mod and go.sum haven't changed)
RUN go mod download

# Generate the certificate chain
RUN cd ./pki && ./create-chain.sh

# Build the Go application (binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o tee-mock-server

FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the previous build stage
COPY --from=build /build/tee-mock-server .

# Set the entrypoint to run the binary
CMD ["/app/tee-mock-server"]