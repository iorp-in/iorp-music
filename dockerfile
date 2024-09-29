# Start from the Alpine base image
FROM golang:1.21-alpine as builder

# Set app directory
WORKDIR /app

# Copy files
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY music.go ./

# Build
RUN go build -o music

# Start from the Alpine base image
FROM alpine:latest

# Install FFmpeg
RUN apk add --no-cache ffmpeg

# Copy your executable library into the container
COPY --from=builder /app/music /usr/local/bin/music

# Set permissions for the executable
RUN chmod +x /usr/local/bin/music

# Set the command to run your application
CMD ["/usr/local/bin/music"]
