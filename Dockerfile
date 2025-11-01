FROM golang:1.24-alpine AS build

# Install C toolchain + sqlite3 headers
RUN apk add --no-cache \
      gcc \
      musl-dev \
      sqlite-dev

WORKDIR /app

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the server
RUN CGO_ENABLED=1 go build -o puush ./cmd/standalone/main.go

FROM alpine AS app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app
COPY --from=build /app/puush /app/puush
COPY --from=build /app/web /app/web

# Create data volume
VOLUME ["/app/.data"]

# Run the compiled binary
ENTRYPOINT ["/app/puush"]