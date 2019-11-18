ARG backend_target=adjust_server
ARG version=unknown

FROM golang:1.13 as builder

WORKDIR /app
COPY . .
ARG backend_target
ARG version

WORKDIR /app/${backend_target}/run
ENV VERSION=${version}
RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-X main.Version=${VERSION}" -o main .



######## Start a new stage from scratch #######
FROM alpine:latest
ARG backend_target

RUN apk --no-cache add ca-certificates

WORKDIR /app/${backend_target}/run

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/${backend_target}/run/main .

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./main"]

