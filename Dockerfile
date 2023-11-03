# Этап 1: Сборка приложения
FROM golang:latest AS build
ENV GOOS=linux GOARCH=amd64
WORKDIR /app

# Копирование и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода приложения
COPY . .

# Компиляция приложения
RUN go build -o main .
RUN chmod +x main

# Этап 2: Перенос исполняемого файла в новый контейнер
FROM alpine

WORKDIR /app

# Копирование исполняемого файла из этапа сборки
COPY --from=build /app/main /app/

# Запуск приложения
CMD ["./main"]