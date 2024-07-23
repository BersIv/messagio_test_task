package producer

import (
	"encoding/json"
	"messagio_test_task/internal/models"
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer() (*Producer, error) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")

	config := sarama.NewConfig()
	config.Producer.Retry.Backoff = 5 * time.Second

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		slog.Error("Error creating new producer", "error", err)
		return nil, err
	}
	slog.Info("Procuder created")
	return &Producer{producer: producer, topic: topic}, nil
}

func (p *Producer) SendMessage(msg *models.Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal message to JSON", "error", err)
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	}

	_, _, err = p.producer.SendMessage(kafkaMsg)
	if err != nil {
		slog.Error("Failed to send message to Kafka", "error", err)
		return err
	}
	slog.Info("Message sent to Kafka", "message", msg)
	return nil
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return err
	}
	return nil
}
