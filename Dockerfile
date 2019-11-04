ARG backend_target=adjust_server

FROM golang:1.13

WORKDIR /app
COPY . .
ARG backend_target

RUN cd ${backend_target}/run

RUN go build -o ${backend_target}
EXPOSE 80

CMD ./${backend_target}
