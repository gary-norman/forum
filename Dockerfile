# Start with a lightweight base image
FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download
# Install SQLite
RUN apk update && apk add sqlite
 
Copy . .
RUN make build run
# Set the default command to open SQLite
# CMD ["sqlite3"]
