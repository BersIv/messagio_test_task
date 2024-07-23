package message

import (
	"context"
	"messagio_test_task/internal/models"
	"time"
)

type Service struct {
	MessageRepository
	timeout time.Duration
}

func NewMessageService(repository MessageRepository) MessageService {
	return &Service{MessageRepository: repository, timeout: time.Duration(2) * time.Second}
}

func (s *Service) createMessage(ctx context.Context, message *models.Message) error {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer func() {
		cancel()
	}()
	if err := s.MessageRepository.createMessage(newCtx, message); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateMessage(ctx context.Context, message *models.Message) error {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer func() {
		cancel()
	}()
	if err := s.MessageRepository.updateMessage(newCtx, message); err != nil {
		return err
	}

	return nil
}

func (s *Service) getPendingMessagesCount(ctx context.Context) (*uint, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer func() {
		cancel()
	}()
	count, err := s.MessageRepository.getPendingMessagesCount(newCtx)
	if err != nil {
		return nil, err
	}

	return count, nil
}

func (s *Service) getProcessedMessagesCount(ctx context.Context) (*uint, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer func() {
		cancel()
	}()
	count, err := s.MessageRepository.getProcessedMessagesCount(newCtx)
	if err != nil {
		return nil, err
	}

	return count, nil
}
