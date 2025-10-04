# -------- Stage 1: Build --------
FROM golang:1.24.2 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app (static binary, no CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nasa-hdf-processor .

# -------- Stage 2: Runtime --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/nasa-hdf-processor .


VOLUME ["/data"]

# Run the application
ENTRYPOINT ["./nasa-hdf-processor"]
