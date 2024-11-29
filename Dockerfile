FROM golang:1.23.3 AS build

WORKDIR /app

# Include project files
COPY . .

# Download dependencies (this layer will be cached if go.mod and go.sum haven't changed)
RUN go mod download

# Generate the certificate chain
RUN ./pki/create-chain.sh

# Build the Go application (binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tee-server-mock

FROM alpine:latest

# Install any necessary runtime dependencies (if needed)
# RUN apk --no-cache add <runtime-dependency>

# Set the working directory
WORKDIR /

# Copy the built binary from the previous build stage
COPY --from=build /app/tee-server-mock .

# Set the entrypoint to run the binary
CMD ["/tee-server-mock"]