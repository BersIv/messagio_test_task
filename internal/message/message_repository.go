package message

import (
	"context"
	"database/sql"
	"log/slog"
	"messagio_test_task/internal/models"
)

type repository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &repository{db: db}
}

func (r *repository) createMessage(ctx context.Context, message *models.Message) error {
	query := "INSERT INTO messages(message) VALUES($1) RETURNING id"
	if err := r.db.QueryRowContext(ctx, query, message.Message).Scan(&message.Id); err != nil {
		slog.Error("Error inserting message", "error", err)
		return err
	}
	return nil
}

func (r *repository) updateMessage(ctx context.Context, message *models.Message) error {
	query := "UPDATE messages SET status = 'processed' WHERE id = $1"
	if _, err := r.db.ExecContext(ctx, query, message.Id); err != nil {
		slog.Error("Error updating message", "error", err)
		return err
	}
	return nil
}

func (r *repository) getPendingMessagesCount(ctx context.Context) (*uint, error) {
	var count uint
	query := "SELECT COUNT(*) FROM messages WHERE status = 'pending'"
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		slog.Error("Error getting pending messages count", "error", err)
		return nil, err
	}
	return &count, nil
}

func (r *repository) getProcessedMessagesCount(ctx context.Context) (*uint, error) {
	var count uint
	query := "SELECT COUNT(*) FROM messages WHERE status = 'processed'"
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		slog.Error("Error getting processed messages count", "error", err)
		return nil, err
	}
	return &count, nil
}
