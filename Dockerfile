FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev rsync && rm -rf /var/lib/apt/lists/*

# Copy all files except hidden directories and bin
COPY . /tmp/app
RUN rsync -a --exclude='.*' --exclude='bin' /tmp/app/ /app && rm -rf /tmp/app

# Run make build and make run
CMD ["make", "build", "run"]
