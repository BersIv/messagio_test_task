package main

import (
	"log"
	"log/slog"
	"messagio_test_task/db"
	"messagio_test_task/internal/message"
	"messagio_test_task/internal/producer"
	"messagio_test_task/router"
	"os"

	"github.com/joho/godotenv"
)

// @title Messagio Test Task API
// @version 1.0
// @description Server to create messages in postgres and kafka
// @host 194.247.187.44:5000
// @BasePath /
func main() {
	file, err := os.Create("logs.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env", "error", err)
	}

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}
	defer dbConn.Close()

	kafkaProducer, err := producer.NewProducer()
	if err != nil {
		log.Fatalf("Could not create kafka produced: %s", err)
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Fatalf("Failed to close Kafka producer: %s", err)
		}
	}()
	messageHandler := message.NewMessageHandler(message.NewMessageService(message.NewMessageRepository(dbConn.GetDB())), kafkaProducer)

	r := router.InitRouter(
		router.MessageRouter(messageHandler),
	)

	slog.Info("Server started")

	if err := router.Start("0.0.0.0:5000", r); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
