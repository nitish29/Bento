# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY bot ./bot

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bento -x

# Stage 2: Runtime
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /bento /bento

# Run
CMD ["/bento"]
