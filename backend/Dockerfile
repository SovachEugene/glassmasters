# Используем образ Golang
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальную часть приложения
COPY . .

# Загружаем и устанавливаем wait-for-it.sh
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Собираем приложение
RUN go build -o main .

# Указываем порт
EXPOSE 8080

# Запускаем приложение через wait-for-it.sh
CMD ["/wait-for-it.sh", "mysql-db:3306", "--", "./main"]
