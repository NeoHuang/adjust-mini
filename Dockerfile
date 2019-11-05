ARG backend_target=adjust_server
ARG version=unknown

FROM golang:1.13

WORKDIR /app
COPY . .
ARG backend_target
ARG version

WORKDIR /app/${backend_target}/run
ENV VERSION=${version}
RUN go build --ldflags "-X main.Version=${VERSION}" -o main .

EXPOSE 80

CMD ["./main"]
