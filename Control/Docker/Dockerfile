FROM golang:latest

WORKDIR /app

# Копируем go.mod и go.sum
COPY ../go.mod ../go.sum ./
RUN go mod download

# Копируем весь исходный код проекта
COPY .. .

# Собираем приложение
RUN go build -o app ./main.go

CMD ["./app"]
