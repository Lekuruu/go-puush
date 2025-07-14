FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 go build -o puush_api ./cmd/api/main.go

FROM gcr.io/distroless/static:nonroot

# Set the user to non-root
USER nonroot:nonroot

WORKDIR /app
COPY --from=build /app/puush /app/puush

# Create data volume
VOLUME ["/app/.data"]

# Run the compiled binary
ENTRYPOINT ["/app/puush"]