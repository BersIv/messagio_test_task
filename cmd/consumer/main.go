package main

import (
	"log"
	"log/slog"
	"messagio_test_task/db"
	"messagio_test_task/internal/consumer"
	"messagio_test_task/internal/message"
	"os"

	"github.com/joho/godotenv"
)

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

	consumer.StartConsumerGroup(message.NewMessageService(message.NewMessageRepository(dbConn.GetDB())))
}
