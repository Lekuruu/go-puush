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

# Build
RUN CGO_ENABLED=1 go build -o puush ./cmd/api/main.go

FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Run as non‑root
RUN addgroup -S app && adduser -S -G app app
USER app

WORKDIR /app
COPY --from=build /app/puush /app/puush

# Create data volume
VOLUME ["/app/.data"]

# Run the compiled binary
ENTRYPOINT ["/app/puush"]