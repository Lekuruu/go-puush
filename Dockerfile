FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build flags for cross-compilation
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Build
RUN go build -o /app/puush_api ./cmd/api/main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=build /app/puush_api /app/puush_api

# Set the user to non-root
USER nonroot:nonroot

# Run the compiled binary
ENTRYPOINT ["./puush_api"]