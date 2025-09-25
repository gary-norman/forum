# ----------------------
# Stage 1: Build the app
# ----------------------
FROM golang:1.21 AS builder

WORKDIR /app

# Install build deps (SQLite headers only)
RUN apt-get update && apt-get install -y \
  libsqlite3-dev \
  && rm -rf /var/lib/apt/lists/*

COPY . /app

# Build codex binary (stripped for smaller size)
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/codex ./cmd/server

# ----------------------
# Stage 2: Runtime (distroless)
# ----------------------
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

# Copy binary & required assets
COPY --from=builder /app/codex /app/codex
COPY --from=builder /app/migrations/001_schema.sql /app/migrations/001_schema.sql
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/robots.txt /app/robots.txt
COPY --from=builder /app/README.md /app/README.md
COPY --from=builder /app/LICENSE /app/LICENSE
COPY --from=builder /app/default-images /app/default-images

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh

USER nonroot:nonroot

EXPOSE 8888

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/codex"]
