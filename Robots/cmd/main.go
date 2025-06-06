// @title RobotService API
// @version 1.0
// @description API for managing robots in the RobotService system.
// @host localhost:8083
// @BasePath /
// @schemes http

package main

import (
	"RobotService/internal/handler"
	"RobotService/internal/metrics"
	"RobotService/internal/rabbit"
	"RobotService/internal/repositories"
	"RobotService/internal/services"
	"RobotService/internal/storage"
	logger "RobotService/pkg/Logger"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "RobotService/cmd/docs" // импорт сгенерированных Swagger-доков

	httpSwagger "github.com/swaggo/http-swagger" // swagger UI
)

func main() {
	log := logger.GetLogger("dev")

	metrics.Register()

	psqlConnectionUrl := storage.MakeURL(storage.ConnectionInfo{
		Username: "postgres",
		Password: "admin",
		Host:     "localhost",
		Port:     "5432",
		DBName:   "robotdatabase",
		SSLMode:  "disable",
	})

	conn, err := storage.CreatePostgresConnection(psqlConnectionUrl)

	if err != nil {
		log.Error("Connection error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	redis := storage.NewClient("localhost:6379")

	log.Info("Success connect to database")

	rabbitPublisher, err := rabbit.NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Error("Failed connected to rabbitmq", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	// Swagger docs route
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8083/swagger/doc.json"),
	))

	expenseRepo := repositories.RobotRepositories{DataBase: conn}
	expenseService := services.RobotService{RobotRepository: expenseRepo, Redis: redis, Rabbit: rabbitPublisher}
	expenseHandler := handler.RobotHandlers{RobotService: expenseService}
	expenseHandler.Register(router)

	log.Info("Server starting...")

	serverPort := ":8083"

	err = http.ListenAndServe(serverPort, router)
	if err != nil {
		log.Error("Starting server error", slog.String("error", err.Error()))
	}
}
