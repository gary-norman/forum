# Use multi-stage build to reduce final image size
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev rsync && rm -rf /var/lib/apt/lists/*

# Copy all files except hidden directories and bin
COPY . /tmp/app
RUN rsync -a --exclude='.*' --exclude='bin' /tmp/app/ /app && rm -rf /tmp/app

# Build the application
RUN make build

# Use a smaller base image for the final stage
FROM debian:stable-slim

# Install only necessary runtime dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app /app

# Run the application
CMD ["make", "run"]
