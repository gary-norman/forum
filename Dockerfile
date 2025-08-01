# Use multi-stage build to reduce final image size
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev rsync && rm -rf /var/lib/apt/lists/*

# Copy all files except hidden directories and bin
COPY . /tmp/app
RUN rsync -a --exclude='audit' --exclude='diagrams' --exclude='Dockerfile' --exclude='identifier.sqlite' --exclude='docker-compose.yml' --exclude='.*' --exclude='bin' /tmp/app/ /app && rm -rf /tmp/app

# Build the application
RUN make build

# Use a smaller base image for the final stage
FROM debian:sid-20250317-slim

# Install only necessary runtime dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/bin/codex /app/codex
COPY --from=builder /app/schema.sql /app/schema.sql
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/robots.txt /app/robots.txt
COPY --from=builder /app/README.md /app/README.md
COPY --from=builder /app/LICENSE /app/LICENSE

# Create the database directory and initialize the database
RUN mkdir -p /app/db && sqlite3 /app/db/forum_database.db < /app/schema.sql

# Create the userdata directory structure
RUN mkdir -p /app/db/userdata/images/channel-images
RUN mkdir /app/db/userdata/images/user-images
RUN mkdir /app/db/userdata/images/post-images
# Copy donkey user image
COPY --from=builder /app/db/userdata/images/user-images/K1xcM1YYVT6xkHoEUICS4g==.png /app/db/userdata/images/user-images/K1xcM1YYVT6xkHoEUICS4g==.png
# Copy codex channel image
COPY --from=builder /app/db/userdata/images/channel-images/-ou7VWQ7pK1JHiOvK7lJiw==.png /app/db/userdata/images/channel-images/-ou7VWQ7pK1JHiOvK7lJiw==.png
# Create a base user with a single post
# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use the script as the container entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Run the application
CMD ["/bin/codex"]
