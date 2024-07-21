package message

import (
	"context"
	"time"
)

type service struct {
	MessageRepository
	timeout time.Duration
}

func NewMessageService(repository MessageRepository) MessageService {
	return &service{MessageRepository: repository, timeout: time.Duration(2) * time.Second}
}

func (s *service) createMessage(ctx context.Context, message *message) (*idRes, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer func() {
		cancel()
	}()
	id, err := s.MessageRepository.createMessage(newCtx, message)
	if err != nil {
		return nil, err
	}

	return id, nil
}
