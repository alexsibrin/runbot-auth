# Сначала укажем образ для сборки
FROM golang:1.21.6 AS build
LABEL authors="alex sibrin"

WORKDIR /temp

# Копируем все файлы проекта в рабочую директорию
COPY . .

RUN go get ./...

# Собираем приложение
RUN go build -o bin/app ./cmd/main.go

# Затем создаем конечный образ
FROM debian:12-slim
LABEL authors="alex sibrin"

WORKDIR /appbin

# Копируем только исполняемый файл из предыдущего образа
COPY --from=build /temp/bin/app /temp/config.yaml ./

# Задаем команду для запуска приложения
CMD ["./app"]