FROM golang:latest

WORKDIR /app

# Копируем зависимости
COPY ../go.mod ../go.sum ./
RUN go mod download

# Копируем исходники
COPY .. .

# Билдим
RUN go build -o app ./cmd/main.go

CMD ["./app"]
