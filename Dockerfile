# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies (cached between builds)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source code
COPY . .

# Build the server with cached build artifacts
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o puush ./cmd/puush/

FROM alpine AS app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /app/puush /app/puush

# Create web directory volume
VOLUME ["/app/web"]

# Create data volume
VOLUME ["/app/.data"]

# Run the compiled binary
ENTRYPOINT ["/app/puush"]