# ----------------------
# Stage 1: Build the app
# ----------------------
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Install dependencies
RUN  apt-get update && apt-get install -y sqlite3 libsqlite3-dev dos2unix rsync && rm -rf /var/lib/apt/lists/*

# Copy all files except hidden directories and bin
COPY . /app
COPY --chmod=755 entrypoint.sh /entrypoint.sh
COPY .env /app/.env
RUN dos2unix /entrypoint.sh

# Build the application
RUN make build

# Use a smaller base image for the final stage
FROM debian:sid-20250317-slim

# Install only necessary runtime dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0  && rm -rf /var/lib/apt/lists/*


# Set the working directory
WORKDIR /app

# ----------------------------------------------------------
# Stage 2: Copy the built application from the builder stage
# ----------------------------------------------------------
COPY --from=builder /app/bin/codex /app/bin/codex
COPY --from=builder /app/migrations/001_schema.sql /app/migrations/001_schema.sql
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/robots.txt /app/robots.txt

# Create the database directory and initialize the database
RUN mkdir -p /var/lib/db-codex && sqlite3 /var/lib/db-codex/forum_database.db < /app/migrations/001_schema.sql

# Create the userdata directory structure
RUN mkdir -p /app/db/userdata/images/{channel-images,user-images,post-images}
# Copy donkey user image
COPY --from=builder /app/default-images/donkey.png /app/db/userdata/images/user-images/donkey.png
# Copy codex channel image
COPY --from=builder /app/default-images/codex.png /app/db/userdata/images/channel-images/codex.png

# Create a base user with a single post
# Copy the entrypoint script
COPY --from=builder /entrypoint.sh /entrypoint.sh


COPY --from=builder /app/.env /app/.env

EXPOSE 8888

# Use the script as the container entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Run the application
CMD ["/app/bin/codex"]
