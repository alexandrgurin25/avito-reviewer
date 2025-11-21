FROM golang:1.25.3-alpine

# Устанавливаем рабочую директорию
WORKDIR /reviwers

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Переходим в директорию с main.go и собираем приложение
RUN cd cmd/app && go build -o reviwers .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./cmd/app/reviwers"]