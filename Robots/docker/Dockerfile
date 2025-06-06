FROM golang:latest

WORKDIR /app

# Копируем go.mod и go.sum из родительской директории
COPY ../go.mod ../go.sum ./
RUN go mod download

# Копируем весь исходный код проекта
COPY .. .

# Собираем приложение (если main.go находится в cmd/, укажи путь)
RUN go build -o app ./cmd/main.go

CMD ["./app"]
