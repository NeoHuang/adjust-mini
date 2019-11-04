ARG backend_target=adjust_server

FROM golang:1.13

WORKDIR /app
COPY . .
ARG backend_target

WORKDIR /app/${backend_target}/run
RUN go build -o main .

EXPOSE 80

CMD ["./main"]
