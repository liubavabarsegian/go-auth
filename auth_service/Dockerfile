# Dockerfile для Go-приложения

# Указываем базовый образ для Go
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы в контейнер
COPY . .

# Сборка приложения
RUN go build -o auth_service .

# Указываем команду для запуска
CMD ["./auth_service"]
