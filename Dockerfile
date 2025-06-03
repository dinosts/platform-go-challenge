# Build stage
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Minimal image
FROM alpine:latest

ENV APP_ENV="dev"
ENV JWT_SECRET_KEY="GWI_CHALLENGE"
ENV HASHING_SALT="SALTY"

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api .

# Expose the port your app listens on (change if needed)
EXPOSE 3008

# Command to run the executable
CMD ["./api"]
