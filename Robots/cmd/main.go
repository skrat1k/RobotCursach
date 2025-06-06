// @title RobotService API
// @version 1.0
// @description API for managing robots in the RobotService system.
// @host localhost:8083
// @BasePath /
// @schemes http

package main

import (
	"context"
	"net/http"
	"os"

	"RobotService/internal/handlers"
	"RobotService/internal/prometheusinfo"
	"RobotService/internal/rabbit"
	"RobotService/internal/repositories"
	"RobotService/internal/services"
	"RobotService/internal/storage"
	logger "RobotService/pkg/Logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Отсюда сваггер подсасывает данные для себя
	_ "RobotService/cmd/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	log := logger.GetLogger("dev")

	// Init metrics
	prometheusinfo.Register()

	// Setup dependencies
	db := setupDatabase(log)
	defer db.Close(context.Background())

	cache := storage.NewClient("redis:6379")
	rmq := setupRabbitMQ(log)

	// Init services
	repo := repositories.RobotRepositories{DataBase: db}
	service := services.RobotService{
		RobotRepository: repo,
		Redis:           cache,
		Rabbit:          rmq,
	}
	ctrl := handlers.RobotHandlers{RobotService: service}

	// Инициализация роутера
	router := buildRouter(ctrl)

	log.Info("RobotService is running at http://localhost:8083")
	if err := http.ListenAndServe(":8083", router); err != nil {
		log.Error("Failed to start HTTP server", "error", err.Error())
	}
}

func setupDatabase(log *slog.Logger) *pgx.Conn {
	dbURL := storage.MakeURL(storage.ConnectionInfo{
		Username: "postgres",
		Password: "admin",
		Host:     "postgres",
		Port:     "5432",
		DBName:   "robotdatabase",
		SSLMode:  "disable",
	})

	conn, err := storage.CreatePostgresConnection(dbURL)
	if err != nil {
		log.Error("Unable to connect to PostgreSQL", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Connected to PostgreSQL")
	return conn
}

func setupRabbitMQ(log *slog.Logger) *rabbit.Publisher {
	publisher, err := rabbit.NewPublisher("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Error("Unable to connect to RabbitMQ", "error", err.Error())
		os.Exit(1)
	}
	return publisher
}

func buildRouter(ctrl handlers.RobotHandlers) *chi.Mux {
	r := chi.NewRouter()

	// Инициализация прометеуса
	r.Handle("/metrics", promhttp.Handler())

	// Инициализация сваггера
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8083/swagger/doc.json"),
	))

	// Регистрация эндпоинтов
	ctrl.Register(r)

	return r
}
