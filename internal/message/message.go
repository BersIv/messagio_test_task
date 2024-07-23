package message

import (
	"context"
	"messagio_test_task/internal/models"
)

type messageReq struct {
	Message string `json:"message"`
}

type MessageRepository interface {
	createMessage(ctx context.Context, message *models.Message) error
	updateMessage(ctx context.Context, message *models.Message) error
	getPendingMessagesCount(ctx context.Context) (*uint, error)
	getProcessedMessagesCount(ctx context.Context) (*uint, error)
}

type MessageService interface {
	createMessage(ctx context.Context, message *models.Message) error
	UpdateMessage(ctx context.Context, message *models.Message) error
	getPendingMessagesCount(ctx context.Context) (*uint, error)
	getProcessedMessagesCount(ctx context.Context) (*uint, error)
}
