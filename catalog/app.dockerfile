# Базовый образ на котором будет строится контейнер - golang:1.13-alpine3.11 (alpine - минималистический образ linux)
FROM golang:1.13-alpine3.11 AS build 

RUN apk --no-cache add gcc g++ make ca-certificates
# Указывается рабочая директория при запуске контейнера, откуда будут выполняться все команды
WORKDIR /go/src/github.com/GrudTrigger/go-grpc-graphql-microservice
# Копирование файлов из проекта в контейнер
COPY go.mod go.sum ./
# Копирование файлов из проекта в контейнер
COPY vendor vendor
# Копирование файлов из проекта в контейнер
COPY catalog catalog

RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
# Порт для контейнера
EXPOSE 8080
# запуск app при запуске контейнера
CMD ["app"]

#Создание образа
# docker build -t (-t позволяет задать имя образа) server.app .(. указывает относительный путь к докер файлу)