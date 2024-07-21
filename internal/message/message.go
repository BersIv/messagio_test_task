package message

import "context"

type message struct {
	Message string `json:"message"`
}

type idRes struct {
	Id uint `json:"id"`
}

type MessageRepository interface {
	createMessage(ctx context.Context, message *message) (*idRes, error)
}

type MessageService interface {
	createMessage(ctx context.Context, message *message) (*idRes, error)
}
