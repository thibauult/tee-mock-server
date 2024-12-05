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

# Env variables
ENV TEE_GOOGLE_SERVICE_ACCOUNT=tee-mock-server@localhost.gserviceaccount.com
ENV TEE_TOKEN_EXPIRATION_IN_MINUTES=5

# Set the working directory
WORKDIR /app

# Copy the built binary from the previous build stage
COPY --from=build /build/tee-mock-server .

# Set the entrypoint to run the binary
CMD ["/app/tee-mock-server"]