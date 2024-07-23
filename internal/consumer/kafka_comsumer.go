package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"messagio_test_task/internal/message"
	"messagio_test_task/internal/models"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	service message.MessageService
}

func NewConsumerGroupHandler(service message.MessageService) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{service: service}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type Consumer struct {
	ready chan bool
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				slog.Info("Message channel was closed")
				return nil
			}
			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				slog.Error("Failed to unmarshal message", "error", err)
				continue
			}
			if err := h.service.UpdateMessage(context.Background(), &message); err != nil {
				continue
			}

			slog.Info("Message processed", "message", message.Message)
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			slog.Info("Stopping consumer")
			return nil
		}
	}
}

func StartConsumerGroup(service message.MessageService) {
	keepRunning := true

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	topics := strings.Split(os.Getenv("KAFKA_TOPIC"), ",")

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	handler := NewConsumerGroupHandler(service)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				slog.Error("Error from consumer", "error", err)
				break
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	slog.Info("Sarama consumer up")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("Terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("Terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		slog.Error("Error closing client", "error", err)
	}
}
